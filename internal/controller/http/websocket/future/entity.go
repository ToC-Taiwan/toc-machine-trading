package future

import "tmt/internal/entity"

type AutomationType int

const (
	AutomationNone AutomationType = iota
	AutomationByBalance
	AutomationByTimePeriod
	AutomationByTimePeriodAndBalance
)

type futureTradeClientMsg struct {
	Code   string               `json:"code"`
	Action entity.OrderAction   `json:"action"`
	Price  float64              `json:"price"`
	Qty    int64                `json:"qty"`
	Option HalfAutomationOption `json:"option"`
}

type HalfAutomationOption struct {
	AutomationType AutomationType `json:"automation_type"`
	ByBalanceHigh  float64        `json:"by_balance_high"`
	ByBalanceLow   float64        `json:"by_balance_low"`
	ByTimePeriod   int64          `json:"by_time_period"`
}

type periodTradeVolume struct {
	FirstPeriod  entity.OutInVolume `json:"first_period"`
	SecondPeriod entity.OutInVolume `json:"second_period"`
	ThirdPeriod  entity.OutInVolume `json:"third_period"`
	FourthPeriod entity.OutInVolume `json:"fourth_period"`
}

type futurePosition struct {
	Position []*entity.FuturePosition `json:"position"`
}

// type assistStatus struct {
// 	Running bool `json:"running"`
// }
