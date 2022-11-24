package main

import (
	"tmt/internal/app"
	"tmt/internal/usecase/modules/config"
)

func main() {
	cfg := config.GetConfig()

	app.MigrateDB(cfg)
	app.Run(cfg)
}
