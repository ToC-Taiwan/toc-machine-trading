package usecase

import (
	"context"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/global"
)

// HistoryUseCase -.
type HistoryUseCase struct {
	repo    HistoryRepo
	grpcapi HistorygRPCAPI
}

// NewHistory -.
func NewHistory(r *repo.HistoryRepo, t *grpcapi.HistorygRPCAPI) {
	uc := &HistoryUseCase{
		repo:    r,
		grpcapi: t,
	}

	bus.SubscribeTopic(topicTargets, uc.FetchHistory)
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

	bus.PublishTopicEvent(topicStreamTargets, ctx, targetArr)
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
		closeMap, err := uc.repo.QueryMutltiStockCloseByDate(context.Background(), stockNumArr, d)
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
	result := make(map[string][]*entity.HistoryTick)
	dataChan := make(chan *entity.HistoryTick)
	wait := make(chan struct{})
	go func() {
		for {
			tick, ok := <-dataChan
			if !ok {
				break
			}
			result[tick.StockNum] = append(result[tick.StockNum], tick)
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
		for _, v := range result {
			go uc.processTickArr(v)
			if err := uc.repo.InsertHistoryTickArr(context.Background(), v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (uc *HistoryUseCase) findExistHistoryTick(fetchTradeDayArr []time.Time, stockNumArr []string) (map[time.Time][]string, error) {
	result := make(map[time.Time][]string)
	for _, d := range fetchTradeDayArr {
		var stockNumArrInDay []string
		tickArrMap, err := uc.repo.QueryMultiStockTickArrByDate(context.Background(), stockNumArr, d)
		if err != nil {
			return nil, err
		}
		for _, s := range stockNumArr {
			if _, ok := tickArrMap[s]; !ok {
				stockNumArrInDay = append(stockNumArrInDay, s)
			} else {
				go uc.processTickArr(tickArrMap[s])
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
	result := make(map[string][]*entity.HistoryKbar)
	dataChan := make(chan *entity.HistoryKbar)
	wait := make(chan struct{})
	go func() {
		for {
			kbar, ok := <-dataChan
			if !ok {
				break
			}
			result[kbar.StockNum] = append(result[kbar.StockNum], kbar)
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
		for _, v := range result {
			go uc.processKbarArr(v)
			if err := uc.repo.InsertHistoryKbarArr(context.Background(), v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (uc *HistoryUseCase) findExistHistoryKbar(fetchTradeDayArr []time.Time, stockNumArr []string) (map[time.Time][]string, error) {
	result := make(map[time.Time][]string)
	for _, d := range fetchTradeDayArr {
		var stockNumArrInDay []string
		kbarArrMap, err := uc.repo.QueryMultiStockKbarArrByDate(context.Background(), stockNumArr, d)
		if err != nil {
			return nil, err
		}
		for _, s := range stockNumArr {
			if _, ok := kbarArrMap[s]; !ok {
				stockNumArrInDay = append(stockNumArrInDay, s)
			} else {
				go uc.processKbarArr(kbarArrMap[s])
			}
		}
		result[d] = stockNumArrInDay
	}
	return result, nil
}

func (uc *HistoryUseCase) processKbarArr(arr []*entity.HistoryKbar) {
	firstKbar := arr[0]
	cc.SetHistoryOpen(firstKbar.StockNum, firstKbar.KbarTime, firstKbar.Open)
}

func (uc *HistoryUseCase) processTickArr(arr []*entity.HistoryTick) {}
