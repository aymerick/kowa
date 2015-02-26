package main

import (
	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/commands"
)

func main() {
	loadLocales()

	commands.InitConf()
	commands.Execute()
}

// load i18n locales
func loadLocales() {
	i18n.MustLoadTranslationFile("./locales/en.json")
	i18n.MustLoadTranslationFile("./locales/fr.json")
}
