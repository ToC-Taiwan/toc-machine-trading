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
	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/mq"
	"tmt/internal/usecase/repo"
	"tmt/pkg/grpc"
	"tmt/pkg/httpserver"
	"tmt/pkg/log"
	"tmt/pkg/postgres"

	"github.com/gin-gonic/gin"
)

var logger = log.Get()

type app struct {
	pg *postgres.Postgres
	sc *grpc.Connection
	fg *grpc.Connection
}

func newApp(cfg *config.Config) *app {
	pg, err := postgres.New(
		fmt.Sprintf("%s%s", cfg.Database.URL, cfg.Database.DBName),
		postgres.MaxPoolSize(cfg.Database.PoolMax),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connecting to sinopac gRPC server")
	sc, err := grpc.New(
		cfg.Sinopac.URL,
		grpc.MaxPoolSize(cfg.Sinopac.PoolMax),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connecting to fugle gRPC server")
	fg, err := grpc.New(
		cfg.Fugle.URL,
		grpc.MaxPoolSize(cfg.Fugle.PoolMax),
	)
	if err != nil {
		logger.Fatal(err)
	}

	return &app{
		pg: pg,
		sc: sc,
		fg: fg,
	}
}

func RunApp(cfg *config.Config) {
	app := newApp(cfg)
	defer app.pg.Close()

	basicUseCase := usecase.NewBasic(repo.NewBasic(app.pg), grpcapi.NewBasic(app.sc), grpcapi.NewBasic(app.fg))
	tradeUseCase := usecase.NewTrade(grpcapi.NewTrade(app.sc), grpcapi.NewTrade(app.fg), repo.NewTrade(app.pg))
	analyzeUseCase := usecase.NewAnalyze(repo.NewHistory(app.pg))
	historyUseCase := usecase.NewHistory(repo.NewHistory(app.pg), grpcapi.NewHistory(app.sc))
	realTimeUseCase := usecase.NewRealTime(repo.NewRealTime(app.pg), grpcapi.NewRealTime(app.sc), grpcapi.NewSubscribe(app.sc), mq.NewRabbit())
	targetUseCase := usecase.NewTarget(repo.NewTarget(app.pg), grpcapi.NewRealTime(app.sc))

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
