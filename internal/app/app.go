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

func RunApp(cfg *config.Config) {
	base := usecase.NewUseCaseBase(cfg)

	// Do not adjust the order
	basicUseCase := base.NewBasic()
	tradeUseCase := base.NewTrade()
	analyzeUseCase := base.NewAnalyze()
	historyUseCase := base.NewHistory()
	realTimeUseCase := base.NewRealTime()
	targetUseCase := base.NewTarget()

	// HTTP Server
	r := router.NewRouter().
		AddV1BasicRoutes(basicUseCase).
		AddV1TradeRoutes(tradeUseCase).
		AddV1RealTimeRoutes(realTimeUseCase, tradeUseCase, historyUseCase).
		AddV1AnalyzeRoutes(analyzeUseCase).
		AddV1HistoryRoutes(historyUseCase).
		AddV1TargetRoutes(targetUseCase)

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
	base.Close()
}
