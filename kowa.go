package main

import (
	"github.com/aymerick/kowa/commands"
	"github.com/aymerick/kowa/core"
)

func main() {
	core.LoadLocales()

	commands.InitConf()
	commands.Execute()
}
