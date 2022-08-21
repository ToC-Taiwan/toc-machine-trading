// Package usecase package usecase
package usecase

import (
	"context"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pb"
	"toc-machine-trading/pkg/eventbus"
	"toc-machine-trading/pkg/logger"
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
	}

	// BasicgRPCAPI -.
	BasicgRPCAPI interface {
		Heartbeat() error
		Terminate() error
		GetAllStockDetail() ([]*pb.StockDetailMessage, error)
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
		UnSubscribeStockAllTick() (*pb.FunctionErr, error)
		SubscribeStockBidAsk(stockNumArr []string) ([]string, error)
		UnSubscribeStockBidAsk(stockNumArr []string) ([]string, error)
		UnSubscribeStockAllBidAsk() (*pb.FunctionErr, error)
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
		GetStockHistoryTick(stockNumArr []string, date string) ([]*pb.StockHistoryTickMessage, error)
		GetStockHistoryKbar(stockNumArr []string, date string) ([]*pb.StockHistoryKbarMessage, error)
		GetStockHistoryClose(stockNumArr []string, date string) ([]*pb.StockHistoryCloseMessage, error)
		GetStockHistoryCloseByDateArr(stockNumArr []string, date []string) ([]*pb.StockHistoryCloseMessage, error)
		GetStockTSEHistoryTick(date string) ([]*pb.StockHistoryTickMessage, error)
		GetStockTSEHistoryKbar(date string) ([]*pb.StockHistoryKbarMessage, error)
		GetStockTSEHistoryClose(date string) ([]*pb.StockHistoryCloseMessage, error)
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
		GetAllStockSnapshot() ([]*pb.StockSnapshotMessage, error)
		GetStockSnapshotByNumArr(stockNumArr []string) ([]*pb.StockSnapshotMessage, error)
		GetStockSnapshotTSE() (*pb.StockSnapshotMessage, error)
		GetFutureSnapshotFIMTX() (*pb.StockSnapshotMessage, error)
	}

	// StreamRabbit -.
	StreamRabbit interface {
		EventConsumer(eventChan chan *entity.SinopacEvent)
		OrderStatusConsumer(orderStatusChan chan *entity.Order)
		TickConsumer(stockNum string, tickChan chan *entity.RealTimeTick)
		BidAskConsumer(stockNum string, bidAskChan chan *entity.RealTimeBidAsk)
	}
)

type (
	// Order -.
	Order interface {
		GetAllOrder(ctx context.Context) ([]*entity.Order, error)
		GetAllTradeBalance(ctx context.Context) ([]*entity.TradeBalance, error)
		CalculateBuyCost(price float64, quantity int64) int64
		CalculateSellCost(price float64, quantity int64) int64
		CalculateTradeDiscount(price float64, quantity int64) int64

		AskOrderUpdate() error
	}

	// OrderRepo -.
	OrderRepo interface {
		InsertOrUpdateTradeBalance(ctx context.Context, t *entity.TradeBalance) error
		QueryTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.TradeBalance, error)
		InsertOrUpdateOrderByOrderID(ctx context.Context, t *entity.Order) error
		QueryOrderByID(ctx context.Context, orderID string) (*entity.Order, error)
		QueryAllOrderByDate(ctx context.Context, date time.Time) ([]*entity.Order, error)
		QueryAllOrder(ctx context.Context) ([]*entity.Order, error)
		QueryAllTradeBalance(ctx context.Context) ([]*entity.TradeBalance, error)
	}

	// OrdergRPCAPI -.
	OrdergRPCAPI interface {
		BuyStock(order *entity.Order, sim bool) (*pb.TradeResult, error)
		SellStock(order *entity.Order, sim bool) (*pb.TradeResult, error)
		SellFirstStock(order *entity.Order, sim bool) (*pb.TradeResult, error)
		CancelStock(orderID string, sim bool) (*pb.TradeResult, error)
		GetOrderStatusByID(orderID string, sim bool) (*pb.TradeResult, error)
		GetOrderStatusArr() ([]*pb.StockOrderStatus, error)
		GetNonBlockOrderStatusArr() (*pb.FunctionErr, error)
	}
)

type (
	// Analyze -.
	Analyze interface {
		GetRebornMap(ctx context.Context) map[time.Time][]entity.Stock
		SimulateOnHistoryTick(ctx context.Context, useDefault bool)
	}
)
