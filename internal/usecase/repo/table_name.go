package repo

import "toc-machine-trading/pkg/logger"

var log = logger.Get()

const (
	tableNameStock    string = "basic_stock"
	tableNameCalendar string = "basic_calendar"
	tableNameEvent    string = "sinopac_event"
	tableNameTarget   string = "basic_targets"
)
