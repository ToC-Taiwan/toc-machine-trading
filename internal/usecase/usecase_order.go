package usecase

import (
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/config"
)

// OrderUseCase -.
type OrderUseCase struct {
	repo     OrderRepo
	gRPCAPI  OrdergRPCAPI
	simTrade bool
}

// NewOrder -.
func NewOrder(r *repo.OrderRepo, t *grpcapi.OrdergRPCAPI) {
	uc := &OrderUseCase{
		repo:    r,
		gRPCAPI: t,
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	uc.simTrade = cfg.TradeSwitch.Simulation

	go func() {
		for range time.NewTicker(1500 * time.Millisecond).C {
			uc.callOrderStatus()
		}
	}()

	if err := bus.SubscribeTopic(topicOrder, uc.processOrder); err != nil {
		log.Panic(err)
	}
}

func (uc *OrderUseCase) processOrder(order *entity.Order) {
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
	if status == entity.StatusCancelled || status == entity.StatusFailed {
		return
	}
	order.OrderID = orderID
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

// CancelOrder -.
func (uc *OrderUseCase) CancelOrder(orderID string) (string, entity.OrderStatus, error) {
	result, err := uc.gRPCAPI.CancelStock(orderID, uc.simTrade)
	if err != nil {
		return "", entity.StatusUnknow, err
	}
	statusMap := entity.StatusListMap
	return result.GetOrderId(), statusMap[result.GetStatus()], nil
}

func (uc *OrderUseCase) callOrderStatus() {
	msg, err := uc.gRPCAPI.GetNonBlockOrderStatusArr()
	if err != nil {
		log.Error(err)
	}

	if errMsg := msg.GetErr(); errMsg != "" {
		log.Error(errMsg)
	}
}
