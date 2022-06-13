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

	if err := uc.bus.SubscribeTopic(topicBuyOrder, nil); err != nil {
		log.Panic(err)
	}
	if err := uc.bus.SubscribeTopic(topicSellOrder, nil); err != nil {
		log.Panic(err)
	}
	if err := uc.bus.SubscribeTopic(topicSellFirstOrder, nil); err != nil {
		log.Panic(err)
	}
	if err := uc.bus.SubscribeTopic(topicBuyOrder, nil); err != nil {
		log.Panic(err)
	}
}
