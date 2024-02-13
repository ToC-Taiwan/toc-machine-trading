package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"tmt/internal/config"

	"tmt/internal/entity"
	"tmt/internal/usecase/grpc"
	"tmt/internal/usecase/modules/calendar"
	"tmt/internal/usecase/modules/quota"
	"tmt/internal/usecase/repo"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"

	"github.com/google/uuid"
)

// TradeUseCase -.
type TradeUseCase struct {
	repo TradeRepo
	sc   TradegRPCAPI

	quota    *quota.Quota
	tradeDay *calendar.Calendar

	stockTradeDay  calendar.TradePeriod
	futureTradeDay calendar.TradePeriod

	finishedStockOrderMap  map[string]*entity.StockOrder
	finishedFutureOrderMap map[string]*entity.FutureOrder
	updateFutureOrderLock  sync.Mutex
	updateStockOrderLock   sync.Mutex

	logger *log.Log
	bus    *eventbus.Bus

	authUserMap     map[string]struct{}
	authUserMapLock sync.RWMutex
}

func NewTrade() Trade {
	cfg := config.Get()
	tradeDay := calendar.Get()
	uc := &TradeUseCase{
		sc:    grpc.NewTrade(cfg.GetSinopacPool(), cfg.Simulation),
		repo:  repo.NewTrade(cfg.GetPostgresPool()),
		quota: quota.NewQuota(cfg.Quota),

		tradeDay:       tradeDay,
		stockTradeDay:  tradeDay.GetStockTradeDay(),
		futureTradeDay: tradeDay.GetFutureTradeDay(),

		finishedStockOrderMap:  make(map[string]*entity.StockOrder),
		finishedFutureOrderMap: make(map[string]*entity.FutureOrder),

		logger: log.Get(),
		bus:    eventbus.Get(),
	}

	uc.bus.SubscribeAsync(topicInsertOrUpdateStockOrder, true, uc.updateStockOrderCacheAndInsertDB)
	uc.bus.SubscribeAsync(topicInsertOrUpdateFutureOrder, true, uc.updateFutureOrderCacheAndInsertDB)
	uc.bus.SubscribeAsync(topicUpdateAuthTradeUser, true, uc.updateAuthUserMap)

	go uc.askOrderStatus(cfg.Simulation)
	go uc.updateAccountDetail()
	go uc.updateAllTradeBalance()

	return uc
}

func (uc *TradeUseCase) updateAccountDetail() {
	for range time.NewTicker(time.Minute).C {
		uc.updateAccountBalance()
		uc.updateAccountSettlement()
		uc.updateStockInventory()
		uc.updateFutureInventory()
	}
}

func (uc *TradeUseCase) updateAuthUserMap(username []string) {
	uc.authUserMapLock.Lock()
	defer uc.authUserMapLock.Unlock()

	uc.authUserMap = make(map[string]struct{})
	for _, v := range username {
		uc.logger.Warnf("auth user: %s", v)
		uc.authUserMap[v] = struct{}{}
	}
}

func (uc *TradeUseCase) IsAuthUser(username string) bool {
	uc.authUserMapLock.RLock()
	defer uc.authUserMapLock.RUnlock()

	if _, ok := uc.authUserMap[username]; ok {
		return true
	}
	return false
}

func (uc *TradeUseCase) updateAccountBalance() {
	margin, err := uc.sc.GetMargin()
	if err != nil {
		uc.logger.Fatal(err)
	}

	accountBalance, er := uc.sc.GetAccountBalance()
	if err != nil {
		uc.logger.Fatal(er)
	}

	balance := &entity.AccountBalance{
		Date:            time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		Balance:         accountBalance.Balance,
		TodayMargin:     margin.TodayBalance,
		AvailableMargin: margin.AvailableMargin,
		YesterdayMargin: margin.YesterdayBalance,
		RiskIndicator:   margin.RiskIndicator,
	}

	err = uc.repo.InsertOrUpdateAccountBalance(context.Background(), balance)
	if err != nil {
		uc.logger.Fatal(err)
	}
}

func (uc *TradeUseCase) updateAccountSettlement() {
	sinopacSettlement, err := uc.sc.GetSettlement()
	if err != nil {
		uc.logger.Fatal(err)
	}

	for _, v := range sinopacSettlement.GetSettlement() {
		dateTime, err := time.ParseInLocation(entity.LongTimeLayout, v.GetDate(), time.Local)
		if err != nil {
			uc.logger.Fatal(err)
		}
		err = uc.repo.InsertOrUpdateAccountSettlement(context.Background(), &entity.Settlement{
			Date:       dateTime,
			Settlement: v.GetAmount(),
		})
		if err != nil {
			uc.logger.Fatal(err)
		}
	}
}

