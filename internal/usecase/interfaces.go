// Package usecase package usecase
package usecase

import (
	"context"
	"time"

	"tmt/internal/entity"
	"tmt/pb"
	"tmt/pkg/eventbus"
	"tmt/pkg/logger"
)

var (
	log = logger.Get()
	bus = eventbus.New()
	cc  = NewCache()
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
		GetTargets(ctx context.Context) []*entity.Target
	}

	// TargetRepo -.
	TargetRepo interface {
		InsertOrUpdateTargetArr(ctx context.Context, t []*entity.Target) error
		QueryTargetsByTradeDay(ctx context.Context, tradeDay time.Time) ([]*entity.Target, error)
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
	}
)

type (
	// History -.
	History interface {
		GetTradeDay() time.Time
		GetDayKbarByStockNumDate(stockNum string, date time.Time) *entity.HistoryKbar
	}

	// HistoryRepo -.
	HistoryRepo interface {
		InsertHistoryCloseArr(ctx context.Context, t []*entity.HistoryClose) error
		QueryMutltiStockCloseByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string]*entity.HistoryClose, error)
		InsertHistoryTickArr(ctx context.Context, t []*entity.HistoryTick) error
		QueryMultiStockTickArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.HistoryTick, error)
		InsertHistoryKbarArr(ctx context.Context, t []*entity.HistoryKbar) error
		QueryMultiStockKbarArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.HistoryKbar, error)
		InsertQuaterMA(ctx context.Context, t *entity.HistoryAnalyze) error
		QueryAllQuaterMAByStockNum(ctx context.Context, stockNum string) (map[time.Time]*entity.HistoryAnalyze, error)
		DeleteHistoryKbarByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error
		DeleteHistoryTickByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error
		DeleteHistoryCloseByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error
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
		ReceiveStreamData(ctx context.Context, targetArr []*entity.Target)
		GetTSESnapshot(ctx context.Context) (*entity.StockSnapShot, error)
		GetStockSnapshotByNumArr(stockNumArr []string) ([]*entity.StockSnapShot, error)
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

		GetFutureSnapshotByCodeArr(codeArr []string) (*pb.SnapshotResponse, error)
	}

	// StreamRabbit -.
	StreamRabbit interface {
		FillAllBasic(allStockMap map[string]*entity.Stock, allFutureMap map[string]*entity.Future)

		EventConsumer(eventChan chan *entity.SinopacEvent)
		OrderStatusConsumer(orderStatusChan chan interface{})
		TickConsumer(stockNum string, tickChan chan *entity.RealTimeTick)
		BidAskConsumer(stockNum string, bidAskChan chan *entity.RealTimeBidAsk)

		FutureTickConsumer(code string, tickChan chan *entity.RealTimeFutureTick)
	}
)

type (
	// Order -.
	Order interface {
		GetAllOrder(ctx context.Context) ([]*entity.StockOrder, error)
		GetAllTradeBalance(ctx context.Context) ([]*entity.TradeBalance, error)
		CalculateBuyCost(price float64, quantity int64) int64
		CalculateSellCost(price float64, quantity int64) int64
		CalculateTradeDiscount(price float64, quantity int64) int64

		AskOrderUpdate() error
	}

	// OrderRepo -.
	OrderRepo interface {
		QueryAllStockTradeBalance(ctx context.Context) ([]*entity.TradeBalance, error)
		QueryStockTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.TradeBalance, error)
		InsertOrUpdateStockTradeBalance(ctx context.Context, t *entity.TradeBalance) error
		QueryFutureTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.TradeBalance, error)
		InsertOrUpdateFutureTradeBalance(ctx context.Context, t *entity.TradeBalance) error

		QueryStockOrderByID(ctx context.Context, orderID string) (*entity.StockOrder, error)
		InsertOrUpdateOrderByOrderID(ctx context.Context, t *entity.StockOrder) error
		QueryAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error)
		QueryAllStockOrderByDate(ctx context.Context, date time.Time) ([]*entity.StockOrder, error)

		QueryFutureOrderByID(ctx context.Context, orderID string) (*entity.FutureOrder, error)
		InsertOrUpdateFutureOrderByOrderID(ctx context.Context, t *entity.FutureOrder) error
		QueryAllFutureOrderByDate(ctx context.Context, date time.Time) ([]*entity.FutureOrder, error)
	}

	// OrdergRPCAPI -.
	OrdergRPCAPI interface {
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
		SimulateOnHistoryTick(ctx context.Context, useDefault bool)
	}
)
