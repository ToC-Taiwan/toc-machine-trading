package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/global"
	"tmt/internal/entity"
	"tmt/internal/usecase/modules/event"
	"tmt/internal/usecase/modules/quota"
	"tmt/internal/usecase/modules/tradeday"
)

// OrderUseCase -.
type OrderUseCase struct {
	gRPCAPI OrdergRPCAPI
	repo    OrderRepo

	quota *quota.Quota

	placeOrderLock       sync.Mutex
	placeFutureOrderLock sync.Mutex

	stockTradeDay  tradeday.TradePeriod
	futureTradeDay tradeday.TradePeriod

	simTrade bool
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

		stockTradeDay:  tradeDay.GetStockTradeDay(),
		futureTradeDay: tradeDay.GetFutureTradeDay(),
	}

	bus.SubscribeTopic(event.TopicPlaceStockOrder, uc.placeStockOrder)
	bus.SubscribeTopic(event.TopicCancelStockOrder, uc.cancelStockOrder)
	bus.SubscribeTopic(event.TopicInsertOrUpdateStockOrder, uc.updateStockOrderCacheAndInsertDB)

	bus.SubscribeTopic(event.TopicPlaceFutureOrder, uc.placeFutureOrder)
	bus.SubscribeTopic(event.TopicCancelFutureOrder, uc.cancelFutureOrder)
	bus.SubscribeTopic(event.TopicInsertOrUpdateFutureOrder, uc.updateFutureOrderCacheAndInsertDB)

	uc.updateAllTradeBalance()
	return uc
}

func (uc *OrderUseCase) updateAllTradeBalance() {
	go func() {
		for range time.NewTicker(20 * time.Second).C {
			stockOrders, err := uc.repo.QueryAllStockOrderByDate(context.Background(), uc.stockTradeDay.ToStartEndArray())
			if err != nil {
				log.Panic(err)
			}
			uc.calculateStockTradeBalance(stockOrders)

			futureOrders, err := uc.repo.QueryAllFutureOrderByDate(context.Background(), uc.futureTradeDay.ToStartEndArray())
			if err != nil {
				log.Panic(err)
			}
			uc.calculateFutureTradeBalance(futureOrders)
		}
	}()
	go func() {
		for range time.NewTicker(3 * time.Second).C {
			err := uc.AskOrderUpdate()
			if err != nil {
				log.Error(err)
			}
		}
	}()
}

// AskOrderUpdate -.
func (uc *OrderUseCase) AskOrderUpdate() error {
	if !uc.simTrade {
		msg, err := uc.gRPCAPI.GetNonBlockOrderStatusArr()
		if err != nil {
			return err
		}

		if errMsg := msg.GetErr(); errMsg != "" {
			return errors.New(errMsg)
		}
	} else {
		orders, err := uc.gRPCAPI.GetOrderStatusArr()
		if err != nil {
			return err
		}
		actionMap := entity.ActionListMap
		statusMap := entity.StatusListMap
		for _, v := range orders {
			orderTime, err := time.ParseInLocation(global.LongTimeLayout, v.GetOrderTime(), time.Local)
			if err != nil {
				return err
			}

			switch {
			case cc.GetOrderByOrderID(v.GetOrderId()) != nil:
				o := &entity.StockOrder{
					StockNum: v.GetCode(),
					BaseOrder: entity.BaseOrder{
						OrderID:   v.GetOrderId(),
						Action:    actionMap[v.GetAction()],
						Price:     v.GetPrice(),
						Quantity:  v.GetQuantity(),
						Status:    statusMap[v.GetStatus()],
						OrderTime: orderTime,
					},
				}
				uc.updateStockOrderCacheAndInsertDB(o)
			case cc.GetFutureOrderByOrderID(v.GetOrderId()) != nil:
				o := &entity.FutureOrder{
					Code: v.GetCode(),
					BaseOrder: entity.BaseOrder{
						OrderID:   v.GetOrderId(),
						Action:    actionMap[v.GetAction()],
						Price:     v.GetPrice(),
						Quantity:  v.GetQuantity(),
						Status:    statusMap[v.GetStatus()],
						OrderTime: orderTime,
					},
				}
				uc.updateFutureOrderCacheAndInsertDB(o)
			}
		}
	}
	return nil
}

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
		log.Error(err)
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

	log.Warnf("Place Order -> Stock: %s, Action: %d, Price: %.2f, Qty: %d, Quota: %d", order.StockNum, order.Action, order.Price, order.Quantity, uc.quota.GetCurrentQuota())
}

