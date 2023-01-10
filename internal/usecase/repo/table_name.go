package repo

import "tmt/pkg/log"

var (
	logger        = log.Get()
	batchSize int = 2000
)

const (
	tableNameCalendar string = "basic_calendar"
	tableNameStock    string = "basic_stock"
	tableNameFuture   string = "basic_future"
	tableNameTarget   string = "basic_targets"

	tableNameHistoryStockAnalyze string = "history_stock_analyze"
	tableNameHistoryStockClose   string = "history_stock_close"
	tableNameHistoryStockKbar    string = "history_stock_kbar"
	tableNameHistoryStockTick    string = "history_stock_tick"

	tableNameHistoryFutureClose string = "history_future_close"
	tableNameHistoryFutureTick  string = "history_future_tick"

	tableNameTradeStockOrder    string = "trade_stock_order"
	tableNameTradeStockBalance  string = "trade_stock_balance"
	tableNameTradeFutureOrder   string = "trade_future_order"
	tableNameFutureTradeBalance string = "trade_future_balance"

	tableNameEvent string = "sinopac_event"
)
