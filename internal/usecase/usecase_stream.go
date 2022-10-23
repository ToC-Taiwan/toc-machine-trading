package usecase

import (
	"context"
	"sort"
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/modules/event"
	"tmt/internal/usecase/modules/target"
	"tmt/internal/usecase/modules/tradeday"
	"tmt/internal/usecase/modules/trader"
)

// StreamUseCase -.
type StreamUseCase struct {
	repo    StreamRepo
	rabbit  StreamRabbit
	grpcapi StreamgRPCAPI

	basic entity.BasicInfo

	stockTradeSwitchCfg  config.StockTradeSwitch
	futureTradeSwitchCfg config.FutureTradeSwitch

	stockAnalyzeCfg  config.StockAnalyze
	futureAnalyzeCfg config.FutureAnalyze

	targetFilter *target.Filter
	tradeDay     *tradeday.TradeDay

	stockTradeInSwitch  bool
	futureTradeInSwitch bool
}

// NewStream -.
func NewStream(r StreamRepo, g StreamgRPCAPI, t StreamRabbit) *StreamUseCase {
	cfg := config.GetConfig()
	uc := &StreamUseCase{
		repo:    r,
		rabbit:  t,
		grpcapi: g,

		stockTradeSwitchCfg:  cfg.StockTradeSwitch,
		futureTradeSwitchCfg: cfg.FutureTradeSwitch,
		stockAnalyzeCfg:      cfg.StockAnalyze,
		futureAnalyzeCfg:     cfg.FutureAnalyze,

		basic:        *cc.GetBasicInfo(),
		targetFilter: target.NewFilter(cfg.TargetCond),
		tradeDay:     tradeday.NewTradeDay(),
	}
	t.FillAllBasic(uc.basic.AllStocks, uc.basic.AllFutures)

	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveOrderStatus(context.Background())

	go func() {
		time.Sleep(time.Until(cc.GetBasicInfo().TradeDay.Add(time.Hour * 9)))
		for range time.NewTicker(time.Second * 60).C {
			if uc.stockTradeInSwitch {
				if err := uc.realTimeAddTargets(); err != nil {
					log.Panic(err)
				}
			}
		}
	}()

	bus.SubscribeTopic(event.TopicStreamTargets, uc.ReceiveStreamData)
	bus.SubscribeTopic(event.TopicStreamFutureTargets, uc.ReceiveFutureStreamData)
	return uc
}

func (uc *StreamUseCase) realTimeAddTargets() error {
	data, err := uc.grpcapi.GetAllStockSnapshot()
	if err != nil {
		return err
	}

	// at least 200 snapshot to rank volume
	if len(data) < 200 {
		log.Warnf("stock snapshot len is not enough: %d", len(data))
		return nil
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
	})
	data = data[:uc.targetFilter.RealTimeRank]

	currentTargets := cc.GetTargets()
	targetsMap := make(map[string]*entity.Target)
	for _, t := range currentTargets {
		targetsMap[t.StockNum] = t
	}

	var newTargets []*entity.Target
	for i, d := range data {
		stock := cc.GetStockDetail(d.GetCode())
		if stock == nil {
			continue
		}

		if !uc.targetFilter.IsTarget(stock, d.GetClose()) {
			continue
		}

		if targetsMap[d.GetCode()] == nil {
			newTargets = append(newTargets, &entity.Target{
				Rank:     100 + i + 1,
				StockNum: d.GetCode(),
				Volume:   d.GetTotalVolume(),
				TradeDay: uc.basic.TradeDay,
				Stock:    stock,
			})
		}
	}

	if len(newTargets) != 0 {
		bus.PublishTopicEvent(event.TopicNewTargets, newTargets)
	}
	return nil
}

