package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/global"
)

// TargetUseCase -.
type TargetUseCase struct {
	repo          TargetRepo
	gRPCAPI       TargetgRPCAPI
	streamgRPCAPI StreamgRPCAPI

	targetCond config.TargetCond
}

// NewTarget -.
func NewTarget(r TargetRepo, t TargetgRPCAPI, s StreamgRPCAPI) *TargetUseCase {
	uc := &TargetUseCase{
		repo:          r,
		gRPCAPI:       t,
		streamgRPCAPI: s,
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}
	uc.targetCond = cfg.TargetCond

	// unsubscriba all first
	if err = uc.UnSubscribeAll(context.Background()); err != nil {
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

		if len(targetArr) == 0 {
			stuck := make(chan struct{})
			log.Error("no targets")
			<-stuck
		}
	}

	uc.publishNewTargets(targetArr)

	// sub events
	bus.SubscribeTopic(topicNewTargets, uc.publishNewTargets)
	bus.SubscribeTopic(topicSubscribeTickTargets, uc.SubscribeStockTick, uc.SubscribeStockBidAsk)
	bus.SubscribeTopic(topicUnSubscribeTickTargets, uc.UnSubscribeStockTick, uc.UnSubscribeStockBidAsk)

	return uc
}

func (uc *TargetUseCase) publishNewTargets(targetArr []*entity.Target) {
	err := uc.repo.InsertOrUpdateTargetArr(context.Background(), targetArr)
	if err != nil {
		log.Panic(err)
	}

	cc.AppendTargets(targetArr)

	bus.PublishTopicEvent(topicFetchHistory, context.Background(), targetArr)
	bus.PublishTopicEvent(topicStreamTargets, context.Background(), targetArr)
}

// GetTargets - get targets from cache
func (uc *TargetUseCase) GetTargets(ctx context.Context) []*entity.Target {
	return cc.GetTargets()
}

// SearchTradeDayTargets - search targets by trade day
func (uc *TargetUseCase) SearchTradeDayTargets(ctx context.Context, tradeDay time.Time) ([]*entity.Target, error) {
	lastTradeDay := cc.GetBasicInfo().LastTradeDay
	t, err := uc.gRPCAPI.GetStockVolumeRank(lastTradeDay.Format(global.ShortTimeLayout))
	if err != nil {
		return nil, err
	}

	if len(t) == 0 && time.Now().Before(cc.GetBasicInfo().OpenTime.Add(-30*time.Minute)) {
		return uc.SearchTradeDayTargetsFromAllSnapshot(tradeDay)
	}

	condition := uc.targetCond
	var result []*entity.Target
	for _, c := range condition.PriceVolumeLimit {
		for _, v := range t {
			stock := cc.GetStockDetail(v.GetCode())
			if stock == nil {
				continue
			}

			if !blackStockFilter(stock.Number, condition) ||
				!blackCatagoryFilter(stock.Category, condition) ||
				!targetFilter(v.GetClose(), v.GetTotalVolume(), c, false) {
				continue
			}

			result = append(result, &entity.Target{
				Rank:     len(result) + 1,
				StockNum: v.GetCode(),
				Volume:   v.GetTotalVolume(),
				TradeDay: tradeDay,
				Stock:    stock,
			})
		}
	}
	return result, nil
}

// SearchTradeDayTargetsFromAllSnapshot -.
func (uc *TargetUseCase) SearchTradeDayTargetsFromAllSnapshot(tradeDay time.Time) ([]*entity.Target, error) {
	data, err := uc.streamgRPCAPI.GetAllStockSnapshot()
	if err != nil {
		return []*entity.Target{}, err
	}

	if len(data) < 200 {
		return []*entity.Target{}, errors.New("no all snapshots")
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
	})

	condition := uc.targetCond
	var result []*entity.Target
	for _, c := range condition.PriceVolumeLimit {
		for _, v := range data[:200] {
			stock := cc.GetStockDetail(v.GetCode())
			if stock == nil {
				continue
			}

			if !blackStockFilter(stock.Number, condition) ||
				!blackCatagoryFilter(stock.Category, condition) ||
				!targetFilter(v.GetClose(), v.GetTotalVolume(), c, false) {
				continue
			}

			result = append(result, &entity.Target{
				Rank:     len(result) + 1,
				StockNum: v.GetCode(),
				Volume:   v.GetTotalVolume(),
				TradeDay: tradeDay,
				Stock:    stock,
			})
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
		subArr = append(subArr, v.StockNum)
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
		subArr = append(subArr, v.StockNum)
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
func (uc *TargetUseCase) UnSubscribeStockTick(stockNum string) error {
	failUnSubNumArr, err := uc.gRPCAPI.UnSubscribeStockTick([]string{stockNum})
	if err != nil {
		return err
	}

	if len(failUnSubNumArr) != 0 {
		return fmt.Errorf("unsubscribe fail %v", failUnSubNumArr)
	}

	return nil
}

// UnSubscribeStockBidAsk -.
func (uc *TargetUseCase) UnSubscribeStockBidAsk(stockNum string) error {
	failUnSubNumArr, err := uc.gRPCAPI.UnSubscribeStockBidAsk([]string{stockNum})
	if err != nil {
		return err
	}

	if len(failUnSubNumArr) != 0 {
		return fmt.Errorf("unsubscribe fail %v", failUnSubNumArr)
	}

	return nil
}
