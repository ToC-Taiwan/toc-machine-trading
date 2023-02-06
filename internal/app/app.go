// Package app package app
package app

import (
	"os"
	"os/signal"
	"syscall"

	"tmt/cmd/config"
	v1 "tmt/internal/controller/http/v1"
	"tmt/internal/usecase"
	"tmt/pkg/httpserver"
	"tmt/pkg/log"

	"github.com/gin-gonic/gin"
)

var logger = log.Get()

func RunApp(cfg *config.Config) {
	base := usecase.NewUseCaseBase(cfg)
	defer base.Close()

	logger.Infof("Simulation Mode: %v", cfg.Simulation)

	// Do not adjust the order
	basicUseCase := base.NewBasic()
	tradeUseCase := base.NewTrade()
	analyzeUseCase := base.NewAnalyze()
	historyUseCase := base.NewHistory()
	realTimeUseCase := base.NewRealTime()
	targetUseCase := base.NewTarget()

	// HTTP Server
	handler := gin.New()
	r := v1.NewRouter(handler)
	{
		r.AddBasicRoutes(handler, basicUseCase)
		r.AddTradeRoutes(handler, tradeUseCase)
		r.AddRealTimeRoutes(handler, realTimeUseCase, tradeUseCase, historyUseCase)
		r.AddAnalyzeRoutes(handler, analyzeUseCase)
		r.AddHistoryRoutes(handler, historyUseCase)
		r.AddTargetRoutes(handler, targetUseCase)
	}
	httpServer := httpserver.New(handler, httpserver.Port(cfg.Server.HTTP))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		logger.Info(s.String())
	case err := <-httpServer.Notify():
		logger.Error(err)
	}

	// Shutdown
	if err := httpServer.Shutdown(); err != nil {
		logger.Error(err)
	}
}
