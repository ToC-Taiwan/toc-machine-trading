package realtime

import (
	"context"

	"tmt/internal/entity"
	"tmt/internal/usecase/cache"
	"tmt/pb"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=realtime_test

type RealTime interface {
	Init(logger *log.Log, cc *cache.Cache, bus *eventbus.Bus) RealTime
	GetStockSnapshotByNumArr(stockNumArr []string) ([]*entity.StockSnapShot, error)
	GetTradeIndex() *entity.TradeIndex
	GetTSESnapshot(ctx context.Context) (*entity.StockSnapShot, error)
	GetOTCSnapshot(ctx context.Context) (*entity.StockSnapShot, error)

	GetMainFuture() *entity.Future
	GetFutureSnapshotByCode(code string) (*entity.FutureSnapShot, error)

	NewFutureRealTimeClient(tickChan chan *entity.RealTimeFutureTick, orderStatusChan chan interface{}, connectionID string)
	DeleteFutureRealTimeClient(connectionID string)
}

type RealTimeRepo interface {
	InsertEvent(ctx context.Context, t *entity.SinopacEvent) error
}

type RealTimegRPCAPI interface {
	GetAllStockSnapshot() ([]*pb.SnapshotMessage, error)
	GetStockSnapshotByNumArr(stockNumArr []string) ([]*pb.SnapshotMessage, error)
	GetStockSnapshotTSE() (*pb.SnapshotMessage, error)
	GetStockSnapshotOTC() (*pb.SnapshotMessage, error)
	GetNasdaq() (*pb.YahooFinancePrice, error)
	GetNasdaqFuture() (*pb.YahooFinancePrice, error)
	GetStockVolumeRank(date string) ([]*pb.StockVolumeRankMessage, error)
	GetFutureSnapshotByCode(code string) (*pb.SnapshotMessage, error)
}

type SubscribegRPCAPI interface {
	SubscribeStockTick(stockNumArr []string, odd bool) ([]string, error)
	UnSubscribeStockTick(stockNumArr []string) ([]string, error)
	UnSubscribeAllTick() (*pb.ErrorMessage, error)
	SubscribeStockBidAsk(stockNumArr []string) ([]string, error)
	UnSubscribeStockBidAsk(stockNumArr []string) ([]string, error)
	UnSubscribeAllBidAsk() (*pb.ErrorMessage, error)
	SubscribeFutureTick(codeArr []string) ([]string, error)
	UnSubscribeFutureTick(codeArr []string) ([]string, error)
	SubscribeFutureBidAsk(codeArr []string) ([]string, error)
	UnSubscribeFutureBidAsk(codeArr []string) ([]string, error)
}

type Rabbit interface {
	FillAllBasic(allStockMap map[string]*entity.Stock, allFutureMap map[string]*entity.Future)

	EventConsumer(eventChan chan *entity.SinopacEvent)
	OrderStatusConsumer(orderStatusChan chan interface{})
	OrderStatusArrConsumer(orderStatusChan chan interface{})
	StockTickConsumer(stockNum string, tickChan chan *entity.RealTimeStockTick)
	StockBidAskConsumer(stockNum string, bidAskChan chan *entity.RealTimeStockBidAsk)
	FutureTickConsumer(code string, tickChan chan *entity.RealTimeFutureTick)
	FutureBidAskConsumer(code string, bidAskChan chan *entity.FutureRealTimeBidAsk)

	Close()
}
