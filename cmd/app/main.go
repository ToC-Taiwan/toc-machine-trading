package main

import (
	"os"
	"path/filepath"

	"toc-machine-trading/internal/app"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/global"
)

var (
	cfg *config.Config
	err error
)

func init() {
	// Configuration
	cfg, err = config.GetConfig()
	if err != nil {
		panic(err)
	}

	// get binary path
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	global.SetBasePath(filepath.Clean(filepath.Dir(ex)))

	// check if env is production or development
	if cfg.Deployment != "prod" {
		global.SetIsDevelopment(true)
	}
}

func main() {
	app.MigrateDB(cfg)
	app.Run(cfg)
}