// ReceiveEvent -.
func (uc *StreamUseCase) ReceiveEvent(ctx context.Context) {
	eventChan := make(chan *entity.SinopacEvent)
	go func() {
		for {
			event := <-eventChan
			if err := uc.repo.InsertEvent(ctx, event); err != nil {
				log.Error(err)
			}

			if event.EventCode != 16 {
				log.Warnf("EventCode: %d, Event: %s, ResoCode: %d, Info: %s", event.EventCode, event.Event, event.Response, event.Info)
			}
		}
	}()
	uc.rabbit.EventConsumer(eventChan)
}

// ReceiveOrderStatus -.
func (uc *StreamUseCase) ReceiveOrderStatus(ctx context.Context) {
	orderStatusChan := make(chan interface{})
	go func() {
		for {
			order := <-orderStatusChan
			switch t := order.(type) {
			case *entity.StockOrder:
				if cc.GetOrderByOrderID(t.OrderID) != nil {
					bus.PublishTopicEvent(event.TopicInsertOrUpdateOrder, t)
				}
			case *entity.FutureOrder:
				if cc.GetFutureOrderByOrderID(t.OrderID) != nil {
					bus.PublishTopicEvent(event.TopicInsertOrUpdateFutureOrder, t)
				}
			}
		}
	}()
	uc.rabbit.OrderStatusConsumer(orderStatusChan)
}

// GetTSESnapshot -.
func (uc *StreamUseCase) GetTSESnapshot(ctx context.Context) (*entity.StockSnapShot, error) {
	body, err := uc.grpcapi.GetStockSnapshotTSE()
	if err != nil {
		return nil, err
	}
	return &entity.StockSnapShot{
		SnapTime:        time.Unix(0, body.GetTs()).Add(-8 * time.Hour),
		Open:            body.GetOpen(),
		High:            body.GetHigh(),
		Low:             body.GetLow(),
		Close:           body.GetClose(),
		TickType:        body.GetTickType(),
		PriceChg:        body.GetChangePrice(),
		PctChg:          body.GetChangeRate(),
		ChgType:         body.GetChangeType(),
		Volume:          body.GetVolume(),
		VolumeSum:       body.GetTotalVolume(),
		Amount:          body.GetAmount(),
		AmountSum:       body.GetTotalAmount(),
		YesterdayVolume: body.GetYesterdayVolume(),
		VolumeRatio:     body.GetVolumeRatio(),
	}, nil
}

// GetStockSnapshotByNumArr -.
func (uc *StreamUseCase) GetStockSnapshotByNumArr(stockNumArr []string) ([]*entity.StockSnapShot, error) {
	var fetchArr, stockNotExist []string
	for _, s := range stockNumArr {
		if cc.GetStockDetail(s) == nil {
			stockNotExist = append(stockNotExist, s)
		} else {
			fetchArr = append(fetchArr, s)
		}
	}

	snapshot, err := uc.grpcapi.GetStockSnapshotByNumArr(fetchArr)
	if err != nil {
		return nil, err
	}

	var result []*entity.StockSnapShot
	for _, body := range snapshot {
		stockNum := body.GetCode()
		result = append(result, &entity.StockSnapShot{
			StockNum:        stockNum,
			StockName:       cc.GetStockDetail(stockNum).Name,
			SnapTime:        time.Unix(0, body.GetTs()).Add(-8 * time.Hour),
			Open:            body.GetOpen(),
			High:            body.GetHigh(),
			Low:             body.GetLow(),
			Close:           body.GetClose(),
			TickType:        body.GetTickType(),
			PriceChg:        body.GetChangePrice(),
			PctChg:          body.GetChangeRate(),
			ChgType:         body.GetChangeType(),
			Volume:          body.GetVolume(),
			VolumeSum:       body.GetTotalVolume(),
			Amount:          body.GetAmount(),
			AmountSum:       body.GetTotalAmount(),
			YesterdayVolume: body.GetYesterdayVolume(),
			VolumeRatio:     body.GetVolumeRatio(),
		})
	}

	for _, v := range stockNotExist {
		result = append(result, &entity.StockSnapShot{
			StockNum: v,
		})
	}

	return result, nil
}

