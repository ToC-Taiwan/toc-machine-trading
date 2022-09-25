package usecase

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/global"
	"tmt/internal/entity"
)

// StreamUseCase -.
type StreamUseCase struct {
	repo    StreamRepo
	rabbit  StreamRabbit
	grpcapi StreamgRPCAPI

	tradeSwitchCfg       config.TradeSwitch
	futureTradeSwitchCfg config.FutureTradeSwitch
	basic                entity.BasicInfo

	stockAnalyzeCfg  config.StockAnalyze
	futureAnalyzeCfg config.FutureAnalyze

	targetFilter *TargetFilter

	tradeInSwitch       bool
	futureTradeInSwitch bool

	// allowForward bool
	// allowReverse bool
}

// NewStream -.
func NewStream(r StreamRepo, g StreamgRPCAPI, t StreamRabbit) *StreamUseCase {
	cfg := config.GetConfig()
	basic := *cc.GetBasicInfo()
	t.FillAllBasic(basic.AllStocks, basic.AllFutures)

	uc := &StreamUseCase{
		repo:                 r,
		rabbit:               t,
		grpcapi:              g,
		tradeSwitchCfg:       cfg.TradeSwitch,
		futureTradeSwitchCfg: cfg.FutureTradeSwitch,
		stockAnalyzeCfg:      cfg.StockAnalyze,
		futureAnalyzeCfg:     cfg.FutureAnalyze,
		basic:                basic,
		targetFilter:         NewTargetFilter(cfg.TargetCond),
	}

	go uc.checkTradeSwitch()
	go uc.checkFutureTradeSwitch()
	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveOrderStatus(context.Background())

	go func() {
		time.Sleep(time.Until(cc.GetBasicInfo().TradeDay.Add(time.Hour * 9)))
		for range time.NewTicker(time.Second * 60).C {
			if uc.tradeInSwitch {
				if err := uc.realTimeAddTargets(); err != nil {
					log.Panic(err)
				}
			}
		}
	}()

	bus.SubscribeTopic(topicStreamTargets, uc.ReceiveStreamData)
	bus.SubscribeTopic(topicStreamFutureTargets, uc.ReceiveFutureStreamData)
	return uc
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
					bus.PublishTopicEvent(topicInsertOrUpdateOrder, t)
				}
			case *entity.FutureOrder:
				if cc.GetFutureOrderByOrderID(t.OrderID) != nil {
					bus.PublishTopicEvent(topicInsertOrUpdateFutureOrder, t)
				}
			}
		}
	}()
	uc.rabbit.OrderStatusConsumer(orderStatusChan)
}

// ReceiveStreamData - receive target data, start goroutine to trade
func (uc *StreamUseCase) ReceiveStreamData(ctx context.Context, targetArr []*entity.Target) {
	agentChan := make(chan *TradeAgent)
	targetMap := make(map[string]*entity.Target)
	mutex := sync.RWMutex{}

	go func() {
		for {
			agent, ok := <-agentChan
			if !ok {
				break
			}
			go uc.tradingRoom(agent)

			// send tick, bidask to trade room's channel
			go uc.rabbit.TickConsumer(agent.stockNum, agent.tickChan)
			go uc.rabbit.BidAskConsumer(agent.stockNum, agent.bidAskChan)

			mutex.RLock()
			target := targetMap[agent.stockNum]
			mutex.RUnlock()

			bus.PublishTopicEvent(topicSubscribeTickTargets, []*entity.Target{target})
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
			agent := NewAgent(target.StockNum, uc.tradeSwitchCfg)
			agentChan <- agent
		}()
	}
	wg.Wait()
	close(agentChan)
}

func (uc *StreamUseCase) tradingRoom(agent *TradeAgent) {
	go func() {
		for {
			agent.lastTick = <-agent.tickChan
			agent.tickArr = append(agent.tickArr, agent.lastTick)
			log.Debugf("%s tick time delay: %s", agent.stockNum, time.Since(agent.lastTick.TickTime).String())

			if agent.waitingOrder != nil || agent.analyzeTickTime.IsZero() || !agent.openPass {
				continue
			}

			order := agent.generateOrder(uc.stockAnalyzeCfg)
			if order == nil {
				continue
			}

			uc.placeOrder(agent, order)
		}
	}()

	go func() {
		for {
			agent.lastBidAsk = <-agent.bidAskChan
			log.Debugf("%s bidask time delay: %s", agent.stockNum, time.Since(agent.lastBidAsk.BidAskTime).String())
		}
	}()
}

func (uc *StreamUseCase) placeOrder(agent *TradeAgent, order *entity.StockOrder) {
	// switch order.Action {
	// case entity.ActionBuy:
	// 	if !uc.allowForward {
	// 		return
	// 	}
	// case entity.ActionSellFirst:
	// 	if !uc.allowReverse {
	// 		return
	// 	}
	// }

	if order.Price == 0 {
		log.Errorf("%s Order price is 0", order.StockNum)
		return
	}

	agent.waitingOrder = order

	// if out of trade in time, return
	if !uc.tradeInSwitch && (order.Action == entity.ActionBuy || order.Action == entity.ActionSellFirst) {
		// avoid stuck in the market
		agent.waitingOrder = nil
		return
	}

	bus.PublishTopicEvent(topicPlaceOrder, order)
	go agent.checkPlaceOrderStatus(order)
}

