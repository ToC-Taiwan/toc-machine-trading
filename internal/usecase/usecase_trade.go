package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/module/quota"
	"tmt/internal/usecase/module/tradeday"
	"tmt/internal/usecase/repo"
	"tmt/internal/usecase/topic"
)

// TradeUseCase -.
type TradeUseCase struct {
	repo TradeRepo

	sc TradegRPCAPI
	fg TradegRPCAPI

	quota    *quota.Quota
	tradeDay *tradeday.TradeDay

	stockTradeDay  tradeday.TradePeriod
	futureTradeDay tradeday.TradePeriod

	updateFutureOrderLock sync.Mutex
	updateStockOrderLock  sync.Mutex
}

func (u *UseCaseBase) NewTrade() Trade {
	cfg := u.cfg
	tradeDay := tradeday.Get()

	uc := &TradeUseCase{
		sc:    grpcapi.NewTrade(u.sc, cfg.Simulation),
		fg:    grpcapi.NewTrade(u.fg, cfg.Simulation),
		repo:  repo.NewTrade(u.pg),
		quota: quota.NewQuota(cfg.Quota),

		tradeDay:       tradeDay,
		stockTradeDay:  tradeDay.GetStockTradeDay(),
		futureTradeDay: tradeDay.GetFutureTradeDay(),
	}

	bus.SubscribeTopic(topic.TopicInsertOrUpdateStockOrder, uc.updateStockOrderCacheAndInsertDB)
	bus.SubscribeTopic(topic.TopicInsertOrUpdateFutureOrder, uc.updateFutureOrderCacheAndInsertDB)

	if cfg.Simulation {
		go uc.askSimulateOrderStatus()
	} else {
		go uc.askOrderStatus()
	}

	go uc.updateAllTradeBalance()
	return uc
}

func (uc *TradeUseCase) updateAllTradeBalance() {
	for range time.NewTicker(time.Second * 20).C {
		if uc.IsStockTradeTime() {
			stockOrders, err := uc.repo.QueryAllStockOrderByDate(context.Background(), uc.stockTradeDay.ToStartEndArray())
			if err != nil {
				logger.Fatal(err)
			}
			uc.calculateStockTradeBalance(stockOrders, uc.stockTradeDay.TradeDay)
		}

		if uc.IsFutureTradeTime() {
			futureOrders, err := uc.repo.QueryAllFutureOrderByDate(context.Background(), uc.futureTradeDay.ToStartEndArray())
			if err != nil {
				logger.Fatal(err)
			}
			uc.calculateFutureTradeBalance(futureOrders, uc.futureTradeDay.TradeDay)
		}
	}
}

// UpdateTradeBalanceByTradeDay -.
func (uc *TradeUseCase) UpdateTradeBalanceByTradeDay(ctx context.Context, date string) error {
	stockTradePeriod, err := uc.tradeDay.GetStockTradePeriodByDate(date)
	if err != nil {
		return err
	}

	futureTradePeriod, err := uc.tradeDay.GetFutureTradePeriodByDate(date)
	if err != nil {
		return err
	}

	stockOrders, err := uc.repo.QueryAllStockOrderByDate(ctx, stockTradePeriod.ToStartEndArray())
	if err != nil {
		return err
	}
	uc.calculateStockTradeBalance(stockOrders, stockTradePeriod.TradeDay)

	futureOrders, err := uc.repo.QueryAllFutureOrderByDate(ctx, futureTradePeriod.ToStartEndArray())
	if err != nil {
		return err
	}
	uc.calculateFutureTradeBalance(futureOrders, futureTradePeriod.TradeDay)

	return nil
}

func (uc *TradeUseCase) MoveStockOrderToLatestTradeDay(ctx context.Context, orderID string) error {
	order, err := uc.repo.QueryStockOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.OrderTime = uc.stockTradeDay.StartTime
	return uc.repo.InsertOrUpdateOrderByOrderID(ctx, order)
}

func (uc *TradeUseCase) MoveFutureOrderToLatestTradeDay(ctx context.Context, orderID string) error {
	order, err := uc.repo.QueryFutureOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.OrderTime = uc.futureTradeDay.StartTime
	return uc.repo.InsertOrUpdateFutureOrderByOrderID(ctx, order)
}

func (uc *TradeUseCase) askOrderStatus() {
	for range time.NewTicker(750 * time.Millisecond).C {
		if !uc.IsFutureTradeTime() && !uc.IsStockTradeTime() {
			continue
		}

		if err := uc.sc.GetLocalOrderStatusArr(); err != nil {
			logger.Error(err)
		}

		if err := uc.fg.GetLocalOrderStatusArr(); err != nil {
			logger.Error(err)
		}
	}
}

func (uc *TradeUseCase) askSimulateOrderStatus() {
	for range time.NewTicker(750 * time.Millisecond).C {
		if !uc.IsFutureTradeTime() && !uc.IsStockTradeTime() {
			continue
		}

		if err := uc.sc.GetSimulateOrderStatusArr(); err != nil {
			logger.Error(err)
		}

		if err := uc.fg.GetSimulateOrderStatusArr(); err != nil {
			logger.Error(err)
		}
	}
}

