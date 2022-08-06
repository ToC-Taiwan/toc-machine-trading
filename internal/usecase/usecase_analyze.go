package usecase

import (
	"context"
	"sort"
	"sync"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/global"
)

// AnalyzeUseCase -.
type AnalyzeUseCase struct {
	repo      HistoryRepo
	targetArr []*entity.Target

	basic       entity.BasicInfo
	tradeSwitch config.TradeSwitch
	quotaCfg    config.Quota
	analyzeCfg  config.Analyze

	beforeHistoryClose map[string]float64
	historyTick        map[string]*[]*entity.HistoryTick
	historyDataLock    sync.RWMutex

	lastBelowMAStock map[string]*entity.HistoryAnalyze
	rebornMap        map[time.Time][]entity.Stock
	rebornLock       sync.Mutex
}

// NewAnalyze -.
func NewAnalyze(r HistoryRepo) *AnalyzeUseCase {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	uc := &AnalyzeUseCase{
		repo:               r,
		basic:              *cc.GetBasicInfo(),
		tradeSwitch:        cfg.TradeSwitch,
		quotaCfg:           cfg.Quota,
		analyzeCfg:         cfg.Analyze,
		beforeHistoryClose: make(map[string]float64),
		historyTick:        make(map[string]*[]*entity.HistoryTick),
		lastBelowMAStock:   make(map[string]*entity.HistoryAnalyze),
		rebornMap:          make(map[time.Time][]entity.Stock),
	}

	bus.SubscribeTopic(topicAnalyzeTargets, uc.AnalyzeAll)
	return uc
}

type simulateResult struct {
	cfg     config.Analyze
	balance *entity.TradeBalance
	orders  []*entity.Order
}

// AnalyzeAll -.
func (uc *AnalyzeUseCase) AnalyzeAll(ctx context.Context, targetArr []*entity.Target) {
	uc.findBelowQuaterMATargets(ctx, targetArr)
}

func (uc *AnalyzeUseCase) findBelowQuaterMATargets(ctx context.Context, targetArr []*entity.Target) {
	defer uc.rebornLock.Unlock()
	uc.rebornLock.Lock()
	uc.targetArr = append(uc.targetArr, targetArr...)

	for _, t := range targetArr {
		maMap, err := uc.repo.QueryAllQuaterMAByStockNum(ctx, t.StockNum)
		if err != nil {
			log.Panic(err)
		}

		basicInfo := cc.GetBasicInfo()
		for _, ma := range maMap {
			tmp := ma
			if close := cc.GetHistoryClose(ma.StockNum, ma.Date); close != 0 && close-ma.QuaterMA > 0 {
				continue
			}
			if nextTradeDay := getAbsNextTradeDayTime(ma.Date); nextTradeDay.Equal(basicInfo.TradeDay) {
				uc.lastBelowMAStock[tmp.StockNum] = tmp
			} else if nextOpen := cc.GetHistoryOpen(ma.StockNum, nextTradeDay); nextOpen != 0 && nextOpen-ma.QuaterMA > 0 {
				uc.rebornMap[ma.Date] = append(uc.rebornMap[ma.Date], *tmp.Stock)
			}
		}
	}
	log.Info("FindBelowQuaterMATargets Done")
}

// GetRebornMap -.
func (uc *AnalyzeUseCase) GetRebornMap(ctx context.Context) map[time.Time][]entity.Stock {
	uc.rebornLock.Lock()
	basicInfo := cc.GetBasicInfo()
	if len(uc.lastBelowMAStock) != 0 {
		for _, s := range uc.lastBelowMAStock {
			if open := cc.GetHistoryOpen(s.Stock.Number, basicInfo.TradeDay); open != 0 {
				if open > s.QuaterMA {
					uc.rebornMap[s.Date] = append(uc.rebornMap[s.Date], *s.Stock)
				}
				delete(uc.lastBelowMAStock, s.Stock.Number)
			}
		}
	}
	uc.rebornLock.Unlock()
	return uc.rebornMap
}

