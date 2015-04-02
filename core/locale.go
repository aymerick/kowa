package core

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n"
)

var (
	Langs []string
)

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
