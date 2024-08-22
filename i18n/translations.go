package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

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

func T(ctx context.Context, key string) string {
	locale := FromCtx(ctx)
	content := translations[locale][key]
	if content == "" {
		return key
	}
	return content
}
