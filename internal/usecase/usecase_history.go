package usecase

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/event"
	"tmt/internal/usecase/modules/tradeday"
	"tmt/internal/usecase/modules/trader"
	"tmt/pkg/common"
	"tmt/pkg/utils"
)

// HistoryUseCase -.
type HistoryUseCase struct {
	repo    HistoryRepo
	grpcapi HistorygRPCAPI

	stockAnalyzeCfg config.StockAnalyze
	basic           entity.BasicInfo

	fetchList map[string]*entity.StockTarget
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
		fetchList: make(map[string]*entity.StockTarget),
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
func (uc *HistoryUseCase) GetDayKbarByStockNumDate(stockNum string, date time.Time) *entity.StockHistoryKbar {
	return cc.GetDaykbar(stockNum, date)
}

// FetchHistory FetchHistory
func (uc *HistoryUseCase) FetchHistory(ctx context.Context, targetArr []*entity.StockTarget) {
	defer uc.mutex.Unlock()
	uc.mutex.Lock()

	var fetchArr []*entity.StockTarget
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
		logger.Panic(err)
	}

	err = uc.fetchHistoryTick(fetchArr)
	if err != nil {
		logger.Panic(err)
	}

	err = uc.fetchHistoryClose(fetchArr)
	if err != nil {
		logger.Panic(err)
	}

	bus.PublishTopicEvent(event.TopicAnalyzeStockTargets, ctx, fetchArr)
}

func (uc *HistoryUseCase) fetchHistoryClose(targetArr []*entity.StockTarget) error {
	fetchTradeDayArr := cc.GetBasicInfo().HistoryCloseRange
	stockNumArr := []string{}
	for _, target := range targetArr {
		stockNumArr = append(stockNumArr, target.StockNum)
	}

	stockNumArrInDayMap, total, err := uc.findExistHistoryClose(fetchTradeDayArr, stockNumArr)
	if err != nil {
		return err
	}
	defer logger.Info("Fetching History Close Done")
	logger.Infof("Fetching History Close -> Count: %d", total)
	result := make(map[string][]*entity.StockHistoryClose)
	dataChan := make(chan *entity.StockHistoryClose)
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
		closeArr, err := uc.grpcapi.GetStockHistoryClose(stockArr, date.Format(common.ShortTimeLayout))
		if err != nil {
			logger.Error(err)
		}
		for _, close := range closeArr {
			dataChan <- &entity.StockHistoryClose{
				StockNum: close.GetCode(),
				HistoryCloseBase: entity.HistoryCloseBase{
					Date:  date,
					Close: close.GetClose(),
				},
			}
		}
	}
	close(dataChan)
	<-wait
	if len(result) != 0 {
		logger.Info("Inserting History Close")
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
	logger.Info("Query Exist History Close")
	result := make(map[time.Time][]string)
	dbCloseMap := make(map[string][]*entity.StockHistoryClose)
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
				logger.Panic(dErr)
			}
			result[d] = stockNumArrInDay
		}
	}

	for _, v := range dbCloseMap {
		uc.processCloseArr(v)
	}
	return result, total, nil
}

func (uc *HistoryUseCase) fetchHistoryTick(targetArr []*entity.StockTarget) error {
	fetchTradeDayArr := cc.GetBasicInfo().HistoryTickRange
	stockNumArr := []string{}
	for _, target := range targetArr {
		stockNumArr = append(stockNumArr, target.StockNum)
	}

	stockNumArrInDayMap, total, err := uc.findExistStockHistoryTick(fetchTradeDayArr, stockNumArr)
	if err != nil {
		return err
	}
	defer logger.Info("Fetching History Tick Done")
	logger.Infof("Fetching History Tick -> Count: %d", total)
	result := make(map[string][]*entity.StockHistoryTick)
	dataChan := make(chan *entity.StockHistoryTick)
	wait := make(chan struct{})
	go func() {
		for {
			tick, ok := <-dataChan
			if !ok {
				break
			}
			key := fmt.Sprintf("%s:%s", tick.StockNum, tick.TickTime.Format(common.ShortTimeLayout))
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
		tickArr, err := uc.grpcapi.GetStockHistoryTick(stockArr, date.Format(common.ShortTimeLayout))
		if err != nil {
			logger.Error(err)
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
		logger.Info("Inserting History Tick")
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
	logger.Info("Query Exist History Tick")
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
				logger.Panic(dErr)
			}
			result[d] = stockNumArrInDay
		}
	}
	return result, total, nil
}

