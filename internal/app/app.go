// Package app package app
package app

import (
	"os"
	"os/signal"
	"syscall"

	"tmt/internal/config"
	"tmt/internal/controller/http/router"
	"tmt/internal/usecase"
	"tmt/pkg/httpserver"
	"tmt/pkg/log"

	"github.com/robfig/cron/v3"
)

func Run() {
	logger := log.Get()
	cfg := config.Get()

	logger.Warn("TMT is running")
	logger.Warnf("Simulation Mode: %v", cfg.Simulation)

	// Do not adjust the order
	basic := usecase.NewBasic()
	trade := usecase.NewTrade()
	analyze := usecase.NewAnalyze()
	history := usecase.NewHistory()
	realTime := usecase.NewRealTime()
	target := usecase.NewTarget()

	// HTTP Server
	r := router.NewRouter().
		AddV1BasicRoutes(basic).
		AddV1TradeRoutes(trade).
		AddV1RealTimeRoutes(realTime, trade, history).
		AddV1AnalyzeRoutes(analyze).
		AddV1HistoryRoutes(history).
		AddV1TargetRoutes(target)

	if e := httpserver.New(
		r.GetHandler(),
		httpserver.Port(cfg.Server.HTTP),
		httpserver.AddLogger(logger),
	).Start(); e != nil {
		logger.Fatalf("API Server error: %s", e)
	}

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	setupCronJob(interrupt)
	<-interrupt

	// post process
	cfg.CloseDB()
	logger.Warn("TMT is shutting down")
}

func setupCronJob(interrupt chan os.Signal) {
	exit := func() {
		interrupt <- os.Interrupt
	}

	job := cron.New()
	if _, e := job.AddFunc("20 8 * * *", exit); e != nil {
		panic(e)
	}
	if _, e := job.AddFunc("40 14 * * *", exit); e != nil {
		panic(e)
	}

	go job.Start()
}
