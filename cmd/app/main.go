package main

import (
	"toc-machine-trading/internal/app"
	"toc-machine-trading/pkg/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	app.MigrateDB(cfg)
	app.Run(cfg)
}
