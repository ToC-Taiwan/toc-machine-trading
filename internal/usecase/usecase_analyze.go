package usecase

import (
	"context"
	"sort"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/pkg/config"
	"tmt/pkg/global"
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

	tradeDay *TradeDay
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
		tradeDay:           NewTradeDay(),
	}

	bus.SubscribeTopic(topicAnalyzeTargets, uc.AnalyzeAll)
	return uc
}

type simulateResult struct {
	cfg     config.Analyze
	balance *entity.TradeBalance
	orders  []*entity.StockOrder
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
			if nextTradeDay := uc.tradeDay.getAbsNextTradeDayTime(ma.Date); nextTradeDay.Equal(basicInfo.TradeDay) {
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
		var bestBalance *entity.TradeBalance
		var orders *[]*entity.StockOrder
		for {
			res, ok := <-resultChan
			if !ok {
				for _, o := range *orders {
					log.Warnf("TradeTime: %s, Stock: %s, Action: %d, Qty: %d, Price: %.2f", o.TradeTime.Format(global.LongTimeLayout), o.StockNum, o.Action, o.Quantity, o.Price)
				}
				break
			}

			if bestBalance == nil || res.balance.Total > bestBalance.Total {
				bestBalance = res.balance
				bestCfg = res.cfg
				orders = &res.orders
				log.Infof("TradeCount: %d, Forward: %d, Reverse: %d, Discount: %d, Total: %d", bestBalance.TradeCount, bestBalance.Forward, bestBalance.Reverse, bestBalance.Discount, bestBalance.Total)
				log.Warnf("RSIMinCount: %d", bestCfg.RSIMinCount)
				log.Warnf("VolumePRLimit: %.1f", bestCfg.VolumePRLimit)
				log.Warnf("AllOutInRatio %.1f", bestCfg.AllOutInRatio)
				log.Warnf("AllInOutRatio: %.1f", bestCfg.AllInOutRatio)
				log.Warnf("TickAnalyzePeriod: %.0f", bestCfg.TickAnalyzePeriod)
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

func (uc *AnalyzeUseCase) getSimulateCond(targetArr []*entity.Target, analyzeCfg config.Analyze) (config.Analyze, *entity.TradeBalance, []*entity.StockOrder) {
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

	var allOrders []*entity.StockOrder
	for i := 0; i < len(agentArr); i++ {
		orders := agentArr[i].getAllOrders()
		if len(orders) != 0 {
			allOrders = append(allOrders, orders...)
		}
	}

	if len(allOrders) == 0 {
		return config.Analyze{}, &entity.TradeBalance{}, []*entity.StockOrder{}
	}

	balancer := NewSimulateBalance(uc.quotaCfg, allOrders)
	tmp, orders := balancer.calculateBalance(allOrders)
	return analyzeCfg, tmp, orders
}

// SimulateBalance -.
type SimulateBalance struct {
	quota     *Quota
	allOrders []*entity.StockOrder
}

// NewSimulateBalance -.
func NewSimulateBalance(quotaCfg config.Quota, allOrders []*entity.StockOrder) *SimulateBalance {
	return &SimulateBalance{
		quota:     NewQuota(quotaCfg),
		allOrders: allOrders,
	}
}

func (uc *SimulateBalance) calculateBalance(allOrders []*entity.StockOrder) (*entity.TradeBalance, []*entity.StockOrder) {
	sort.Slice(allOrders, func(i, j int) bool {
		return allOrders[i].TradeTime.Before(allOrders[j].TradeTime)
	})

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

	var orders []*entity.StockOrder
	orders = append(orders, forwardOrder...)
	orders = append(orders, reverseOrder...)

	tmp := &entity.TradeBalance{
		TradeDay:        cc.GetBasicInfo().TradeDay,
		TradeCount:      tradeCount,
		Forward:         forwardBalance,
		Reverse:         revereBalance,
		OriginalBalance: forwardBalance + revereBalance,
		Discount:        discount,
		Total:           forwardBalance + revereBalance + discount,
	}

	return tmp, orders
}

func (uc *SimulateBalance) splitOrdersByQuota(allOrders []*entity.StockOrder) ([]*entity.StockOrder, []*entity.StockOrder) {
	var forwardOrder, reverseOrder []*entity.StockOrder
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

func (uc *SimulateBalance) splitOrdersByAction(allOrders []*entity.StockOrder) ([]*entity.StockOrder, []*entity.StockOrder) {
	orderMap := make(map[string][]*entity.StockOrder)
	for _, v := range allOrders {
		orderMap[v.GroupID] = append(orderMap[v.GroupID], v)
	}

	var tempForwardOrder, tempReverseOrder []*entity.StockOrder
	forwardOrder, reverseOrder := uc.splitOrdersByQuota(allOrders)
	for _, v := range forwardOrder {
		tempForwardOrder = append(tempForwardOrder, orderMap[v.GroupID]...)
	}

	for _, v := range reverseOrder {
		tempReverseOrder = append(tempReverseOrder, orderMap[v.GroupID]...)
	}

	return tempForwardOrder, tempReverseOrder
}

func generateAnalyzeCfg(useDefault bool) []config.Analyze {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	if useDefault {
		return []config.Analyze{{
			MaxHoldTime:          cfg.Analyze.MaxHoldTime,
			CloseChangeRatioLow:  cfg.Analyze.CloseChangeRatioLow,
			CloseChangeRatioHigh: cfg.Analyze.CloseChangeRatioHigh,
			VolumePRLimit:        cfg.Analyze.VolumePRLimit,
			TickAnalyzePeriod:    cfg.Analyze.TickAnalyzePeriod,
			RSIMinCount:          cfg.Analyze.RSIMinCount,
			AllOutInRatio:        cfg.Analyze.AllOutInRatio,
			AllInOutRatio:        cfg.Analyze.AllInOutRatio,
		}}
	}

	base := []config.Analyze{
		{
			RSIMinCount:   50,
			VolumePRLimit: 99,
			AllOutInRatio: 50,
			AllInOutRatio: 50,

			MaxHoldTime:          cfg.Analyze.MaxHoldTime,
			TickAnalyzePeriod:    cfg.Analyze.TickAnalyzePeriod,
			CloseChangeRatioLow:  cfg.Analyze.CloseChangeRatioLow,
			CloseChangeRatioHigh: cfg.Analyze.CloseChangeRatioHigh,
		},
	}

	AppendRSICountVar(&base)
	AppendVolumePRLimitVar(&base)
	AppendAllOutInRatioVar(&base)
	AppendAllInOutRatioVar(&base)

	log.Warnf("Total analyze times: %d", len(base))

	return base
}

// AppendRSICountVar -.
func AppendRSICountVar(cfgArr *[]config.Analyze) {
	var appendCfg []config.Analyze
	for _, v := range *cfgArr {
		for {
			if v.RSIMinCount >= 200 {
				break
			}
			v.RSIMinCount += 50
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
			if v.VolumePRLimit <= 95 {
				break
			}
			v.VolumePRLimit--
			appendCfg = append(appendCfg, v)
		}
	}
	*cfgArr = append(*cfgArr, appendCfg...)
}

// AppendAllOutInRatioVar -.
func AppendAllOutInRatioVar(cfgArr *[]config.Analyze) {
	var appendCfg []config.Analyze
	for _, v := range *cfgArr {
		for {
			if v.AllOutInRatio >= 90 {
				break
			}
			v.AllOutInRatio += 10
			appendCfg = append(appendCfg, v)
		}
	}
	*cfgArr = append(*cfgArr, appendCfg...)
}

// AppendAllInOutRatioVar -.
func AppendAllInOutRatioVar(cfgArr *[]config.Analyze) {
	var appendCfg []config.Analyze
	for _, v := range *cfgArr {
		for {
			if v.AllInOutRatio >= 90 {
				break
			}
			v.AllInOutRatio += 10
			appendCfg = append(appendCfg, v)
		}
	}
	*cfgArr = append(*cfgArr, appendCfg...)
}
