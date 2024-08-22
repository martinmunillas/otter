package i18n

var defaultLocale string

// Set Defaults changes the default locale, the default locale is the first one added by default
func SetDefault(locale string) {
	defaultLocale = locale
}
