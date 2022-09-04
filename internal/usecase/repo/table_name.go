package repo

import "tmt/pkg/logger"

var (
	log           = logger.Get()
	batchSize int = 2000
)

const (
	tableNameCalendar         string = "basic_calendar"
	tableNameStock            string = "basic_stock"
	tableNameFuture           string = "basic_future"
	tableNameTarget           string = "basic_targets"
	tableNameHistoryAnalyze   string = "history_analyze"
	tableNameHistoryClose     string = "history_close"
	tableNameHistoryKbar      string = "history_kbar"
	tableNameHistoryTick      string = "history_tick"
	tableNameTradeOrder       string = "trade_order"
	tableNameTradeFutureOrder string = "trade_future_order"
	tableNameEvent            string = "sinopac_event"
	tableNameTradeBalance     string = "trade_balance"
)