// FillHistoryTick -.
func (uc *AnalyzeUseCase) FillHistoryTick(targetArr []*entity.Target) {
	for _, t := range targetArr {
		tickArr := cc.GetHistoryTickArr(t.StockNum, uc.basic.LastTradeDay)[1:]
		beforeLastTradeDayClose := cc.GetHistoryClose(t.StockNum, uc.basic.BefroeLastTradeDay)

		uc.historyDataLock.Lock()
		uc.historyTick[t.StockNum] = &tickArr
		uc.beforeHistoryClose[t.StockNum] = beforeLastTradeDayClose
		uc.historyDataLock.Unlock()
	}
}

// SimulateOnHistoryTick -.
func (uc *AnalyzeUseCase) SimulateOnHistoryTick(ctx context.Context, useDefault bool) {
	if len(uc.targetArr) == 0 {
		return
	}

	uc.FillHistoryTick(uc.targetArr)
	resultChan := make(chan simulateResult)

	go func() {
		var bestCfg config.Analyze
		var bestBalance entity.TradeBalance
		for {
			res, ok := <-resultChan
			if !ok {
				break
			}

			if res.balance.Total > bestBalance.Total {
				bestBalance = *res.balance
				bestCfg = res.cfg
				log.Infof("TradeCount: %d, Forward: %d, Reverse: %d, Discount: %d, Total: %d", bestBalance.TradeCount, bestBalance.Forward, bestBalance.Reverse, bestBalance.Discount, bestBalance.Total)
				log.Warnf("OutInRatio %.0f", bestCfg.OutInRatio)
				log.Warnf("InOutRatio: %.0f", bestCfg.InOutRatio)
				log.Warnf("VolumePRLimit: %.0f", bestCfg.VolumePRLimit)
				log.Warnf("TickAnalyzePeriod: %.0f", bestCfg.TickAnalyzePeriod)
				log.Warnf("RSIMinCount: %d", bestCfg.RSIMinCount)
				log.Warnf("RSIHigh: %.1f", bestCfg.RSIHigh)
				log.Warnf("RSILow: %.1f", bestCfg.RSILow)
				for _, o := range res.orders {
					log.Warnf("TradeTime: %s, Stock: %s, Action: %d, Qty: %d, Price: %.2f", o.TradeTime.Format(global.LongTimeLayout), o.StockNum, o.Action, o.Quantity, o.Price)
				}
			}
		}
	}()

	for _, cfg := range generateAnalyzeCfg(useDefault) {
		simCfg, balance, orders := uc.getSimulateCond(uc.targetArr, cfg)
		resultChan <- simulateResult{
			cfg:     simCfg,
			balance: balance,
			orders:  orders,
		}
	}
	close(resultChan)
	log.Info("Simulate Done")
}

func (uc *AnalyzeUseCase) getSimulateCond(targetArr []*entity.Target, analyzeCfg config.Analyze) (config.Analyze, *entity.TradeBalance, []*entity.Order) {
	var wg sync.WaitGroup
	var agentArr []*SimulateTradeAgent
	var agentLock sync.Mutex
	for _, t := range targetArr {
		stock := t

		wg.Add(1)
		go func() {
			defer wg.Done()
			uc.historyDataLock.RLock()
			tickArr := *uc.historyTick[stock.StockNum]
			beforeLastTradeDayClose := uc.beforeHistoryClose[stock.StockNum]
			uc.historyDataLock.RUnlock()

			openChangeRatio := 100 * (tickArr[0].Close - beforeLastTradeDayClose) / beforeLastTradeDayClose
			if openChangeRatio < uc.tradeSwitch.OpenCloseChangeRatioLow || openChangeRatio > uc.tradeSwitch.OpenCloseChangeRatioHigh {
				return
			}

			simulateAgent := NewSimulateAgent(stock.StockNum)
			simulateAgent.analyzeTickTime = tickArr[0].TickTime
			simulateAgent.tradeSwitch = uc.tradeSwitch

			agentLock.Lock()
			agentArr = append(agentArr, simulateAgent)
			agentLock.Unlock()

			simulateAgent.searchOrder(analyzeCfg, &tickArr, beforeLastTradeDayClose)
		}()
	}
	wg.Wait()

	var allOrders []*entity.Order
	for i := 0; i < len(agentArr); i++ {
		orders := agentArr[i].getAllOrders()
		if len(orders) != 0 {
			allOrders = append(allOrders, orders...)
		}
	}

	if len(allOrders) == 0 {
		return config.Analyze{}, &entity.TradeBalance{}, []*entity.Order{}
	}

	balancer := NewSimulateBalance(uc.quotaCfg, allOrders)
	tmp, orders := balancer.calculateBalance(allOrders)
	return analyzeCfg, tmp, orders
}

