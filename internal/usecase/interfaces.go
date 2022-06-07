// Package usecase package usecase
package usecase

import (
	"context"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pb"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test
type (
	// Basic -.
	Basic interface {
		GetAllSinopacStockAndUpdateRepo(ctx context.Context) ([]*entity.Stock, error)
		GetAllRepoStock(ctx context.Context) ([]*entity.Stock, error)
	}

	// BasicRepo -.
	BasicRepo interface {
		QueryAllStock(ctx context.Context) ([]*entity.Stock, error)
		InserOrUpdatetStockArr(ctx context.Context, t []*entity.Stock) error
		QueryAllCalendar(ctx context.Context) ([]*entity.CalendarDate, error)
		InserOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error
	}

	// BasicgRPCAPI -.
	BasicgRPCAPI interface {
		GetServerToken() (string, error)
		GetAllStockDetail() ([]*pb.StockDetailMessage, error)
		GetAllStockSnapshot() ([]*pb.StockSnapshotMessage, error)
		GetStockSnapshotByNumArr(stockNumArr []string) ([]*pb.StockSnapshotMessage, error)
		GetStockSnapshotTSE() ([]*pb.StockSnapshotMessage, error)
	}
)

type (
	// Target -.
	Target interface{}

	// TargetRepo -.
	TargetRepo interface {
		InsertTargetArr(ctx context.Context, t []*entity.Target) error
		QueryTargetsByTradeDay(ctx context.Context, tradeDay time.Time) ([]*entity.Target, error)
	}

	// TargetRPCAPI -.
	TargetRPCAPI interface {
		GetStockVolumeRank(date string) ([]*pb.StockVolumeRankMessage, error)
		SubscribeStockTick(stockNumArr []string) ([]string, error)
		UnSubscribeStockTick(stockNumArr []string) ([]string, error)
		SubscribeStockBidAsk(stockNumArr []string) ([]string, error)
		UnSubscribeStockBidAsk(stockNumArr []string) ([]string, error)
	}
)

type (
	// History -.
	History interface{}

	// HistoryRepo -.
	HistoryRepo interface{}

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
		ReceiveEvent(ctx context.Context) error
		ReceiveTicks(ctx context.Context) error
		ReceiveBidAsk(ctx context.Context) error
		ReceiveOrderStatus(ctx context.Context) error
	}

	// StreamRepo -.
	StreamRepo interface {
		InsertEvent(ctx context.Context, t *entity.SinopacEvent) error
	}

	// StreamgRPCAPI -.
	StreamgRPCAPI interface {
		EventChannel(eventChan chan *entity.SinopacEvent) error
		BidAskChannel(bidAskChan chan *entity.RealTimeBidAsk) error
		TickChannel(tickChan chan *entity.RealTimeTick) error
		OrderStatusChannel(orderStatusChan chan *entity.OrderStatus) error
	}
)

type (
	// Order -.
	Order interface{}

	// OrderRepo -.
	OrderRepo interface{}

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