func (uc *TradeUseCase) updateStockInventory() {
	inv := []*entity.InventoryStock{}
	timeNowZero := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	queryData, err := uc.sc.GetStockPosition()
	if err != nil {
		uc.logger.Fatal(err)
	}

	dbInvMap, err := uc.repo.QueryInventoryUUIDStockByDate(context.Background(), timeNowZero)
	if err != nil {
		uc.logger.Fatal(err)
	}
	for _, s := range queryData.GetPositionArr() {
		invID := dbInvMap[s.GetCode()]
		if invID == "" {
			invID = uuid.NewString()
		} else {
			delete(dbInvMap, s.GetCode())
		}
		lot, share := int(s.GetQuantity())/1000, int(s.GetQuantity())%1000
		position := []*entity.PositionStock{}
		for _, p := range s.GetDetailArr() {
			date, pErr := time.ParseInLocation(entity.ShortTimeLayout, p.GetDate(), time.Local)
			if pErr != nil {
				continue
			}
			position = append(position, &entity.PositionStock{
				StockNum:  s.GetCode(),
				Date:      date,
				Quantity:  int(p.GetQuantity()),
				Price:     p.GetPrice(),
				LastPrice: p.GetLastPrice(),
				Dseq:      p.GetDseq(),
				Direction: p.GetDirection(),
				Pnl:       p.GetPnl(),
				Fee:       p.GetFee(),
				InvID:     invID,
			})
		}
		inv = append(inv, &entity.InventoryStock{
			InventoryBase: entity.InventoryBase{
				UUID:     invID,
				AvgPrice: s.GetPrice(),
				Date:     timeNowZero,
			},
			StockNum: s.GetCode(),
			Lot:      lot,
			Share:    share,
			Position: position,
		})
	}
	if len(inv) == 0 {
		return
	}
	if err = uc.repo.InsertOrUpdateInventoryStock(context.Background(), inv); err != nil {
		uc.logger.Fatal(err)
	}
	for _, v := range dbInvMap {
		_ = uc.repo.ClearInventoryStockByUUID(context.Background(), v)
	}
}

func (uc *TradeUseCase) updateFutureInventory() {}

func (uc *TradeUseCase) updateAllTradeBalance() {
	for range time.NewTicker(time.Second * 20).C {
		if uc.IsStockTradeTime() {
			stockOrders, err := uc.repo.QueryAllStockOrderByDate(context.Background(), uc.stockTradeDay.ToStartEndArray())
			if err != nil {
				uc.logger.Fatal(err)
			}
			uc.calculateStockTradeBalance(stockOrders, uc.stockTradeDay.TradeDay)
		}

		if uc.IsFutureTradeTime() {
			futureOrders, err := uc.repo.QueryAllFutureOrderByDate(context.Background(), uc.futureTradeDay.ToStartEndArray())
			if err != nil {
				uc.logger.Fatal(err)
			}
			uc.calculateFutureTradeBalance(futureOrders, uc.futureTradeDay.TradeDay)
		}
	}
}

func (uc *TradeUseCase) askOrderStatus(sim bool) {
	scFn := uc.sc.GetLocalOrderStatusArr
	if sim {
		scFn = uc.sc.GetSimulateOrderStatusArr
	}

	for range time.NewTicker(750 * time.Millisecond).C {
		if !uc.IsFutureTradeTime() && !uc.IsStockTradeTime() {
			continue
		}

		if err := scFn(); err != nil {
			uc.logger.Error(err)
		}
	}
}

func (uc *TradeUseCase) updateStockOrderCacheAndInsertDB(order *entity.StockOrder) {
	defer uc.updateStockOrderLock.Unlock()
	uc.updateStockOrderLock.Lock()
	if _, ok := uc.finishedStockOrderMap[order.OrderID]; ok {
		return
	}

	// insert or update order to db
	if err := uc.repo.InsertOrUpdateOrderByOrderID(context.Background(), order); err != nil {
		uc.logger.Fatal(err)
	}

	if !order.Cancellable() {
		uc.finishedStockOrderMap[order.OrderID] = order
	}
}

