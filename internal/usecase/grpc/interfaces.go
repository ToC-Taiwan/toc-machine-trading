package grpc

import (
	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
	"github.com/toc-taiwan/toc-trade-protobuf/golang/pb"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=grpc

type BasicgRPCAPI interface {
	CreateLongConnection() error
	Login() error
	CheckUsage() (*pb.ShioajiUsage, error)
	GetAllStockDetail() ([]*pb.StockDetailMessage, error)
	GetAllFutureDetail() ([]*pb.FutureDetailMessage, error)
	GetAllOptionDetail() ([]*pb.OptionDetailMessage, error)
}

type HistorygRPCAPI interface {
	GetStockHistoryTick(stockNumArr []string, date string) ([]*pb.HistoryTickMessage, error)
	GetStockHistoryKbar(stockNumArr []string, date string) ([]*pb.HistoryKbarMessage, error)
	GetStockHistoryClose(stockNumArr []string, date string) ([]*pb.HistoryCloseMessage, error)
	GetFutureHistoryKbar(codeArr []string, date string) (*pb.HistoryKbarResponse, error)
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
	GetStockVolumeRankPB(date string) (*pb.StockVolumeRankResponse, error)
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

type TradegRPCAPI interface {
	CancelOrder(orderID string) (*pb.TradeResult, error)

	BuyStock(order *entity.StockOrder) (*pb.TradeResult, error)
	SellStock(order *entity.StockOrder) (*pb.TradeResult, error)
	SellFirstStock(order *entity.StockOrder) (*pb.TradeResult, error)

	BuyOddStock(order *entity.StockOrder) (*pb.TradeResult, error)
	SellOddStock(order *entity.StockOrder) (*pb.TradeResult, error)

	BuyFuture(order *entity.FutureOrder) (*pb.TradeResult, error)
	SellFuture(order *entity.FutureOrder) (*pb.TradeResult, error)
	SellFirstFuture(order *entity.FutureOrder) (*pb.TradeResult, error)

	GetLocalOrderStatusArr() error
	GetSimulateOrderStatusArr() error

	GetFuturePosition() (*pb.FuturePositionArr, error)
	GetStockPosition() (*pb.StockPositionArr, error)
	GetSettlement() (*pb.SettlementList, error)
	GetAccountBalance() (*pb.AccountBalance, error)
	GetMargin() (*pb.Margin, error)
}
