package acceptlanguage

// Key to use when setting the accept language.
type ctxKeyAcceptLanguage string

const (
	AcceptLanguageKey ctxKeyAcceptLanguage = "acceptLanguage"
)

const (
	IFTELangSubtagEnglish = "en"
	IFTELangSubtagRussian = "ru"
	IFTELangSubtagKazakh  = "kk"
)

// Validate - валидация сабтега языка если это не Русский, Казахский, Английский
// то вернётся сабтег русского языка.
func Validate(subtag string) string {
	switch subtag {
	default:
		fallthrough
	case IFTELangSubtagRussian:
		return IFTELangSubtagRussian
	case IFTELangSubtagKazakh:
		return IFTELangSubtagKazakh
	case IFTELangSubtagEnglish:
		return IFTELangSubtagEnglish
	}
}
