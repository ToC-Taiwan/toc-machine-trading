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
	usecase.StartHealthCheck()
	// Do not adjust the order
	basicUseCase := usecase.NewBasic()
	tradeUseCase := usecase.NewTrade()
	analyzeUseCase := usecase.NewAnalyze()
	historyUseCase := usecase.NewHistory()
	realTimeUseCase := usecase.NewRealTime()
	targetUseCase := usecase.NewTarget()

	// HTTP Server
	r := router.NewRouter().
		AddV1BasicRoutes(basicUseCase).
		AddV1TradeRoutes(tradeUseCase).
		AddV1RealTimeRoutes(realTimeUseCase, tradeUseCase, historyUseCase).
		AddV1AnalyzeRoutes(analyzeUseCase).
		AddV1HistoryRoutes(historyUseCase).
		AddV1TargetRoutes(targetUseCase)

	cfg := config.Get()
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
	cfg.Close()
}
