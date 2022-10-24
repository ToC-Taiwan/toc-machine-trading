// Package event package event
package event

const (
	// TopicNewTargets -.
	TopicNewTargets string = "new_targets"
)

const (
	// TopicFetchStockHistory -.
	TopicFetchStockHistory string = "fetch_stock_history"
	// TopicAnalyzeStockTargets -.
	TopicAnalyzeStockTargets string = "analyze_stock_targets"
	// TopicStreamStockTargets -.
	TopicStreamStockTargets string = "stream_stock_targets"
	// TopicSubscribeStockTickTargets -.
	TopicSubscribeStockTickTargets string = "subscribe_stock_tick_targets"
	// TopicUnSubscribeStockTickTargets -.
	TopicUnSubscribeStockTickTargets string = "unsubscribe_stock_tick_targets"
	// TopicPlaceStockOrder -.
	TopicPlaceStockOrder string = "place_stock_order"
	// TopicCancelStockOrder -.
	TopicCancelStockOrder string = "cancel_stock_order"
	// TopicInsertOrUpdateStockOrder -.
	TopicInsertOrUpdateStockOrder string = "insert_or_update_order"
	// TopicUpdateStockTradeSwitch -.
	TopicUpdateStockTradeSwitch string = "update_stock_trade_switch"
)

const (
	// TopicStreamFutureTargets -.
	TopicStreamFutureTargets string = "stream_future_targets"
	// TopicSubscribeFutureTickTargets -.
	TopicSubscribeFutureTickTargets string = "subscribe_future_targets"
	// TopicPlaceFutureOrder -.
	TopicPlaceFutureOrder string = "place_future_order"
	// TopicCancelFutureOrder -.
	TopicCancelFutureOrder string = "cancel_future_order"
	// TopicInsertOrUpdateFutureOrder -.
	TopicInsertOrUpdateFutureOrder string = "insert_or_update_future_order"
	// TopicUpdateFutureTradeSwitch -.
	TopicUpdateFutureTradeSwitch string = "update_future_trade_switch"
	// TopicQueryMonitorFutureCode -.
	TopicQueryMonitorFutureCode string = "query_monitor_future_code"
	// TopicMonitorFutureCode -.
	TopicMonitorFutureCode string = "monitor_future_code"
)
