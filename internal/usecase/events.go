// Package usecase package usecase
package usecase

const (
	topicFetchStockHistory string = "fetch_stock_history"

	topicAnalyzeStockTargets string = "analyze_stock_targets"

	topicFetchFutureHistory string = "fetch_future_history"
)

const (
	topicInsertOrUpdateStockOrder string = "insert_or_update_order"
)

const (
	topicInsertOrUpdateFutureOrder string = "insert_or_update_future_order"

	topicSubscribeFutureTickTargets string = "subscribe_future_targets"

	// topicUnSubscribeFutureTickTargets string = "unsubscribe_future_targets"
)

const (
	topicUpdateAuthTradeUser string = "update_auth_trade_user"

	topicUpdatePushUser string = "update_push_user"
)
