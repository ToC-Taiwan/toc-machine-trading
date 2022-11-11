// Package app package app
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tmt/cmd/config"
	"tmt/internal/usecase"
	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/rabbit"
	"tmt/internal/usecase/repo"
	"tmt/pkg/grpc"
	"tmt/pkg/httpserver"
	"tmt/pkg/logger"
	"tmt/pkg/postgres"

	v1 "tmt/internal/controller/http/v1"

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

	sc, err := grpc.New(
		cfg.Sinopac.URL,
		grpc.MaxPoolSize(cfg.Sinopac.PoolMax),
	)
	if err != nil {
		log.Panic(err)
	}

	basicUseCase := usecase.NewBasic(repo.NewBasic(pg), grpcapi.NewBasic(sc))
	orderUseCase := usecase.NewOrder(grpcapi.NewOrder(sc), repo.NewOrder(pg))
	streamUseCase := usecase.NewStream(repo.NewStream(pg), grpcapi.NewStream(sc), rabbit.NewStream())
	analyzeUseCase := usecase.NewAnalyze(repo.NewHistory(pg))
	historyUseCase := usecase.NewHistory(repo.NewHistory(pg), grpcapi.NewHistory(sc))
	targetUseCase := usecase.NewTarget(repo.NewTarget(pg), grpcapi.NewTarget(sc), grpcapi.NewStream(sc))

	// HTTP Server
	handler := gin.New()
	r := v1.NewRouter(handler)
	{
		r.AddBasicRoutes(handler, basicUseCase)
		r.AddOrderRoutes(handler, orderUseCase)
		r.AddStreamRoutes(handler, streamUseCase, orderUseCase)
		r.AddAnalyzeRoutes(handler, analyzeUseCase)
		r.AddHistoryRoutes(handler, historyUseCase)
		r.AddTargetRoutes(handler, targetUseCase)
	}
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
