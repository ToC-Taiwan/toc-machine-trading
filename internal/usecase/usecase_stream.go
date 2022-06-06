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
