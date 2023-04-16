// Package usecase package usecase
package usecase

import (
	"context"
	"time"

	"tmt/internal/entity"
)

//go:generate mockgen -source=interfaces_repo.go -destination=./mocks_repo_test.go -package=usecase_test

type BasicRepo interface {
	QueryAllStock(ctx context.Context) (map[string]*entity.Stock, error)
	InsertOrUpdatetStockArr(ctx context.Context, t []*entity.Stock) error
	UpdateAllStockDayTradeToNo(ctx context.Context) error

	QueryAllCalendar(ctx context.Context) (map[time.Time]*entity.CalendarDate, error)
	InsertOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error

	QueryAllFuture(ctx context.Context) (map[string]*entity.Future, error)
	QueryFutureByLikeName(ctx context.Context, name string) ([]*entity.Future, error)
	InsertOrUpdatetFutureArr(ctx context.Context, t []*entity.Future) error

	QueryAllOption(ctx context.Context) (map[string]*entity.Option, error)
	QueryOptionByLikeName(ctx context.Context, name string) ([]*entity.Option, error)
	InsertOrUpdatetOptionArr(ctx context.Context, t []*entity.Option) error
}

type TargetRepo interface {
	InsertOrUpdateTargetArr(ctx context.Context, t []*entity.StockTarget) error
	QueryTargetsByTradeDay(ctx context.Context, tradeDay time.Time) ([]*entity.StockTarget, error)
	QueryAllMXFFuture(ctx context.Context) ([]*entity.Future, error)
}

type HistoryRepo interface {
	InsertHistoryCloseArr(ctx context.Context, t []*entity.StockHistoryClose) error
	QueryMutltiStockCloseByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string]*entity.StockHistoryClose, error)
	InsertHistoryTickArr(ctx context.Context, t []*entity.StockHistoryTick) error
	QueryMultiStockTickArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.StockHistoryTick, error)
	InsertHistoryKbarArr(ctx context.Context, t []*entity.StockHistoryKbar) error
	QueryMultiStockKbarArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.StockHistoryKbar, error)
	InsertQuaterMA(ctx context.Context, t *entity.StockHistoryAnalyze) error
	QueryAllQuaterMAByStockNum(ctx context.Context, stockNum string) (map[time.Time]*entity.StockHistoryAnalyze, error)
	DeleteHistoryKbarByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error
	DeleteHistoryTickByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error
	DeleteHistoryCloseByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error
	InsertFutureHistoryTickArr(ctx context.Context, t []*entity.FutureHistoryTick) error
	QueryFutureHistoryTickArrByTime(ctx context.Context, code string, startTime, endTime time.Time) ([]*entity.FutureHistoryTick, error)
	InsertFutureHistoryClose(ctx context.Context, c *entity.FutureHistoryClose) error
	QueryFutureHistoryCloseByDate(ctx context.Context, code string, tradeDay time.Time) (*entity.FutureHistoryClose, error)
}

type RealTimeRepo interface {
	InsertEvent(ctx context.Context, t *entity.SinopacEvent) error
}

type TradeRepo interface {
	QueryAllStockTradeBalance(ctx context.Context) ([]*entity.StockTradeBalance, error)
	QueryLastStockTradeBalance(ctx context.Context) (*entity.StockTradeBalance, error)
	QueryStockTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.StockTradeBalance, error)
	InsertOrUpdateStockTradeBalance(ctx context.Context, t *entity.StockTradeBalance) error
	QueryAllFutureTradeBalance(ctx context.Context) ([]*entity.FutureTradeBalance, error)
	QueryLastFutureTradeBalance(ctx context.Context) (*entity.FutureTradeBalance, error)
	QueryFutureTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.FutureTradeBalance, error)
	InsertOrUpdateFutureTradeBalance(ctx context.Context, t *entity.FutureTradeBalance) error
	QueryStockOrderByID(ctx context.Context, orderID string) (*entity.StockOrder, error)
	InsertOrUpdateOrderByOrderID(ctx context.Context, t *entity.StockOrder) error
	QueryAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error)
	QueryAllStockOrderByDate(ctx context.Context, timeTange []time.Time) ([]*entity.StockOrder, error)
	QueryFutureOrderByID(ctx context.Context, orderID string) (*entity.FutureOrder, error)
	InsertOrUpdateFutureOrderByOrderID(ctx context.Context, t *entity.FutureOrder) error
	QueryAllFutureOrder(ctx context.Context) ([]*entity.FutureOrder, error)
	QueryAllFutureOrderByDate(ctx context.Context, timeTange []time.Time) ([]*entity.FutureOrder, error)
	QueryAllLastAccountBalance(ctx context.Context, bankIDArr []int) ([]*entity.AccountBalance, error)
	QueryAccountBalanceByDateAndBankID(ctx context.Context, date time.Time, bankID int) (*entity.AccountBalance, error)
	InsertOrUpdateAccountBalance(ctx context.Context, t *entity.AccountBalance) error
	QueryAccountSettlementByDate(ctx context.Context, date time.Time) (*entity.Settlement, error)
	InsertOrUpdateAccountSettlement(ctx context.Context, t *entity.Settlement) error
	QueryInventoryStockByDate(ctx context.Context, date time.Time) ([]*entity.InventoryStock, error)
	DeleteInventoryStockByDate(ctx context.Context, date time.Time) error
	InsertInventoryStock(ctx context.Context, t []*entity.InventoryStock) error
}
