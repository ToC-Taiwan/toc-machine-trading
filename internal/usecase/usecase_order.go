package usecase

import (
	"context"
	"errors"
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

	quota    *Quota
	simTrade bool
}

// NewOrder -.
func NewOrder(t *grpcapi.OrdergRPCAPI, r *repo.OrderRepo) *OrderUseCase {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	uc := &OrderUseCase{
		gRPCAPI: t,
		repo:    r,
		quota:   NewQuota(cfg.Quota),
	}

	uc.simTrade = cfg.TradeSwitch.Simulation

	bus.SubscribeTopic(topicPlaceOrder, uc.placeOrder)
	bus.SubscribeTopic(topicCancelOrder, uc.cancelOrder)
	bus.SubscribeTopic(topicAllOrders, uc.CalculateTradeBalance)

	go func() {
		for range time.NewTicker(1500 * time.Millisecond).C {
			uc.updateOrderStatusCache()
		}
	}()

	return uc
}

func (uc *OrderUseCase) placeOrder(order *entity.Order) {
	var orderID string
	var status entity.OrderStatus
	var err error

	if uc.quota.quota < uc.quota.CalculateOrderCost(order) {
		return
	}

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
		return
	}

	if status == entity.StatusFailed {
		return
	}

	uc.quota.quota -= uc.quota.CalculateOrderCost(order)
	order.OrderID = orderID
	order.Status = status
	cc.SetOrderByOrderID(order)
}

func (uc *OrderUseCase) cancelOrder(orderID string) {
	orderID, status, err := uc.CancelOrderID(orderID)
	if err != nil {
		log.Error(err)
		return
	}

	order := cc.GetOrderByOrderID(orderID)
	if order == nil {
		log.Error("Order not found")
		return
	}

	uc.quota.quota += uc.quota.CalculateOrderCost(order)
	order.Status = status
	cc.SetOrderByOrderID(order)
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

func (uc *OrderUseCase) updateOrderStatusCache() {
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
			bus.PublishTopicEvent(topicUpdateOrderStatus, context.Background(), o)
		}
	}
}

// CalculateTradeBalance -.
func (uc *OrderUseCase) CalculateTradeBalance(allOrders []*entity.Order) {
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

	if len(forwardOrder)/2 == 0 && len(reverseOrder)/2 == 0 {
		return
	}

	var forwardBalance, revereBalance, discount int64
	for _, v := range forwardOrder {
		switch v.Action {
		case entity.ActionBuy:
			forwardBalance -= uc.quota.GetStockBuyCost(v.Price, v.Quantity)
		case entity.ActionSell:
			forwardBalance += uc.quota.GetStockSellCost(v.Price, v.Quantity)
		}
		discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Quantity)
	}

	for _, v := range reverseOrder {
		switch v.Action {
		case entity.ActionSellFirst:
			revereBalance += uc.quota.GetStockSellCost(v.Price, v.Quantity)
		case entity.ActionBuyLater:
			revereBalance -= uc.quota.GetStockBuyCost(v.Price, v.Quantity)
		}
		discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Quantity)
	}

	tmp := &entity.TradeBalance{
		TradeDay:        cc.GetBasicInfo().TradeDay,
		TradeCount:      int64(len(forwardOrder)/2 + len(reverseOrder)/2),
		Forward:         forwardBalance,
		Reverse:         revereBalance,
		OriginalBalance: forwardBalance + revereBalance,
		Discount:        discount,
		Total:           forwardBalance + revereBalance + discount,
	}

	err := uc.repo.InserOrUpdateTradeBalance(context.Background(), tmp)
	if err != nil {
		log.Panic(err)
	}
}
