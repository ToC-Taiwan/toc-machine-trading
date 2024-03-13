package usecase

import (
	"context"
	"errors"
	"sort"
	"time"

	"tmt/internal/config"
	"tmt/pb"

	"tmt/internal/entity"
	"tmt/internal/usecase/cache"
	"tmt/internal/usecase/grpc"
	"tmt/internal/usecase/modules/calendar"
	"tmt/internal/usecase/repo"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"
)

// TargetUseCase -.
type TargetUseCase struct {
	repo    repo.TargetRepo
	gRPCAPI grpc.RealTimegRPCAPI

	cfg      *config.Config
	tradeDay *calendar.Calendar

	logger *log.Log
	cc     *cache.Cache
	bus    *eventbus.Bus

	rankFromSnapshot *pb.StockVolumeRankResponse
}

func NewTarget() Target {
	cfg := config.Get()
	uc := &TargetUseCase{
		repo:     repo.NewTarget(cfg.GetPostgresPool()),
		gRPCAPI:  grpc.NewRealTime(cfg.GetSinopacPool()),
		cfg:      cfg,
		tradeDay: calendar.Get(),
		logger:   log.Get(),
		cc:       cache.Get(),
		bus:      eventbus.Get(),
	}

	// query targets from db
	tDay := uc.tradeDay.GetStockTradeDay().TradeDay
	targetArr, err := uc.searchTradeDayTargetsFromDB(tDay)
	if err != nil {
		uc.logger.Fatal(err)
	}

	// db has no targets, find targets from gRPC
	if len(targetArr) == 0 {
		targetArr, err = uc.searchTradeDayTargets(tDay)
		if err != nil {
			uc.logger.Fatal(err)
		}

		if len(targetArr) == 0 {
			stuck := make(chan struct{})
			uc.logger.Error("no targets")
			<-stuck
		}
	}

	uc.cc.AppendStockTargets(targetArr)
	uc.publishNewStockTargets(targetArr)

	return uc
}

// GetTargets - get targets from cache
func (uc *TargetUseCase) GetTargets(ctx context.Context) []*entity.StockTarget {
	return uc.cc.GetStockTargets()
}

func (uc *TargetUseCase) publishNewStockTargets(targetArr []*entity.StockTarget) {
	if err := uc.repo.InsertOrUpdateTargetArr(context.Background(), targetArr); err != nil {
		uc.logger.Fatal(err)
	}
	uc.bus.PublishTopicEvent(topicFetchStockHistory, targetArr)
}

func (uc *TargetUseCase) searchTradeDayTargetsFromDB(tradeDay time.Time) ([]*entity.StockTarget, error) {
	targetArr, err := uc.repo.QueryTargetsByTradeDay(context.Background(), tradeDay)
	if err != nil {
		return nil, err
	}

	var result []*entity.StockTarget
	for _, v := range targetArr {
		stock := uc.cc.GetStockDetail(v.StockNum)
		if stock == nil {
			continue
		}

		v.Stock = stock
		result = append(result, v)
	}

	return result, nil
}

