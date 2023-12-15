package main

import (
	"os"
	"path/filepath"

	"tmt/internal/app"
	cfg "tmt/internal/config"

	"github.com/joho/godotenv"
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
