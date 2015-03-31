package main

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/commands"
	"github.com/aymerick/kowa/core"
)

func main() {
	loadLocales()

	commands.InitConf()
	commands.Execute()
}

// load i18n locales
func loadLocales() {
	for _, lang := range core.Langs {
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
