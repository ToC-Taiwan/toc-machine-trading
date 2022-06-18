package usecase

import (
	"context"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/eventbus"
	"toc-machine-trading/pkg/global"
)

// HistoryUseCase -.
type HistoryUseCase struct {
	repo    HistoryRepo
	grpcapi HistorygRPCAPI
	bus     *eventbus.Bus
}

// NewHistory -.
func NewHistory(r *repo.HistoryRepo, t *grpcapi.HistorygRPCAPI, bus *eventbus.Bus) {
	uc := &HistoryUseCase{
		repo:    r,
		grpcapi: t,
		bus:     bus,
	}

	if err := uc.bus.SubscribeTopic(topicTargets, uc.FetchHistory); err != nil {
		log.Panic(err)
	}
}

// FetchHistory FetchHistory
func (uc *HistoryUseCase) FetchHistory(ctx context.Context, targetArr []*entity.Target) {
	err := uc.fetchHistoryClose(targetArr)
	if err != nil {
		log.Panic(err)
	}

	err = uc.fetchHistoryTick(targetArr)
	if err != nil {
		log.Panic(err)
	}

	err = uc.fetchHistoryKbar(targetArr)
	if err != nil {
		log.Panic(err)
	}

	uc.bus.PublishTopicEvent(topicStreamTickTargets, ctx, targetArr)
	uc.bus.PublishTopicEvent(topicStreamBidAskTargets, ctx, targetArr)
}

func (uc *HistoryUseCase) fetchHistoryClose(targetArr []*entity.Target) error {
	fetchTradeDayArr := cc.GetBasicInfo().HistoryCloseRange
	stockNumArr := []string{}
	for _, target := range targetArr {
		stockNumArr = append(stockNumArr, target.StockNum)
	}

	stockNumArrInDayMap, err := uc.findExistHistoryClose(fetchTradeDayArr, stockNumArr)
	if err != nil {
		return err
	}
	defer log.Info("Fetching History Close Done")
	result := []*entity.HistoryClose{}
	dataChan := make(chan *entity.HistoryClose)
	wait := make(chan struct{})
	go func() {
		for {
			close, ok := <-dataChan
			if !ok {
				break
			}
			result = append(result, close)
		}
		close(wait)
	}()
	for d, s := range stockNumArrInDayMap {
		if len(s) == 0 {
			continue
		}
		stockArr := s
		date := d
		log.Infof("Fetching History Close -> StockCount: %d, Date: %s", len(stockArr), date.Format(global.ShortTimeLayout))
		closeArr, err := uc.grpcapi.GetStockHistoryClose(stockArr, date.Format(global.ShortTimeLayout))
		if err != nil {
			log.Error(err)
		}
		for _, close := range closeArr {
			if close.GetClose() != 0 {
				dataChan <- &entity.HistoryClose{
					Date:     date,
					StockNum: close.GetCode(),
					Close:    close.GetClose(),
				}
			}
		}
	}
	close(dataChan)
	<-wait
	if len(result) != 0 {
		if err := uc.repo.InsertHistoryCloseArr(context.Background(), result); err != nil {
			return err
		}
	}
	return nil
}

func (uc *HistoryUseCase) findExistHistoryClose(fetchTradeDayArr []time.Time, stockNumArr []string) (map[time.Time][]string, error) {
	result := make(map[time.Time][]string)
	for _, d := range fetchTradeDayArr {
		var stockNumArrInDay []string
		closeMap, err := uc.repo.QueryHistoryCloseByMutltiStockNumDate(context.Background(), stockNumArr, d)
		if err != nil {
			return nil, err
		}
		for _, s := range stockNumArr {
			if close := closeMap[s]; (close != nil && close.Close == 0) || close == nil {
				stockNumArrInDay = append(stockNumArrInDay, s)
			}
		}
		result[d] = stockNumArrInDay
	}
	return result, nil
}

