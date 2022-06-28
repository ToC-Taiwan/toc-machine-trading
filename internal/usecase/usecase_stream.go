package usecase

import (
	"context"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/rabbit"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/config"
)

// StreamUseCase -.
type StreamUseCase struct {
	repo   StreamRepo
	rabbit StreamRabbit

	tradeSwitchCfg config.TradeSwitch
	analyzeCfg     config.Analyze
	basic          entity.BasicInfo

	tradeInSwitch bool
}

// NewStream -.
func NewStream(r *repo.StreamRepo, t *rabbit.StreamRabbit) *StreamUseCase {
	uc := &StreamUseCase{
		repo:   r,
		rabbit: t,
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	uc.tradeSwitchCfg = cfg.TradeSwitch
	uc.analyzeCfg = cfg.Analyze
	uc.basic = *cc.GetBasicInfo()

	go uc.checkTradeSwitch()

	bus.SubscribeTopic(topicStreamTargets, uc.ReceiveStreamData)

	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveOrderStatus(context.Background())

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
			stockNum: t.StockNum,
			orderMap: make(map[entity.OrderAction][]*entity.Order),
			// quantity should decide by bisrate
			orderQuantity: 1,
			tickChan:      make(chan *entity.RealTimeTick),
			bidAskChan:    make(chan *entity.RealTimeBidAsk),
		}
		data.setHistoryTickAnalyze(cc.GetHistoryTickAnalyze(t.StockNum))

		finishChan := make(chan struct{})
		go uc.tradeAgent(data, finishChan)
		for {
			_, ok := <-finishChan
			if !ok {
				break
			}
		}
		go uc.rabbit.TickConsumer(t.StockNum, data.tickChan)
		go uc.rabbit.BidAskConsumer(t.StockNum, data.bidAskChan)
	}
	bus.PublishTopicEvent(topicSubscribeTickTargets, ctx, targetArr)
}

func (uc *StreamUseCase) tradeAgent(data *RealTimeData, finishChan chan struct{}) {
	go func() {
		for {
			tick := <-data.tickChan
			data.tickArr = append(data.tickArr, tick)
			if data.bidAsk == nil {
				continue
			}

			order := data.generateOrder(uc.analyzeCfg)
			if order == nil {
				continue
			}

			bus.PublishTopicEvent(topicPlaceOrder, order)
			data.waitingOrder = order

			var timeout time.Duration
			switch order.Action {
			case entity.ActionBuy, entity.ActionSellFirst:
				if !uc.tradeInSwitch {
					continue
				}

				timeout = time.Duration(uc.tradeSwitchCfg.TradeInEndTime) * time.Second
			case entity.ActionSell, entity.ActionBuyLater:
				timeout = time.Duration(uc.tradeSwitchCfg.TradeOutEndTime) * time.Second
			}
			go data.checkOrderStatus(order, timeout)
		}
	}()
	go func() {
		for {
			data.bidAsk = <-data.bidAskChan
		}
	}()

	// close channel to start receive data from rabbitmq
	close(finishChan)
}

func (uc *StreamUseCase) checkTradeSwitch() {
	openTime := uc.basic.OpenTime

	for range time.Tick(5 * time.Second) {
		if uc.basic.TradeDay.After(openTime) && uc.basic.TradeDay.Before(openTime.Add(time.Duration(uc.tradeSwitchCfg.TradeInEndTime)*time.Hour)) {
			uc.tradeInSwitch = true
		} else {
			uc.tradeInSwitch = false
		}
	}
}
