package usecase

import (
	"context"
	"sync"
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

	uc.bus.PublishTopicEvent(topicStreamTickTargets, ctx, targetArr)
	uc.bus.PublishTopicEvent(topicStreamBidAskTargets, ctx, targetArr)
}

func (uc *HistoryUseCase) fetchHistoryClose(targetArr []*entity.Target) error {
	fetchTradeDayArr := GetLastNTradeDayByDate(20, CacheGetTradeDay())
	stockNumArr := []string{}
	for _, target := range targetArr {
		stockNumArr = append(stockNumArr, target.StockNum)
	}

	stockNumArrInDayMap, err := uc.findExistHistoryClose(fetchTradeDayArr, stockNumArr)
	if err != nil {
		result := make(map[time.Time][]string)
		for _, d := range fetchTradeDayArr {
			result[d] = stockNumArr
		}
		stockNumArrInDayMap = result
	}

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
	var wg sync.WaitGroup
	for d, s := range stockNumArrInDayMap {
		stockArr := s
		date := d
		wg.Add(1)
		go func() {
			defer wg.Done()
			closeArr, err := uc.grpcapi.GetStockHistoryClose(stockArr, date.Format(global.ShortTimeLayout))
			if err != nil {
				log.Error(err)
			}
			for _, close := range closeArr {
				dataChan <- &entity.HistoryClose{
					Date:     date,
					StockNum: close.GetCode(),
					Close:    close.GetClose(),
				}
			}
		}()
	}
	wg.Wait()
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
		closeArr, err := uc.repo.QueryHistoryCloseByMutltiStockNumDate(context.Background(), stockNumArr, d)
		if err != nil {
			return nil, err
		}
		for _, s := range stockNumArr {
			if close := closeArr[s]; (close != nil && close.Close == 0) || close == nil {
				stockNumArrInDay = append(stockNumArrInDay, s)
			}
		}
		result[d] = stockNumArrInDay
	}
	return result, nil
}
