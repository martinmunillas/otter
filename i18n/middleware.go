package i18n

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/martinmunillas/otter/utils"
)

const cookieName = "otter-lang"

type localeKeyType string

var localeKey localeKeyType = "locale"

func parseAcceptLanguage(header string) []string {
	langs := strings.Split(header, ",")

	for i, l := range langs {
		if strings.Contains(l, ";") {
			langs[i] = strings.Split(l, ";")[0]
		}
	}
	return langs
}

func Middleware(next http.Handler) http.Handler {
	if len(supportedLocales) == 0 {
		utils.Throw("invalid i18n middleware initialization, before initializing the middleware make sure to add your locales with i18n.AddLocale()")
	}
	defaultLocale := defaultLocale
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/set-locale" {
			locale := r.FormValue("locale")
			w.Header().Add("HX-Refresh", "true")
			SetLocale(w, r, locale)
			return
		}
		locale := "INVALID"
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			header := r.Header["Accept-Language"]
			if len(header) == 0 {
				locale = defaultLocale
			} else {
				langs := parseAcceptLanguage(header[0])

				for _, lang := range langs {
					if lang == "*" {
						locale = defaultLocale
						break
					}

					for _, l := range supportedLocales {
						if lang == l {
							locale = lang
						}
					}
				}

				if locale == "INVALID" {
					locale = defaultLocale
				}
			}
		} else {
			lang := cookie.Value
			if lang == "*" {
				locale = defaultLocale
			}

			for _, l := range supportedLocales {
				if lang == l {
					locale = lang
				}
			}
		}

		ctx := context.WithValue(r.Context(), localeKey, locale)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetLocale(w http.ResponseWriter, r *http.Request, lang string) {
	cookie := http.Cookie{
		Name:  cookieName,
		Value: lang,
	}

	if err := cookie.Valid(); err != nil {
		slog.Error(err.Error())
	}
	http.SetCookie(w, &cookie)
}

func FromCtx(ctx context.Context) string {
	l := ctx.Value(localeKey)
	if l == nil {
		return defaultLocale
	}
	locale, ok := l.(string)
	if !ok {
		return defaultLocale
	}
	return locale
}
