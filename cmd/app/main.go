package main

import (
	"tmt/cmd/config"
	"tmt/internal/app"
)

func main() {
	cfg := config.GetConfig()

	app.MigrateDB(cfg)
	app.Run(cfg)
}
