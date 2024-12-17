package handlers

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/innovay-software/famapp-main/app/utils"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	zhTrans  ut.Translator
	enTrans  ut.Translator
)

func HandlerInit() {
	validate = validator.New()
	en := en.New()
	zh := zh.New()
	uni = ut.New(en, zh)
	zhTrans, _ = uni.GetTranslator("zh")
	zh_translations.RegisterDefaultTranslations(validate, zhTrans)
	enTrans, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, enTrans)

	utils.LogSuccess("Handlers.Init completed")
}
