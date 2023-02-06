package config

// TradeTimeRange -.
type TradeTimeRange struct {
	FirstPartDuration  int64 `json:"first_part_duration"  yaml:"first_part_duration"`
	SecondPartDuration int64 `json:"second_part_duration" yaml:"second_part_duration"`
}

// History -.
type History struct {
	HistoryClosePeriod int64 `json:"history_close_period" env-required:"true" yaml:"history_close_period"`
	HistoryTickPeriod  int64 `json:"history_tick_period"  env-required:"true" yaml:"history_tick_period"`
	HistoryKbarPeriod  int64 `json:"history_kbar_period"  env-required:"true" yaml:"history_kbar_period"`
}

// Quota -.
type Quota struct {
	StockTradeQuota  int64   `json:"stock_trade_quota"  env-required:"true" yaml:"stock_trade_quota"`
	StockFeeDiscount float64 `json:"stock_fee_discount" env-required:"true" yaml:"stock_fee_discount"`
	FutureTradeFee   int64   `json:"future_trade_fee"   env-required:"true" yaml:"future_trade_fee"`
}

// TargetStock -.
type TargetStock struct {
	BlackStock    []string     `json:"black_stock"    env-required:"true" yaml:"black_stock"`
	BlackCategory []string     `json:"black_category" env-required:"true" yaml:"black_category"`
	RealTimeRank  int64        `json:"real_time_rank" env-required:"true" yaml:"real_time_rank"`
	LimitVolume   int64        `json:"limit_volume"   env-required:"true" yaml:"limit_volume"`
	PriceLimit    []PriceLimit `json:"price_limit"    env-required:"true" yaml:"price_limit"`
}

// PriceLimit -.
type PriceLimit struct {
	Low  float64 `json:"low"  env-required:"true" yaml:"low"`
	High float64 `json:"high" env-required:"true" yaml:"high"`
}

// AnalyzeStock -.
type AnalyzeStock struct {
	MaxHoldTime          float64 `json:"max_hold_time"           env-required:"true" yaml:"max_hold_time"`
	CloseChangeRatioLow  float64 `json:"close_change_ratio_low"  env-required:"true" yaml:"close_change_ratio_low"`
	CloseChangeRatioHigh float64 `json:"close_change_ratio_high" env-required:"true" yaml:"close_change_ratio_high"`
	AllOutInRatio        float64 `json:"all_out_in_ratio"        env-required:"true" yaml:"all_out_in_ratio"`
	AllInOutRatio        float64 `json:"all_in_out_ratio"        env-required:"true" yaml:"all_in_out_ratio"`
	VolumePRLimit        float64 `json:"volume_pr_limit"         env-required:"true" yaml:"volume_pr_limit"`
	TickAnalyzePeriod    float64 `json:"tick_analyze_period"     env-required:"true" yaml:"tick_analyze_period"`
	RSIMinCount          int     `json:"rsi_min_count"           env-required:"true" yaml:"rsi_min_count"`
	MAPeriod             int64   `json:"ma_period"               env-required:"true" yaml:"ma_period"`
}

// TradeStock -.
type TradeStock struct {
	AllowTrade bool `json:"allow_trade"         yaml:"allow_trade"`
	Subscribe  bool `json:"subscribe"           yaml:"subscribe"`

	HoldTimeFromOpen float64 `json:"hold_time_from_open" env-required:"true" yaml:"hold_time_from_open"`
	TotalOpenTime    float64 `json:"total_open_time"     env-required:"true" yaml:"total_open_time"`
	TradeInEndTime   float64 `json:"trade_in_end_time"   env-required:"true" yaml:"trade_in_end_time"`
	TradeInWaitTime  int64   `json:"trade_in_wait_time"  env-required:"true" yaml:"trade_in_wait_time"`
	TradeOutWaitTime int64   `json:"trade_out_wait_time" env-required:"true" yaml:"trade_out_wait_time"`
	CancelWaitTime   int64   `json:"cancel_wait_time"    env-required:"true" yaml:"cancel_wait_time"`
}

// TradeFuture -.
type TradeFuture struct {
	AllowTrade bool `json:"allow_trade"         yaml:"allow_trade"`
	Subscribe  bool `json:"subscribe"           yaml:"subscribe"`

	TradeInWaitTime  int64 `json:"trade_in_wait_time"  env-required:"true" yaml:"trade_in_wait_time"`
	TradeOutWaitTime int64 `json:"trade_out_wait_time" env-required:"true" yaml:"trade_out_wait_time"`
	CancelWaitTime   int64 `json:"cancel_wait_time"    env-required:"true" yaml:"cancel_wait_time"`

	Quantity          int64   `json:"quantity"            env-required:"true" yaml:"quantity"`
	TargetBalanceHigh float64 `json:"target_balance_high" env-required:"true" yaml:"target_balance_high"`
	TargetBalanceLow  float64 `json:"target_balance_low"  env-required:"true" yaml:"target_balance_low"`

	TradeTimeRange TradeTimeRange `json:"trade_time_range"    env-required:"true" yaml:"trade_time_range"`
}
