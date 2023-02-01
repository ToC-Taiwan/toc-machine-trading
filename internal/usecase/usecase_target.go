package usecase

import (
	"context"
	"errors"
	"sort"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/module/target"
	"tmt/internal/usecase/module/tradeday"
	"tmt/internal/usecase/topic"
	"tmt/pkg/common"
)

// TargetUseCase -.
type TargetUseCase struct {
	repo    TargetRepo
	gRPCAPI RealTimegRPCAPI

	targetFilter *target.Filter
	cfg          *config.Config
	basic        *entity.BasicInfo
	tradeDay     *tradeday.TradeDay

	stockTradeInSwitch  bool
	futureTradeInSwitch bool
}

func (uc *TargetUseCase) checkStockTradeSwitch() {
	if !uc.cfg.StockTradeSwitch.AllowTrade {
		return
	}

	openTime := uc.basic.OpenTime
	tradeInEndTime := uc.basic.TradeInEndTime

	for range time.NewTicker(2500 * time.Millisecond).C {
		now := time.Now()
		var tempSwitch bool
		switch {
		case now.Before(openTime) || now.After(tradeInEndTime):
			tempSwitch = false
		case now.After(openTime) && now.Before(tradeInEndTime):
			tempSwitch = true
		}

		if uc.stockTradeInSwitch != tempSwitch {
			uc.stockTradeInSwitch = tempSwitch
			bus.PublishTopicEvent(topic.TopicUpdateStockTradeSwitch, uc.stockTradeInSwitch)
		}
	}
}

func (uc *TargetUseCase) checkFutureTradeSwitch() {
	if !uc.cfg.FutureTradeSwitch.AllowTrade {
		return
	}

	futureTradeDay := uc.tradeDay.GetFutureTradeDay()
	timeRange := [][]time.Time{}
	firstStart := futureTradeDay.StartTime
	secondStart := futureTradeDay.EndTime.Add(-300 * time.Minute)

	timeRange = append(timeRange, []time.Time{firstStart, firstStart.Add(time.Duration(uc.cfg.FutureTradeSwitch.TradeTimeRange.FirstPartDuration) * time.Minute)})
	timeRange = append(timeRange, []time.Time{secondStart, secondStart.Add(time.Duration(uc.cfg.FutureTradeSwitch.TradeTimeRange.SecondPartDuration) * time.Minute)})

	for range time.NewTicker(2500 * time.Millisecond).C {
		now := time.Now()
		var tempSwitch bool
		for _, rangeTime := range timeRange {
			if now.After(rangeTime[0]) && now.Before(rangeTime[1]) {
				tempSwitch = true
			}
		}

		if uc.futureTradeInSwitch != tempSwitch {
			uc.futureTradeInSwitch = tempSwitch
			bus.PublishTopicEvent(topic.TopicUpdateFutureTradeSwitch, uc.futureTradeInSwitch)
		}
	}
}

// GetTargets - get targets from cache
func (uc *TargetUseCase) GetTargets(ctx context.Context) []*entity.StockTarget {
	return cc.GetStockTargets()
}

func (uc *TargetUseCase) publishNewStockTargets(targetArr []*entity.StockTarget) {
	if err := uc.repo.InsertOrUpdateTargetArr(context.Background(), targetArr); err != nil {
		logger.Fatal(err)
	}
	bus.PublishTopicEvent(topic.TopicFetchStockHistory, targetArr)
}

func (uc *TargetUseCase) publishNewFutureTargets() {
	if futureTarget, err := uc.getFutureTarget(); err != nil {
		logger.Fatal(err)
	} else {
		bus.PublishTopicEvent(topic.TopicSubscribeFutureTickTargets, futureTarget)
	}
}

func (uc *TargetUseCase) getFutureTarget() (string, error) {
	futures, err := uc.repo.QueryAllMXFFuture(context.Background())
	if err != nil {
		return "", err
	}

	for _, v := range futures {
		if v.Code == "MXFR1" || v.Code == "MXFR2" {
			continue
		}

		if time.Now().Before(v.DeliveryDate) {
			return v.Code, nil
		}
	}

	return "", errors.New("no future")
}