func (uc *TradeUseCase) updateStockOrderCacheAndInsertDB(order *entity.StockOrder) {
	defer uc.updateStockOrderLock.Unlock()
	uc.updateStockOrderLock.Lock()

	// insert or update order to db
	if err := uc.repo.InsertOrUpdateOrderByOrderID(context.Background(), order); err != nil {
		logger.Fatal(err)
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
			qtyMap[v.StockNum] += v.Quantity
		case entity.ActionSell:
			if qtyMap[v.StockNum] > 0 {
				forward = append(forward, v)
			} else {
				reverse = append(reverse, v)
			}
			qtyMap[v.StockNum] -= v.Quantity
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
		logger.Fatal(err)
	}
}

func (uc *TradeUseCase) calculateForwardStockBalance(forward []*entity.StockOrder) (int64, int64, int64) {
	var forwardBalance, tradeCount, discount int64
	var qty int64
	for _, v := range forward {
		tradeCount++

		switch v.Action {
		case entity.ActionBuy:
			qty += v.Quantity
			forwardBalance -= uc.quota.GetStockBuyCost(v.Price, v.Quantity)
		case entity.ActionSell:
			qty -= v.Quantity
			forwardBalance += uc.quota.GetStockSellCost(v.Price, v.Quantity)
		}
		discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Quantity)
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
			qty -= v.Quantity
			revereBalance += uc.quota.GetStockSellCost(v.Price, v.Quantity)
		case entity.ActionBuy:
			qty += v.Quantity
			revereBalance -= uc.quota.GetStockBuyCost(v.Price, v.Quantity)
		}
		discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Quantity)
	}

	if qty != 0 {
		return 0, tradeCount, 0
	}

	return revereBalance, tradeCount, discount
}

func (uc *TradeUseCase) updateFutureOrderCacheAndInsertDB(order *entity.FutureOrder) {
	defer uc.updateFutureOrderLock.Unlock()
	uc.updateFutureOrderLock.Lock()

	// insert or update order to db
	if err := uc.repo.InsertOrUpdateFutureOrderByOrderID(context.Background(), order); err != nil {
		logger.Fatal(err)
	}
}

func (uc *TradeUseCase) ManualInsertFutureOrder(ctx context.Context, order *entity.FutureOrder) error {
	defer uc.updateFutureOrderLock.Unlock()
	uc.updateFutureOrderLock.Lock()

	// insert or update order to db
	if err := uc.repo.InsertOrUpdateFutureOrderByOrderID(context.Background(), order); err != nil {
		return err
	}

	return nil
}

// BuyFuture -.
func (uc *TradeUseCase) BuyFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error) {
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
			qtyMap[v.Code] += v.Quantity
		case entity.ActionSell:
			if qtyMap[v.Code] > 0 {
				forward = append(forward, v)
			} else {
				reverse = append(reverse, v)
			}
			qtyMap[v.Code] -= v.Quantity
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
		logger.Fatal(err)
	}
}

func (uc *TradeUseCase) calculateForwardFutureBalance(forward []*entity.FutureOrder) (int64, int64) {
	var forwardBalance, tradeCount int64
	var qty int64
	for _, v := range forward {
		tradeCount++

		switch v.Action {
		case entity.ActionBuy:
			qty += v.Quantity
			forwardBalance -= uc.quota.GetFutureBuyCost(v.Price, v.Quantity)
		case entity.ActionSell:
			qty -= v.Quantity
			forwardBalance += uc.quota.GetFutureSellCost(v.Price, v.Quantity)
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
			qty -= v.Quantity
			reverseBalance += uc.quota.GetFutureSellCost(v.Price, v.Quantity)
		case entity.ActionBuy:
			qty += v.Quantity
			reverseBalance -= uc.quota.GetFutureBuyCost(v.Price, v.Quantity)
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

func (uc *TradeUseCase) GetLastStockTradeBalance(ctx context.Context) (*entity.StockTradeBalance, error) {
	return uc.repo.QueryLastStockTradeBalance(ctx)
}

func (uc *TradeUseCase) GetLastFutureTradeBalance(ctx context.Context) (*entity.FutureTradeBalance, error) {
	return uc.repo.QueryLastFutureTradeBalance(ctx)
}

// CalculateBuyCost -.
func (uc *TradeUseCase) CalculateBuyCost(price float64, quantity int64) int64 {
	return uc.quota.GetStockBuyCost(price, quantity)
}

// CalculateSellCost -.
func (uc *TradeUseCase) CalculateSellCost(price float64, quantity int64) int64 {
	return uc.quota.GetStockSellCost(price, quantity)
}

// CalculateTradeDiscount -.
func (uc *TradeUseCase) CalculateTradeDiscount(price float64, quantity int64) int64 {
	return uc.quota.GetStockTradeFeeDiscount(price, quantity)
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
			Quantity:  int64(v.GetQuantity()),
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
