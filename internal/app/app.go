// Package app package app
package app

import (
	"os"
	"os/signal"
	"syscall"

	"tmt/cmd/config"
	"tmt/internal/controller/http/router"
	"tmt/internal/usecase/cache"
	"tmt/internal/usecase/cases/analyze"
	"tmt/internal/usecase/cases/basic"
	"tmt/internal/usecase/cases/history"
	"tmt/internal/usecase/cases/realtime"
	"tmt/internal/usecase/cases/target"
	"tmt/internal/usecase/cases/trade"
	"tmt/pkg/eventbus"
	"tmt/pkg/httpserver"
	"tmt/pkg/log"
)

func RunApp() {
	logger := log.Get()
	cfg := config.Get()

	cc := cache.New()
	bus := eventbus.New()

	logger.Warn("TMT is running")
	logger.Warnf("Simulation Mode: %v", cfg.Simulation)

	// Do not adjust the order
	basic := basic.NewBasic().Init(logger, cc)
	trade := trade.NewTrade().Init(logger, bus)
	analyze := analyze.NewAnalyze().Init(logger, cc, bus)
	history := history.NewHistory().Init(logger, cc, bus)
	realTime := realtime.NewRealTime().Init(logger, cc, bus)
	target := target.NewTarget().Init(logger, cc, bus)

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
