// Package app package app
package app

import (
	"os"
	"os/signal"
	"syscall"

	"tmt/internal/config"
	"tmt/internal/controller/http/router"
	"tmt/internal/usecase"
	"tmt/pkg/eventbus"
	"tmt/pkg/httpserver"
	"tmt/pkg/log"
)

func RunApp() {
	logger := log.Get()
	cfg := config.Get()

	bus := eventbus.New()
	cc := usecase.New()

	logger.Warn("TMT is running")
	logger.Warnf("Simulation Mode: %v", cfg.Simulation)

	// Do not adjust the order
	basic := usecase.NewBasic().Init(logger, cc)
	trade := usecase.NewTrade().Init(logger, bus)
	analyze := usecase.NewAnalyze().Init(logger, cc, bus)
	history := usecase.NewHistory().Init(logger, cc, bus)
	realTime := usecase.NewRealTime().Init(logger, cc, bus)
	target := usecase.NewTarget().Init(logger, cc, bus)

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
	cfg.CloseDB()
}
