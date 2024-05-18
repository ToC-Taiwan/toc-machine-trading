package main

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/toc-taiwan/toc-machine-trading/internal/app"
	cfg "github.com/toc-taiwan/toc-machine-trading/internal/config"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	if err := godotenv.Load(filepath.Join(filepath.Dir(ex), ".env")); err != nil {
		panic(err)
	}
	cfg.Init()
	app.Run()
}
