package usecase

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/logger"
)

// StreamUseCase -.
type StreamUseCase struct {
	repo    StreamRepo
	gRPCAPI StreamgRPCAPI
}

// NewStream -.
func NewStream(r *repo.StreamRepo, t *grpcapi.StreamgRPCAPI) {
	uc := &StreamUseCase{
		repo:    r,
		gRPCAPI: t,
	}
	go func() {
		if err := uc.ReceiveEvent(context.Background()); err != nil {
			logger.Get().Panic(err)
		}
	}()
}

// ReceiveEvent -.
func (uc *StreamUseCase) ReceiveEvent(ctx context.Context) error {
	eventChan := make(chan *entity.SinopacEvent)
	if err := uc.gRPCAPI.EventChannel(eventChan); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-eventChan:
			if err := uc.repo.InsertEvent(ctx, event); err != nil {
				return err
			}
		}
	}
}

// ReceiveTicks -.
func (uc *StreamUseCase) ReceiveTicks(ctx context.Context) error {
	tickChan := make(chan *entity.RealTimeTick)
	if err := uc.gRPCAPI.TickChannel(tickChan); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case tick := <-tickChan:
			// if err := uc.repo.InsertEvent(ctx, event); err != nil {
			// 	return err
			// }
			logger.Get().Info(tick)
		}
	}
}

// ReceiveBidAsk -.
func (uc *StreamUseCase) ReceiveBidAsk(ctx context.Context) error {
	bidAskChan := make(chan *entity.RealTimeBidAsk)
	if err := uc.gRPCAPI.BidAskChannel(bidAskChan); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case bidAsk := <-bidAskChan:
			// if err := uc.repo.InsertEvent(ctx, event); err != nil {
			// 	return err
			// }
			logger.Get().Info(bidAsk)
		}
	}
}

// ReceiveOrderStatus -.
func (uc *StreamUseCase) ReceiveOrderStatus(ctx context.Context) error {
	orderStatusChan := make(chan *entity.OrderStatus)
	if err := uc.gRPCAPI.OrderStatusChannel(orderStatusChan); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case orderStatus := <-orderStatusChan:
			// if err := uc.repo.InsertEvent(ctx, event); err != nil {
			// 	return err
			// }
			logger.Get().Info(orderStatus)
		}
	}
}