func (uc *HistoryUseCase) fetchHistoryKbar(targetArr []*entity.StockTarget) error {
	fetchTradeDayArr := cc.GetBasicInfo().HistoryKbarRange
	stockNumArr := []string{}
	for _, target := range targetArr {
		stockNumArr = append(stockNumArr, target.StockNum)
	}

	stockNumArrInDayMap, total, err := uc.findExistHistoryKbar(fetchTradeDayArr, stockNumArr)
	if err != nil {
		return err
	}
	defer logger.Info("Fetching History Kbar Done")
	logger.Infof("Fetching History Kbar -> Count: %d", total)
	result := make(map[string][]*entity.StockHistoryKbar)
	dataChan := make(chan *entity.StockHistoryKbar)
	wait := make(chan struct{})
	go func() {
		for {
			kbar, ok := <-dataChan
			if !ok {
				break
			}
			key := fmt.Sprintf("%s:%s", kbar.StockNum, kbar.KbarTime.Format(common.ShortTimeLayout))
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
		tickArr, err := uc.grpcapi.GetStockHistoryKbar(stockArr, date.Format(common.ShortTimeLayout))
		if err != nil {
			logger.Error(err)
		}
		for _, t := range tickArr {
			dataChan <- &entity.StockHistoryKbar{
				StockNum: t.GetCode(),
				HistoryKbarBase: entity.HistoryKbarBase{
					KbarTime: time.Unix(0, t.GetTs()).Add(-8 * time.Hour),
					Open:     t.GetOpen(), High: t.GetHigh(), Low: t.GetLow(),
					Close: t.GetClose(), Volume: t.GetVolume(),
				},
			}
		}
	}
	close(dataChan)
	<-wait
	if len(result) != 0 {
		logger.Info("Inserting History Kbar")
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
	logger.Info("Query Exist History Kbar")
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
				logger.Panic(dErr)
			}
			result[d] = stockNumArrInDay
		}
	}
	return result, total, nil
}

