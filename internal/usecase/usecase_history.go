package usecase

import (
	"context"
	"fmt"
	"sort"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/global"
	"toc-machine-trading/pkg/utils"
)

// HistoryUseCase -.
type HistoryUseCase struct {
	repo    HistoryRepo
	grpcapi HistorygRPCAPI

	analyzeCfg config.Analyze
}

// NewHistory -.
func NewHistory(r *repo.HistoryRepo, t *grpcapi.HistorygRPCAPI) {
	uc := &HistoryUseCase{
		repo:    r,
		grpcapi: t,
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}
	uc.analyzeCfg = cfg.Analyze

	bus.SubscribeTopic(topicTargets, uc.FetchHistory)
}

// FetchHistory FetchHistory
func (uc *HistoryUseCase) FetchHistory(ctx context.Context, targetArr []*entity.Target) {
	err := uc.fetchHistoryKbar(targetArr)
	if err != nil {
		log.Panic(err)
	}

	err = uc.fetchHistoryTick(targetArr)
	if err != nil {
		log.Panic(err)
	}

	err = uc.fetchHistoryClose(targetArr)
	if err != nil {
		log.Panic(err)
	}

	bus.PublishTopicEvent(topicAnalyzeTargets, ctx, targetArr)
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
	result := make(map[string][]*entity.HistoryClose)
	dataChan := make(chan *entity.HistoryClose)
	wait := make(chan struct{})
	go func() {
		for {
			close, ok := <-dataChan
			if !ok {
				break
			}
			result[close.StockNum] = append(result[close.StockNum], close)
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
			dataChan <- &entity.HistoryClose{
				Date:     date,
				StockNum: close.GetCode(),
				Close:    close.GetClose(),
			}
		}
	}
	close(dataChan)
	<-wait
	if len(result) != 0 {
		log.Info("Inserting History Close")
		for _, v := range result {
			uc.processCloseArr(v)
			if err := uc.repo.InsertHistoryCloseArr(context.Background(), v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (uc *HistoryUseCase) findExistHistoryClose(fetchTradeDayArr []time.Time, stockNumArr []string) (map[time.Time][]string, error) {
	result := make(map[time.Time][]string)
	dbCloseMap := make(map[string][]*entity.HistoryClose)
	for _, d := range fetchTradeDayArr {
		closeMap, err := uc.repo.QueryMutltiStockCloseByDate(context.Background(), stockNumArr, d)
		if err != nil {
			return nil, err
		}

		var stockNumArrInDay []string
		for _, s := range stockNumArr {
			if c := closeMap[s]; c != nil && c.Close != 0 {
				dbCloseMap[s] = append(dbCloseMap[s], c)
				continue
			}
			stockNumArrInDay = append(stockNumArrInDay, s)
		}

		if len(stockNumArrInDay) != 0 {
			dErr := uc.repo.DeleteHistoryCloseByStockAndDate(context.Background(), stockNumArrInDay, d)
			if dErr != nil {
				log.Panic(dErr)
			}
			result[d] = stockNumArrInDay
		}
	}

	for _, v := range dbCloseMap {
		uc.processCloseArr(v)
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
			key := fmt.Sprintf("%s:%s", tick.StockNum, tick.TickTime.Format(global.ShortTimeLayout))
			result[key] = append(result[key], tick)
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
		log.Info("Inserting History Tick")
		for _, v := range result {
			uc.processTickArr(v)
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
		tickArrMap, err := uc.repo.QueryMultiStockTickArrByDate(context.Background(), stockNumArr, d)
		if err != nil {
			return nil, err
		}

		var stockNumArrInDay []string
		for _, s := range stockNumArr {
			if _, ok := tickArrMap[s]; !ok {
				stockNumArrInDay = append(stockNumArrInDay, s)
			} else {
				uc.processTickArr(tickArrMap[s])
			}
		}

		if len(stockNumArrInDay) != 0 {
			dErr := uc.repo.DeleteHistoryTickByStockAndDate(context.Background(), stockNumArrInDay, d)
			if dErr != nil {
				log.Panic(dErr)
			}
			result[d] = stockNumArrInDay
		}
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
			key := fmt.Sprintf("%s:%s", kbar.StockNum, kbar.KbarTime.Format(global.ShortTimeLayout))
			result[key] = append(result[key], kbar)
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
		log.Info("Inserting History Kbar")
		for _, v := range result {
			uc.processKbarArr(v)
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
		kbarArrMap, err := uc.repo.QueryMultiStockKbarArrByDate(context.Background(), stockNumArr, d)
		if err != nil {
			return nil, err
		}

		var stockNumArrInDay []string
		for _, s := range stockNumArr {
			if _, ok := kbarArrMap[s]; !ok {
				stockNumArrInDay = append(stockNumArrInDay, s)
			} else {
				uc.processKbarArr(kbarArrMap[s])
			}
		}

		if len(stockNumArrInDay) != 0 {
			dErr := uc.repo.DeleteHistoryKbarByStockAndDate(context.Background(), stockNumArrInDay, d)
			if dErr != nil {
				log.Panic(dErr)
			}
			result[d] = stockNumArrInDay
		}
	}
	return result, nil
}

func (uc *HistoryUseCase) processCloseArr(arr []*entity.HistoryClose) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Date.After(arr[j].Date)
	})

	stockNum := arr[0].StockNum
	closeArr := []float64{}
	for _, v := range arr {
		closeArr = append(closeArr, v.Close)
		cc.SetHistoryClose(stockNum, v.Date, v.Close)
	}

	biasRate, err := utils.GetBiasRateByCloseArr(closeArr)
	if err != nil {
		return
	}
	cc.SetBiasRate(stockNum, biasRate)

	i := 0
	for {
		if i+int(uc.analyzeCfg.MAPeriod) > len(closeArr) {
			break
		}
		tmp := closeArr[i : i+int(uc.analyzeCfg.MAPeriod)]
		ma := utils.GenerareMAByCloseArr(tmp)
		if err := uc.repo.InsertQuaterMA(context.Background(), &entity.HistoryAnalyze{
			Date:     arr[i].Date,
			StockNum: stockNum,
			QuaterMA: ma,
		}); err != nil {
			log.Error(err)
		}
		i++
	}
}

func (uc *HistoryUseCase) processTickArr(arr []*entity.HistoryTick) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].TickTime.Before(arr[j].TickTime)
	})
}

func (uc *HistoryUseCase) processKbarArr(arr []*entity.HistoryKbar) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].KbarTime.Before(arr[j].KbarTime)
	})
	firstKbar := arr[0]
	cc.SetHistoryOpen(firstKbar.StockNum, firstKbar.KbarTime, firstKbar.Open)
}
