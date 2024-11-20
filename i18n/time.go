package i18n

import (
	"context"
	"time"

	"github.com/goodsign/monday"
)

type DateStyle string

const (
	DateStyleFull     DateStyle = "full"
	DateStyleLong     DateStyle = "long"
	DateStyleMedium   DateStyle = "medium"
	DateStyleShort    DateStyle = "short"
	DateStyleDateTime DateStyle = "datetime"
	DateStyleTime     DateStyle = "time"
)

var layout = map[string]map[DateStyle]string{
	"es": {
		"full":     monday.DefaultFormatEsESFull,
		"long":     monday.DefaultFormatEsESLong,
		"medium":   monday.DefaultFormatEsESMedium,
		"short":    monday.DefaultFormatEsESShort,
		"datetime": monday.DefaultFormatEsESDateTime,
		"time":     monday.DefaultFormatEsESTime,
	},
	"it": {
		"full":     monday.DefaultFormatItITFull,
		"long":     monday.DefaultFormatItITLong,
		"medium":   monday.DefaultFormatItITMedium,
		"short":    monday.DefaultFormatItITShort,
		"datetime": monday.DefaultFormatItITDateTime,
		"time":     monday.DefaultFormatItITTime,
	},
	"fr": {
		"full":     monday.DefaultFormatFrFRFull,
		"long":     monday.DefaultFormatFrFRLong,
		"medium":   monday.DefaultFormatFrFRMedium,
		"short":    monday.DefaultFormatFrFRShort,
		"datetime": monday.DefaultFormatFrFRDateTime,
		"time":     monday.DefaultFormatFrFRTime,
	},
	"en": {
		"full":     monday.DefaultFormatEnUSFull,
		"long":     monday.DefaultFormatEnUSLong,
		"medium":   monday.DefaultFormatEnUSMedium,
		"short":    monday.DefaultFormatEnUSShort,
		"datetime": monday.DefaultFormatEnUSDateTime,
		"time":     monday.DefaultFormatEnUSTime,
	},
}
var mondayLocale = map[string]string{
	"es": monday.LocaleEsES,
	"it": monday.LocaleItIT,
	"fr": monday.LocaleFrFR,
	"en": monday.LocaleEnUS,
}

func DateTime(ctx context.Context, t time.Time, style DateStyle) string {
	return DateTimeLocale(FromCtx(ctx), t, style)
}
func DateTimeLocale(locale string, t time.Time, style DateStyle) string {
	return monday.Format(t, layout[locale][style], monday.Locale(mondayLocale[locale]))
}
