package main

import (
	"github.com/im-kulikov/doc-bot/misc"
	"github.com/im-kulikov/doc-bot/modules/bot"
	"github.com/im-kulikov/helium"
)

func main() {
	h, err := helium.New(&helium.Settings{
		File:         misc.Config,
		Name:         misc.Name,
		Prefix:       misc.Prefix,
		BuildTime:    misc.Build,
		BuildVersion: misc.Version,
	}, bot.Module)
	helium.CatchTrace(err)
	helium.CatchTrace(h.Run())
}
