package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/event"
	"tmt/internal/usecase/module/quota"
	"tmt/internal/usecase/module/tradeday"
)

// OrderUseCase -.
type OrderUseCase struct {
	gRPCAPI OrdergRPCAPI
	repo    OrderRepo

	quota *quota.Quota

	placeOrderLock       sync.Mutex
	placeFutureOrderLock sync.Mutex

	tradeDay       *tradeday.TradeDay
	stockTradeDay  tradeday.TradePeriod
	futureTradeDay tradeday.TradePeriod

	simTrade              bool
	updateFutureOrderLock sync.Mutex
	updateStockOrderLock  sync.Mutex
}

// NewOrder -.
func NewOrder(t OrdergRPCAPI, r OrderRepo) *OrderUseCase {
	cfg := config.GetConfig()
	tradeDay := tradeday.NewTradeDay()

	uc := &OrderUseCase{
		simTrade: cfg.Simulation,

		gRPCAPI: t,
		repo:    r,
		quota:   quota.NewQuota(cfg.Quota),

		tradeDay:       tradeDay,
		stockTradeDay:  tradeDay.GetStockTradeDay(),
		futureTradeDay: tradeDay.GetFutureTradeDay(),
	}

	bus.SubscribeTopic(event.TopicPlaceStockOrder, uc.placeStockOrder)
	bus.SubscribeTopic(event.TopicCancelStockOrder, uc.cancelStockOrder)
	bus.SubscribeTopic(event.TopicInsertOrUpdateStockOrder, uc.updateStockOrderCacheAndInsertDB)

	bus.SubscribeTopic(event.TopicPlaceFutureOrder, uc.placeFutureOrder)
	bus.SubscribeTopic(event.TopicCancelFutureOrder, uc.cancelFutureOrder)
	bus.SubscribeTopic(event.TopicInsertOrUpdateFutureOrder, uc.updateFutureOrderCacheAndInsertDB)

	if uc.simTrade {
		go uc.askSimulateOrderStatus()
	} else {
		go uc.askOrderStatus()
	}
	go uc.updateAllTradeBalance()
	return uc
}

func (uc *OrderUseCase) updateAllTradeBalance() {
	for range time.NewTicker(time.Second * 20).C {
		if uc.IsStockTradeTime() {
			stockOrders, err := uc.repo.QueryAllStockOrderByDate(context.Background(), uc.stockTradeDay.ToStartEndArray())
			if err != nil {
				logger.Panic(err)
			}
			uc.calculateStockTradeBalance(stockOrders, uc.stockTradeDay.TradeDay)
		}

		if uc.IsFutureTradeTime() {
			futureOrders, err := uc.repo.QueryAllFutureOrderByDate(context.Background(), uc.futureTradeDay.ToStartEndArray())
			if err != nil {
				logger.Panic(err)
			}
			uc.calculateFutureTradeBalance(futureOrders, uc.futureTradeDay.TradeDay)
		}
	}
}