func (uc *OrderUseCase) cancelStockOrder(order *entity.StockOrder) {
	defer uc.placeOrderLock.Unlock()
	uc.placeOrderLock.Lock()

	order.TradeTime = time.Now()
	log.Warnf("Cancel Order -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)

	// result will return instantly
	_, _, err := uc.CancelOrderID(order.OrderID)
	if err != nil {
		log.Error(err)
		return
	}

	if cosumeQuota := uc.quota.CalculateOriginalOrderCost(order); cosumeQuota > 0 {
		uc.quota.BackQuota(cosumeQuota)
		log.Warnf("Quota Back: %d", uc.quota.GetCurrentQuota())
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

	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
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

	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
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

	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
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

	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
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

	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
}

func (uc *OrderUseCase) updateStockOrderCacheAndInsertDB(order *entity.StockOrder) {
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
		log.Panic(err)
	}
}

// calculateStockTradeBalance -.
func (uc *OrderUseCase) calculateStockTradeBalance(allOrders []*entity.StockOrder) {
	var forwardOrder, reverseOrder []*entity.StockOrder
	for _, v := range allOrders {
		if v.Status != entity.StatusFilled {
			continue
		}

		switch v.Action {
		case entity.ActionBuy, entity.ActionSell:
			forwardOrder = append(forwardOrder, v)
		case entity.ActionSellFirst, entity.ActionBuyLater:
			reverseOrder = append(reverseOrder, v)
		}
	}

	forwardOrder = forwardOrder[:2*(len(forwardOrder)/2)]
	reverseOrder = reverseOrder[:2*(len(reverseOrder)/2)]

	var forwardBalance, revereBalance, discount, tradeCount int64
	for _, v := range forwardOrder {
		switch v.Action {
		case entity.ActionBuy:
			tradeCount++
			forwardBalance -= uc.quota.GetStockBuyCost(v.Price, v.Quantity)
		case entity.ActionSell:
			forwardBalance += uc.quota.GetStockSellCost(v.Price, v.Quantity)
		}
		discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Quantity)
	}

	for _, v := range reverseOrder {
		switch v.Action {
		case entity.ActionSellFirst:
			tradeCount++
			revereBalance += uc.quota.GetStockSellCost(v.Price, v.Quantity)
		case entity.ActionBuyLater:
			revereBalance -= uc.quota.GetStockBuyCost(v.Price, v.Quantity)
		}
		discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Quantity)
	}

	tmp := &entity.TradeBalance{
		TradeDay:        cc.GetBasicInfo().TradeDay,
		TradeCount:      tradeCount,
		Forward:         forwardBalance,
		Reverse:         revereBalance,
		OriginalBalance: forwardBalance + revereBalance,
		Discount:        discount,
		Total:           forwardBalance + revereBalance + discount,
	}

	err := uc.repo.InsertOrUpdateStockTradeBalance(context.Background(), tmp)
	if err != nil {
		log.Panic(err)
	}
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
		log.Error(err)
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

	log.Warnf("Place Future Order -> Future: %s, Action: %d, Price: %.0f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)
}