func (uc *HistoryUseCase) processCloseArr(arr []*entity.StockHistoryClose) {
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
		if err := uc.repo.InsertQuaterMA(context.Background(), &entity.StockHistoryAnalyze{
			StockNum: stockNum,
			HistoryAnalyzeBase: entity.HistoryAnalyzeBase{
				Date:     arr[i].Date,
				QuaterMA: utils.Round(ma, 2),
			},
		}); err != nil {
			logger.Error(err)
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

func (uc *HistoryUseCase) processKbarArr(arr []*entity.StockHistoryKbar) {
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
	cc.SetDaykbar(firstKbar.StockNum, firstKbar.KbarTime, &entity.StockHistoryKbar{
		StockNum: firstKbar.StockNum,
		Stock:    cc.GetStockDetail(firstKbar.StockNum),
		HistoryKbarBase: entity.HistoryKbarBase{
			KbarTime: lastKbarTime,
			Open:     open,
			High:     high,
			Low:      low,
			Close:    close,
			Volume:   volume,
		},
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
	tickArr, err := uc.grpcapi.GetFutureHistoryTick([]string{code}, date.Format(common.ShortTimeLayout))
	if err != nil {
		logger.Error(err)
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
	dbTickArr, err := uc.repo.QueryFutureHistoryTickArrByTime(context.Background(), code, date.StartTime, date.EndTime)
	if err != nil {
		return nil, err
	}

	if len(dbTickArr) == 0 {
		dbTickArr = uc.FetchFutureHistoryTick(code, date.TradeDay)
		if dbTickArr == nil {
			return nil, fmt.Errorf("fetch %s tick failed", date.TradeDay.Format(common.ShortTimeLayout))
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

func (uc *HistoryUseCase) fetchFutureHistoryClose(code string, date time.Time) *entity.FutureHistoryClose {
	result := make(map[string]*entity.FutureHistoryClose)
	dataChan := make(chan *entity.FutureHistoryClose)
	wait := make(chan struct{})
	go func() {
		for {
			close, ok := <-dataChan
			if !ok {
				break
			}
			result[close.Code] = close
		}
		close(wait)
	}()
	closeArr, err := uc.grpcapi.GetFutureHistoryClose([]string{code}, date.Format(common.ShortTimeLayout))
	if err != nil {
		logger.Error(err)
	}
	for _, close := range closeArr {
		dataChan <- &entity.FutureHistoryClose{
			Code: close.GetCode(),
			HistoryCloseBase: entity.HistoryCloseBase{
				Date:  date,
				Close: close.GetClose(),
			},
		}
	}
	close(dataChan)
	<-wait
	if len(result) != 0 {
		return result[code]
	}
	return nil
}

func (uc *HistoryUseCase) findExistFutureHistoryClose(date tradeday.TradePeriod, code string) (*entity.FutureHistoryClose, error) {
	dbClose, err := uc.repo.QueryFutureHistoryCloseByDate(context.Background(), code, date.TradeDay)
	if err != nil {
		return nil, err
	}

	if dbClose.Close == 0 {
		dbClose = uc.fetchFutureHistoryClose(code, date.TradeDay)
		if dbClose == nil {
			return nil, fmt.Errorf("fetch %s close failed", date.TradeDay.Format(common.ShortTimeLayout))
		}

		err = uc.repo.InsertFutureHistoryClose(context.Background(), dbClose)
		if err != nil {
			return nil, err
		}
	}
	return dbClose, nil
}

func (uc *HistoryUseCase) GetFutureTradeCond(days int) trader.SimulateBalance {
	simulateDateArr := uc.tradeDay.GetLastNFutureTradeDay(days)
	var balanceArr []trader.SimulateBalance
	for _, date := range simulateDateArr {
		dbTickArrArr, err := uc.findExistFutureHistoryTick(date, uc.simulateFutureCode)
		if err != nil {
			logger.Error(err)
			continue
		}

		lastPeriod := date.GetLastFutureTradePeriod()
		dbClose, err := uc.findExistFutureHistoryClose(lastPeriod, uc.simulateFutureCode)
		if err != nil {
			logger.Error(err)
			continue
		}

		cond := config.FutureAnalyze{
			MaxHoldTime: 20,
		}

		logger.Infof("Simulating %s %s, last close: %.0f", uc.simulateFutureCode, date.TradeDay.Format(common.ShortTimeLayout), dbClose.Close)
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
					PctChg:   utils.Round(100*(tick.Close-dbClose.Close)/dbClose.Close, 2),
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

	return trader.SimulateBalance{
		Count:   totalCount,
		Balance: totalBalance,
	}
}

// FetchFutureHistoryKbar -.
func (uc *HistoryUseCase) FetchFutureHistoryKbar(code string, date time.Time) ([]*entity.FutureHistoryKbar, error) {
	result := make(map[string][]*entity.FutureHistoryKbar)
	dataChan := make(chan *entity.FutureHistoryKbar)
	wait := make(chan struct{})
	go func() {
		for {
			kbar, ok := <-dataChan
			if !ok {
				break
			}
			// key := fmt.Sprintf("%s:%s", kbar.StockNum, kbar.KbarTime.Format(common.ShortTimeLayout))
			result[code] = append(result[code], kbar)
		}
		close(wait)
	}()

	kbarArr, err := uc.grpcapi.GetFutureHistoryKbar([]string{code}, date.Format(common.ShortTimeLayout))
	if err != nil {
		return nil, err
	}

	for _, t := range kbarArr {
		dataChan <- &entity.FutureHistoryKbar{
			Code: t.GetCode(),
			HistoryKbarBase: entity.HistoryKbarBase{
				KbarTime: time.Unix(0, t.GetTs()).Add(-8 * time.Hour),
				Open:     t.GetOpen(), High: t.GetHigh(), Low: t.GetLow(),
				Close: t.GetClose(), Volume: t.GetVolume(),
			},
		}
	}
	close(dataChan)
	<-wait
	return result[code], nil
}
