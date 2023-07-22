// Package app package app
package app

import (
	"os"
	"os/signal"
	"syscall"

	"tmt/cmd/config"
	"tmt/internal/controller/http/router"
	"tmt/internal/usecase"
	"tmt/pkg/httpserver"
	"tmt/pkg/log"
)

var logger = log.Get()

func RunApp() {
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
	<-interrupt
	basic.LogoutAll()
	cfg.Close()
}
