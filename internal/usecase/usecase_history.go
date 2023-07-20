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
	"tmt/internal/usecase/event"
	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/module/quota"
	"tmt/internal/usecase/module/simulator"
	"tmt/internal/usecase/module/tradeday"
	"tmt/internal/usecase/repo"
	"tmt/internal/utils"
)

// HistoryUseCase -.
type HistoryUseCase struct {
	repo    HistoryRepo
	grpcapi HistorygRPCAPI

	analyzeStockCfg config.AnalyzeStock

	fetchList map[string]*entity.StockTarget
	mutex     sync.Mutex

	tradeDay *tradeday.TradeDay
	cfg      *config.Config

	simulateFutureTicks entity.RealTimeFutureTickArr
	simulateDate        tradeday.TradePeriod
	simulateCode        string

	slackMsgChan chan string
}

// NewHistory -.
func NewHistory() History {
	cfg := config.Get()
	uc := &HistoryUseCase{
		repo:            repo.NewHistory(cfg.GetPostgresPool()),
		grpcapi:         grpcapi.NewHistory(cfg.GetSinopacPool()),
		fetchList:       make(map[string]*entity.StockTarget),
		tradeDay:        tradeday.Get(),
		analyzeStockCfg: cfg.AnalyzeStock,
		cfg:             cfg,
		slackMsgChan:    make(chan string),
	}

	go uc.SendMessage()

	bus.SubscribeAsync(event.TopicFetchStockHistory, true, uc.FetchStockHistory)
	bus.SubscribeAsync(event.TopicFetchFutureHistory, true, uc.FetchFutureHistory)

	return uc
}

func (uc *HistoryUseCase) SendMessage() {
	for {
		msg := <-uc.slackMsgChan
		logger.Warn(msg)
	}
}

// GetTradeDay -.
func (uc *HistoryUseCase) GetTradeDay() time.Time {
	return uc.tradeDay.GetStockTradeDay().TradeDay
}

// GetDayKbarByStockNumDate -.
func (uc *HistoryUseCase) GetDayKbarByStockNumDate(stockNum string, date time.Time) *entity.StockHistoryKbar {
	return cc.GetDaykbar(stockNum, date)
}

// FetchStockHistory -.
func (uc *HistoryUseCase) FetchStockHistory(targetArr []*entity.StockTarget) {
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
		logger.Fatal(err)
	}

	err = uc.fetchHistoryTick(fetchArr)
	if err != nil {
		logger.Fatal(err)
	}

	err = uc.fetchHistoryClose(fetchArr)
	if err != nil {
		logger.Fatal(err)
	}

	bus.PublishTopicEvent(event.TopicAnalyzeStockTargets, fetchArr)
}

func (uc *HistoryUseCase) FetchFutureHistory(code string) {
	td := uc.tradeDay.GetLastNFutureTradeDay(1)
	ticks, err := uc.FetchFutureHistoryTick(code, td[0])
	if err != nil {
		logger.Error(err)
		return
	}

	if len(ticks) == 0 {
		logger.Error("Fetch Future History Tick Failed, Cancel subscribe future tick")
		return
	}

	uc.simulateDate = td[0]
	uc.simulateCode = code
	for _, tick := range ticks {
		uc.simulateFutureTicks = append(uc.simulateFutureTicks, tick.ToRealTimeTick())
	}

	bus.PublishTopicEvent(event.TopicSubscribeFutureTickTargets, code)
}

func (uc *HistoryUseCase) SimulateOne(cond *config.TradeFuture) *simulator.SimulateBalance {
	s := simulator.NewSimulatorFuture(simulator.SimulatorFutureTarget{
		Code:           uc.simulateCode,
		Ticks:          uc.simulateFutureTicks,
		TradePeriod:    uc.simulateDate,
		Quota:          quota.NewQuota(uc.cfg.Quota),
		TradeConfigArr: []*config.TradeFuture{cond},
	})

	result := s.OneCond()
	return result
}

func (uc *HistoryUseCase) SimulateMulti(cond []*config.TradeFuture) {
	simulator.NewSimulatorFuture(simulator.SimulatorFutureTarget{
		Code:           uc.simulateCode,
		Ticks:          uc.simulateFutureTicks,
		TradePeriod:    uc.simulateDate,
		Quota:          quota.NewQuota(uc.cfg.Quota),
		TradeConfigArr: cond,
	}).AllConds(uc.slackMsgChan)
}

