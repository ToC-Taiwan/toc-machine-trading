// Package app package app
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	v1 "toc-machine-trading/internal/controller/http/v1"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/httpserver"
	"toc-machine-trading/pkg/logger"
	"toc-machine-trading/pkg/postgres"

	"github.com/gin-gonic/gin"
)

// Run -.
func Run(cfg *config.Config) {
	dbPath := fmt.Sprintf("%s%s", cfg.PG.URL, cfg.PG.DBName)

	pg, err := postgres.New(dbPath, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		logger.Get().Panic(err)
	}
	defer pg.Close()

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler)
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

	// stuck := make(chan struct{})
	// sinopac.HealthCheck()
	// logger.New("INOF").Warn("message string")
	// go sinopac.EventChannel()
	// go sinopac.TickChannel()
	// go sinopac.BidAskChannel()
	// // time.Sleep(time.Second * 3)
	// sinopac.SubscribeStockTick()
	// sinopac.UnSubscribeStockTick()
	// sinopac.SubscribeStockBidAsk()
	// sinopac.UnSubscribeStockBidAsk()
	// sinopac.GetAllStockDetail()
	// sinopac.GetAllSnapshot()
	// sinopac.GetStockSnapshotByStockNumArr()
	// sinopac.GetAllSnapshotTSE()
	// sinopac.GetStockHistoryTick()
	// sinopac.GetStockHistoryKbar()
	// sinopac.GetStockHistoryClose()
	// sinopac.GetStockVolumeRank()
	// sinopac.GetStockTSEHistoryTick()
	// sinopac.GetStockTSEHistoryClose()
	// <-stuck
}