func (uc *TargetUseCase) searchTradeDayTargets(tradeDay time.Time) ([]*entity.StockTarget, error) {
	lastTradeDay := cc.GetBasicInfo().LastTradeDay
	t, err := uc.gRPCAPI.GetStockVolumeRank(lastTradeDay.Format(common.ShortTimeLayout))
	if err != nil {
		return nil, err
	}

	if len(t) == 0 && time.Now().Before(cc.GetBasicInfo().TradeDay.Add(8*time.Hour)) {
		logger.Warn("VolumeRank is empty, search from all snapshot")
		return uc.searchTradeDayTargetsFromAllSnapshot(tradeDay)
	}

	var result []*entity.StockTarget
	for _, v := range t {
		stock := cc.GetStockDetail(v.GetCode())
		if stock == nil {
			continue
		}

		if !uc.targetFilter.CheckVolume(v.GetTotalVolume()) || !uc.targetFilter.IsTarget(stock, v.GetClose()) {
			continue
		}

		result = append(result, &entity.StockTarget{
			Rank:     len(result) + 1,
			StockNum: v.GetCode(),
			Volume:   v.GetTotalVolume(),
			TradeDay: tradeDay,
			Stock:    stock,
		})
	}
	return result, nil
}

func (uc *TargetUseCase) searchTradeDayTargetsFromAllSnapshot(tradeDay time.Time) ([]*entity.StockTarget, error) {
	data, err := uc.gRPCAPI.GetAllStockSnapshot()
	if err != nil {
		return []*entity.StockTarget{}, err
	}

	if len(data) < 200 {
		return []*entity.StockTarget{}, errors.New("no all snapshots")
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
	})

	var result []*entity.StockTarget
	for _, v := range data[:200] {
		stock := cc.GetStockDetail(v.GetCode())
		if stock == nil {
			continue
		}

		if !uc.targetFilter.CheckVolume(v.GetTotalVolume()) || !uc.targetFilter.IsTarget(stock, v.GetClose()) {
			continue
		}

		result = append(result, &entity.StockTarget{
			Rank:     len(result) + 1,
			StockNum: v.GetCode(),
			Volume:   v.GetTotalVolume(),
			TradeDay: tradeDay,
			Stock:    stock,
		})
	}
	return result, nil
}

func (uc *TargetUseCase) realTimeAddTargets() error {
	data, err := uc.gRPCAPI.GetAllStockSnapshot()
	if err != nil {
		return err
	}

	// at least 200 snapshot to rank volume
	if len(data) < 200 {
		logger.Warnf("stock snapshot len is not enough: %d", len(data))
		return nil
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
	})
	data = data[:uc.targetFilter.RealTimeRank]

	currentTargets := cc.GetStockTargets()
	targetsMap := make(map[string]*entity.StockTarget)
	for _, t := range currentTargets {
		targetsMap[t.StockNum] = t
	}

	var newTargets []*entity.StockTarget
	for i, d := range data {
		stock := cc.GetStockDetail(d.GetCode())
		if stock == nil {
			continue
		}

		if !uc.targetFilter.IsTarget(stock, d.GetClose()) {
			continue
		}

		if targetsMap[d.GetCode()] == nil {
			newTargets = append(newTargets, &entity.StockTarget{
				Rank:     100 + i + 1,
				StockNum: d.GetCode(),
				Volume:   d.GetTotalVolume(),
				TradeDay: uc.basic.TradeDay,
				Stock:    stock,
			})
		}
	}

	if len(newTargets) != 0 {
		cc.AppendStockTargets(newTargets)
		uc.publishNewStockTargets(newTargets)
		for _, t := range newTargets {
			logger.Infof("New target: %s", t.StockNum)
		}
	}
	return nil
}
