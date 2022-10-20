// Package config package config
package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	globalConfig *Config
	once         sync.Once
)

// GetConfig -.
func GetConfig() *Config {
	if globalConfig != nil {
		return globalConfig
	}

	once.Do(parseConfigFile)
	return globalConfig
}

func parseConfigFile() {
	newConfig := Config{}
	err := cleanenv.ReadConfig("./configs/config.yml", &newConfig)
	if err != nil {
		panic(err)
	}

	err = cleanenv.ReadEnv(&newConfig)
	if err != nil {
		panic(err)
	}

	globalConfig = &newConfig
}

// Config -.
type Config struct {
	HTTP              HTTP              `json:"http"                env-required:"true" yaml:"http"`
	Postgres          Postgres          `json:"postgres"            env-required:"true" yaml:"postgres"`
	Sinopac           Sinopac           `json:"sinopac"             env-required:"true" yaml:"sinopac"`
	RabbitMQ          RabbitMQ          `json:"rabbitmq"            env-required:"true" yaml:"rabbitmq"`
	Simulation        bool              `json:"simulation"          yaml:"simulation"`
	TradeSwitch       TradeSwitch       `json:"trade_switch"        env-required:"true" yaml:"trade_switch"`
	FutureTradeSwitch FutureTradeSwitch `json:"future_trade_switch" env-required:"true" yaml:"future_trade_switch"`
	History           History           `json:"history"             env-required:"true" yaml:"history"`
	Quota             Quota             `json:"quota"               env-required:"true" yaml:"quota"`
	TargetCond        TargetCond        `json:"target_cond"         env-required:"true" yaml:"target_cond"`
	StockAnalyze      StockAnalyze      `json:"stock_analyze"       env-required:"true" yaml:"stock_analyze"`
	FutureAnalyze     FutureAnalyze     `json:"future_analyze"      env-required:"true" yaml:"future_analyze"`
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
	AllowTrade       bool    `json:"allow_trade"         yaml:"allow_trade"`
	HoldTimeFromOpen float64 `json:"hold_time_from_open" env-required:"true" yaml:"hold_time_from_open"`
	TotalOpenTime    float64 `json:"total_open_time"     env-required:"true" yaml:"total_open_time"`
	TradeInEndTime   float64 `json:"trade_in_end_time"   env-required:"true" yaml:"trade_in_end_time"`
	TradeInWaitTime  int64   `json:"trade_in_wait_time"  env-required:"true" yaml:"trade_in_wait_time"`
	TradeOutWaitTime int64   `json:"trade_out_wait_time" env-required:"true" yaml:"trade_out_wait_time"`
	CancelWaitTime   int64   `json:"cancel_wait_time"    env-required:"true" yaml:"cancel_wait_time"`
}

// FutureTradeSwitch -.
type FutureTradeSwitch struct {
	AllowTrade       bool             `json:"allow_trade"         yaml:"allow_trade"`
	Quantity         int64            `json:"quantity"            env-required:"true" yaml:"quantity"`
	TradeInWaitTime  int64            `json:"trade_in_wait_time"  env-required:"true" yaml:"trade_in_wait_time"`
	TradeOutWaitTime int64            `json:"trade_out_wait_time" env-required:"true" yaml:"trade_out_wait_time"`
	CancelWaitTime   int64            `json:"cancel_wait_time"    env-required:"true" yaml:"cancel_wait_time"`
	TradeTimeRange   []TradeTimeRange `json:"trade_time_range"    env-required:"true" yaml:"trade_time_range"`
}

// TradeTimeRange -.
type TradeTimeRange struct {
	StartTime string `json:"start_time" env-required:"true" yaml:"start_time"`
	Duration  int64  `json:"duration"   env-required:"true" yaml:"duration"`
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

// StockAnalyze -.
type StockAnalyze struct {
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

// FutureAnalyze -.
type FutureAnalyze struct {
	MaxHoldTime       float64 `json:"max_hold_time"       env-required:"true" yaml:"max_hold_time"`
	AllOutInRatio     float64 `json:"all_out_in_ratio"    env-required:"true" yaml:"all_out_in_ratio"`
	AllInOutRatio     float64 `json:"all_in_out_ratio"    env-required:"true" yaml:"all_in_out_ratio"`
	TickAnalyzePeriod float64 `json:"tick_analyze_period" env-required:"true" yaml:"tick_analyze_period"`
	RSIMinCount       int     `json:"rsi_min_count"       env-required:"true" yaml:"rsi_min_count"`
}
