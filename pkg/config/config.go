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
	HTTP        HTTP         `json:"http"         env-required:"true" yaml:"http"`
	Postgres    Postgres     `json:"postgres"     env-required:"true" yaml:"postgres"`
	Sinopac     Sinopac      `json:"sinopac"      env-required:"true" yaml:"sinopac"`
	RabbitMQ    RabbitMQ     `json:"rabbitmq"     env-required:"true" yaml:"rabbitmq"`
	TradeSwitch TradeSwitch  `json:"trade_switch" env-required:"true" yaml:"trade_switch"`
	History     History      `json:"history"      env-required:"true" yaml:"history"`
	Quota       Quota        `json:"quota"        env-required:"true" yaml:"quota"`
	TargetCond  []TargetCond `json:"target_cond"  env-required:"true" yaml:"target_cond"`
	Analyze     Analyze      `json:"analyze"      env-required:"true" yaml:"analyze"`

	Deployment string `json:"deployment" env-required:"true" env:"DEPLOYMENT"`
}

// HTTP -.
type HTTP struct {
	Port string `json:"port" env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

// Postgres -.
type Postgres struct {
	PoolMax int    `json:"pool_max" env-required:"true" yaml:"pool_max"`
	URL     string `json:"url"      env-required:"true" env:"PG_URL"`
	DBName  string `json:"db_name"  env-required:"true" env:"DB_NAME"`
}

// Sinopac -.
type Sinopac struct {
	PoolMax int    `json:"pool_max" env-required:"true" yaml:"pool_max"`
	URL     string `json:"url"      env-required:"true" env:"SINOPAC_URL"`
}

// RabbitMQ -.
type RabbitMQ struct {
	URL      string `json:"url"       env-required:"true" env:"RABBITMQ_URL"`
	Exchange string `json:"exchange"  env-required:"true" env:"RABBITMQ_EXCHANGE"`
	WaitTime int64  `json:"wait_time" env-required:"true" yaml:"wait_time"`
	Attempts int    `json:"attempts"  env-required:"true" yaml:"attempts"`
}

// TradeSwitch -.
type TradeSwitch struct {
	Simulation bool `json:"simulation" env-required:"true" yaml:"simulation"`

	Buy       bool `json:"buy"        env-required:"true" yaml:"buy"`
	Sell      bool `json:"sell"       env-required:"true" yaml:"sell"`
	SellFirst bool `json:"sell_first" env-required:"true" yaml:"sell_first"`
	BuyLater  bool `json:"buy_later"  env-required:"true" yaml:"buy_later"`

	HoldTimeFromOpen float64 `json:"hold_time_from_open" env-required:"true" yaml:"hold_time_from_open"`
	TotalOpenTime    float64 `json:"total_open_time"     env-required:"true" yaml:"total_open_time"`

	TradeInWaitTime  int64   `json:"trade_in_wait_time"  env-required:"true" yaml:"trade_in_wait_time"`
	TradeOutWaitTime int64   `json:"trade_out_wait_time" env-required:"true" yaml:"trade_out_wait_time"`
	TradeInEndTime   float64 `json:"trade_in_end_time"   env-required:"true" yaml:"trade_in_end_time"`
	TradeOutEndTime  float64 `json:"trade_out_end_time"  env-required:"true" yaml:"trade_out_end_time"`

	MeanTimeForward int64 `json:"mean_time_forward" env-required:"true" yaml:"mean_time_forward"`
	MeanTimeReverse int64 `json:"mean_time_reverse" env-required:"true" yaml:"mean_time_reverse"`

	ForwardMax int64 `json:"forward_max" env-required:"true" yaml:"forward_max"`
	ReverseMax int64 `json:"reverse_max" env-required:"true" yaml:"reverse_max"`
}

// History -.
type History struct {
	HistoryClosePeriod int64 `json:"history_close_period" env-required:"true" yaml:"history_close_period"`
	HistoryTickPeriod  int64 `json:"history_tick_period"  env-required:"true" yaml:"history_tick_period"`
	HistoryKbarPeriod  int64 `json:"history_kbar_period"  env-required:"true" yaml:"history_kbar_period"`
}

// Quota -.
type Quota struct {
	TradeQuota    int64   `json:"trade_quota"     env-required:"true" yaml:"trade_quota"`
	TradeTaxRatio float64 `json:"trade_tax_ratio" env-required:"true" yaml:"trade_tax_ratio"`
	TradeFeeRatio float64 `json:"trade_fee_ratio" env-required:"true" yaml:"trade_fee_ratio"`
	FeeDiscount   float64 `json:"fee_discount"    env-required:"true" yaml:"fee_discount"`
}

// TargetCond -.
type TargetCond struct {
	LimitPriceLow  float64 `json:"limit_price_low"  env-required:"true" yaml:"limit_price_low"`
	LimitPriceHigh float64 `json:"limit_price_high" env-required:"true" yaml:"limit_price_high"`
	LimitVolume    int64   `json:"limit_volume"     env-required:"true" yaml:"limit_volume"`
	Subscribe      bool    `json:"subscribe"        env-required:"true" yaml:"subscribe"`
}

// Analyze -.
type Analyze struct {
	CloseChangeRatioLow      float64 `json:"close_change_ratio_low"       env-required:"true" yaml:"close_change_ratio_low"`
	CloseChangeRatioHigh     float64 `json:"close_change_ratio_high"      env-required:"true" yaml:"close_change_ratio_high"`
	OpenCloseChangeRatioLow  float64 `json:"open_close_change_ratio_low"  env-required:"true" yaml:"open_close_change_ratio_low"`
	OpenCloseChangeRatioHigh float64 `json:"open_close_change_ratio_high" env-required:"true" yaml:"open_close_change_ratio_high"`
	OutInRatio               float64 `json:"out_in_ratio"                 env-required:"true" yaml:"out_in_ratio"`
	InOutRatio               float64 `json:"in_out_ratio"                 env-required:"true" yaml:"in_out_ratio"`
	VolumePRLow              float64 `json:"volume_pr_low"                env-required:"true" yaml:"volume_pr_low"`
	VolumePRHigh             float64 `json:"volume_pr_high"               env-required:"true" yaml:"volume_pr_high"`
	TickAnalyzePeriod        float64 `json:"tick_analyze_period"          env-required:"true" yaml:"tick_analyze_period"`
	RSIMinCount              int     `json:"rsi_min_count"                env-required:"true" yaml:"rsi_min_count"`
	RSIHigh                  float64 `json:"rsi_high"                     env-required:"true" yaml:"rsi_high"`
	RSILow                   float64 `json:"rsi_low"                      env-required:"true" yaml:"rsi_low"`
	MaxLoss                  float64 `json:"max_loss"                     env-required:"true" yaml:"max_loss"`
	MAPeriod                 int64   `json:"ma_period"                    env-required:"true" yaml:"ma_period"`
}