// calculateStockTradeBalance -.
func (uc *TradeUseCase) calculateStockTradeBalance(allOrders []*entity.StockOrder, tradeDay time.Time) {
	var forward, reverse []*entity.StockOrder
	qtyMap := make(map[string]int64)
	for _, v := range allOrders {
		if v.Status != entity.StatusFilled {
			continue
		}

		switch v.Action {
		case entity.ActionBuy:
			if qtyMap[v.StockNum] >= 0 {
				forward = append(forward, v)
			} else {
				reverse = append(reverse, v)
			}
			qtyMap[v.StockNum] += v.Lot
		case entity.ActionSell:
			if qtyMap[v.StockNum] > 0 {
				forward = append(forward, v)
			} else {
				reverse = append(reverse, v)
			}
			qtyMap[v.StockNum] -= v.Lot
		}
	}

	forwardBalance, fDiscount, fTradeCount := uc.calculateForwardStockBalance(forward)
	revereBalance, rDiscount, rTradeCount := uc.calculateReverseStockBalance(reverse)
	tmp := &entity.StockTradeBalance{
		TradeDay:        tradeDay,
		TradeCount:      fTradeCount + rTradeCount,
		Forward:         forwardBalance,
		Reverse:         revereBalance,
		OriginalBalance: forwardBalance + revereBalance,
		Discount:        fDiscount + rDiscount,
		Total:           forwardBalance + revereBalance + fDiscount + rDiscount,
	}

	err := uc.repo.InsertOrUpdateStockTradeBalance(context.Background(), tmp)
	if err != nil {
		uc.logger.Fatal(err)
	}
}

func (uc *TradeUseCase) calculateForwardStockBalance(forward []*entity.StockOrder) (int64, int64, int64) {
	var forwardBalance, tradeCount, discount int64
	var qty int64
	for _, v := range forward {
		tradeCount++

		switch v.Action {
		case entity.ActionBuy:
			qty += v.Lot
			forwardBalance -= uc.quota.GetStockBuyCost(v.Price, v.Lot, v.Share)
		case entity.ActionSell:
			qty -= v.Lot
			forwardBalance += uc.quota.GetStockSellCost(v.Price, v.Lot, v.Share)
		}
		discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Lot, v.Share)
	}

	if qty != 0 {
		return 0, tradeCount, 0
	}

	return forwardBalance, tradeCount, discount
}

func (uc *TradeUseCase) calculateReverseStockBalance(reverse []*entity.StockOrder) (int64, int64, int64) {
	var revereBalance, tradeCount, discount int64
	var qty int64
	for _, v := range reverse {
		tradeCount++

		switch v.Action {
		case entity.ActionSell:
			qty -= v.Lot
			revereBalance += uc.quota.GetStockSellCost(v.Price, v.Lot, v.Share)
		case entity.ActionBuy:
			qty += v.Lot
			revereBalance -= uc.quota.GetStockBuyCost(v.Price, v.Lot, v.Share)
		}
		discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Lot, v.Share)
	}

	if qty != 0 {
		return 0, tradeCount, 0
	}

	return revereBalance, tradeCount, discount
}

func (uc *TradeUseCase) updateFutureOrderCacheAndInsertDB(order *entity.FutureOrder) {
	defer uc.updateFutureOrderLock.Unlock()
	uc.updateFutureOrderLock.Lock()
	if _, ok := uc.finishedFutureOrderMap[order.OrderID]; ok {
		return
	}

	// insert or update order to db
	if err := uc.repo.InsertOrUpdateFutureOrderByOrderID(context.Background(), order); err != nil {
		uc.logger.Fatal(err)
	}

	if !order.Cancellable() {
		uc.finishedFutureOrderMap[order.OrderID] = order
	}
}