// SimulateBalance -.
type SimulateBalance struct {
	quota     *Quota
	allOrders []*entity.Order
}

// NewSimulateBalance -.
func NewSimulateBalance(quotaCfg config.Quota, allOrders []*entity.Order) *SimulateBalance {
	return &SimulateBalance{
		quota:     NewQuota(quotaCfg),
		allOrders: allOrders,
	}
}

func (uc *SimulateBalance) calculateBalance(allOrders []*entity.Order) (*entity.TradeBalance, []*entity.Order) {
	forwardOrder, reverseOrder := uc.splitOrdersByAction(allOrders)
	var forwardBalance, revereBalance, discount, tradeCount int64
	for _, v := range forwardOrder {
		switch v.Action {
		case entity.ActionBuy:
			tradeCount++
			forwardBalance -= uc.quota.GetStockBuyCost(v.Price, v.Quantity)
		case entity.ActionSell:
			forwardBalance += uc.quota.GetStockSellCost(v.Price, v.Quantity)
		}
		discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Quantity)
		// log.Warnf("TradeTime: %s, Stock: %s, Action: %d, Qty: %d, Price: %.2f", v.TradeTime.Format(global.LongTimeLayout), v.StockNum, v.Action, v.Quantity, v.Price)
	}

	for _, v := range reverseOrder {
		switch v.Action {
		case entity.ActionSellFirst:
			tradeCount++
			revereBalance += uc.quota.GetStockSellCost(v.Price, v.Quantity)
		case entity.ActionBuyLater:
			revereBalance -= uc.quota.GetStockBuyCost(v.Price, v.Quantity)
		}
		discount += uc.quota.GetStockTradeFeeDiscount(v.Price, v.Quantity)
		// log.Warnf("TradeTime: %s, Stock: %s, Action: %d, Qty: %d, Price: %.2f", v.TradeTime.Format(global.LongTimeLayout), v.StockNum, v.Action, v.Quantity, v.Price)
	}

	tmp := &entity.TradeBalance{
		TradeDay:        cc.GetBasicInfo().TradeDay,
		TradeCount:      tradeCount,
		Forward:         forwardBalance,
		Reverse:         revereBalance,
		OriginalBalance: forwardBalance + revereBalance,
		Discount:        discount,
		Total:           forwardBalance + revereBalance + discount,
	}

	return tmp, allOrders
}

func (uc *SimulateBalance) splitOrdersByQuota(allOrders []*entity.Order) ([]*entity.Order, []*entity.Order) {
	var forwardOrder, reverseOrder []*entity.Order

	sort.Slice(allOrders, func(i, j int) bool {
		return allOrders[i].TradeTime.Before(allOrders[j].TradeTime)
	})

	for _, v := range allOrders {
		consumeQuota := uc.quota.calculateOriginalOrderCost(v)
		if uc.quota.quota-consumeQuota < 0 {
			break
		}
		uc.quota.quota -= consumeQuota
		switch v.Action {
		case entity.ActionBuy:
			forwardOrder = append(forwardOrder, v)
		case entity.ActionSellFirst:
			reverseOrder = append(reverseOrder, v)
		}
	}

	return forwardOrder, reverseOrder
}

func (uc *SimulateBalance) splitOrdersByAction(allOrders []*entity.Order) ([]*entity.Order, []*entity.Order) {
	forwardOrder, reverseOrder := uc.splitOrdersByQuota(allOrders)
	var tempForwardOrder []*entity.Order
	for _, v := range forwardOrder {
		for _, a := range allOrders {
			if a.Action == entity.ActionSell && a.StockNum == v.StockNum {
				tempForwardOrder = append(tempForwardOrder, a)
			}
		}
	}
	forwardOrder = append(forwardOrder, tempForwardOrder...)

	var tempReverseOrder []*entity.Order
	for _, v := range reverseOrder {
		for _, a := range allOrders {
			if a.Action == entity.ActionBuyLater && a.StockNum == v.StockNum {
				tempReverseOrder = append(tempReverseOrder, a)
			}
		}
	}
	reverseOrder = append(reverseOrder, tempReverseOrder...)
	return forwardOrder, reverseOrder
}

