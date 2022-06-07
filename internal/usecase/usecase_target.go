package usecase

import (
	"context"
	"fmt"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/eventbus"
	"toc-machine-trading/pkg/global"
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

	targetArr, err := uc.repo.QueryTargetsByTradeDay(context.Background(), CacheGetTradeDay())
	if err != nil {
		logger.Get().Panic(err)
	}

	if len(targetArr) == 0 {
		targetArr, err = uc.SearchTradeDayTargets(context.Background(), CacheGetTradeDay())
		if err != nil {
			logger.Get().Panic(err)
		}

		if len(targetArr) != 0 {
			if err = uc.repo.InsertTargetArr(context.Background(), targetArr); err != nil {
				logger.Get().Panic(err)
			}
		}
	}

	for _, v := range targetArr {
		logger.Get().Warnf("%+v\n", v)
	}

	bus.PublishTopicEvent(targetsTopic, targetArr)
}

// SearchTradeDayTargets -.
func (uc *TargetUseCase) SearchTradeDayTargets(ctx context.Context, dateTime time.Time) ([]*entity.Target, error) {
	lastTradeDay := GetLastNTradeDayByDate(1, dateTime)[0]
	t, err := uc.gRPCAPI.GetStockVolumeRank(lastTradeDay.Format(global.ShortTimeLayout))
	if err != nil {
		return nil, err
	}
	var result []*entity.Target
	for i, v := range t {
		result = append(result, &entity.Target{
			StockNum:    v.GetCode(),
			TradeDay:    dateTime,
			Rank:        i + 1,
			Volume:      v.GetTotalVolume(),
			Subscribe:   false,
			RealTimeAdd: false,
		})
	}
	return result, nil
}

// SubscribeStockTick -.
func (uc *TargetUseCase) SubscribeStockTick(ctx context.Context, targetArr []*entity.Target) error {
	var tmp []string
	for _, v := range targetArr {
		tmp = append(tmp, v.StockNum)
	}

	fail, err := uc.gRPCAPI.SubscribeStockTick(tmp)
	if err != nil {
		return err
	}

	if len(fail) != 0 {
		return fmt.Errorf("subscribe fail %+v", fail)
	}

	return nil
}
