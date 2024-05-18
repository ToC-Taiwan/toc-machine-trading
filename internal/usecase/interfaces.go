package usecase

import (
	"context"
	"time"

	"tmt/internal/entity"

	"github.com/toc-taiwan/toc-trade-protobuf/src/golang/pb"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase

type Analyze interface {
	GetRebornMap(ctx context.Context) map[time.Time][]entity.Stock
}

type Basic interface {
	GetStockDetail(stockNum string) *entity.Stock
	GetFutureDetail(futureCode string) *entity.Future
	GetShioajiUsage() (*entity.ShioajiUsage, error)
	CreateStockSearchRoom(com chan string, dataChan chan []*entity.Stock)
	CreateFutureSearchRoom(com chan string, dataChan chan []*entity.Future)
}

type History interface {
	GetDayKbarByStockNumMultiDate(stockNum string, date time.Time, interval int64) ([]*entity.StockHistoryKbar, error)
	GetFutureHistoryPBKbarByDate(code string, date time.Time) (*pb.HistoryKbarResponse, error)
}

type RealTime interface {
	GetStockSnapshotByNumArr(stockNumArr []string) ([]*entity.StockSnapShot, error)
	GetTradeIndex() *entity.TradeIndex
	GetFutureSnapshotByCode(code string) (*pb.SnapshotMessage, error)
	DeleteRealTimeClient(connectionID string)
	CreateRealTimePick(connectionID string, odd bool, com chan *pb.PickRealMap, tickChan chan []byte)
	CreateRealTimePickFuture(ctx context.Context, code string, tickChan chan *pb.FutureRealTimeTickMessage)
}

type Target interface {
	GetTargets(ctx context.Context) []*entity.StockTarget
	GetCurrentVolumeRank() (*pb.StockVolumeRankResponse, error)
}

type Trade interface {
	GetAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error)
	GetAllFutureOrder(ctx context.Context) ([]*entity.FutureOrder, error)
	GetAllStockTradeBalance(ctx context.Context) ([]*entity.StockTradeBalance, error)
	GetAllFutureTradeBalance(ctx context.Context) ([]*entity.FutureTradeBalance, error)
	GetFutureOrderByTradeDay(ctx context.Context, tradeDay string) ([]*entity.FutureOrder, error)
	BuyOddStock(num string, price float64, share int64) (string, entity.OrderStatus, error)
	SelloddStock(num string, price float64, share int64) (string, entity.OrderStatus, error)
	BuyFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error)
	SellFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error)
	CancelOrderByID(orderID string) (string, entity.OrderStatus, error)
	GetFuturePosition() ([]*entity.FuturePosition, error)
	IsFutureTradeTime() bool
	IsAuthUser(username string) bool
	GetLatestInventoryStock() ([]*entity.InventoryStock, error)
}

type System interface {
	AddUser(ctx context.Context, t *entity.NewUser) error
	InsertPushToken(ctx context.Context, token, username string, enabled bool) error
	Login(ctx context.Context, username, password string) error
	VerifyEmail(ctx context.Context, username, code string) error
	UpdateAuthTradeUser()
	DeleteAllPushTokens(ctx context.Context) error
	IsPushTokenEnabled(ctx context.Context, token string) (bool, error)
	GetUserInfo(ctx context.Context, username string) (*entity.User, error)
	GetLastJWT(ctx context.Context) (string, error)
	InsertJWT(ctx context.Context, jwt string) error
}

type FCM interface {
	AnnounceMessage(msg string) error
	PushNotification(title, msg string) error
}