func (uc *TargetUseCase) searchTradeDayTargets(tradeDay time.Time) ([]*entity.StockTarget, error) {
	lastTradeDay := uc.tradeDay.GetLastNStockTradeDay(1)[0]
	t, err := uc.gRPCAPI.GetStockVolumeRank(lastTradeDay.Format(entity.ShortTimeLayout))
	if err != nil {
		return nil, err
	}

	if len(t) == 0 && time.Now().Before(uc.tradeDay.GetStockTradeDay().TradeDay.Add(8*time.Hour)) {
		uc.logger.Warn("VolumeRank is empty, search from all snapshot")
		return uc.searchTradeDayTargetsFromAllSnapshot(tradeDay)
	}

	var result []*entity.StockTarget
	for _, v := range t {
		if len(result) >= 25 {
			break
		}

		stock := uc.cc.GetStockDetail(v.GetCode())
		if stock == nil {
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

func (uc *TargetUseCase) GetCurrentVolumeRank() (*pb.StockVolumeRankResponse, error) {
	tradeStart := uc.tradeDay.GetStockTradeDay().StartTime
	tradeEnd := uc.tradeDay.GetStockTradeDay().EndTime
	if time.Now().Before(tradeStart) || time.Now().After(tradeEnd) {
		return uc.getVolumeRankFromSnapshot()
	}

	t, err := uc.gRPCAPI.GetStockVolumeRankPB(time.Now().Format(entity.ShortTimeLayout))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (uc *TargetUseCase) searchTradeDayTargetsFromAllSnapshot(tradeDay time.Time) ([]*entity.StockTarget, error) {
	data, err := uc.gRPCAPI.GetAllStockSnapshot()
	if err != nil {
		return []*entity.StockTarget{}, err
	}

	if len(data) < 200 {
		return []*entity.StockTarget{}, errors.New("no all snapshots")
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
	})

	var result []*entity.StockTarget
	for _, v := range data[:200] {
		stock := uc.cc.GetStockDetail(v.GetCode())
		if stock == nil {
			continue
		}

		if len(result) >= 25 {
			break
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

func (uc *TargetUseCase) getVolumeRankFromSnapshot() (*pb.StockVolumeRankResponse, error) {
	if uc.rankFromSnapshot != nil {
		return uc.rankFromSnapshot, nil
	}

	snapshots, err := uc.gRPCAPI.GetAllStockSnapshot()
	if err != nil {
		return nil, err
	}

	sort.SliceStable(snapshots, func(i, j int) bool {
		return snapshots[i].GetTotalVolume() > snapshots[j].GetTotalVolume()
	})

	result := &pb.StockVolumeRankResponse{}
	for _, v := range snapshots {
		if len(result.GetData()) >= 200 {
			break
		}
		stock := uc.cc.GetStockDetail(v.GetCode())
		if stock == nil {
			continue
		}

		tickType := 1
		if v.GetTickType() == "Buy" {
			tickType = 2
		}
		date := time.Unix(0, v.GetTs()).Format(entity.ShortTimeLayout)
		result.Data = append(result.Data, &pb.StockVolumeRankMessage{
			Date:            date,
			Code:            v.GetCode(),
			Name:            stock.Name,
			Ts:              v.GetTs(),
			Open:            v.GetOpen(),
			High:            v.GetHigh(),
			Low:             v.GetLow(),
			Close:           v.GetClose(),
			TickType:        int64(tickType),
			ChangePrice:     v.GetChangePrice(),
			AveragePrice:    v.GetAveragePrice(),
			Volume:          v.GetVolume(),
			TotalVolume:     v.GetTotalVolume(),
			Amount:          v.GetAmount(),
			TotalAmount:     v.GetTotalAmount(),
			YesterdayVolume: int64(v.GetYesterdayVolume()),
			VolumeRatio:     v.GetVolumeRatio(),
			BuyPrice:        v.GetBuyPrice(),
			BuyVolume:       int64(v.GetBuyVolume()),
			SellPrice:       v.GetSellPrice(),
			SellVolume:      v.GetSellVolume(),
		})
	}
	uc.rankFromSnapshot = result
	return result, nil
}

// func (uc *TargetUseCase) realTimeAddTargets() error {
// 	data, err := uc.gRPCAPI.GetAllStockSnapshot()
// 	if err != nil {
// 		return err
// 	}

// 	// at least 200 snapshot to rank volume
// 	if len(data) < 200 {
// 		uc.logger.Warnf("stock snapshot len is not enough: %d", len(data))
// 		return nil
// 	}

// 	sort.SliceStable(data, func(i, j int) bool {
// 		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
// 	})
// 	data = data[:uc.targetFilter.RealTimeRank]

// 	currentTargets := uc.cc.GetStockTargets()
// 	targetsMap := make(map[string]*entity.StockTarget)
// 	for _, t := range currentTargets {
// 		targetsMap[t.StockNum] = t
// 	}

// 	var newTargets []*entity.StockTarget
// 	for i, d := range data {
// 		stock := uc.cc.GetStockDetail(d.GetCode())
// 		if stock == nil {
// 			continue
// 		}

// 		if !uc.targetFilter.IsTarget(stock, d.GetClose()) {
// 			continue
// 		}

// 		if targetsMap[d.GetCode()] == nil {
// 			newTargets = append(newTargets, &entity.StockTarget{
// 				Rank:     100 + i + 1,
// 				StockNum: d.GetCode(),
// 				Volume:   d.GetTotalVolume(),
// 				TradeDay: uc.basic.TradeDay,
// 				Stock:    stock,
// 			})
// 		}
// 	}

// 	if len(newTargets) != 0 {
// 		uc.cc.AppendStockTargets(newTargets)
// 		uc.publishNewStockTargets(newTargets)
// 		for _, t := range newTargets {
// 			uc.logger.Infof("New target: %s", t.StockNum)
// 		}
// 	}
// 	return nil
// }
