package app

import (
	"os"

	"github.com/robfig/cron/v3"
)

func SetupCronJob() {
	job := cron.New()
	if _, e := job.AddFunc("20 8 * * *", exit); e != nil {
		panic(e)
	}
	if _, e := job.AddFunc("40 14 * * *", exit); e != nil {
		panic(e)
	}
	job.Start()
}

func exit() {
	os.Exit(0)
}
