package utils

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// Translate validator error into target language
func TranslateValidatorError(
	errs validator.ValidationErrors, language string, enTrans, zhTrans ut.Translator,
) string {
	trans := enTrans
	if strings.HasPrefix(language, "zh") {
		trans = zhTrans
	}
	translatedMessages := []string{}
	for _, e := range errs {
		// can translate each error one at a time.
		translatedMessages = append(translatedMessages, e.Translate(trans))
	}
	return strings.Join(translatedMessages, ". ")
}