func (uc *HistoryUseCase) fetchHistoryClose(targetArr []*entity.StockTarget) error {
	fetchTradeDayArr := uc.tradeDay.GetLastNStockTradeDay(uc.cfg.History.HistoryClosePeriod)
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
		closeArr, err := uc.grpcapi.GetStockHistoryClose(stockArr, date.Format(global.ShortTimeLayout))
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
				logger.Fatal(dErr)
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
	fetchTradeDayArr := uc.tradeDay.GetLastNStockTradeDay(uc.cfg.History.HistoryTickPeriod)
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
				logger.Fatal(dErr)
			}
			result[d] = stockNumArrInDay
		}
	}
	return result, total, nil
}

func (uc *HistoryUseCase) fetchHistoryKbar(targetArr []*entity.StockTarget) error {
	fetchTradeDayArr := uc.tradeDay.GetLastNStockTradeDay(uc.cfg.History.HistoryKbarPeriod)
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
				logger.Fatal(dErr)
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

	i := 0
	for {
		if i+int(uc.analyzeStockCfg.MAPeriod) > len(closeArr) {
			break
		}
		tmp := closeArr[i : i+int(uc.analyzeStockCfg.MAPeriod)]
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
	if uc.tradeDay.GetLastNStockTradeDay(1)[0].Equal(tickTradeDay) {
		cc.SetHistoryTickArr(stockNum, tickTradeDay, arr)
	}

	minPeriod := time.Duration(uc.analyzeStockCfg.TickAnalyzePeriod) * time.Millisecond
	maxPeriod := time.Duration(uc.analyzeStockCfg.TickAnalyzePeriod*1.1) * time.Millisecond

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
func (uc *HistoryUseCase) FetchFutureHistoryTick(code string, date tradeday.TradePeriod) ([]*entity.FutureHistoryTick, error) {
	dbTicks, err := uc.repo.QueryFutureHistoryTickArrByTime(context.Background(), code, date.StartTime, date.EndTime)
	if err != nil {
		return []*entity.FutureHistoryTick{}, err
	}

	if len(dbTicks) > 0 {
		return dbTicks, nil
	}

	result := []*entity.FutureHistoryTick{}
	tickArr, err := uc.grpcapi.GetFutureHistoryTick([]string{code}, date.TradeDay.Format(global.ShortTimeLayout))
	if err != nil {
		return nil, err
	}

	for _, t := range tickArr {
		result = append(result, &entity.FutureHistoryTick{
			Code: t.GetCode(),
			HistoryTickBase: entity.HistoryTickBase{
				TickTime: time.Unix(0, t.GetTs()).Add(-8 * time.Hour), Close: t.GetClose(),
				TickType: t.GetTickType(), Volume: t.GetVolume(),
				BidPrice: t.GetBidPrice(), BidVolume: t.GetBidVolume(),
				AskPrice: t.GetAskPrice(), AskVolume: t.GetAskVolume(),
			},
		})
	}

	e := uc.repo.InsertFutureHistoryTickArr(context.Background(), result)
	if e != nil {
		logger.Error(e)
		return nil, e
	}

	return result, nil
}

// FetchFutureHistoryKbar -.
func (uc *HistoryUseCase) FetchFutureHistoryKbar(code string, date time.Time) ([]*entity.FutureHistoryKbar, error) {
	result := []*entity.FutureHistoryKbar{}
	kbarArr, err := uc.grpcapi.GetFutureHistoryKbar([]string{code}, date.Format(global.ShortTimeLayout))
	if err != nil {
		return nil, err
	}

	for _, t := range kbarArr {
		result = append(result, &entity.FutureHistoryKbar{
			Code: t.GetCode(),
			HistoryKbarBase: entity.HistoryKbarBase{
				KbarTime: time.Unix(0, t.GetTs()).Add(-8 * time.Hour),
				Open:     t.GetOpen(), High: t.GetHigh(), Low: t.GetLow(),
				Close: t.GetClose(), Volume: t.GetVolume(),
			},
		})
	}
	return result, nil
}
