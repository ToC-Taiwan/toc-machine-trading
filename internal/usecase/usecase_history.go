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
	"tmt/internal/usecase/modules/event"
	"tmt/internal/usecase/modules/tradeday"
	"tmt/internal/usecase/modules/trader"
	"tmt/pkg/utils"
)

// HistoryUseCase -.
type HistoryUseCase struct {
	repo    HistoryRepo
	grpcapi HistorygRPCAPI

	stockAnalyzeCfg config.StockAnalyze
	basic           entity.BasicInfo

	fetchList map[string]*entity.Target
	mutex     sync.Mutex

	biasRateArr []float64

	tradeDay           *tradeday.TradeDay
	simulateFutureCode string
}

// NewHistory -.
func NewHistory(r HistoryRepo, t HistorygRPCAPI) *HistoryUseCase {
	uc := &HistoryUseCase{
		repo:      r,
		grpcapi:   t,
		fetchList: make(map[string]*entity.Target),
		tradeDay:  tradeday.NewTradeDay(),
	}

	cfg := config.GetConfig()
	uc.stockAnalyzeCfg = cfg.StockAnalyze
	uc.basic = *cc.GetBasicInfo()

	bus.SubscribeTopic(event.TopicFetchStockHistory, uc.FetchHistory)
	bus.SubscribeTopic(event.TopicMonitorFutureCode, uc.updateSimulateFutureCode)
	return uc
}

func (uc *HistoryUseCase) updateSimulateFutureCode(future *entity.Future) {
	uc.simulateFutureCode = future.Code
}

// GetTradeDay -.
func (uc *HistoryUseCase) GetTradeDay() time.Time {
	return uc.basic.TradeDay
}

// GetDayKbarByStockNumDate -.
func (uc *HistoryUseCase) GetDayKbarByStockNumDate(stockNum string, date time.Time) *entity.HistoryKbar {
	return cc.GetDaykbar(stockNum, date)
}

// FetchHistory FetchHistory
func (uc *HistoryUseCase) FetchHistory(ctx context.Context, targetArr []*entity.Target) {
	defer uc.mutex.Unlock()
	uc.mutex.Lock()

	var fetchArr []*entity.Target
	for _, v := range targetArr {
		if _, ok := uc.fetchList[v.StockNum]; !ok {
			uc.fetchList[v.StockNum] = v
			fetchArr = append(fetchArr, v)
		}
	}

	if len(fetchArr) == 0 {
		return
	}

	err := uc.fetchHistoryKbar(fetchArr)
	if err != nil {
		log.Panic(err)
	}

	err = uc.fetchHistoryTick(fetchArr)
	if err != nil {
		log.Panic(err)
	}

	err = uc.fetchHistoryClose(fetchArr)
	if err != nil {
		log.Panic(err)
	}

	bus.PublishTopicEvent(event.TopicAnalyzeStockTargets, ctx, fetchArr)
}

func (uc *HistoryUseCase) fetchHistoryClose(targetArr []*entity.Target) error {
	fetchTradeDayArr := cc.GetBasicInfo().HistoryCloseRange
	stockNumArr := []string{}
	for _, target := range targetArr {
		stockNumArr = append(stockNumArr, target.StockNum)
	}

	stockNumArrInDayMap, total, err := uc.findExistHistoryClose(fetchTradeDayArr, stockNumArr)
	if err != nil {
		return err
	}
	defer log.Info("Fetching History Close Done")
	log.Infof("Fetching History Close -> Count: %d", total)
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
	uc.processBiasRate()
	return nil
}

