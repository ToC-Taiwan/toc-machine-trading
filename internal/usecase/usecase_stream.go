package usecase

import (
	"context"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/rabbit"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/config"

	"github.com/google/go-cmp/cmp"
)

// StreamUseCase -.
type StreamUseCase struct {
	repo   StreamRepo
	rabbit StreamRabbit

	tradeSwitchCfg config.TradeSwitch
	analyzeCfg     config.Analyze
	basic          entity.BasicInfo
}

// NewStream -.
func NewStream(r *repo.StreamRepo, t *rabbit.StreamRabbit) {
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

	bus.SubscribeTopic(topicStreamTargets, uc.ReceiveStreamData)
	bus.SubscribeTopic(topicUpdateOrderStatus, uc.updateOrderSatusCache)

	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveOrderStatus(context.Background())
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
			uc.updateOrderSatusCache(ctx, order)
		}
	}()
	uc.rabbit.OrderStatusConsumer(orderStatusChan)
}

func (uc *StreamUseCase) updateOrderSatusCache(ctx context.Context, order *entity.Order) {
	cacheOrder := cc.GetOrderByOrderID(order.OrderID)
	order.TradeTime = cacheOrder.TradeTime
	if !cmp.Equal(order, cacheOrder) {
		cc.SetOrderByOrderID(order)
	}

	if err := uc.repo.InserOrUpdatetOrder(ctx, order); err != nil {
		log.Error(err)
	}
}

// ReceiveStreamData -.
func (uc *StreamUseCase) ReceiveStreamData(ctx context.Context, targetArr []*entity.Target) {
	for _, t := range targetArr {
		target := t
		data := &RealTimeData{
			stockNum: target.StockNum,
			orderMap: make(map[entity.OrderAction][]*entity.Order),
			// quantity should decide by bisrate
			orderQuantity: 1,
			tickChan:      make(chan *entity.RealTimeTick),
			bidAskChan:    make(chan *entity.RealTimeBidAsk),
		}
		finishChan := make(chan struct{})
		go uc.tradeAgent(data, finishChan)
		for {
			_, ok := <-finishChan
			if !ok {
				break
			}
		}
		go uc.rabbit.TickConsumer(target.StockNum, data.tickChan)
		go uc.rabbit.BidAskConsumer(target.StockNum, data.bidAskChan)
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
				if time.Now().After(uc.basic.TradeDay.Add(time.Duration(uc.tradeSwitchCfg.TradeInEndTime) * time.Hour)) {
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
