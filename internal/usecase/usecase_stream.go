package usecase

import (
	"context"
	"sort"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/rabbit"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/config"
)

// StreamUseCase -.
type StreamUseCase struct {
	repo    StreamRepo
	rabbit  StreamRabbit
	grpcapi StreamgRPCAPI

	tradeSwitchCfg config.TradeSwitch
	analyzeCfg     config.Analyze
	basic          entity.BasicInfo
	targetCond     []config.TargetCond

	tradeInSwitch bool
	clearAll      bool
}

// NewStream -.
func NewStream(r *repo.StreamRepo, g *grpcapi.StreamgRPCAPI, t *rabbit.StreamRabbit) *StreamUseCase {
	uc := &StreamUseCase{
		repo:    r,
		rabbit:  t,
		grpcapi: g,
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	uc.tradeSwitchCfg = cfg.TradeSwitch
	uc.analyzeCfg = cfg.Analyze
	uc.targetCond = cfg.TargetCond
	uc.basic = *cc.GetBasicInfo()

	go uc.checkTradeSwitch()
	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveOrderStatus(context.Background())

	go func() {
		for range time.NewTicker(60 * time.Second).C {
			if uc.tradeInSwitch {
				if err := uc.realTimeAddTargets(context.Background()); err != nil {
					log.Panic(err)
				}
			}
		}
	}()

	bus.SubscribeTopic(topicStreamTargets, uc.ReceiveStreamData)
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
	orderStatusChan := make(chan *entity.Order)
	go func() {
		for {
			order := <-orderStatusChan
			bus.PublishTopicEvent(topicInsertOrUpdateOrder, order)
		}
	}()
	uc.rabbit.OrderStatusConsumer(orderStatusChan)
}

// ReceiveStreamData -.
func (uc *StreamUseCase) ReceiveStreamData(ctx context.Context, targetArr []*entity.Target) {
	for _, t := range targetArr {
		data := &RealTimeData{
			stockNum:      t.StockNum,
			orderMap:      make(map[entity.OrderAction][]*entity.Order),
			orderQuantity: 1,
			tickChan:      make(chan *entity.RealTimeTick),
			bidAskChan:    make(chan *entity.RealTimeBidAsk),
		}
		data.setHistoryTickAnalyze(cc.GetHistoryTickAnalyze(t.StockNum))
		if biasRate := cc.GetBiasRate(t.StockNum); biasRate > 4 || biasRate < -4 {
			data.orderQuantity = 2
		}
		go data.checkFirstTickArrive()

		go uc.tradeAgent(data)
		go uc.rabbit.TickConsumer(t.StockNum, data.tickChan)
		go uc.rabbit.BidAskConsumer(t.StockNum, data.bidAskChan)
	}
	bus.PublishTopicEvent(topicSubscribeTickTargets, targetArr)
}

func (uc *StreamUseCase) tradeAgent(data *RealTimeData) {
	go func() {
		for {
			tick := <-data.tickChan
			if tick.PctChg < uc.analyzeCfg.CloseChangeRatioLow || tick.PctChg > uc.analyzeCfg.CloseChangeRatioHigh {
				continue
			}
			data.tickArr = append(data.tickArr, tick)
			order := data.generateOrder(uc.analyzeCfg, uc.clearAll)
			if order == nil {
				continue
			}
			uc.placeOrder(data, order)
		}
	}()

	go func() {
		for {
			data.bidAsk = <-data.bidAskChan
		}
	}()

	for {
		time.Sleep(15 * time.Second)
		if uc.clearAll {
			order := data.clearUnfinishedOrder()
			if order != nil {
				uc.placeOrder(data, order)
			}
		}
	}
}

func (uc *StreamUseCase) placeOrder(data *RealTimeData, order *entity.Order) {
	var timeout time.Duration
	switch order.Action {
	case entity.ActionBuy, entity.ActionSellFirst:
		if !uc.tradeInSwitch {
			return
		}
		timeout = time.Duration(uc.tradeSwitchCfg.TradeInWaitTime) * time.Second
	case entity.ActionSell, entity.ActionBuyLater:
		timeout = time.Duration(uc.tradeSwitchCfg.TradeOutWaitTime) * time.Second
	}

	bus.PublishTopicEvent(topicPlaceOrder, order)
	data.waitingOrder = order

	log.Warnf("Place Order -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)
	go data.checkPlaceOrderStatus(order, timeout)
}

func (uc *StreamUseCase) checkTradeSwitch() {
	openTime := uc.basic.OpenTime
	endTime := uc.basic.EndTime
	tradeInEndTime := uc.basic.TradeInEndTime
	tradeOutEndTime := uc.basic.TradeOutEndTime

	for range time.NewTicker(5 * time.Second).C {
		now := time.Now()
		switch {
		case now.Before(openTime) || now.After(endTime) || (now.After(tradeInEndTime) && now.Before(tradeOutEndTime)):
			uc.tradeInSwitch = false
			uc.clearAll = false
		case now.After(openTime) && now.Before(tradeInEndTime):
			uc.tradeInSwitch = true
			uc.clearAll = false
		case now.After(tradeOutEndTime) && now.Before(endTime):
			uc.tradeInSwitch = false
			uc.clearAll = true
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

func (uc *StreamUseCase) realTimeAddTargets(ctx context.Context) error {
	if !uc.tradeInSwitch {
		return nil
	}

	data, err := uc.grpcapi.GetAllStockSnapshot()
	if err != nil {
		return err
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
	})
	data = data[:10]

	currentTargets := cc.GetTargets()
	targetsMap := make(map[string]*entity.Target)
	for _, t := range currentTargets {
		targetsMap[t.StockNum] = t
	}

	var newTargets []*entity.Target
	for _, c := range uc.targetCond {
		for i, d := range data {
			if targetFilter(d.GetClose(), d.GetTotalVolume(), c, true) {
				if stock := cc.GetStockDetail(d.GetCode()); stock != nil && targetsMap[d.GetCode()] == nil {
					newTargets = append(newTargets, &entity.Target{
						Rank:        100 + i + 1,
						StockNum:    d.GetCode(),
						Volume:      d.GetTotalVolume(),
						Subscribe:   c.Subscribe,
						RealTimeAdd: true,
						TradeDay:    uc.basic.TradeDay,
						Stock:       stock,
					})
				}
			}
		}
	}

	if len(newTargets) != 0 {
		cc.AppendTargets(newTargets)
		bus.PublishTopicEvent(topicRealTimeTargets, ctx, newTargets)
		bus.PublishTopicEvent(topicTargets, ctx, newTargets)
	}
	return nil
}

// GetStockSnapshotByNumArr -.
func (uc *StreamUseCase) GetStockSnapshotByNumArr(stockNumArr []string) ([]*entity.StockSnapShot, error) {
	snapshot, err := uc.grpcapi.GetStockSnapshotByNumArr(stockNumArr)
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
	return result, nil
}