func (uc *HistoryUseCase) findExistHistoryClose(fetchTradeDayArr []time.Time, stockNumArr []string) (map[time.Time][]string, int64, error) {
	log.Info("Query Exist History Close")
	result := make(map[time.Time][]string)
	dbCloseMap := make(map[string][]*entity.HistoryClose)
	var total int64
	for _, d := range fetchTradeDayArr {
		closeMap, err := uc.repo.QueryMutltiStockCloseByDate(context.Background(), stockNumArr, d)
		if err != nil {
			return nil, 0, err
		}

		var stockNumArrInDay []string
		for _, s := range stockNumArr {
			if c := closeMap[s]; c != nil && c.Close != 0 {
				dbCloseMap[s] = append(dbCloseMap[s], c)
				continue
			}
			total++
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
	return result, total, nil
}

func (uc *HistoryUseCase) fetchHistoryTick(targetArr []*entity.Target) error {
	fetchTradeDayArr := cc.GetBasicInfo().HistoryTickRange
	stockNumArr := []string{}
	for _, target := range targetArr {
		stockNumArr = append(stockNumArr, target.StockNum)
	}

	stockNumArrInDayMap, total, err := uc.findExistStockHistoryTick(fetchTradeDayArr, stockNumArr)
	if err != nil {
		return err
	}
	defer log.Info("Fetching History Tick Done")
	log.Infof("Fetching History Tick -> Count: %d", total)
	result := make(map[string][]*entity.StockHistoryTick)
	dataChan := make(chan *entity.StockHistoryTick)
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
		tickArr, err := uc.grpcapi.GetStockHistoryTick(stockArr, date.Format(global.ShortTimeLayout))
		if err != nil {
			log.Error(err)
		}
		for _, t := range tickArr {
			dataChan <- &entity.StockHistoryTick{
				StockNum: t.GetCode(),
				HistoryTickBase: entity.HistoryTickBase{
					TickTime: time.Unix(0, t.GetTs()).Add(-8 * time.Hour), Close: t.GetClose(),
					TickType: t.GetTickType(), Volume: t.GetVolume(),
					BidPrice: t.GetBidPrice(), BidVolume: t.GetBidVolume(),
					AskPrice: t.GetAskPrice(), AskVolume: t.GetAskVolume(),
				},
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

func (uc *HistoryUseCase) findExistStockHistoryTick(fetchTradeDayArr []time.Time, stockNumArr []string) (map[time.Time][]string, int64, error) {
	log.Info("Query Exist History Tick")
	result := make(map[time.Time][]string)
	var total int64
	for _, d := range fetchTradeDayArr {
		tickArrMap, err := uc.repo.QueryMultiStockTickArrByDate(context.Background(), stockNumArr, d)
		if err != nil {
			return nil, 0, err
		}

		var stockNumArrInDay []string
		for _, s := range stockNumArr {
			if _, ok := tickArrMap[s]; !ok {
				total++
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
	return result, total, nil
}

func (uc *HistoryUseCase) fetchHistoryKbar(targetArr []*entity.Target) error {
	fetchTradeDayArr := cc.GetBasicInfo().HistoryKbarRange
	stockNumArr := []string{}
	for _, target := range targetArr {
		stockNumArr = append(stockNumArr, target.StockNum)
	}

	stockNumArrInDayMap, total, err := uc.findExistHistoryKbar(fetchTradeDayArr, stockNumArr)
	if err != nil {
		return err
	}
	defer log.Info("Fetching History Kbar Done")
	log.Infof("Fetching History Kbar -> Count: %d", total)
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
		tickArr, err := uc.grpcapi.GetStockHistoryKbar(stockArr, date.Format(global.ShortTimeLayout))
		if err != nil {
			log.Error(err)
		}
		for _, t := range tickArr {
			dataChan <- &entity.HistoryKbar{
				StockNum: t.GetCode(),
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

func (uc *HistoryUseCase) findExistHistoryKbar(fetchTradeDayArr []time.Time, stockNumArr []string) (map[time.Time][]string, int64, error) {
	log.Info("Query Exist History Kbar")
	result := make(map[time.Time][]string)
	var total int64
	for _, d := range fetchTradeDayArr {
		kbarArrMap, err := uc.repo.QueryMultiStockKbarArrByDate(context.Background(), stockNumArr, d)
		if err != nil {
			return nil, 0, err
		}

		var stockNumArrInDay []string
		for _, s := range stockNumArr {
			if _, ok := kbarArrMap[s]; !ok {
				total++
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
	return result, total, nil
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

	if biasRate != 0 {
		uc.biasRateArr = append(uc.biasRateArr, biasRate)
		cc.SetBiasRate(stockNum, biasRate)
	}

	i := 0
	for {
		if i+int(uc.stockAnalyzeCfg.MAPeriod) > len(closeArr) {
			break
		}
		tmp := closeArr[i : i+int(uc.stockAnalyzeCfg.MAPeriod)]
		ma := utils.GenerareMAByCloseArr(tmp)
		if err := uc.repo.InsertQuaterMA(context.Background(), &entity.HistoryAnalyze{
			Date:     arr[i].Date,
			StockNum: stockNum,
			QuaterMA: utils.Round(ma, 2),
		}); err != nil {
			log.Error(err)
		}
		i++
	}
}

func (uc *HistoryUseCase) processBiasRate() {
	sort.Slice(uc.biasRateArr, func(i, j int) bool {
		return uc.biasRateArr[i] > uc.biasRateArr[j]
	})

	total := len(uc.biasRateArr)
	cc.SetHighBiasRate(uc.biasRateArr[total/4])
	cc.SetLowBiasRate(uc.biasRateArr[3*total/4])
}

func (uc *HistoryUseCase) processTickArr(arr []*entity.StockHistoryTick) {
	if len(arr) < 2 {
		return
	}

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].TickTime.Before(arr[j].TickTime)
	})

	stockNum := arr[0].StockNum
	firsTickTime := arr[0].TickTime
	tickTradeDay := time.Date(firsTickTime.Year(), firsTickTime.Month(), firsTickTime.Day(), 0, 0, 0, 0, time.Local)
	if uc.basic.LastTradeDay.Equal(tickTradeDay) {
		cc.SetHistoryTickArr(stockNum, tickTradeDay, arr)
	}

	minPeriod := time.Duration(uc.stockAnalyzeCfg.TickAnalyzePeriod) * time.Millisecond
	maxPeriod := time.Duration(uc.stockAnalyzeCfg.TickAnalyzePeriod*1.1) * time.Millisecond

	var volumeArr []int64
	var periodVolume int64

	startTime := arr[1].TickTime
	for _, tick := range arr[1:] {
		if tick.TickTime.Sub(startTime) > maxPeriod {
			periodVolume = tick.Volume
			startTime = tick.TickTime
			continue
		}

		if tick.TickTime.Sub(startTime) < minPeriod {
			periodVolume += tick.Volume
		} else {
			volumeArr = append(volumeArr, periodVolume)

			periodVolume = tick.Volume
			startTime = tick.TickTime
		}
	}
	cc.AppendHistoryTickAnalyze(stockNum, volumeArr)
}

func (uc *HistoryUseCase) processKbarArr(arr []*entity.HistoryKbar) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].KbarTime.Before(arr[j].KbarTime)
	})
	firstKbar := arr[0]
	cc.SetHistoryOpen(firstKbar.StockNum, firstKbar.KbarTime, firstKbar.Open)
	var close, open, high, low float64
	var volume int64
	var lastKbarTime time.Time
	for i, kbar := range arr {
		if i == 0 {
			open = kbar.Open
		}
		if i == len(arr)-1 {
			close = kbar.Close
			lastKbarTime = kbar.KbarTime
		}
		if high == 0 {
			high = kbar.High
		} else if kbar.High > high {
			high = kbar.High
		}
		if low == 0 {
			low = kbar.Low
		} else if kbar.Low < low {
			low = kbar.Low
		}
		volume += kbar.Volume
	}
	cc.SetDaykbar(firstKbar.StockNum, firstKbar.KbarTime, &entity.HistoryKbar{
		StockNum: firstKbar.StockNum,
		KbarTime: lastKbarTime,
		Open:     open,
		High:     high,
		Low:      low,
		Close:    close,
		Volume:   volume,
		Stock:    cc.GetStockDetail(firstKbar.StockNum),
	})
}

// FetchFutureHistoryTick -.
func (uc *HistoryUseCase) FetchFutureHistoryTick(code string, date time.Time) []*entity.FutureHistoryTick {
	result := make(map[string][]*entity.FutureHistoryTick)
	dataChan := make(chan *entity.FutureHistoryTick)
	wait := make(chan struct{})
	go func() {
		for {
			tick, ok := <-dataChan
			if !ok {
				break
			}
			result[tick.Code] = append(result[tick.Code], tick)
		}
		close(wait)
	}()
	tickArr, err := uc.grpcapi.GetFutureHistoryTick([]string{code}, date.Format(global.ShortTimeLayout))
	if err != nil {
		log.Error(err)
	}
	for _, t := range tickArr {
		dataChan <- &entity.FutureHistoryTick{
			Code: t.GetCode(),
			HistoryTickBase: entity.HistoryTickBase{
				TickTime: time.Unix(0, t.GetTs()).Add(-8 * time.Hour), Close: t.GetClose(),
				TickType: t.GetTickType(), Volume: t.GetVolume(),
				BidPrice: t.GetBidPrice(), BidVolume: t.GetBidVolume(),
				AskPrice: t.GetAskPrice(), AskVolume: t.GetAskVolume(),
			},
		}
	}
	close(dataChan)
	<-wait
	return result[code]
}

func (uc *HistoryUseCase) findExistFutureHistoryTick(date tradeday.TradePeriod, code string) ([][]*entity.FutureHistoryTick, error) {
	log.Infof("Fetching %s tick %s", code, date.TradeDay.Format(global.ShortTimeLayout))
	dbTickArr, err := uc.repo.QueryFutureTickArrByTime(context.Background(), code, date.StartTime, date.EndTime)
	if err != nil {
		return nil, err
	}

	if len(dbTickArr) == 0 {
		dbTickArr = uc.FetchFutureHistoryTick(code, date.TradeDay)
		if dbTickArr == nil {
			return nil, fmt.Errorf("fetch %s tick failed", date.TradeDay.Format(global.ShortTimeLayout))
		}

		err = uc.repo.InsertFutureHistoryTickArr(context.Background(), dbTickArr)
		if err != nil {
			return nil, err
		}
	}

	var cut int
	for i, v := range dbTickArr {
		if v.TickTime.After(date.StartTime.Add(14 * time.Hour)) {
			cut = i
			break
		}
	}

	firstPart := dbTickArr[:cut]
	secondPart := dbTickArr[cut:]

	return [][]*entity.FutureHistoryTick{firstPart, secondPart}, nil
}

func (uc *HistoryUseCase) GetFutureTradeCond(days int) trader.TradeBalance {
	simulateDateArr := uc.tradeDay.GetLastNFutureTradeDay(days)
	var balanceArr []trader.TradeBalance
	for _, date := range simulateDateArr {
		dbTickArrArr, err := uc.findExistFutureHistoryTick(date, uc.simulateFutureCode)
		if err != nil {
			log.Error(err)
			continue
		}

		cond := config.FutureAnalyze{
			MaxHoldTime: 5,
		}

		for _, dbTickArr := range dbTickArrArr {
			simulator := trader.NewFutureSimulator(uc.simulateFutureCode, cond, date)
			tickChan := simulator.GetTickChan()
			for _, tick := range dbTickArr {
				tickChan <- &entity.RealTimeFutureTick{
					Code:     uc.simulateFutureCode,
					TickTime: tick.TickTime,
					Close:    tick.Close,
					Volume:   tick.Volume,
					TickType: tick.TickType,
				}
			}
			close(tickChan)
			balanceArr = append(balanceArr, simulator.CalculateFutureTradeBalance())
		}
	}

	var totalCount, totalBalance int64
	for _, balance := range balanceArr {
		totalCount += balance.Count
		totalBalance += balance.Balance
	}

	return trader.TradeBalance{
		Count:   totalCount,
		Balance: totalBalance,
	}
}
