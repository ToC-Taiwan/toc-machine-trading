// Package topic package topic
package topic

const (
	TopicFetchStockHistory   string = "fetch_stock_history"
	TopicAnalyzeStockTargets string = "analyze_stock_targets"

	TopicSubscribeStockTickTargets   string = "subscribe_stock_tick_targets"
	TopicUnSubscribeStockTickTargets string = "unsubscribe_stock_tick_targets"

	TopicPlaceStockOrder          string = "place_stock_order"
	TopicCancelStockOrder         string = "cancel_stock_order"
	TopicInsertOrUpdateStockOrder string = "insert_or_update_order"
	TopicUpdateStockTradeSwitch   string = "update_stock_trade_switch"
)

const (
	TopicSubscribeFutureTickTargets string = "subscribe_future_targets"

	TopicPlaceFutureOrder          string = "place_future_order"
	TopicCancelFutureOrder         string = "cancel_future_order"
	TopicInsertOrUpdateFutureOrder string = "insert_or_update_future_order"
	TopicUpdateFutureTradeSwitch   string = "update_future_trade_switch"
)
