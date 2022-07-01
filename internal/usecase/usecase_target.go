package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/global"
)

// TargetUseCase -.
type TargetUseCase struct {
	repo    TargetRepo
	gRPCAPI TargetRPCAPI
}

// NewTarget -.
func NewTarget(r *repo.TargetRepo, t *grpcapi.TargetgRPCAPI) *TargetUseCase {
	uc := &TargetUseCase{
		repo:    r,
		gRPCAPI: t,
	}

	// unsubscriba all first
	if err := uc.UnSubscribeAll(context.Background()); err != nil {
		log.Panic("unsubscribe all fail")
	}

	// query targets from db
	tradeDay := cc.GetBasicInfo().TradeDay
	targetArr, err := uc.repo.QueryTargetsByTradeDay(context.Background(), tradeDay)
	if err != nil {
		log.Panic(err)
	}

	// db has no targets, find targets from gRPC
	if len(targetArr) == 0 {
		targetArr, err = uc.SearchTradeDayTargets(context.Background(), tradeDay)
		if err != nil {
			log.Panic(err)
		}

		if len(targetArr) != 0 {
			if err = uc.InsertTargets(context.Background(), targetArr); err != nil {
				log.Panic(err)
			}
		} else {
			stuck := make(chan struct{})
			log.Error("no targets")
			<-stuck
		}
	}

	// sub events
	bus.SubscribeTopic(topicRealTimeTargets, uc.InsertTargets)
	bus.SubscribeTopic(topicSubscribeTickTargets, uc.SubscribeStockTick, uc.SubscribeStockBidAsk)
	bus.SubscribeTopic(topicUnSubscribeTickTargets, uc.UnSubscribeStockTick, uc.UnSubscribeStockBidAsk)

	// save to cache
	cc.AppendTargets(targetArr)
	// pub events
	bus.PublishTopicEvent(topicTargets, context.Background(), targetArr)
	return uc
}

// GetTargets - get targets from cache
func (uc *TargetUseCase) GetTargets(ctx context.Context) []*entity.Target {
	return cc.GetTargets()
}

// InsertTargets - insert targets to db
func (uc *TargetUseCase) InsertTargets(ctx context.Context, targetArr []*entity.Target) error {
	if err := uc.repo.InsertTargetArr(ctx, targetArr); err != nil {
		return err
	}
	return nil
}

// SearchTradeDayTargets - search targets by trade day
func (uc *TargetUseCase) SearchTradeDayTargets(ctx context.Context, tradeDay time.Time) ([]*entity.Target, error) {
	lastTradeDay := cc.GetBasicInfo().LastTradeDay
	t, err := uc.gRPCAPI.GetStockVolumeRank(lastTradeDay.Format(global.ShortTimeLayout))
	if err != nil {
		return nil, err
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	condition := cfg.TargetCond

	var result []*entity.Target
	for _, c := range condition {
		for _, v := range t {
			if !targetFilter(v.GetClose(), v.GetTotalVolume(), c, false) {
				continue
			}

			if stock := cc.GetStockDetail(v.GetCode()); stock != nil {
				result = append(result, &entity.Target{
					Rank:        len(result) + 1,
					StockNum:    v.GetCode(),
					Volume:      v.GetTotalVolume(),
					Subscribe:   c.Subscribe,
					RealTimeAdd: false,
					TradeDay:    tradeDay,
					Stock:       stock,
				})
			}
		}
	}
	return result, nil
}

// UnSubscribeAll -.
func (uc *TargetUseCase) UnSubscribeAll(ctx context.Context) error {
	failMessge, err := uc.gRPCAPI.UnSubscribeStockAllTick()
	if err != nil {
		return err
	}

	if m := failMessge.GetErr(); m != "" {
		return errors.New(m)
	}

	failMessge, err = uc.gRPCAPI.UnSubscribeStockAllBidAsk()
	if err != nil {
		return err
	}

	if m := failMessge.GetErr(); m != "" {
		return errors.New(m)
	}

	return nil
}

// SubscribeStockTick -.
func (uc *TargetUseCase) SubscribeStockTick(targetArr []*entity.Target) error {
	var subArr []string
	for _, v := range targetArr {
		if v.Subscribe {
			subArr = append(subArr, v.StockNum)
		}
	}

	failSubNumArr, err := uc.gRPCAPI.SubscribeStockTick(subArr)
	if err != nil {
		return err
	}

	if len(failSubNumArr) != 0 {
		return fmt.Errorf("subscribe fail %v", failSubNumArr)
	}

	return nil
}

// SubscribeStockBidAsk -.
func (uc *TargetUseCase) SubscribeStockBidAsk(targetArr []*entity.Target) error {
	var subArr []string
	for _, v := range targetArr {
		if v.Subscribe {
			subArr = append(subArr, v.StockNum)
		}
	}

	failSubNumArr, err := uc.gRPCAPI.SubscribeStockBidAsk(subArr)
	if err != nil {
		return err
	}

	if len(failSubNumArr) != 0 {
		return fmt.Errorf("subscribe fail %v", failSubNumArr)
	}

	return nil
}

// UnSubscribeStockTick -.
func (uc *TargetUseCase) UnSubscribeStockTick(target *entity.Target) error {
	failUnSubNumArr, err := uc.gRPCAPI.UnSubscribeStockTick([]string{target.StockNum})
	if err != nil {
		return err
	}

	if len(failUnSubNumArr) != 0 {
		return fmt.Errorf("unsubscribe fail %v", failUnSubNumArr)
	}

	return nil
}

// UnSubscribeStockBidAsk -.
func (uc *TargetUseCase) UnSubscribeStockBidAsk(target *entity.Target) error {
	failUnSubNumArr, err := uc.gRPCAPI.UnSubscribeStockBidAsk([]string{target.StockNum})
	if err != nil {
		return err
	}

	if len(failUnSubNumArr) != 0 {
		return fmt.Errorf("unsubscribe fail %v", failUnSubNumArr)
	}

	return nil
}
