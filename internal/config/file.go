package config

// History -.
type History struct {
	HistoryClosePeriod int64 `yaml:"HistoryClosePeriod"`
	HistoryTickPeriod  int64 `yaml:"HistoryTickPeriod"`
	HistoryKbarPeriod  int64 `yaml:"HistoryKbarPeriod"`
}

// Quota -.
type Quota struct {
	StockTradeQuota  int64   `yaml:"StockTradeQuota"`
	StockFeeDiscount float64 `yaml:"StockFeeDiscount"`
	FutureTradeFee   int64   `yaml:"FutureTradeFee"`
}

// TargetStock -.
type TargetStock struct {
	BlackStock    []string     `yaml:"BlackStock"`
	BlackCategory []string     `yaml:"BlackCategory"`
	RealTimeRank  int64        `yaml:"RealTimeRank"`
	LimitVolume   int64        `yaml:"LimitVolume"`
	PriceLimit    []PriceLimit `yaml:"PriceLimit"`
}

// PriceLimit -.
type PriceLimit struct {
	Low  float64 `yaml:"Low"`
	High float64 `yaml:"High"`
}

// AnalyzeStock -.
type AnalyzeStock struct {
	MaxHoldTime          float64 `yaml:"MaxHoldTime"`
	CloseChangeRatioLow  float64 `yaml:"CloseChangeRatioLow"`
	CloseChangeRatioHigh float64 `yaml:"CloseChangeRatioHigh"`
	AllOutInRatio        float64 `yaml:"AllOutInRatio"`
	AllInOutRatio        float64 `yaml:"AllInOutRatio"`
	VolumePRLimit        float64 `yaml:"VolumePRLimit"`
	TickAnalyzePeriod    float64 `yaml:"TickAnalyzePeriod"`
	RSIMinCount          int     `yaml:"RSIMinCount"`
	MAPeriod             int64   `yaml:"MAPeriod"`
}

// TradeStock -.
type TradeStock struct {
	AllowTrade bool `yaml:"AllowTrade"`
	Subscribe  bool `yaml:"Subscribe"`
	Odd        bool `yaml:"Odd"`

	HoldTimeFromOpen float64 `yaml:"HoldTimeFromOpen"`
	TotalOpenTime    float64 `yaml:"TotalOpenTime"`
	TradeInEndTime   float64 `yaml:"TradeInEndTime"`
	TradeInWaitTime  int64   `yaml:"TradeInWaitTime"`
	TradeOutWaitTime int64   `yaml:"TradeOutWaitTime"`
	CancelWaitTime   int64   `yaml:"CancelWaitTime"`
}

// TradeFuture -.
type TradeFuture struct {
	AllowTrade bool `yaml:"AllowTrade"`
	Subscribe  bool `yaml:"Subscribe"`

	BuySellWaitTime int64 `yaml:"BuySellWaitTime"`

	Quantity          int64   `yaml:"Quantity"`
	TargetBalanceHigh float64 `yaml:"TargetBalanceHigh"`
	TargetBalanceLow  float64 `yaml:"TargetBalanceLow"`

	TradeTimeRange TradeTimeRange `yaml:"TradeTimeRange"`
	MaxHoldTime    int64          `yaml:"MaxHoldTime"`

	TickInterval    int64   `yaml:"TickInterval"`
	RateLimit       float64 `yaml:"RateLimit"`
	RateChangeRatio float64 `yaml:"RateChangeRatio"`
	OutInRatio      float64 `yaml:"OutInRatio"`
	InOutRatio      float64 `yaml:"InOutRatio"`

	TradeOutWaitTimes int64 `yaml:"TradeOutWaitTimes"`
}

// TradeTimeRange -.
type TradeTimeRange struct {
	FirstPartDuration  int64 `yaml:"FirstPartDuration"`
	SecondPartDuration int64 `yaml:"SecondPartDuration"`
}
