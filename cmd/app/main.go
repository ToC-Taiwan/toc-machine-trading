package main

import (
	"os"
	"path/filepath"

	"tmt/internal/app"
	"tmt/internal/config"

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
	config.Init()
	app.InitDB()
	app.SetupCronJob()
	app.Run()
}
