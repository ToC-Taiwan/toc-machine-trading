// Package config package config
package config

import (
	"sync"

	"toc-machine-trading/pkg/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

var log = logger.Get()

var (
	globalConfig *Config
	once         sync.Once
)

// GetConfig -.
func GetConfig() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	once.Do(parseConfigFile)
	return globalConfig, nil
}

func parseConfigFile() {
	newConfig := Config{}
	err := cleanenv.ReadConfig("./configs/config.yml", &newConfig)
	if err != nil {
		log.Panic(err)
	}

	err = cleanenv.ReadEnv(&newConfig)
	if err != nil {
		log.Panic(err)
	}

	globalConfig = &newConfig
}

// Config -.
type Config struct {
	HTTP        `env-required:"true" yaml:"http"`
	Postgres    `env-required:"true" yaml:"postgres"`
	Sinopac     `env-required:"true" yaml:"sinopac"`
	RabbitMQ    `env-required:"true" yaml:"rabbitmq"`
	TradeSwitch `env-required:"true" yaml:"trade_switch"`
	History     `env-required:"true" yaml:"history"`
	Quota       `env-required:"true" yaml:"quota"`
	TargetCond  `env-required:"true" yaml:"target_cond"`
	Analyze     `env-required:"true" yaml:"analyze"`

	Deployment string `env-required:"true" env:"DEPLOYMENT"`
}

// HTTP -.
type HTTP struct {
	Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

// Postgres -.
type Postgres struct {
	PoolMax int    `env-required:"true" yaml:"pool_max"`
	URL     string `env-required:"true" env:"PG_URL"`
	DBName  string `env-required:"true" env:"DB_NAME"`
}

// Sinopac -.
type Sinopac struct {
	PoolMax int    `env-required:"true" yaml:"pool_max"`
	URL     string `env-required:"true" env:"SINOPAC_URL"`
}

// RabbitMQ -.
type RabbitMQ struct {
	URL      string `env-required:"true" env:"RABBITMQ_URL"`
	Exchange string `env-required:"true" env:"RABBITMQ_EXCHANGE"`
	WaitTime int64  `env-required:"true" yaml:"wait_time"`
	Attempts int    `env-required:"true" yaml:"attempts"`
}

// TradeSwitch -.
type TradeSwitch struct {
	Simulation bool `env-required:"true" yaml:"simulation"`

	Buy       bool `env-required:"true" yaml:"buy"`
	Sell      bool `env-required:"true" yaml:"sell"`
	SellFirst bool `env-required:"true" yaml:"sell_first"`
	BuyLater  bool `env-required:"true" yaml:"buy_later"`

	HoldTimeFromOpen float64 `env-required:"true" yaml:"hold_time_from_open"`
	TotalOpenTime    float64 `env-required:"true" yaml:"total_open_time"`

	TradeInWaitTime  int64   `env-required:"true" yaml:"trade_in_wait_time"`
	TradeOutWaitTime int64   `env-required:"true" yaml:"trade_out_wait_time"`
	TradeInEndTime   float64 `env-required:"true" yaml:"trade_in_end_time"`
	TradeOutEndTime  float64 `env-required:"true" yaml:"trade_out_end_time"`

	MeanTimeForward int64 `env-required:"true" yaml:"mean_time_forward"`
	MeanTimeReverse int64 `env-required:"true" yaml:"mean_time_reverse"`

	ForwardMax int64 `env-required:"true" yaml:"forward_max"`
	ReverseMax int64 `env-required:"true" yaml:"reverse_max"`
}

// History -.
type History struct {
	HistoryClosePeriod int64 `env-required:"true" yaml:"history_close_period"`
	HistoryTickPeriod  int64 `env-required:"true" yaml:"history_tick_period"`
	HistoryKbarPeriod  int64 `env-required:"true" yaml:"history_kbar_period"`
}

// Quota -.
type Quota struct {
	TradeQuota    int64   `env-required:"true" yaml:"trade_quota"`
	TradeTaxRatio float64 `env-required:"true" yaml:"trade_tax_ratio"`
	TradeFeeRatio float64 `env-required:"true" yaml:"trade_fee_ratio"`
	FeeDiscount   float64 `env-required:"true" yaml:"fee_discount"`
}

// TargetCond -.
type TargetCond struct {
	LimitPriceLow  float64 `env-required:"true" yaml:"limit_price_low"`
	LimitPriceHigh float64 `env-required:"true" yaml:"limit_price_high"`
	LimitVolume    int64   `env-required:"true" yaml:"limit_volume"`
}

// Analyze -.
type Analyze struct {
	CloseChangeRatioLow      float64 `env-required:"true" yaml:"close_change_ratio_low"`
	CloseChangeRatioHigh     float64 `env-required:"true" yaml:"close_change_ratio_high"`
	OpenCloseChangeRatioLow  float64 `env-required:"true" yaml:"open_close_change_ratio_low"`
	OpenCloseChangeRatioHigh float64 `env-required:"true" yaml:"open_close_change_ratio_high"`
	OutInRatio               float64 `env-required:"true" yaml:"out_in_ratio"`
	InOutRatio               float64 `env-required:"true" yaml:"in_out_ratio"`
	VolumePRLow              float64 `env-required:"true" yaml:"volume_pr_low"`
	VolumePRHigh             float64 `env-required:"true" yaml:"volume_pr_high"`
	TickAnalyzeMinPeriod     float64 `env-required:"true" yaml:"tick_analyze_min_period"`
	TickAnalyzeMaxPeriod     float64 `env-required:"true" yaml:"tick_analyze_max_period"`
	RSIMinCount              int     `env-required:"true" yaml:"rsi_min_count"`
	RSIHigh                  float64 `env-required:"true" yaml:"rsi_high"`
	RSILow                   float64 `env-required:"true" yaml:"rsi_low"`
	MaxLoss                  float64 `env-required:"true" yaml:"max_loss"`
	MAPeriod                 int64   `env-required:"true" yaml:"ma_period"`
}
