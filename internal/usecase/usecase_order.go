package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/pkg/config"
	"tmt/pkg/global"
)

// OrderUseCase -.
type OrderUseCase struct {
	gRPCAPI OrdergRPCAPI
	repo    OrderRepo

	basicInfo            entity.BasicInfo
	quota                *Quota
	simTrade             bool
	placeOrderLock       sync.Mutex
	placeFutureOrderLock sync.Mutex

	futureTradeDay time.Time
}

// NewOrder -.
func NewOrder(t OrdergRPCAPI, r OrderRepo) *OrderUseCase {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	futureTradeDay, err := time.ParseInLocation(global.ShortTimeLayout, time.Now().Format(global.ShortTimeLayout), time.Local)
	if err != nil {
		log.Panic(err)
	}

	uc := &OrderUseCase{
		gRPCAPI:        t,
		repo:           r,
		quota:          NewQuota(cfg.Quota),
		simTrade:       cfg.Simulation,
		placeOrderLock: sync.Mutex{},
		basicInfo:      *cc.GetBasicInfo(),
		futureTradeDay: futureTradeDay,
	}

	bus.SubscribeTopic(topicPlaceOrder, uc.placeOrder)
	bus.SubscribeTopic(topicCancelOrder, uc.cancelOrder)
	bus.SubscribeTopic(topicInsertOrUpdateOrder, uc.updateCacheAndInsertDB)

	bus.SubscribeTopic(topicPlaceFutureOrder, uc.placeFutureOrder)
	bus.SubscribeTopic(topicCancelFutureOrder, uc.cancelFutureOrder)
	bus.SubscribeTopic(topicInsertOrUpdateFutureOrder, uc.updateCacheAndInsertFutureDB)

	go func() {
		for range time.NewTicker(time.Minute).C {
			stockOrders, err := uc.repo.QueryAllStockOrderByDate(context.Background(), uc.basicInfo.TradeDay)
			if err != nil {
				log.Panic(err)
			}
			uc.calculateStockTradeBalance(stockOrders)

			futureOrders, err := uc.repo.QueryAllFutureOrderByDate(context.Background(), uc.futureTradeDay)
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

	return uc
}

func (uc *OrderUseCase) placeOrder(order *entity.StockOrder) {
	defer uc.placeOrderLock.Unlock()
	uc.placeOrderLock.Lock()

	cosumeQuota := uc.quota.calculateOriginalOrderCost(order)
	if cosumeQuota != 0 && uc.quota.quota-cosumeQuota < 0 {
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
	uc.quota.quota -= cosumeQuota

	// modify order and save to cache
	order.OrderID = orderID
	order.Status = status
	order.TradeTime = time.Now()
	cc.SetOrderByOrderID(order)

	log.Warnf("Place Order -> Stock: %s, Action: %d, Price: %.2f, Qty: %d, Quota: %d", order.StockNum, order.Action, order.Price, order.Quantity, uc.quota.quota)
}

func (uc *OrderUseCase) cancelOrder(order *entity.StockOrder) {
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

	if cosumeQuota := uc.quota.calculateOriginalOrderCost(order); cosumeQuota > 0 {
		uc.quota.quota += cosumeQuota
		log.Warnf("Quota Back: %d", uc.quota.quota)
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
				uc.updateCacheAndInsertDB(o)
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
				uc.updateCacheAndInsertFutureDB(o)
			}
		}
	}
	return nil
}

func (uc *OrderUseCase) updateCacheAndInsertDB(order *entity.StockOrder) {
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

	if len(forwardOrder) == 0 && len(reverseOrder) == 0 {
		return
	}

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

	if len(forwardOrder) == 0 && len(reverseOrder) == 0 {
		return
	}

	var forwardBalance, revereBalance, tradeCount int64
	for _, v := range forwardOrder {
		switch v.Action {
		case entity.ActionBuy:
			tradeCount++
			forwardBalance -= 50*int64(v.Price) + 15 + 30
		case entity.ActionSell:
			forwardBalance += 50*int64(v.Price) - 15 - 30
		}
	}

	for _, v := range reverseOrder {
		switch v.Action {
		case entity.ActionSellFirst:
			tradeCount++
			revereBalance += 50*int64(v.Price) - 15 - 30
		case entity.ActionBuyLater:
			revereBalance -= 50*int64(v.Price) + 15 + 30
		}
	}

	tmp := &entity.TradeBalance{
		TradeDay:   uc.futureTradeDay,
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

// GetAllOrder -.
func (uc *OrderUseCase) GetAllOrder(ctx context.Context) ([]*entity.StockOrder, error) {
	orderArr, err := uc.repo.QueryAllStockOrder(ctx)
	if err != nil {
		return nil, err
	}
	return orderArr, nil
}

// GetAllTradeBalance -.
func (uc *OrderUseCase) GetAllTradeBalance(ctx context.Context) ([]*entity.TradeBalance, error) {
	tradeBalanceArr, err := uc.repo.QueryAllStockTradeBalance(ctx)
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
