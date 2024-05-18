package repo

import (
	"context"
	"time"

	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=repo

type BasicRepo interface {
	UpdateAllStockDayTradeToNo(ctx context.Context) error
	InsertOrUpdatetStockArr(ctx context.Context, t []*entity.Stock) error
	InsertOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error
	InsertOrUpdatetFutureArr(ctx context.Context, t []*entity.Future) error
	InsertOrUpdatetOptionArr(ctx context.Context, t []*entity.Option) error
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
}

type RealTimeRepo interface {
	InsertEvent(ctx context.Context, t *entity.SinopacEvent) error
}

type SystemRepo interface {
	EmailVerification(ctx context.Context, username string) error
	InsertUser(ctx context.Context, t *entity.NewUser) error
	QueryAllUser(ctx context.Context) ([]*entity.User, error)
	QueryUserByUsername(ctx context.Context, username string) (*entity.User, error)
	InsertOrUpdatePushToken(ctx context.Context, token, username string, enabled bool) error
	GetAllPushTokens(ctx context.Context) ([]string, error)
	GetPushToken(ctx context.Context, token string) (*entity.PushToken, error)
	DeleteAllPushTokens(ctx context.Context) error
	GetLastJWT(ctx context.Context) (string, error)
	InsertJWT(ctx context.Context, jwt string) error
}

type TargetRepo interface {
	InsertOrUpdateTargetArr(ctx context.Context, t []*entity.StockTarget) error
	QueryTargetsByTradeDay(ctx context.Context, tradeDay time.Time) ([]*entity.StockTarget, error)
}

type TradeRepo interface {
	QueryAllStockTradeBalance(ctx context.Context) ([]*entity.StockTradeBalance, error)
	InsertOrUpdateStockTradeBalance(ctx context.Context, t *entity.StockTradeBalance) error
	QueryAllFutureTradeBalance(ctx context.Context) ([]*entity.FutureTradeBalance, error)
	InsertOrUpdateFutureTradeBalance(ctx context.Context, t *entity.FutureTradeBalance) error
	InsertOrUpdateOrderByOrderID(ctx context.Context, t *entity.StockOrder) error
	QueryAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error)
	QueryAllStockOrderByDate(ctx context.Context, timeTange []time.Time) ([]*entity.StockOrder, error)
	InsertOrUpdateFutureOrderByOrderID(ctx context.Context, t *entity.FutureOrder) error
	QueryAllFutureOrder(ctx context.Context) ([]*entity.FutureOrder, error)
	QueryAllFutureOrderByDate(ctx context.Context, timeTange []time.Time) ([]*entity.FutureOrder, error)
	QueryLastAccountBalance(ctx context.Context) (*entity.AccountBalance, error)
	InsertOrUpdateAccountBalance(ctx context.Context, t *entity.AccountBalance) error
	InsertOrUpdateAccountSettlement(ctx context.Context, t *entity.Settlement) error
	QueryInventoryUUIDStockByDate(ctx context.Context, date time.Time) (map[string]string, error)
	InsertOrUpdateInventoryStock(ctx context.Context, t []*entity.InventoryStock) error
	ClearInventoryStockByUUID(ctx context.Context, uuid string) error
	QueryInventoryStockByDate(ctx context.Context, date time.Time) ([]*entity.InventoryStock, error)
}
