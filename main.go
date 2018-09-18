package main

import (
	"github.com/im-kulikov/doc-bot/app"
	"github.com/im-kulikov/helium"
)

var (
	BuildTime    = "now"
	BuildVersion = "dev"
)

func main() {
	h, err := helium.New(&helium.Settings{
		File:         "config.yml",
		Name:         "bot",
		BuildTime:    BuildTime,
		BuildVersion: BuildVersion,
	}, app.Module)
	helium.Catch(err)
	helium.Catch(h.Run())
}