func (uc *StreamUseCase) checkTradeSwitch() {
	if !uc.tradeSwitchCfg.AllowTrade {
		return
	}

	openTime := uc.basic.OpenTime
	tradeInEndTime := uc.basic.TradeInEndTime

	for range time.NewTicker(2500 * time.Millisecond).C {
		now := time.Now()
		switch {
		case now.Before(openTime) || now.After(tradeInEndTime):
			uc.tradeInSwitch = false
		case now.After(openTime) && now.Before(tradeInEndTime):
			uc.tradeInSwitch = true
		}
	}
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

func (uc *StreamUseCase) realTimeAddTargets() error {
	data, err := uc.grpcapi.GetAllStockSnapshot()
	if err != nil {
		return err
	}

	// at least 200 snapshot to rank volume
	if len(data) < 200 {
		return nil
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
	})
	data = data[:uc.targetFilter.realTimeRank]

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

		if !uc.targetFilter.isTarget(stock, d.GetClose()) {
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
		bus.PublishTopicEvent(topicNewTargets, newTargets)
	}
	return nil
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

// ReceiveFutureStreamData -.
func (uc *StreamUseCase) ReceiveFutureStreamData(ctx context.Context, code string) {
	agent := NewFutureAgent(code, uc.futureTradeSwitchCfg, uc.futureAnalyzeCfg)

	go uc.futureTradingRoom(agent)
	// go uc.checkFirstFutureTick(agent)
	go uc.rabbit.FutureTickConsumer(code, agent.tickChan)

	bus.PublishTopicEvent(topicSubscribeFutureTickTargets, code)
}

func (uc *StreamUseCase) futureTradingRoom(agent *FutureTradeAgent) {
	for {
		agent.lastTick = <-agent.tickChan
		agent.tickArr = append(agent.tickArr, agent.lastTick)

		log.Debugf("TickTime: %s, Code: %s, Close: %.0f, TickType: %d, Volume: %3d, PriceChg: %.0f", agent.lastTick.TickTime.Format(global.LongTimeLayout), agent.lastTick.Code, agent.lastTick.Close, agent.lastTick.TickType, agent.lastTick.Volume, agent.lastTick.PriceChg)
		if agent.waitingOrder != nil || agent.analyzeTickTime.IsZero() {
			continue
		}

		if order := agent.generateOrder(); order == nil {
			continue
		} else {
			uc.placeFutureOrder(agent, order)
		}
	}
}

// func (uc *StreamUseCase) checkFirstFutureTick(agent *FutureTradeAgent) {
// 	for {
// 		time.Sleep(time.Second)
// 		dayMarketLastTick := cc.GetFutureHistoryTick(agent.code)
// 		if agent.lastTick == nil || dayMarketLastTick == nil {
// 			continue
// 		}
// 		agent.analyzeTickTime = agent.lastTick.TickTime

// 		if agent.lastTick.TickTime.Hour() != 8 {
// 			log.Warn("Not at stock trading time")
// 			log.Warnf("DayMarketLastTickTime: %s, Close: %.0f", dayMarketLastTick.TickTime.Format(global.LongTimeLayout), dayMarketLastTick.Close)
// 			log.Warnf("CurrentTickTime %s, Close: %.0f", agent.lastTick.TickTime.Format(global.LongTimeLayout), agent.lastTick.Close)
// 			break
// 		}

// 		if gap := agent.lastTick.Close - dayMarketLastTick.Close; gap >= 0 {
// 			uc.allowForward = true
// 		} else if gap != 0 {
// 			uc.allowReverse = true
// 		}

// 		break
// 	}
// }

func (uc *StreamUseCase) placeFutureOrder(agent *FutureTradeAgent, order *entity.FutureOrder) {
	if order.Price == 0 {
		log.Errorf("%s Future Order price is 0", order.Code)
		return
	}

	agent.waitingOrder = order

	// if out of trade in time, return
	if !uc.futureTradeInSwitch && (order.Action == entity.ActionBuy || order.Action == entity.ActionSellFirst) {
		// avoid stuck in the market
		agent.waitingOrder = nil
		return
	}

	bus.PublishTopicEvent(topicPlaceFutureOrder, order)
	go agent.checkPlaceOrderStatus(order)
}

func (uc *StreamUseCase) checkFutureTradeSwitch() {
	if !uc.futureTradeSwitchCfg.AllowTrade {
		return
	}

	for range time.NewTicker(2500 * time.Millisecond).C {
		now := time.Now()
		var tempSwitch bool
		for _, v := range uc.futureTradeSwitchCfg.TradeTimeRange {
			start, err := time.ParseInLocation(global.LongTimeLayout, fmt.Sprintf("%s %s", now.Format(global.ShortTimeLayout), v.StartTime), time.Local)
			if err != nil {
				log.Panic(err)
			}
			end := start.Add(time.Duration(v.Duration) * time.Minute)
			if now.After(start) && now.Before(end) {
				tempSwitch = true
			}
		}
		uc.futureTradeInSwitch = tempSwitch
	}
}
