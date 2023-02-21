// Package usecase package usecase
package usecase

import (
	"tmt/internal/entity"

	"tmt/pb"
)

//go:generate mockgen -source=interfaces_grpcapi.go -destination=./mocks_grpcapi_test.go -package=usecase_test

type BasicgRPCAPI interface {
	Heartbeat() error
	Terminate() error
	GetAllStockDetail() ([]*pb.StockDetailMessage, error)
	GetAllFutureDetail() ([]*pb.FutureDetailMessage, error)
	GetAllOptionDetail() ([]*pb.OptionDetailMessage, error)
}

type SubscribegRPCAPI interface {
	SubscribeStockTick(stockNumArr []string) ([]string, error)
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

type HistorygRPCAPI interface {
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

type TradegRPCAPI interface {
	BuyStock(order *entity.StockOrder) (*pb.TradeResult, error)
	SellStock(order *entity.StockOrder) (*pb.TradeResult, error)
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

	NotifyToSlack(message string) error
}