func (uc *OrderUseCase) cancelFutureOrder(order *entity.FutureOrder) {
	defer uc.placeFutureOrderLock.Unlock()
	uc.placeFutureOrderLock.Lock()

	order.TradeTime = time.Now()
	log.Warnf("Cancel Future Order -> Future: %s, Action: %d, Price: %.0f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)

	// result will return instantly
	_, _, err := uc.CancelFutureOrderID(order.OrderID)
	if err != nil {
		log.Error(err)
		return
	}
}

func (uc *OrderUseCase) updateFutureOrderCacheAndInsertDB(order *entity.FutureOrder) {
	// get order from cache
	cacheOrder := cc.GetFutureOrderByOrderID(order.OrderID)
	if cacheOrder == nil {
		return
	}

	cacheOrder.Status = order.Status
	cacheOrder.OrderTime = order.OrderTime

	// qty may not filled with original order, change it by return quantity
	cacheOrder.Quantity = order.Quantity

	// update cache
	cc.SetFutureOrderByOrderID(cacheOrder)

	// insert or update order to db
	if err := uc.repo.InsertOrUpdateFutureOrderByOrderID(context.Background(), cacheOrder); err != nil {
		log.Panic(err)
	}
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

	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
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

	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
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

	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
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

	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
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

	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
}

// calculateFutureTradeBalance -.
func (uc *OrderUseCase) calculateFutureTradeBalance(allOrders []*entity.FutureOrder) {
	var forwardOrder, reverseOrder []*entity.FutureOrder
	for _, v := range allOrders {
		if v.Status != entity.StatusFilled {
			continue
		}

		switch v.Action {
		case entity.ActionBuy, entity.ActionSell:
			forwardOrder = append(forwardOrder, v)
		case entity.ActionSellFirst, entity.ActionBuyLater:
			reverseOrder = append(reverseOrder, v)
		}
	}

	forwardOrder = forwardOrder[:2*(len(forwardOrder)/2)]
	reverseOrder = reverseOrder[:2*(len(reverseOrder)/2)]

	var forwardBalance, revereBalance, tradeCount int64
	for _, v := range forwardOrder {
		switch v.Action {
		case entity.ActionBuy:
			tradeCount++
			forwardBalance -= uc.quota.GetFutureBuyCost(v.Price, v.Quantity)
		case entity.ActionSell:
			forwardBalance += uc.quota.GetFutureSellCost(v.Price, v.Quantity)
		}
	}

	for _, v := range reverseOrder {
		switch v.Action {
		case entity.ActionSellFirst:
			tradeCount++
			revereBalance += uc.quota.GetFutureSellCost(v.Price, v.Quantity)
		case entity.ActionBuyLater:
			revereBalance -= uc.quota.GetFutureBuyCost(v.Price, v.Quantity)
		}
	}

	tmp := &entity.TradeBalance{
		TradeDay:   uc.futureTradeDay.TradeDay,
		TradeCount: tradeCount,
		Forward:    forwardBalance,
		Reverse:    revereBalance,
		Total:      forwardBalance + revereBalance,
	}

	err := uc.repo.InsertOrUpdateFutureTradeBalance(context.Background(), tmp)
	if err != nil {
		log.Panic(err)
	}
}

//
// Usecase below
//

// GetAllStockOrder -.
func (uc *OrderUseCase) GetAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error) {
	orderArr, err := uc.repo.QueryAllStockOrder(ctx)
	if err != nil {
		return nil, err
	}
	return orderArr, nil
}

// GetAllStockTradeBalance -.
func (uc *OrderUseCase) GetAllStockTradeBalance(ctx context.Context) ([]*entity.TradeBalance, error) {
	tradeBalanceArr, err := uc.repo.QueryAllStockTradeBalance(ctx)
	if err != nil {
		return nil, err
	}
	return tradeBalanceArr, nil
}

// GetAllFutureTradeBalance -.
func (uc *OrderUseCase) GetAllFutureTradeBalance(ctx context.Context) ([]*entity.TradeBalance, error) {
	tradeBalanceArr, err := uc.repo.QueryAllFutureTradeBalance(ctx)
	if err != nil {
		return nil, err
	}
	return tradeBalanceArr, nil
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