func generateAnalyzeCfg(useDefault bool) []config.Analyze {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	if useDefault {
		return []config.Analyze{{
			CloseChangeRatioLow:  cfg.Analyze.CloseChangeRatioLow,
			CloseChangeRatioHigh: cfg.Analyze.CloseChangeRatioHigh,
			OutInRatio:           cfg.Analyze.OutInRatio,
			InOutRatio:           cfg.Analyze.InOutRatio,
			VolumePRLimit:        cfg.Analyze.VolumePRLimit,
			TickAnalyzePeriod:    cfg.Analyze.TickAnalyzePeriod,
			RSIMinCount:          cfg.Analyze.RSIMinCount,
			RSIHigh:              cfg.Analyze.RSIHigh,
			RSILow:               cfg.Analyze.RSILow,
			MAPeriod:             cfg.Analyze.MAPeriod,
		}}
	}

	base := []config.Analyze{
		{
			OutInRatio:        95,
			InOutRatio:        95,
			VolumePRLimit:     99,
			TickAnalyzePeriod: cfg.Analyze.TickAnalyzePeriod,
			RSIMinCount:       150,
			RSIHigh:           50,
			RSILow:            50,
		},
	}

	AppendRSICountVar(&base)
	AppendRSIHighVar(&base)
	AppendRSILowVar(&base)
	AppendVolumePRLimitVar(&base)
	AppendOutInRatioVar(&base)
	AppendInOutRatioVar(&base)

	log.Warnf("Total analyze times: %d", len(base))

	return base
}

// AppendOutInRatioVar -.
func AppendOutInRatioVar(cfgArr *[]config.Analyze) {
	var appendCfg []config.Analyze
	for _, v := range *cfgArr {
		for {
			if v.OutInRatio <= 75 {
				break
			}
			v.OutInRatio -= 5
			appendCfg = append(appendCfg, v)
		}
	}
	*cfgArr = append(*cfgArr, appendCfg...)
}

// AppendInOutRatioVar -.
func AppendInOutRatioVar(cfgArr *[]config.Analyze) {
	var appendCfg []config.Analyze
	for _, v := range *cfgArr {
		for {
			if v.InOutRatio <= 95 {
				break
			}
			v.InOutRatio -= 5
			appendCfg = append(appendCfg, v)
		}
	}
	*cfgArr = append(*cfgArr, appendCfg...)
}

// AppendVolumePRLimitVar -.
func AppendVolumePRLimitVar(cfgArr *[]config.Analyze) {
	var appendCfg []config.Analyze
	for _, v := range *cfgArr {
		for {
			if v.VolumePRLimit <= 90 {
				break
			}
			v.VolumePRLimit--

			appendCfg = append(appendCfg, v)
		}
	}
	*cfgArr = append(*cfgArr, appendCfg...)
}

// AppendRSICountVar -.
func AppendRSICountVar(cfgArr *[]config.Analyze) {
	var appendCfg []config.Analyze
	for _, v := range *cfgArr {
		for {
			if v.RSIMinCount >= 600 {
				break
			}
			v.RSIMinCount += 150
			appendCfg = append(appendCfg, v)
		}
	}
	*cfgArr = append(*cfgArr, appendCfg...)
}

// AppendRSIHighVar -.
func AppendRSIHighVar(cfgArr *[]config.Analyze) {
	var appendCfg []config.Analyze
	for _, v := range *cfgArr {
		for {
			if v.RSIHigh >= 50.1 {
				break
			}
			v.RSIHigh += 0.1
			if v.RSIHigh >= v.RSILow {
				appendCfg = append(appendCfg, v)
			}
		}
	}
	*cfgArr = append(*cfgArr, appendCfg...)
}

// AppendRSILowVar -.
func AppendRSILowVar(cfgArr *[]config.Analyze) {
	var appendCfg []config.Analyze
	for _, v := range *cfgArr {
		for {
			if v.RSILow <= 49.9 {
				break
			}
			v.RSILow -= 0.1
			if v.RSIHigh >= v.RSILow {
				appendCfg = append(appendCfg, v)
			}
		}
	}
	*cfgArr = append(*cfgArr, appendCfg...)
}
