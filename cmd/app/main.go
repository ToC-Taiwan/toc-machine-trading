package main

import (
	"tmt/internal/app"
)

func main() {
	app.InitDB()
	app.RunApp()
}
