// Package usecase package usecase
package usecase

import (
	"context"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase/modules/cache"
	"tmt/internal/usecase/modules/event"
	"tmt/internal/usecase/modules/trader"

	"tmt/pb"
	"tmt/pkg/logger"
)

var (
	log = logger.Get()
	cc  = cache.Get()
	bus = event.Get()
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test
type (
	// Basic -.
	Basic interface {
		GetAllSinopacStockAndUpdateRepo(ctx context.Context) ([]*entity.Stock, error)
		GetAllRepoStock(ctx context.Context) ([]*entity.Stock, error)
		TerminateSinopac(ctx context.Context) error
	}

	// BasicRepo -.
	BasicRepo interface {
		QueryAllStock(ctx context.Context) (map[string]*entity.Stock, error)
		InsertOrUpdatetStockArr(ctx context.Context, t []*entity.Stock) error
		QueryAllCalendar(ctx context.Context) (map[time.Time]*entity.CalendarDate, error)
		InsertOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error
		UpdateAllStockDayTradeToNo(ctx context.Context) error

		InsertOrUpdatetFutureArr(ctx context.Context, t []*entity.Future) error
		QueryAllFuture(ctx context.Context) (map[string]*entity.Future, error)
		QueryAllMXFFuture(ctx context.Context) ([]*entity.Future, error)
	}

	// BasicgRPCAPI -.
	BasicgRPCAPI interface {
		Heartbeat() error
		Terminate() error
		GetAllStockDetail() ([]*pb.StockDetailMessage, error)
		GetAllFutureDetail() ([]*pb.FutureDetailMessage, error)
	}
)

type (
	// Target -.
	Target interface {
		GetTargets(ctx context.Context) []*entity.StockTarget
	}

	// TargetRepo -.
	TargetRepo interface {
		InsertOrUpdateTargetArr(ctx context.Context, t []*entity.StockTarget) error
		QueryTargetsByTradeDay(ctx context.Context, tradeDay time.Time) ([]*entity.StockTarget, error)
	}

	// TargetgRPCAPI -.
	TargetgRPCAPI interface {
		GetStockVolumeRank(date string) ([]*pb.StockVolumeRankMessage, error)
		SubscribeStockTick(stockNumArr []string) ([]string, error)
		UnSubscribeStockTick(stockNumArr []string) ([]string, error)
		UnSubscribeStockAllTick() (*pb.ErrorMessage, error)
		SubscribeStockBidAsk(stockNumArr []string) ([]string, error)
		UnSubscribeStockBidAsk(stockNumArr []string) ([]string, error)
		UnSubscribeStockAllBidAsk() (*pb.ErrorMessage, error)

		SubscribeFutureTick(codeArr []string) ([]string, error)
		UnSubscribeFutureTick(codeArr []string) ([]string, error)
		SubscribeFutureBidAsk(codeArr []string) ([]string, error)
		UnSubscribeFutureBidAsk(codeArr []string) ([]string, error)
	}
)

type (
	// History -.
	History interface {
		GetTradeDay() time.Time
		GetDayKbarByStockNumDate(stockNum string, date time.Time) *entity.StockHistoryKbar

		GetFutureTradeCond(days int) trader.TradeBalance
	}

	// HistoryRepo -.
	HistoryRepo interface {
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

	// HistorygRPCAPI -.
	HistorygRPCAPI interface {
		GetStockHistoryTick(stockNumArr []string, date string) ([]*pb.HistoryTickMessage, error)
		GetStockHistoryKbar(stockNumArr []string, date string) ([]*pb.HistoryKbarMessage, error)
		GetStockHistoryClose(stockNumArr []string, date string) ([]*pb.HistoryCloseMessage, error)
		GetStockHistoryCloseByDateArr(stockNumArr []string, date []string) ([]*pb.HistoryCloseMessage, error)

		GetStockTSEHistoryTick(date string) ([]*pb.HistoryTickMessage, error)
		GetStockTSEHistoryKbar(date string) ([]*pb.HistoryKbarMessage, error)
		GetStockTSEHistoryClose(date string) ([]*pb.HistoryCloseMessage, error)

		GetFutureHistoryTick(codeArr []string, date string) ([]*pb.HistoryTickMessage, error)
		GetFutureHistoryKbar(codeArr []string, date string) ([]*pb.HistoryKbarMessage, error)
		GetFutureHistoryClose(codeArr []string, date string) ([]*pb.HistoryCloseMessage, error)
	}
)

type (
	// Stream -.
	Stream interface {
		ReceiveEvent(ctx context.Context)
		ReceiveOrderStatus(ctx context.Context)
		ReceiveStreamData(ctx context.Context, targetArr []*entity.StockTarget)
		GetStockSnapshotByNumArr(stockNumArr []string) ([]*entity.StockSnapShot, error)

		GetTSESnapshot(ctx context.Context) (*entity.StockSnapShot, error)
		GetOTCSnapshot(ctx context.Context) (*entity.StockSnapShot, error)

		GetNasdaqClose() (*entity.YahooPrice, error)
		GetNasdaqFutureClose() (*entity.YahooPrice, error)

		GetMainFutureCode() string
		GetFutureSnapshotByCode(code string) (*entity.FutureSnapShot, error)

		NewFutureRealTimeConnection(tickChan chan *entity.RealTimeFutureTick, connectionID string)
		DeleteFutureRealTimeConnection(connectionID string)
		NewOrderStatusConnection(orderStatusChan chan interface{}, connectionID string)
		DeleteOrderStatusConnection(connectionID string)
	}

	// StreamRepo -.
	StreamRepo interface {
		InsertEvent(ctx context.Context, t *entity.SinopacEvent) error
	}

	// StreamgRPCAPI -.
	StreamgRPCAPI interface {
		GetAllStockSnapshot() ([]*pb.SnapshotMessage, error)
		GetStockSnapshotByNumArr(stockNumArr []string) ([]*pb.SnapshotMessage, error)
		GetStockSnapshotTSE() (*pb.SnapshotMessage, error)
		GetStockSnapshotOTC() (*pb.SnapshotMessage, error)

		GetNasdaq() (*pb.YahooFinancePrice, error)
		GetNasdaqFuture() (*pb.YahooFinancePrice, error)

		GetFutureSnapshotByCode(code string) (*pb.SnapshotMessage, error)
	}

	// StreamRabbit -.
	StreamRabbit interface {
		FillAllBasic(allStockMap map[string]*entity.Stock, allFutureMap map[string]*entity.Future)

		EventConsumer(eventChan chan *entity.SinopacEvent)
		OrderStatusConsumer(orderStatusChan chan interface{})
		TickConsumer(stockNum string, tickChan chan *entity.RealTimeStockTick)
		StockBidAskConsumer(stockNum string, bidAskChan chan *entity.RealTimeStockBidAsk)

		FutureTickConsumer(code string, tickChan chan *entity.RealTimeFutureTick)
		FutureBidAskConsumer(code string, bidAskChan chan *entity.FutureRealTimeBidAsk)

		AddFutureTickChan(tickChan chan *entity.RealTimeFutureTick, connectionID string)
		RemoveFutureTickChan(connectionID string)
		AddOrderStatusChan(orderStatusChan chan interface{}, connectionID string)
		RemoveOrderStatusChan(connectionID string)
	}
)

type (
	// Order -.
	Order interface {
		GetAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error)
		GetAllFutureOrder(ctx context.Context) ([]*entity.FutureOrder, error)

		GetAllStockTradeBalance(ctx context.Context) ([]*entity.StockTradeBalance, error)
		GetAllFutureTradeBalance(ctx context.Context) ([]*entity.FutureTradeBalance, error)

		CalculateBuyCost(price float64, quantity int64) int64
		CalculateSellCost(price float64, quantity int64) int64
		CalculateTradeDiscount(price float64, quantity int64) int64

		AskOrderUpdate() error

		BuyFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error)
		SellFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error)
		SellFirstFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error)
		BuyLaterFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error)
		CancelFutureOrderID(orderID string) (string, entity.OrderStatus, error)

		GetFutureOrderStatusByID(orderID string) (*entity.FutureOrder, error)

		GetFuturePosition() ([]*entity.FuturePosition, error)
		IsFutureTradeTime() bool
	}

	// OrderRepo -.
	OrderRepo interface {
		QueryAllStockTradeBalance(ctx context.Context) ([]*entity.StockTradeBalance, error)
		QueryStockTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.StockTradeBalance, error)
		InsertOrUpdateStockTradeBalance(ctx context.Context, t *entity.StockTradeBalance) error

		QueryAllFutureTradeBalance(ctx context.Context) ([]*entity.FutureTradeBalance, error)
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
	}

	// OrdergRPCAPI -.
	OrdergRPCAPI interface {
		GetFuturePosition() (*pb.FuturePositionArr, error)

		BuyStock(order *entity.StockOrder, sim bool) (*pb.TradeResult, error)
		SellStock(order *entity.StockOrder, sim bool) (*pb.TradeResult, error)
		SellFirstStock(order *entity.StockOrder, sim bool) (*pb.TradeResult, error)
		CancelStock(orderID string, sim bool) (*pb.TradeResult, error)
		GetOrderStatusByID(orderID string, sim bool) (*pb.TradeResult, error)
		GetOrderStatusArr() ([]*pb.StockOrderStatus, error)
		GetNonBlockOrderStatusArr() (*pb.ErrorMessage, error)

		BuyFuture(order *entity.FutureOrder, sim bool) (*pb.TradeResult, error)
		SellFuture(order *entity.FutureOrder, sim bool) (*pb.TradeResult, error)
		SellFirstFuture(order *entity.FutureOrder, sim bool) (*pb.TradeResult, error)
		CancelFuture(orderID string, sim bool) (*pb.TradeResult, error)
	}
)

type (
	// Analyze -.
	Analyze interface {
		GetRebornMap(ctx context.Context) map[time.Time][]entity.Stock
	}
)