// UpdateTradeBalanceByTradeDay -.
func (uc *OrderUseCase) UpdateTradeBalanceByTradeDay(ctx context.Context, date string) error {
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

func (uc *OrderUseCase) MoveStockOrderToLatestTradeDay(ctx context.Context, orderID string) error {
	order, err := uc.repo.QueryStockOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.OrderTime = uc.stockTradeDay.StartTime
	return uc.repo.InsertOrUpdateOrderByOrderID(ctx, order)
}

func (uc *OrderUseCase) MoveFutureOrderToLatestTradeDay(ctx context.Context, orderID string) error {
	order, err := uc.repo.QueryFutureOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.OrderTime = uc.futureTradeDay.StartTime
	return uc.repo.InsertOrUpdateFutureOrderByOrderID(ctx, order)
}

func (uc *OrderUseCase) askOrderStatus() {
	for range time.NewTicker(750 * time.Millisecond).C {
		if !uc.IsFutureTradeTime() && !uc.IsStockTradeTime() {
			continue
		}

		if err := uc.gRPCAPI.GetLocalOrderStatusArr(); err != nil {
			logger.Error(err)
			continue
		}
	}
}

func (uc *OrderUseCase) askSimulateOrderStatus() {
	for range time.NewTicker(750 * time.Millisecond).C {
		if !uc.IsFutureTradeTime() && !uc.IsStockTradeTime() {
			continue
		}

		if err := uc.gRPCAPI.GetSimulateOrderStatusArr(); err != nil {
			logger.Error(err)
			continue
		}
	}
}

// AskOrderUpdate -.
// func (uc *OrderUseCase) askOrderStatusProd() {
// 	for range time.NewTicker(750 * time.Millisecond).C {
// 		if !uc.IsFutureTradeTime() && !uc.IsStockTradeTime() {
// 			continue
// 		}

// 		msg, err := uc.gRPCAPI.GetNonBlockOrderStatusArr()
// 		if err != nil {
// 			logger.Error(err)
// 			continue
// 		}

// 		if errMsg := msg.GetErr(); errMsg != "" {
// 			logger.Error(errMsg)
// 			continue
// 		}
// 	}
// }

func (uc *OrderUseCase) placeStockOrder(order *entity.StockOrder) {
	defer uc.placeOrderLock.Unlock()
	uc.placeOrderLock.Lock()

	cosumeQuota := uc.quota.CalculateOriginalOrderCost(order)
	if cosumeQuota != 0 && uc.quota.IsEnough(cosumeQuota) {
		order.Status = entity.StatusAborted
		return
	}

	var orderID string
	var status entity.OrderStatus
	var err error
	switch order.Action {
	case entity.ActionBuy, entity.ActionBuyLater:
		orderID, status, err = uc.BuyStock(order)
	case entity.ActionSell:
		orderID, status, err = uc.SellStock(order)
	case entity.ActionSellFirst:
		orderID, status, err = uc.SellFirstStock(order)
	}
	if err != nil {
		logger.Error(err)
		order.Status = entity.StatusFailed
		return
	}

	if status == entity.StatusFailed || orderID == "" {
		order.Status = entity.StatusFailed
		return
	}

	// count quota
	uc.quota.CosumeQuota(cosumeQuota)

	// modify order and save to cache
	order.OrderID = orderID
	order.Status = status
	order.TradeTime = time.Now()
	cc.SetOrderByOrderID(order)

	logger.Infof("Place Stock Order -> Stock: %s, Action: %d, Price: %.2f, Qty: %d, Quota: %d", order.StockNum, order.Action, order.Price, order.Quantity, uc.quota.GetCurrentQuota())
}

func (uc *OrderUseCase) cancelStockOrder(order *entity.StockOrder) {
	defer uc.placeOrderLock.Unlock()
	uc.placeOrderLock.Lock()

	order.TradeTime = time.Now()
	logger.Infof("Cancel Stock Order -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)

	// result will return instantly
	_, _, err := uc.CancelOrderID(order.OrderID)
	if err != nil {
		logger.Error(err)
		return
	}

	if cosumeQuota := uc.quota.CalculateOriginalOrderCost(order); cosumeQuota > 0 {
		uc.quota.BackQuota(cosumeQuota)
		logger.Infof("Quota Back: %d", uc.quota.GetCurrentQuota())
	}
}

// BuyStock -.
func (uc *OrderUseCase) BuyStock(order *entity.StockOrder) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.BuyStock(order, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// SellStock -.
func (uc *OrderUseCase) SellStock(order *entity.StockOrder) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.SellStock(order, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// SellFirstStock -.
func (uc *OrderUseCase) SellFirstStock(order *entity.StockOrder) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.SellFirstStock(order, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// BuyLaterStock -.
func (uc *OrderUseCase) BuyLaterStock(order *entity.StockOrder) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.BuyStock(order, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// CancelOrderID -.
func (uc *OrderUseCase) CancelOrderID(orderID string) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.CancelStock(orderID, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

func (uc *OrderUseCase) updateStockOrderCacheAndInsertDB(order *entity.StockOrder) {
	defer uc.updateStockOrderLock.Unlock()
	uc.updateStockOrderLock.Lock()

	// get order from cache
	cacheOrder := cc.GetOrderByOrderID(order.OrderID)
	if cacheOrder == nil {
		return
	}

	cacheOrder.Status = order.Status
	cacheOrder.OrderTime = order.OrderTime

	// qty may not filled with original order, change it by return quantity
	cacheOrder.Quantity = order.Quantity

	// update cache
	cc.SetOrderByOrderID(cacheOrder)

	// insert or update order to db
	if err := uc.repo.InsertOrUpdateOrderByOrderID(context.Background(), cacheOrder); err != nil {
		logger.Panic(err)
	}
}

// calculateStockTradeBalance -.
func (uc *OrderUseCase) calculateStockTradeBalance(allOrders []*entity.StockOrder, tradeDay time.Time) {
	var forward, reverse entity.StockOrderArr
	for _, v := range allOrders {
		if v.Status != entity.StatusFilled {
			continue
		}

		switch v.Action {
		case entity.ActionBuy, entity.ActionSell:
			forward = append(forward, v)
		case entity.ActionSellFirst, entity.ActionBuyLater:
			reverse = append(reverse, v)
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
		logger.Panic(err)
	}
}

func (uc *OrderUseCase) calculateForwardStockBalance(forward entity.StockOrderArr) (int64, int64, int64) {
	var forwardBalance, tradeCount, discount int64
	groupOrder, manual := forward.SplitManualAndGroupID()
	for _, v := range groupOrder {
		if len(v) != 2 {
			continue
		}

		tradeCount += 2
		forwardBalance -= uc.quota.GetStockBuyCost(v[0].Price, v[0].Quantity)
		discount += uc.quota.GetStockTradeFeeDiscount(v[0].Price, v[0].Quantity)

		forwardBalance += uc.quota.GetStockSellCost(v[1].Price, v[1].Quantity)
		discount += uc.quota.GetStockTradeFeeDiscount(v[1].Price, v[1].Quantity)
	}

	if manual.IsAllDone() {
		for _, v := range manual {
			tradeCount++

			switch v.Action {
			case entity.ActionBuy:
				forwardBalance -= uc.quota.GetStockBuyCost(v.Price, v.Quantity)
			case entity.ActionSell:
				forwardBalance += uc.quota.GetStockSellCost(v.Price, v.Quantity)
			}
			discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Quantity)
		}
	}
	return forwardBalance, tradeCount, discount
}

func (uc *OrderUseCase) calculateReverseStockBalance(reverse entity.StockOrderArr) (int64, int64, int64) {
	var revereBalance, tradeCount, discount int64
	groupOrder, manual := reverse.SplitManualAndGroupID()
	for _, v := range groupOrder {
		if len(v) != 2 {
			continue
		}

		tradeCount += 2
		revereBalance += uc.quota.GetStockSellCost(v[0].Price, v[0].Quantity)
		discount += uc.quota.GetStockTradeFeeDiscount(v[0].Price, v[0].Quantity)

		revereBalance -= uc.quota.GetStockBuyCost(v[1].Price, v[1].Quantity)
		discount += uc.quota.GetStockTradeFeeDiscount(v[1].Price, v[1].Quantity)
	}

	if manual.IsAllDone() {
		for _, v := range manual {
			tradeCount++

			switch v.Action {
			case entity.ActionSellFirst:
				revereBalance += uc.quota.GetStockSellCost(v.Price, v.Quantity)
			case entity.ActionBuy:
				revereBalance -= uc.quota.GetStockBuyCost(v.Price, v.Quantity)
			}
			discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Quantity)
		}
	}
	return revereBalance, tradeCount, discount
}

//
// below is future trade
//

func (uc *OrderUseCase) placeFutureOrder(order *entity.FutureOrder) {
	defer uc.placeFutureOrderLock.Unlock()
	uc.placeFutureOrderLock.Lock()

	var orderID string
	var status entity.OrderStatus
	var err error
	switch order.Action {
	case entity.ActionBuy, entity.ActionBuyLater:
		orderID, status, err = uc.BuyFuture(order)
	case entity.ActionSell:
		orderID, status, err = uc.SellFuture(order)
	case entity.ActionSellFirst:
		orderID, status, err = uc.SellFirstFuture(order)
	}
	if err != nil {
		logger.Error(err)
		order.Status = entity.StatusFailed
		return
	}

	if status == entity.StatusFailed || orderID == "" {
		order.Status = entity.StatusFailed
		return
	}

	// modify order and save to cache
	order.OrderID = orderID
	order.Status = status
	order.TradeTime = time.Now()
	cc.SetFutureOrderByOrderID(order)

	logger.Infof("Place Future Order -> Future: %s, Action: %d, Price: %.0f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)
}

func (uc *OrderUseCase) cancelFutureOrder(order *entity.FutureOrder) {
	defer uc.placeFutureOrderLock.Unlock()
	uc.placeFutureOrderLock.Lock()

	order.TradeTime = time.Now()
	logger.Infof("Cancel Future Order -> Future: %s, Action: %d, Price: %.0f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)

	// result will return instantly
	_, _, err := uc.CancelFutureOrderID(order.OrderID)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (uc *OrderUseCase) updateFutureOrderCacheAndInsertDB(order *entity.FutureOrder) {
	defer uc.updateFutureOrderLock.Unlock()
	uc.updateFutureOrderLock.Lock()

	// get order from cache
	cacheOrder := cc.GetFutureOrderByOrderID(order.OrderID)
	if cacheOrder == nil {
		return
	}

	cacheOrder.Status = order.Status
	if cacheOrder.OrderTime.IsZero() {
		cacheOrder.OrderTime = order.OrderTime
	}

	// qty may not filled with original order, change it by return quantity
	cacheOrder.Quantity = order.Quantity

	// update cache
	cc.SetFutureOrderByOrderID(cacheOrder)

	// insert or update order to db
	if err := uc.repo.InsertOrUpdateFutureOrderByOrderID(context.Background(), cacheOrder); err != nil {
		logger.Panic(err)
	}
}

func (uc *OrderUseCase) ManualInsertFutureOrder(ctx context.Context, order *entity.FutureOrder) error {
	defer uc.updateFutureOrderLock.Unlock()
	uc.updateFutureOrderLock.Lock()

	// insert or update order to db
	if err := uc.repo.InsertOrUpdateFutureOrderByOrderID(context.Background(), order); err != nil {
		return err
	}

	return nil
}

// BuyFuture -.
func (uc *OrderUseCase) BuyFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.BuyFuture(order, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// SellFuture -.
func (uc *OrderUseCase) SellFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.SellFuture(order, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// SellFirstFuture -.
func (uc *OrderUseCase) SellFirstFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.SellFirstFuture(order, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// BuyLaterFuture -.
func (uc *OrderUseCase) BuyLaterFuture(order *entity.FutureOrder) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.BuyFuture(order, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// CancelFutureOrderID -.
func (uc *OrderUseCase) CancelFutureOrderID(orderID string) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.CancelFuture(orderID, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}

	if e := result.GetError(); e != "" {
		return "", entity.StatusUnknow, errors.New(e)
	}

	return result.GetOrderId(), entity.StringToOrderStatus(result.GetStatus()), nil
}

// calculateFutureTradeBalance -.
func (uc *OrderUseCase) calculateFutureTradeBalance(allOrders []*entity.FutureOrder, tradeDay time.Time) {
	var forward, reverse entity.FutureOrderArr
	for _, v := range allOrders {
		if v.Status != entity.StatusFilled {
			continue
		}

		switch v.Action {
		case entity.ActionBuy, entity.ActionSell:
			forward = append(forward, v)
		case entity.ActionSellFirst, entity.ActionBuyLater:
			reverse = append(reverse, v)
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
		logger.Panic(err)
	}
}

func (uc *OrderUseCase) calculateForwardFutureBalance(forward entity.FutureOrderArr) (int64, int64) {
	var forwardBalance, tradeCount int64
	groupOrder, manual := forward.SplitManualAndGroupID()
	for _, v := range groupOrder {
		if len(v) != 2 {
			continue
		}

		tradeCount += 2
		forwardBalance -= uc.quota.GetFutureBuyCost(v[0].Price, v[0].Quantity)
		forwardBalance += uc.quota.GetFutureSellCost(v[1].Price, v[1].Quantity)
	}

	if manual.IsAllDone() {
		for _, v := range manual {
			tradeCount++

			switch v.Action {
			case entity.ActionBuy:
				forwardBalance -= uc.quota.GetFutureBuyCost(v.Price, v.Quantity)
			case entity.ActionSell:
				forwardBalance += uc.quota.GetFutureSellCost(v.Price, v.Quantity)
			}
		}
	}
	return forwardBalance, tradeCount
}

func (uc *OrderUseCase) calculateReverseFutureBalance(reverse entity.FutureOrderArr) (int64, int64) {
	var reverseBalance, tradeCount int64
	groupOrder, manual := reverse.SplitManualAndGroupID()
	for _, v := range groupOrder {
		if len(v) != 2 {
			continue
		}

		tradeCount += 2
		reverseBalance += uc.quota.GetFutureSellCost(v[0].Price, v[0].Quantity)
		reverseBalance -= uc.quota.GetFutureBuyCost(v[1].Price, v[1].Quantity)
	}

	if manual.IsAllDone() {
		for _, v := range manual {
			tradeCount++

			switch v.Action {
			case entity.ActionSellFirst:
				reverseBalance += uc.quota.GetFutureSellCost(v.Price, v.Quantity)
			case entity.ActionBuyLater:
				reverseBalance -= uc.quota.GetFutureBuyCost(v.Price, v.Quantity)
			}
		}
	}
	return reverseBalance, tradeCount
}

//
// Usecase below
//

// GetAllStockOrder -.
func (uc *OrderUseCase) GetAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error) {
	return uc.repo.QueryAllStockOrder(ctx)
}

func (uc *OrderUseCase) GetAllFutureOrder(ctx context.Context) ([]*entity.FutureOrder, error) {
	return uc.repo.QueryAllFutureOrder(ctx)
}

// GetAllStockTradeBalance -.
func (uc *OrderUseCase) GetAllStockTradeBalance(ctx context.Context) ([]*entity.StockTradeBalance, error) {
	tradeBalanceArr, err := uc.repo.QueryAllStockTradeBalance(ctx)
	if err != nil {
		return nil, err
	}
	return tradeBalanceArr, nil
}

// GetAllFutureTradeBalance -.
func (uc *OrderUseCase) GetAllFutureTradeBalance(ctx context.Context) ([]*entity.FutureTradeBalance, error) {
	tradeBalanceArr, err := uc.repo.QueryAllFutureTradeBalance(ctx)
	if err != nil {
		return nil, err
	}
	return tradeBalanceArr, nil
}

func (uc *OrderUseCase) GetLastStockTradeBalance(ctx context.Context) (*entity.StockTradeBalance, error) {
	return uc.repo.QueryLastStockTradeBalance(ctx)
}

func (uc *OrderUseCase) GetLastFutureTradeBalance(ctx context.Context) (*entity.FutureTradeBalance, error) {
	return uc.repo.QueryLastFutureTradeBalance(ctx)
}

// CalculateBuyCost -.
func (uc *OrderUseCase) CalculateBuyCost(price float64, quantity int64) int64 {
	return uc.quota.GetStockBuyCost(price, quantity)
}

// CalculateSellCost -.
func (uc *OrderUseCase) CalculateSellCost(price float64, quantity int64) int64 {
	return uc.quota.GetStockSellCost(price, quantity)
}

// CalculateTradeDiscount -.
func (uc *OrderUseCase) CalculateTradeDiscount(price float64, quantity int64) int64 {
	return uc.quota.GetStockTradeFeeDiscount(price, quantity)
}

// GetFuturePosition .
func (uc *OrderUseCase) GetFuturePosition() ([]*entity.FuturePosition, error) {
	query, err := uc.gRPCAPI.GetFuturePosition()
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

func (uc *OrderUseCase) IsStockTradeTime() bool {
	return uc.stockTradeDay.IsStockMarketOpenNow()
}

func (uc *OrderUseCase) IsFutureTradeTime() bool {
	return uc.futureTradeDay.IsFutureMarketOpenNow()
}

func (uc *OrderUseCase) GetFutureOrderByTradeDay(ctx context.Context, tradeDay string) ([]*entity.FutureOrder, error) {
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