func (uc *HistoryUseCase) fetchHistoryTick(targetArr []*entity.Target) error {
	fetchTradeDayArr := cc.GetBasicInfo().HistoryTickRange
	stockNumArr := []string{}
	for _, target := range targetArr {
		stockNumArr = append(stockNumArr, target.StockNum)
	}

	stockNumArrInDayMap, err := uc.findExistHistoryTick(fetchTradeDayArr, stockNumArr)
	if err != nil {
		return err
	}
	defer log.Info("Fetching History Tick Done")
	result := []*entity.HistoryTick{}
	dataChan := make(chan *entity.HistoryTick)
	wait := make(chan struct{})
	go func() {
		for {
			close, ok := <-dataChan
			if !ok {
				break
			}
			result = append(result, close)
		}
		close(wait)
	}()
	for d, s := range stockNumArrInDayMap {
		if len(s) == 0 {
			continue
		}
		stockArr := s
		date := d
		log.Infof("Fetching History Tick -> StockCount: %d, Date: %s", len(stockArr), date.Format(global.ShortTimeLayout))
		tickArr, err := uc.grpcapi.GetStockHistoryTick(stockArr, date.Format(global.ShortTimeLayout))
		if err != nil {
			log.Error(err)
		}
		for _, t := range tickArr {
			dataChan <- &entity.HistoryTick{
				StockNum: t.GetStockNum(),
				TickTime: time.Unix(0, t.GetTs()).Add(-8 * time.Hour), Close: t.GetClose(),
				TickType: t.GetTickType(), Volume: t.GetVolume(),
				BidPrice: t.GetBidPrice(), BidVolume: t.GetBidVolume(),
				AskPrice: t.GetAskPrice(), AskVolume: t.GetAskVolume(),
			}
		}
	}
	close(dataChan)
	<-wait
	if len(result) != 0 {
		if err := uc.repo.InsertHistoryTickArr(context.Background(), result); err != nil {
			return err
		}
	}
	return nil
}

func (uc *HistoryUseCase) findExistHistoryTick(fetchTradeDayArr []time.Time, stockNumArr []string) (map[time.Time][]string, error) {
	result := make(map[time.Time][]string)
	for _, d := range fetchTradeDayArr {
		var stockNumArrInDay []string
		for _, s := range stockNumArr {
			exist, err := uc.repo.CheckHistoryTickExist(context.Background(), s, d)
			if err != nil {
				return nil, err
			}
			if !exist {
				stockNumArrInDay = append(stockNumArrInDay, s)
			}
		}
		result[d] = stockNumArrInDay
	}
	return result, nil
}

func (uc *HistoryUseCase) fetchHistoryKbar(targetArr []*entity.Target) error {
	fetchTradeDayArr := cc.GetBasicInfo().HistoryKbarRange
	stockNumArr := []string{}
	for _, target := range targetArr {
		stockNumArr = append(stockNumArr, target.StockNum)
	}

	stockNumArrInDayMap, err := uc.findExistHistoryKbar(fetchTradeDayArr, stockNumArr)
	if err != nil {
		return err
	}
	defer log.Info("Fetching History Kbar Done")
	result := []*entity.HistoryKbar{}
	dataChan := make(chan *entity.HistoryKbar)
	wait := make(chan struct{})
	go func() {
		for {
			close, ok := <-dataChan
			if !ok {
				break
			}
			result = append(result, close)
		}
		close(wait)
	}()
	for d, s := range stockNumArrInDayMap {
		if len(s) == 0 {
			continue
		}
		stockArr := s
		date := d
		log.Infof("Fetching History Kbar -> StockCount: %d, Date: %s", len(stockArr), date.Format(global.ShortTimeLayout))
		tickArr, err := uc.grpcapi.GetStockHistoryKbar(stockArr, date.Format(global.ShortTimeLayout))
		if err != nil {
			log.Error(err)
		}
		for _, t := range tickArr {
			dataChan <- &entity.HistoryKbar{
				StockNum: t.GetStockNum(),
				KbarTime: time.Unix(0, t.GetTs()).Add(-8 * time.Hour),
				Open:     t.GetOpen(),
				High:     t.GetHigh(),
				Low:      t.GetLow(),
				Close:    t.GetClose(),
				Volume:   t.GetVolume(),
			}
		}
	}
	close(dataChan)
	<-wait
	if len(result) != 0 {
		if err := uc.repo.InsertHistoryKbarArr(context.Background(), result); err != nil {
			return err
		}
	}
	return nil
}

func (uc *HistoryUseCase) findExistHistoryKbar(fetchTradeDayArr []time.Time, stockNumArr []string) (map[time.Time][]string, error) {
	result := make(map[time.Time][]string)
	for _, d := range fetchTradeDayArr {
		var stockNumArrInDay []string
		for _, s := range stockNumArr {
			exist, err := uc.repo.CheckHistoryKbarExist(context.Background(), s, d)
			if err != nil {
				return nil, err
			}
			if !exist {
				stockNumArrInDay = append(stockNumArrInDay, s)
			}
		}
		result[d] = stockNumArrInDay
	}
	return result, nil
}
