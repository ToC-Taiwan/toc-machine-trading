package config

// History -.
type History struct {
	HistoryClosePeriod int64 `json:"HistoryClosePeriod" yaml:"HistoryClosePeriod"`
	HistoryTickPeriod  int64 `json:"HistoryTickPeriod" yaml:"HistoryTickPeriod"`
	HistoryKbarPeriod  int64 `json:"HistoryKbarPeriod" yaml:"HistoryKbarPeriod"`
}

// Quota -.
type Quota struct {
	StockTradeQuota  int64   `json:"StockTradeQuota" yaml:"StockTradeQuota"`
	StockFeeDiscount float64 `json:"StockFeeDiscount" yaml:"StockFeeDiscount"`
	FutureTradeFee   int64   `json:"FutureTradeFee" yaml:"FutureTradeFee"`
}

// TargetStock -.
type TargetStock struct {
	BlackStock    []string     `json:"BlackStock" yaml:"BlackStock"`
	BlackCategory []string     `json:"BlackCategory" yaml:"BlackCategory"`
	RealTimeRank  int64        `json:"RealTimeRank" yaml:"RealTimeRank"`
	LimitVolume   int64        `json:"LimitVolume" yaml:"LimitVolume"`
	PriceLimit    []PriceLimit `json:"PriceLimit" yaml:"PriceLimit"`
}

// PriceLimit -.
type PriceLimit struct {
	Low  float64 `json:"Low" yaml:"Low"`
	High float64 `json:"High" yaml:"High"`
}

// AnalyzeStock -.
type AnalyzeStock struct {
	MaxHoldTime          float64 `json:"MaxHoldTime" yaml:"MaxHoldTime"`
	CloseChangeRatioLow  float64 `json:"CloseChangeRatioLow" yaml:"CloseChangeRatioLow"`
	CloseChangeRatioHigh float64 `json:"CloseChangeRatioHigh" yaml:"CloseChangeRatioHigh"`
	AllOutInRatio        float64 `json:"AllOutInRatio" yaml:"AllOutInRatio"`
	AllInOutRatio        float64 `json:"AllInOutRatio" yaml:"AllInOutRatio"`
	VolumePRLimit        float64 `json:"VolumePRLimit" yaml:"VolumePRLimit"`
	TickAnalyzePeriod    float64 `json:"TickAnalyzePeriod" yaml:"TickAnalyzePeriod"`
	RSIMinCount          int     `json:"RSIMinCount" yaml:"RSIMinCount"`
	MAPeriod             int64   `json:"MAPeriod" yaml:"MAPeriod"`
}

// TradeStock -.
type TradeStock struct {
	AllowTrade bool `json:"AllowTrade" yaml:"AllowTrade"`
	Odd        bool `json:"Odd" yaml:"Odd"`

	HoldTimeFromOpen float64 `json:"HoldTimeFromOpen" yaml:"HoldTimeFromOpen"`
	TotalOpenTime    float64 `json:"TotalOpenTime" yaml:"TotalOpenTime"`
	TradeInEndTime   float64 `json:"TradeInEndTime" yaml:"TradeInEndTime"`
	TradeInWaitTime  int64   `json:"TradeInWaitTime" yaml:"TradeInWaitTime"`
	TradeOutWaitTime int64   `json:"TradeOutWaitTime" yaml:"TradeOutWaitTime"`
	CancelWaitTime   int64   `json:"CancelWaitTime" yaml:"CancelWaitTime"`
}

// TradeFuture -.
type TradeFuture struct {
	AllowTrade bool `json:"AllowTrade" yaml:"AllowTrade"`

	BuySellWaitTime int64 `json:"BuySellWaitTime" yaml:"BuySellWaitTime"`

	Quantity          int64   `json:"Quantity" yaml:"Quantity"`
	TargetBalanceHigh float64 `json:"TargetBalanceHigh" yaml:"TargetBalanceHigh"`
	TargetBalanceLow  float64 `json:"TargetBalanceLow" yaml:"TargetBalanceLow"`

	TradeTimeRange TradeTimeRange `json:"TradeTimeRange" yaml:"TradeTimeRange"`
	MaxHoldTime    int64          `json:"MaxHoldTime" yaml:"MaxHoldTime"`

	TickInterval    int64   `json:"TickInterval" yaml:"TickInterval"`
	RateLimit       float64 `json:"RateLimit" yaml:"RateLimit"`
	RateChangeRatio float64 `json:"RateChangeRatio" yaml:"RateChangeRatio"`
	OutInRatio      float64 `json:"OutInRatio" yaml:"OutInRatio"`
	InOutRatio      float64 `json:"InOutRatio" yaml:"InOutRatio"`

	TradeOutWaitTimes int64 `json:"TradeOutWaitTimes" yaml:"TradeOutWaitTimes"`
}

// TradeTimeRange -.
type TradeTimeRange struct {
	FirstPartDuration  int64 `json:"FirstPartDuration" yaml:"FirstPartDuration"`
	SecondPartDuration int64 `json:"SecondPartDuration" yaml:"SecondPartDuration"`
}
