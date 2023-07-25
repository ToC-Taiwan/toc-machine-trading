package trade

import (
	"context"
	"time"

	"tmt/internal/entity"
	"tmt/pb"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=trade_test

type Trade interface {
	Init(logger *log.Log, bus *eventbus.Bus) Trade

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

type TradegRPCAPI interface {
	BuyStock(order *entity.StockOrder) (*pb.TradeResult, error)
	SellStock(order *entity.StockOrder) (*pb.TradeResult, error)
	BuyOddStock(order *entity.StockOrder) (*pb.TradeResult, error)
	SellOddStock(order *entity.StockOrder) (*pb.TradeResult, error)
	SellFirstStock(order *entity.StockOrder) (*pb.TradeResult, error)
	CancelStock(orderID string) (*pb.TradeResult, error)

	BuyFuture(order *entity.FutureOrder) (*pb.TradeResult, error)
	SellFuture(order *entity.FutureOrder) (*pb.TradeResult, error)
	SellFirstFuture(order *entity.FutureOrder) (*pb.TradeResult, error)
	CancelFuture(orderID string) (*pb.TradeResult, error)

	GetOrderStatusByID(orderID string) (*pb.TradeResult, error)
	GetLocalOrderStatusArr() error
	GetSimulateOrderStatusArr() error

	GetNonBlockOrderStatusArr() (*pb.ErrorMessage, error)
	GetFuturePosition() (*pb.FuturePositionArr, error)
	GetStockPosition() (*pb.StockPositionArr, error)
	GetSettlement() (*pb.SettlementList, error)
	GetAccountBalance() (*pb.AccountBalance, error)
	GetMargin() (*pb.Margin, error)
}
