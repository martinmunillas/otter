package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/a-h/templ"
	"github.com/martinmunillas/otter/utils"
)

// https://github.com/opral/monorepo/blob/main/inlang/source-code/plugins/t-function-matcher/marketplace-manifest.json
func flattenJson(input map[string]interface{}) (map[string]string, error) {
	flatMap := make(map[string]string)

	var flatten func(map[string]interface{}, string) error
	flatten = func(data map[string]interface{}, prefix string) error {
		for key, value := range data {
			fullKey := key
			if prefix != "" {
				fullKey = prefix + "." + key
			}

			// Type switch to handle nested maps
			switch v := value.(type) {
			case map[string]interface{}:
				err := flatten(v, fullKey)
				if err != nil {
					return err
				}
			case string:
				flatMap[fullKey] = v
			default:
				return fmt.Errorf("invalid translation %s of type %t", key, v)
			}
		}
		return nil
	}

	err := flatten(input, "")
	return flatMap, err
}

var translations = make(map[string]map[string]string, 2)
var supportedLocales = make([]string, 0, 2)

func processLang(r io.Reader) (map[string]string, error) {
	m := map[string]interface{}{}
	err := json.NewDecoder(r).Decode(&m)
	if err != nil {
		return nil, err
	}
	translation, err := flattenJson(m)
	if err != nil {
		return nil, err
	}

	return translation, nil

}
func AddLocale(locale string, r io.Reader) {
	translation, err := processLang(r)
	if err != nil {
		utils.Throw(err.Error())
	}
	supportedLocales = append(supportedLocales, locale)
	translations[locale] = translation
	if defaultLocale == "" {
		defaultLocale = locale
	}

}

type Replacements = map[string]any

func strChunk(str string, raw bool, _ ...error) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		var err error
		if raw {
			_, err = io.WriteString(w, str)
		} else {
			_, err = io.WriteString(w, templ.EscapeString(str))
		}
		return err
	})
}

func chunksRender(chunks []templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, chunk := range chunks {
			err := chunk.Render(ctx, w)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func errorThrower(err error) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return err
	})
}

func t(ctx context.Context, key string, raw bool, replacements ...Replacements) templ.Component {
	str := Translation(ctx, key)
	if len(replacements) == 0 || str == key {
		return strChunk(str, raw)
	}
	if len(replacements) > 1 {
		return errorThrower(fmt.Errorf("invalid translation \"%s\" call: more than one replacements map provided", key))
	}
	runes := []rune(str)

	chunks := make([]templ.Component, 0, 1)

	currentStr := ""
	currentVarName := ""
	isCollectingVarName := false
	for i, c := range runes {
		isEscaped := i > 0 && runes[i-1] == '\\'
		if c == '{' && !isEscaped {
			if isCollectingVarName {
				return errorThrower(fmt.Errorf("invalid translation \"%s\" format: opening variable before closing previous", key))
			}
			if currentStr != "" {
				chunks = append(chunks, strChunk(currentStr, raw))
				currentStr = ""
			}
			isCollectingVarName = true
			continue
		}
		if c == '}' && !isEscaped {
			if !isCollectingVarName {
				return errorThrower(fmt.Errorf("invalid translation \"%s\" format: closing variable before opening one", key))
			}
			if currentVarName == "" {
				return errorThrower(fmt.Errorf("invalid translation \"%s\" format: missing variable name between {}", key))
			}
			val, ok := replacements[0][currentVarName]
			if !ok {
				return errorThrower(fmt.Errorf("invalid translation \"%s\" call: missing variable \"%s\" value", key, currentVarName))
			}
			switch v := val.(type) {
			case templ.Component:
				chunks = append(chunks, v)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				chunks = append(chunks, strChunk(fmt.Sprintf("%d", v), raw))
			case string, []rune, []byte:
				chunks = append(chunks, strChunk(fmt.Sprintf("%s", v), raw))
			default:
				return errorThrower(fmt.Errorf("invalid translation \"%s\" call: variable \"%s\" of type %t not supported", key, currentVarName, v))

			}
			currentVarName = ""
			isCollectingVarName = false
			continue
		}

		if isCollectingVarName {
			currentVarName += string(c)
		} else {
			currentStr += string(c)
		}
	}
	if currentStr != "" {
		chunks = append(chunks, strChunk(currentStr, raw))
	}
	return chunksRender(chunks)
}

func T(ctx context.Context, key string, replacements ...Replacements) templ.Component {
	return t(ctx, key, false, replacements...)
}

func RawT(ctx context.Context, key string, replacements ...Replacements) templ.Component {
	return t(ctx, key, true, replacements...)
}

func Translation(ctx context.Context, key string) string {
	locale := FromCtx(ctx)
	content := translations[locale][key]
	if content == "" {
		return key
	}
	return content
}
