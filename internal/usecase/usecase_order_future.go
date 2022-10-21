package usecase

import (
	"context"
	"errors"
	"time"

	"tmt/internal/entity"
)

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

func (uc *OrderUseCase) updateCacheAndInsertFutureDB(order *entity.FutureOrder) {
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
