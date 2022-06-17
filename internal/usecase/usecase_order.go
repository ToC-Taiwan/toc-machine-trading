package usecase

import (
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/eventbus"
)

// OrderUseCase -.
type OrderUseCase struct {
	repo    OrderRepo
	gRPCAPI OrdergRPCAPI
	bus     *eventbus.Bus
}

// NewOrder -.
func NewOrder(r *repo.OrderRepo, t *grpcapi.OrdergRPCAPI, bus *eventbus.Bus) {
	uc := &OrderUseCase{
		repo:    r,
		gRPCAPI: t,
		bus:     bus,
	}

	if err := uc.bus.SubscribeTopic(topicBuyOrder, uc.buyOrder); err != nil {
		log.Panic(err)
	}
	if err := uc.bus.SubscribeTopic(topicSellOrder, uc.sellOrder); err != nil {
		log.Panic(err)
	}
	if err := uc.bus.SubscribeTopic(topicSellFirstOrder, uc.sellFirstOrder); err != nil {
		log.Panic(err)
	}
	if err := uc.bus.SubscribeTopic(topicBuyOrder, uc.buyLaterOrder); err != nil {
		log.Panic(err)
	}
	if err := uc.bus.SubscribeTopic(topicCancelOrder, uc.cancelOrder); err != nil {
		log.Panic(err)
	}
}

func (uc *OrderUseCase) buyOrder() {}

func (uc *OrderUseCase) sellOrder() {}

func (uc *OrderUseCase) sellFirstOrder() {}

func (uc *OrderUseCase) buyLaterOrder() {}

func (uc *OrderUseCase) cancelOrder() {}
