package repo

import "toc-machine-trading/pkg/logger"

var log = logger.Get()

const (
	tableNameCalendar     string = "basic_calendar"
	tableNameStock        string = "basic_stock"
	tableNameTarget       string = "basic_targets"
	tableNameHistoryClose string = "history_close"
	tableNameHistoryKbar  string = "history_kbar"
	tableNameHistoryTick  string = "history_tick"
	tableNameOrderStatus  string = "order_status"
	tableNameEvent        string = "sinopac_event"
	tableNameTradeBalance string = "trade_balance"
)

var batchSize int = 2000