// BuyFuture -.
func (uc *TradeUseCase) BuyFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error) {
	if order.Code == "" {
		return "", entity.StatusUnknow, errors.New("empty code")
	}

	result, err := uc.sc.BuyFuture(order)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// SellFuture -.
func (uc *TradeUseCase) SellFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error) {
	result, err := uc.sc.SellFuture(order)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// CancelFutureOrderByID -.
func (uc *TradeUseCase) CancelFutureOrderByID(orderID string) (string, entity.OrderStatus, error) {
	result, err := uc.sc.CancelFuture(orderID)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

func (uc *TradeUseCase) BuyOddStock(num string, price float64, share int64) (string, entity.OrderStatus, error) {
	if num == "" {
		return "", entity.StatusUnknow, errors.New("empty stock num")
	}

	result, err := uc.sc.BuyOddStock(&entity.StockOrder{
		OrderDetail: entity.OrderDetail{
			Price: price,
		},
		Share:    share,
		StockNum: num,
	})
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// calculateFutureTradeBalance -.
func (uc *TradeUseCase) calculateFutureTradeBalance(allOrders []*entity.FutureOrder, tradeDay time.Time) {
	var forward, reverse []*entity.FutureOrder
	qtyMap := make(map[string]int64)
	for _, v := range allOrders {
		if v.Status != entity.StatusFilled {
			continue
		}

		switch v.Action {
		case entity.ActionBuy:
			if qtyMap[v.Code] >= 0 {
				forward = append(forward, v)
			} else {
				reverse = append(reverse, v)
			}
			qtyMap[v.Code] += v.Position
		case entity.ActionSell:
			if qtyMap[v.Code] > 0 {
				forward = append(forward, v)
			} else {
				reverse = append(reverse, v)
			}
			qtyMap[v.Code] -= v.Position
		}
	}

	forwardBalance, forwardCount := uc.calculateForwardFutureBalance(forward)
	revereBalance, reverseCount := uc.calculateReverseFutureBalance(reverse)
	tmp := &entity.FutureTradeBalance{
		TradeDay:   tradeDay,
		TradeCount: forwardCount + reverseCount,
		Forward:    forwardBalance,
		Reverse:    revereBalance,
		Total:      forwardBalance + revereBalance,
	}

	err := uc.repo.InsertOrUpdateFutureTradeBalance(context.Background(), tmp)
	if err != nil {
		uc.logger.Fatal(err)
	}
}

func (uc *TradeUseCase) calculateForwardFutureBalance(forward []*entity.FutureOrder) (int64, int64) {
	var forwardBalance, tradeCount int64
	var qty int64
	for _, v := range forward {
		tradeCount++

		switch v.Action {
		case entity.ActionBuy:
			qty += v.Position
			forwardBalance -= uc.quota.GetFutureBuyCost(v.Price, v.Position)
		case entity.ActionSell:
			qty -= v.Position
			forwardBalance += uc.quota.GetFutureSellCost(v.Price, v.Position)
		}
	}

	if qty != 0 {
		return 0, tradeCount
	}

	return forwardBalance, tradeCount
}

func (uc *TradeUseCase) calculateReverseFutureBalance(reverse []*entity.FutureOrder) (int64, int64) {
	var reverseBalance, tradeCount int64
	var qty int64
	for _, v := range reverse {
		tradeCount++

		switch v.Action {
		case entity.ActionSell:
			qty -= v.Position
			reverseBalance += uc.quota.GetFutureSellCost(v.Price, v.Position)
		case entity.ActionBuy:
			qty += v.Position
			reverseBalance -= uc.quota.GetFutureBuyCost(v.Price, v.Position)
		}
	}

	if qty != 0 {
		return 0, tradeCount
	}

	return reverseBalance, tradeCount
}

// GetAllStockOrder -.
func (uc *TradeUseCase) GetAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error) {
	return uc.repo.QueryAllStockOrder(ctx)
}

func (uc *TradeUseCase) GetAllFutureOrder(ctx context.Context) ([]*entity.FutureOrder, error) {
	return uc.repo.QueryAllFutureOrder(ctx)
}

// GetAllStockTradeBalance -.
func (uc *TradeUseCase) GetAllStockTradeBalance(ctx context.Context) ([]*entity.StockTradeBalance, error) {
	tradeBalanceArr, err := uc.repo.QueryAllStockTradeBalance(ctx)
	if err != nil {
		return nil, err
	}
	return tradeBalanceArr, nil
}

// GetAllFutureTradeBalance -.
func (uc *TradeUseCase) GetAllFutureTradeBalance(ctx context.Context) ([]*entity.FutureTradeBalance, error) {
	tradeBalanceArr, err := uc.repo.QueryAllFutureTradeBalance(ctx)
	if err != nil {
		return nil, err
	}
	return tradeBalanceArr, nil
}

// GetFuturePosition .
func (uc *TradeUseCase) GetFuturePosition() ([]*entity.FuturePosition, error) {
	query, err := uc.sc.GetFuturePosition()
	if err != nil {
		return nil, err
	}
	var result []*entity.FuturePosition
	for _, v := range query.GetPositionArr() {
		result = append(result, &entity.FuturePosition{
			Code:      v.GetCode(),
			Direction: v.GetDirection(),
			Position:  int64(v.GetQuantity()),
			Price:     v.GetPrice(),
			LastPrice: v.GetLastPrice(),
			Pnl:       v.GetPnl(),
		})
	}
	return result, nil
}

func (uc *TradeUseCase) IsStockTradeTime() bool {
	return uc.stockTradeDay.IsStockMarketOpenNow()
}

func (uc *TradeUseCase) IsFutureTradeTime() bool {
	return uc.futureTradeDay.IsFutureMarketOpenNow()
}

func (uc *TradeUseCase) GetFutureOrderByTradeDay(ctx context.Context, tradeDay string) ([]*entity.FutureOrder, error) {
	period, err := uc.tradeDay.GetFutureTradePeriodByDate(tradeDay)
	if err != nil {
		return nil, err
	}

	orders, err := uc.repo.QueryAllFutureOrderByDate(ctx, []time.Time{period.StartTime, period.EndTime})
	if err != nil {
		return nil, err
	}

	var filledOrder []*entity.FutureOrder
	for _, v := range orders {
		if v.Status == entity.StatusFilled {
			filledOrder = append(filledOrder, v)
		}
	}

	return filledOrder, nil
}

func (uc *TradeUseCase) GetAccountBalance(ctx context.Context) (*entity.AccountBalance, error) {
	data, err := uc.repo.QueryLastAccountBalance(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}
