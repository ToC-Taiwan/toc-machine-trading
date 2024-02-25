package usecase

import (
	"context"
	"time"

	"tmt/internal/entity"
	"tmt/pb"
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

type BasicRepo interface {
	UpdateAllStockDayTradeToNo(ctx context.Context) error
	InsertOrUpdatetStockArr(ctx context.Context, t []*entity.Stock) error
	InsertOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error
	InsertOrUpdatetFutureArr(ctx context.Context, t []*entity.Future) error
	InsertOrUpdatetOptionArr(ctx context.Context, t []*entity.Option) error
}

type BasicgRPCAPI interface {
	CreateLongConnection() error
	Login() error
	CheckUsage() (*pb.ShioajiUsage, error)
	GetAllStockDetail() ([]*pb.StockDetailMessage, error)
	GetAllFutureDetail() ([]*pb.FutureDetailMessage, error)
	GetAllOptionDetail() ([]*pb.OptionDetailMessage, error)
}

type History interface {
	GetDayKbarByStockNumMultiDate(stockNum string, date time.Time, interval int64) ([]*entity.StockHistoryKbar, error)
	GetFutureHistoryPBKbarByDate(code string, date time.Time) (*pb.HistoryKbarResponse, error)
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

type HistorygRPCAPI interface {
	GetStockHistoryTick(stockNumArr []string, date string) ([]*pb.HistoryTickMessage, error)
	GetStockHistoryKbar(stockNumArr []string, date string) ([]*pb.HistoryKbarMessage, error)
	GetStockHistoryClose(stockNumArr []string, date string) ([]*pb.HistoryCloseMessage, error)
	GetFutureHistoryKbar(codeArr []string, date string) (*pb.HistoryKbarResponse, error)
}

type RealTime interface {
	GetStockSnapshotByNumArr(stockNumArr []string) ([]*entity.StockSnapShot, error)
	GetTradeIndex() *entity.TradeIndex
	GetFutureSnapshotByCode(code string) (*pb.SnapshotMessage, error)
	DeleteRealTimeClient(connectionID string)
	CreateRealTimePick(connectionID string, odd bool, com chan *pb.PickRealMap, tickChan chan []byte)
	CreateRealTimePickFuture(ctx context.Context, code string, tickChan chan *pb.FutureRealTimeTickMessage)
}

type RealTimeRepo interface {
	InsertEvent(ctx context.Context, t *entity.SinopacEvent) error
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

type Rabbit interface {
	EventConsumer(eventChan chan *entity.SinopacEvent)
	OrderStatusConsumer(orderStatusChan chan interface{})
	OrderStatusArrConsumer(orderStatusChan chan interface{})
	StockTickPbConsumer(ctx context.Context, stockNum string, tickChan chan []byte)
	StockTickOddsPbConsumer(ctx context.Context, stockNum string, tickChan chan []byte)
	FutureTickConsumer(code string, tickChan chan *entity.RealTimeFutureTick)
	FutureTickPbConsumer(ctx context.Context, code string, tickChan chan *pb.FutureRealTimeTickMessage)
	Close()
}

type Target interface {
	GetTargets(ctx context.Context) []*entity.StockTarget
	GetCurrentVolumeRank() (*pb.StockVolumeRankResponse, error)
}

type TargetRepo interface {
	InsertOrUpdateTargetArr(ctx context.Context, t []*entity.StockTarget) error
	QueryTargetsByTradeDay(ctx context.Context, tradeDay time.Time) ([]*entity.StockTarget, error)
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

type FCM interface {
	AnnounceMessage(msg string) error
	PushNotification(title, msg string) error
}
