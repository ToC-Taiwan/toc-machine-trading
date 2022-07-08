package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/global"
)

// OrderUseCase -.
type OrderUseCase struct {
	gRPCAPI OrdergRPCAPI
	repo    OrderRepo

	quota          *Quota
	simTrade       bool
	placeOrderLock sync.Mutex
}

// NewOrder -.
func NewOrder(t *grpcapi.OrdergRPCAPI, r *repo.OrderRepo) *OrderUseCase {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	uc := &OrderUseCase{
		gRPCAPI:        t,
		repo:           r,
		quota:          NewQuota(cfg.Quota),
		simTrade:       cfg.TradeSwitch.Simulation,
		placeOrderLock: sync.Mutex{},
	}

	bus.SubscribeTopic(topicPlaceOrder, uc.placeOrder)
	bus.SubscribeTopic(topicCancelOrder, uc.cancelOrder)
	bus.SubscribeTopic(topicInsertOrUpdateOrder, uc.updateCacheAndInsertDB)

	go func() {
		for range time.NewTicker(time.Minute).C {
			orders, err := uc.repo.QueryAllOrderByDate(context.Background(), cc.GetBasicInfo().TradeDay)
			if err != nil {
				log.Panic(err)
			}
			uc.calculateTradeBalance(orders)
		}
	}()

	go func() {
		for range time.NewTicker(1500 * time.Millisecond).C {
			uc.askOrderUpdate()
		}
	}()

	return uc
}

func (uc *OrderUseCase) placeOrder(order *entity.Order) {
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

func (uc *OrderUseCase) cancelOrder(orderID string) {
	defer uc.placeOrderLock.Unlock()
	uc.placeOrderLock.Lock()

	cacheOrder := cc.GetOrderByOrderID(orderID)
	cacheOrder.TradeTime = time.Now()
	log.Warnf("Cancel Order -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", cacheOrder.StockNum, cacheOrder.Action, cacheOrder.Price, cacheOrder.Quantity)

	// result will return instantly
	resOrderID, status, err := uc.CancelOrderID(orderID)
	if err != nil {
		log.Error(err)
		return
	}

	if resOrderID != orderID {
		log.Error("OrderID not match")
		return
	}

	if cosumeQuota := uc.quota.calculateOriginalOrderCost(cacheOrder); cosumeQuota > 0 {
		uc.quota.quota += cosumeQuota
		log.Warnf("Quota Back: %d", uc.quota.quota)
	}

	cacheOrder.Status = status
	cc.SetOrderByOrderID(cacheOrder)
}

// BuyStock -.
func (uc *OrderUseCase) BuyStock(order *entity.Order) (string, entity.OrderStatus, error) {
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
func (uc *OrderUseCase) SellStock(order *entity.Order) (string, entity.OrderStatus, error) {
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
func (uc *OrderUseCase) SellFirstStock(order *entity.Order) (string, entity.OrderStatus, error) {
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
func (uc *OrderUseCase) BuyLaterStock(order *entity.Order) (string, entity.OrderStatus, error) {
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

func (uc *OrderUseCase) askOrderUpdate() {
	if !uc.simTrade {
		msg, err := uc.gRPCAPI.GetNonBlockOrderStatusArr()
		if err != nil {
			log.Panic(err)
		}

		if errMsg := msg.GetErr(); errMsg != "" {
			log.Panic(errMsg)
		}
	} else {
		orders, err := uc.gRPCAPI.GetOrderStatusArr()
		if err != nil {
			log.Panic(err)
		}
		actionMap := entity.ActionListMap
		statusMap := entity.StatusListMap
		for _, v := range orders {
			orderTime, err := time.ParseInLocation(global.LongTimeLayout, v.GetOrderTime(), time.Local)
			if err != nil {
				log.Error(err)
				continue
			}
			o := &entity.Order{
				StockNum:  v.GetCode(),
				OrderID:   v.GetOrderId(),
				Action:    actionMap[v.GetAction()],
				Price:     v.GetPrice(),
				Quantity:  v.GetQuantity(),
				Status:    statusMap[v.GetStatus()],
				OrderTime: orderTime,
			}
			uc.updateCacheAndInsertDB(o)
		}
	}
}

func (uc *OrderUseCase) updateCacheAndInsertDB(order *entity.Order) {
	// get order from cache
	cacheOrder := cc.GetOrderByOrderID(order.OrderID)
	if cacheOrder == nil {
		return
	}

	if cacheOrder.Status != order.Status || !cacheOrder.OrderTime.Equal(order.OrderTime) {
		cacheOrder.Status = order.Status
		cacheOrder.OrderTime = order.OrderTime

		// update cache
		cc.SetOrderByOrderID(cacheOrder)

		// insert or update order to db
		if err := uc.repo.InsertOrUpdateOrderByOrderID(context.Background(), cacheOrder); err != nil {
			log.Panic(err)
		}
	}
}

// CalculateTradeBalance -.
func (uc *OrderUseCase) calculateTradeBalance(allOrders []*entity.Order) {
	var forwardOrder, reverseOrder []*entity.Order
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

	err := uc.repo.InsertOrUpdateTradeBalance(context.Background(), tmp)
	if err != nil {
		log.Panic(err)
	}
}

// GetAllOrder -.
func (uc *OrderUseCase) GetAllOrder(ctx context.Context) ([]*entity.Order, error) {
	orderArr, err := uc.repo.QueryAllOrder(ctx)
	if err != nil {
		return nil, err
	}
	return orderArr, nil
}

// GetAllTradeBalance -.
func (uc *OrderUseCase) GetAllTradeBalance(ctx context.Context) ([]*entity.TradeBalance, error) {
	tradeBalanceArr, err := uc.repo.QueryAllTradeBalance(ctx)
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
