// Package usecase package usecase
package usecase

const (
	TopicFetchStockHistory   string = "fetch_stock_history"
	TopicAnalyzeStockTargets string = "analyze_stock_targets"

	TopicFetchFutureHistory string = "fetch_future_history"
)

const (
	TopicSubscribeStockTickTargets   string = "subscribe_stock_tick_targets"
	TopicUnSubscribeStockTickTargets string = "unsubscribe_stock_tick_targets"
	TopicInsertOrUpdateStockOrder    string = "insert_or_update_order"
)

const (
	TopicSubscribeFutureTickTargets   string = "subscribe_future_targets"
	TopicUnSubscribeFutureTickTargets string = "unsubscribe_future_targets"
	TopicInsertOrUpdateFutureOrder    string = "insert_or_update_future_order"
)
