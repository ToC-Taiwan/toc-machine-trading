package main

import (
	"tmt/internal/app"
	"tmt/pkg/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	app.MigrateDB(cfg)
	app.Run(cfg)
}
