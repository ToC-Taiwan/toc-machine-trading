// Package usecase package usecase
package usecase

import (
	"context"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/module/simulator"
)

//go:generate mockgen -source=interfaces_usecase.go -destination=./mocks_usecase_test.go -package=usecase_test

type Basic interface {
	GetAllRepoStock(ctx context.Context) ([]*entity.Stock, error)
	TerminateSinopac() error
	TerminateFugle() error
	GetConfig() *config.Config
}

type Target interface {
	GetTargets(ctx context.Context) []*entity.StockTarget
}

type History interface {
	GetTradeDay() time.Time

	GetDayKbarByStockNumDate(stockNum string, date time.Time) *entity.StockHistoryKbar
	FetchFutureHistoryKbar(code string, date time.Time) ([]*entity.FutureHistoryKbar, error)

	SimulateOne(cond *config.TradeFuture) *simulator.SimulateBalance
	SimulateMulti(cond []*config.TradeFuture)
}

type RealTime interface {
	GetStockSnapshotByNumArr(stockNumArr []string) ([]*entity.StockSnapShot, error)
	GetTradeIndex() *entity.TradeIndex
	GetTSESnapshot(ctx context.Context) (*entity.StockSnapShot, error)
	GetOTCSnapshot(ctx context.Context) (*entity.StockSnapShot, error)

	GetMainFuture() *entity.Future
	GetFutureSnapshotByCode(code string) (*entity.FutureSnapShot, error)

	NewFutureRealTimeClient(tickChan chan *entity.RealTimeFutureTick, orderStatusChan chan interface{}, connectionID string)
	DeleteFutureRealTimeClient(connectionID string)
}

type Trade interface {
	CalculateBuyCost(price float64, quantity int64) int64
	CalculateSellCost(price float64, quantity int64) int64
	CalculateTradeDiscount(price float64, quantity int64) int64

	GetAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error)
	GetAllFutureOrder(ctx context.Context) ([]*entity.FutureOrder, error)
	GetAllStockTradeBalance(ctx context.Context) ([]*entity.StockTradeBalance, error)
	GetAllFutureTradeBalance(ctx context.Context) ([]*entity.FutureTradeBalance, error)
	GetLastStockTradeBalance(ctx context.Context) (*entity.StockTradeBalance, error)
	GetLastFutureTradeBalance(ctx context.Context) (*entity.FutureTradeBalance, error)

	GetFutureOrderByTradeDay(ctx context.Context, tradeDay string) ([]*entity.FutureOrder, error)

	BuyFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error)
	SellFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error)
	CancelFutureOrderByID(orderID string) (string, entity.OrderStatus, error)

	GetFuturePosition() ([]*entity.FuturePosition, error)
	IsFutureTradeTime() bool

	ManualInsertFutureOrder(ctx context.Context, order *entity.FutureOrder) error
	UpdateTradeBalanceByTradeDay(ctx context.Context, date string) error
	MoveStockOrderToLatestTradeDay(ctx context.Context, orderID string) error
	MoveFutureOrderToLatestTradeDay(ctx context.Context, orderID string) error

	GetAccountBalance(ctx context.Context) ([]*entity.AccountBalance, error)
}

type Analyze interface {
	GetRebornMap(ctx context.Context) map[time.Time][]entity.Stock
}
