// Package app package app
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	v1 "toc-machine-trading/internal/controller/http/v1"
	"toc-machine-trading/internal/usecase"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/eventbus"
	"toc-machine-trading/pkg/httpserver"
	"toc-machine-trading/pkg/logger"
	"toc-machine-trading/pkg/postgres"
	"toc-machine-trading/pkg/sinopac"

	"github.com/gin-gonic/gin"
)

// Run -.
func Run(cfg *config.Config) {
	dbPath := fmt.Sprintf("%s%s", cfg.Postgres.URL, cfg.Postgres.DBName)
	pg, err := postgres.New(dbPath, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
	if err != nil {
		logger.Get().Panic(err)
	}
	defer pg.Close()

	sc, err := sinopac.New(cfg.Sinopac.URL, sinopac.MaxPoolSize(cfg.Sinopac.PoolMax))
	if err != nil {
		logger.Get().Panic(err)
	}

	bus := eventbus.New()
	basicUseCase := usecase.NewBasic(repo.New(pg), grpcapi.New(sc), bus)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, basicUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		logger.Get().Info(s.String())
	case err = <-httpServer.Notify():
		logger.Get().Error(err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Get().Error(err)
	}
}
