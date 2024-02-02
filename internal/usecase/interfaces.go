package usecase

import (
	"context"
	"time"

	"tmt/internal/config"
	"tmt/internal/entity"
	"tmt/pb"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase

type Analyze interface {
	GetRebornMap(ctx context.Context) map[time.Time][]entity.Stock
}

type Basic interface {
	GetStockDetail(stockNum string) *entity.Stock
	GetFutureDetail(code string) *entity.Future
	GetConfig() *config.Config
	GetShioajiUsage() (*entity.ShioajiUsage, error)
}

type BasicRepo interface {
	QueryAllStock(ctx context.Context) (map[string]*entity.Stock, error)
	InsertOrUpdatetStockArr(ctx context.Context, t []*entity.Stock) error
	UpdateAllStockDayTradeToNo(ctx context.Context) error

	QueryAllCalendar(ctx context.Context) (map[time.Time]*entity.CalendarDate, error)
	InsertOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error

	QueryAllFuture(ctx context.Context) (map[string]*entity.Future, error)
	QueryFutureByLikeName(ctx context.Context, name string) ([]*entity.Future, error)
	InsertOrUpdatetFutureArr(ctx context.Context, t []*entity.Future) error

	QueryAllOption(ctx context.Context) (map[string]*entity.Option, error)
	QueryOptionByLikeName(ctx context.Context, name string) ([]*entity.Option, error)
	InsertOrUpdatetOptionArr(ctx context.Context, t []*entity.Option) error
}

type BasicgRPCAPI interface {
	CreateLongConnection() error
	Terminate() error
	Login() error
	CheckUsage() (*pb.ShioajiUsage, error)
	GetAllStockDetail() ([]*pb.StockDetailMessage, error)
	GetAllFutureDetail() ([]*pb.FutureDetailMessage, error)
	GetAllOptionDetail() ([]*pb.OptionDetailMessage, error)
}

type History interface {
	GetTradeDay() time.Time

	GetDayKbarByStockNumDate(stockNum string, date time.Time) *entity.StockHistoryKbar
	GetDayKbarByStockNumMultiDate(stockNum string, date time.Time, interval int64) ([]*entity.StockHistoryKbar, error)
	FetchFutureHistoryKbar(code string, date time.Time) ([]*entity.FutureHistoryKbar, error)
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
	InsertFutureHistoryTickArr(ctx context.Context, t []*entity.FutureHistoryTick) error
	QueryFutureHistoryTickArrByTime(ctx context.Context, code string, startTime, endTime time.Time) ([]*entity.FutureHistoryTick, error)
	InsertFutureHistoryClose(ctx context.Context, c *entity.FutureHistoryClose) error
	QueryFutureHistoryCloseByDate(ctx context.Context, code string, tradeDay time.Time) (*entity.FutureHistoryClose, error)
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

type RealTime interface {
	GetStockSnapshotByNumArr(stockNumArr []string) ([]*entity.StockSnapShot, error)
	GetTradeIndex() *entity.TradeIndex
	GetTSESnapshot(ctx context.Context) (*entity.StockSnapShot, error)
	GetOTCSnapshot(ctx context.Context) (*entity.StockSnapShot, error)

	GetMainFuture() *entity.Future
	GetFutureSnapshotByCode(code string) (*entity.FutureSnapShot, error)

	NewFutureRealTimeClient(tickChan chan *entity.RealTimeFutureTick, orderStatusChan chan interface{}, connectionID string)
	DeleteFutureRealTimeClient(connectionID string)

	CreateRealTimePick(connectionID string, com chan *pb.PickRealMap, tickChan chan []byte)
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
	StockTickConsumer(stockNum string, tickChan chan *entity.RealTimeStockTick)
	StockBidAskConsumer(stockNum string, bidAskChan chan *entity.RealTimeStockBidAsk)
	FutureTickConsumer(code string, tickChan chan *entity.RealTimeFutureTick)
	FutureBidAskConsumer(code string, bidAskChan chan *entity.FutureRealTimeBidAsk)

	StockTickPbConsumer(ctx context.Context, stockNum string, tickChan chan []byte)

	Close()
}

type Target interface {
	GetTargets(ctx context.Context) []*entity.StockTarget
	GetCurrentVolumeRank() (*pb.StockVolumeRankResponse, error)
}

type TargetRepo interface {
	InsertOrUpdateTargetArr(ctx context.Context, t []*entity.StockTarget) error
	QueryTargetsByTradeDay(ctx context.Context, tradeDay time.Time) ([]*entity.StockTarget, error)
	QueryAllMXFFuture(ctx context.Context) ([]*entity.Future, error)
}

type Trade interface {
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

	BuyLotStock(num string, price float64, lot int64) (string, entity.OrderStatus, error)
	BuyOddStock(num string, price float64, share int64) (string, entity.OrderStatus, error)

	GetFuturePosition() ([]*entity.FuturePosition, error)
	IsFutureTradeTime() bool

	ManualInsertFutureOrder(ctx context.Context, order *entity.FutureOrder) error
	UpdateTradeBalanceByTradeDay(ctx context.Context, date string) error
	MoveStockOrderToLatestTradeDay(ctx context.Context, orderID string) error
	MoveFutureOrderToLatestTradeDay(ctx context.Context, orderID string) error

	GetAccountBalance(ctx context.Context) ([]*entity.AccountBalance, error)
	IsAuthUser(username string) bool
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

type System interface {
	AddUser(ctx context.Context, t *entity.NewUser) error
	InsertPushToken(ctx context.Context, token, username string, enabled bool) error
	Login(ctx context.Context, username, password string) error
	VerifyEmail(ctx context.Context, username, code string) error
	UpdateAuthTradeUser()
	DeleteAllPushTokens(ctx context.Context) error
	IsPushTokenEnabled(ctx context.Context, token string) (bool, error)
	GetUserInfo(ctx context.Context, username string) (*entity.User, error)
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
}

type FCM interface {
	AnnounceMessage(msg string) error
	PushNotification(title, msg string) error
	SendTargets() error
}
