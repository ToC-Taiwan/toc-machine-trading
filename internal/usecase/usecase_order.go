package usecase

import (
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/pkg/config"
)

// OrderUseCase -.
type OrderUseCase struct {
	gRPCAPI  OrdergRPCAPI
	simTrade bool
}

// NewOrder -.
func NewOrder(t *grpcapi.OrdergRPCAPI) {
	uc := &OrderUseCase{
		gRPCAPI: t,
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	uc.simTrade = cfg.TradeSwitch.Simulation

	go func() {
		for range time.NewTicker(1500 * time.Millisecond).C {
			msg, err := uc.gRPCAPI.GetNonBlockOrderStatusArr()
			if err != nil {
				log.Panic(err)
			}

			if errMsg := msg.GetErr(); errMsg != "" {
				log.Panic(errMsg)
			}
		}
	}()

	if err := bus.SubscribeTopic(topicPlaceOrder, uc.placeOrder); err != nil {
		log.Panic(err)
	}
	if err := bus.SubscribeTopic(topicCancelOrder, uc.cancelOrder); err != nil {
		log.Panic(err)
	}
}

func (uc *OrderUseCase) placeOrder(order *entity.Order) {
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
		return
	}

	if status == entity.StatusFailed {
		return
	}

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

	order.Status = status
	cc.SetOrderByOrderID(order)
}

// BuyStock -.
func (uc *OrderUseCase) BuyStock(order *entity.Order) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.BuyStock(order, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
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
	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
}

// SellFirstStock -.
func (uc *OrderUseCase) SellFirstStock(order *entity.Order) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.SellFirstStock(order, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
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
	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
}

// CancelOrderID -.
func (uc *OrderUseCase) CancelOrderID(orderID string) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.CancelStock(orderID, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}
	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
}