// NewFutureRealTimeConnection -.
func (uc *StreamUseCase) NewFutureRealTimeConnection(timestamp int64, tickChan chan *entity.RealTimeFutureTick) {
	uc.rabbit.AddFutureTickChan(timestamp, tickChan)
}

// DeleteFutureRealTimeConnection -.
func (uc *StreamUseCase) DeleteFutureRealTimeConnection(timestamp int64) {
	uc.rabbit.RemoveFutureTickChan(timestamp)
}

// ReceiveStreamData - receive target data, start goroutine to trade
func (uc *StreamUseCase) ReceiveStreamData(ctx context.Context, targetArr []*entity.Target) {
	agentChan := make(chan *trader.StockTrader)
	targetMap := make(map[string]*entity.Target)
	mutex := sync.RWMutex{}

	go func() {
		for {
			agent, ok := <-agentChan
			if !ok {
				break
			}
			go agent.TradingRoom()

			// send tick, bidask to trade room's channel
			go uc.rabbit.TickConsumer(agent.GetStockNum(), agent.GetTickChan())
			go uc.rabbit.StockBidAskConsumer(agent.GetStockNum(), agent.GetBidAskChan())

			mutex.RLock()
			target := targetMap[agent.GetStockNum()]
			mutex.RUnlock()

			bus.PublishTopicEvent(event.TopicSubscribeTickTargets, []*entity.Target{target})
		}
	}()

	var wg sync.WaitGroup
	for _, t := range targetArr {
		target := t
		mutex.Lock()
		targetMap[target.StockNum] = target
		mutex.Unlock()

		wg.Add(1)
		go func() {
			defer wg.Done()
			agent := trader.NewStockTrader(target.StockNum, uc.stockTradeSwitchCfg, uc.stockAnalyzeCfg)
			agentChan <- agent
		}()
	}

	wg.Wait()
	close(agentChan)

	go uc.checkStockTradeSwitch()
}

// ReceiveFutureStreamData -.
func (uc *StreamUseCase) ReceiveFutureStreamData(ctx context.Context, code string) {
	agent := trader.NewFutureTrader(code, uc.futureTradeSwitchCfg, uc.futureAnalyzeCfg)

	go agent.TradingRoom()

	go uc.rabbit.FutureTickConsumer(code, agent.GetTickChan())
	go uc.rabbit.FutureBidAskConsumer(code, agent.GetBidAskChan())

	bus.PublishTopicEvent(event.TopicSubscribeFutureTickTargets, code)

	go uc.checkFutureTradeSwitch()
}

func (uc *StreamUseCase) checkStockTradeSwitch() {
	if !uc.stockTradeSwitchCfg.AllowTrade {
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
			bus.PublishTopicEvent(event.TopicUpdateStockTradeSwitch, uc.stockTradeInSwitch)
		}
	}
}

func (uc *StreamUseCase) checkFutureTradeSwitch() {
	if !uc.futureTradeSwitchCfg.AllowTrade {
		return
	}

	futureTradeDay := uc.tradeDay.GetFutureTradeDay()
	timeRange := [][]time.Time{}
	firstStart := futureTradeDay.StartTime
	secondStart := futureTradeDay.EndTime.Add(-300 * time.Minute)

	timeRange = append(timeRange, []time.Time{firstStart, firstStart.Add(time.Duration(uc.futureTradeSwitchCfg.TradeTimeRange.FirstPartDuration) * time.Minute)})
	timeRange = append(timeRange, []time.Time{secondStart, secondStart.Add(time.Duration(uc.futureTradeSwitchCfg.TradeTimeRange.SecondPartDuration) * time.Minute)})

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
			bus.PublishTopicEvent(event.TopicUpdateFutureTradeSwitch, uc.futureTradeInSwitch)
		}
	}
}
