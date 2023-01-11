package main

import (
	"tmt/cmd/config"
	"tmt/internal/app"
)

func main() {
	cfg := config.GetConfig()

	app.InitDB(cfg.Database)
	app.RunApp(cfg)
}
