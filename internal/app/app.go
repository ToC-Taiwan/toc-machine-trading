// Package app package app
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tmt/cmd/config"
	v1 "tmt/internal/controller/http/v1"
	"tmt/internal/usecase"
	"tmt/internal/usecase/grpcapi/sinopac"
	"tmt/internal/usecase/rabbit"
	"tmt/internal/usecase/repo"
	"tmt/pkg/grpc"
	"tmt/pkg/httpserver"
	"tmt/pkg/log"
	"tmt/pkg/postgres"

	"github.com/gin-gonic/gin"
)

var logger = log.Get()

func RunApp(cfg *config.Config) {
	pg, err := postgres.New(
		fmt.Sprintf("%s%s", cfg.Database.URL, cfg.Database.DBName),
		postgres.MaxPoolSize(cfg.Database.PoolMax),
	)
	if err != nil {
		logger.Panic(err)
	}
	defer pg.Close()

	sc, err := grpc.New(
		cfg.Sinopac.URL,
		grpc.MaxPoolSize(cfg.Sinopac.PoolMax),
	)
	if err != nil {
		logger.Panic(err)
	}

	basicUseCase := usecase.NewBasic(repo.NewBasic(pg), sinopac.NewBasic(sc))
	orderUseCase := usecase.NewOrder(sinopac.NewOrder(sc), repo.NewOrder(pg))
	streamUseCase := usecase.NewStream(repo.NewStream(pg), sinopac.NewStream(sc), rabbit.NewStream())
	analyzeUseCase := usecase.NewAnalyze(repo.NewHistory(pg))
	historyUseCase := usecase.NewHistory(repo.NewHistory(pg), sinopac.NewHistory(sc))
	targetUseCase := usecase.NewTarget(repo.NewTarget(pg), sinopac.NewTarget(sc), sinopac.NewStream(sc))

	// HTTP Server
	handler := gin.New()
	r := v1.NewRouter(handler)
	{
		r.AddBasicRoutes(handler, basicUseCase)
		r.AddOrderRoutes(handler, orderUseCase)
		r.AddStreamRoutes(handler, streamUseCase, orderUseCase, historyUseCase)
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
	case err = <-httpServer.Notify():
		logger.Error(err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(err)
	}
}
