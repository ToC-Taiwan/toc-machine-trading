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
	"toc-machine-trading/internal/usecase/rabbit"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/httpserver"
	"toc-machine-trading/pkg/logger"
	"toc-machine-trading/pkg/postgres"
	"toc-machine-trading/pkg/sinopac"

	"github.com/gin-gonic/gin"
)

var log = logger.Get()

// Run -.
func Run(cfg *config.Config) {
	pg, err := postgres.New(
		fmt.Sprintf("%s%s", cfg.Postgres.URL, cfg.Postgres.DBName),
		postgres.MaxPoolSize(cfg.Postgres.PoolMax),
	)
	if err != nil {
		log.Panic(err)
	}
	defer pg.Close()

	sc, err := sinopac.New(
		cfg.Sinopac.URL,
		sinopac.MaxPoolSize(cfg.Sinopac.PoolMax),
	)
	if err != nil {
		log.Panic(err)
	}

	// basic
	basicUseCase := usecase.NewBasic(repo.NewBasic(pg), grpcapi.NewBasic(sc))

	// order
	usecase.NewOrder(grpcapi.NewOrder(sc))

	// stream
	usecase.NewStream(repo.NewStream(pg), rabbit.NewStream())

	// analyze
	analyzeUseCase := usecase.NewAnalyze(repo.NewHistory(pg))

	// history
	usecase.NewHistory(repo.NewHistory(pg), grpcapi.NewHistory(sc))

	// target
	usecase.NewTarget(repo.NewTarget(pg), grpcapi.NewTarget(sc))

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, basicUseCase, analyzeUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		log.Info(s.String())
	case err = <-httpServer.Notify():
		log.Error(err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(err)
	}
}
