package i18n

var defaultLocale string

// SetDefault changes the default locale, by default the default locale is the first one added
func SetDefault(locale string) {
	defaultLocale = locale
}
