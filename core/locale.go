package core

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n"
)

var (
	Langs []string
)

// sugar to write core.P{} instead of map[string]interface{} for i18n parameters
type P map[string]interface{}

type TranslateFunc func(translationID string, args ...interface{}) string

func init() {
	Langs = []string{"en", "fr"}
}

// load i18n locales
func LoadLocales() {
	for _, lang := range Langs {
		filePath := fmt.Sprintf("locales/%s.json", lang)

		// fetch file from embedded assets
		data, err := Asset(filePath)
		if err != nil {
			panic("Failed to load translation files for language: " + lang)
		}

		// load translations
		i18n.ParseTranslationFileBytes(filePath, data)
	}
}

// Our own wrapper around i18n.MustTfunc that converts core.P arguments into map[string]interface{} to please i18n lib
func MustTfunc(lang string) TranslateFunc {
	f := i18n.MustTfunc(lang)

	return TranslateFunc(func(translationID string, args ...interface{}) string {
		var convArgs []interface{}

		if len(args) > 0 {
			arg0, arg0isP := args[0].(P)
			if arg0isP {
				convArgs = append(convArgs, map[string]interface{}(arg0))
			} else {
				convArgs = append(convArgs, args[0])
			}

			if len(args) > 1 {
				arg1, arg1isP := args[1].(P)
				if arg1isP {
					convArgs = append(convArgs, map[string]interface{}(arg1))
				} else {
					convArgs = append(convArgs, args[1])
				}
			}
		}

		return f(translationID, convArgs...)
	})
}
