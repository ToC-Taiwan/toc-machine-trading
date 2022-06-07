package usecase

import (
	"context"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/eventbus"
	"toc-machine-trading/pkg/logger"
)

// TargetUseCase -.
type TargetUseCase struct {
	repo    TargetRepo
	gRPCAPI TargetRPCAPI
	bus     *eventbus.Bus
}

// NewTarget -.
func NewTarget(r *repo.TargetRepo, t *grpcapi.TargetgRPCAPI, bus *eventbus.Bus) {
	uc := &TargetUseCase{
		repo:    r,
		gRPCAPI: t,
		bus:     bus,
	}

	targetArr, err := uc.SearchTargets(context.Background())
	if err != nil {
		logger.Get().Panic(err)
	}

	bus.PublishTopicEvent(targetsTopic, targetArr)
}

// SearchTargets -.
func (uc *TargetUseCase) SearchTargets(ctx context.Context) ([]*entity.Target, error) {
	t, err := uc.gRPCAPI.GetStockVolumeRank("2022-06-06")
	if err != nil {
		return nil, err
	}
	var result []*entity.Target
	for i, v := range t {
		result = append(result, &entity.Target{
			StockNum:    v.GetCode(),
			TradeDay:    time.Now(),
			Rank:        i + 1,
			Volume:      v.GetTotalVolume(),
			Subscribe:   false,
			RealTimeAdd: false,
		})
	}
	return result, nil
}
