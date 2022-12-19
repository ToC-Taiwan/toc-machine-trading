package future

import (
	"fmt"
)

type futureTradeError struct {
	ErrCode  int    `json:"err_code"`
	Response string `json:"response"`
}

func (f *futureTradeError) Error() string {
	return fmt.Sprintf("Code: %d, Err: %s", f.ErrCode, f.Response)
}

var (
	errNotTradeTime          = &futureTradeError{-1, "not trade time"}
	errNotFilled             = &futureTradeError{-2, "please wait for previous order to be filled"}
	errAssistNotSupport      = &futureTradeError{-3, "assist only support qty <= 4"}
	errUnmarshal             = &futureTradeError{-4, "unmarshal error"}
	errGetSnapshot           = &futureTradeError{-5, "get snapshot error"}
	errGetPosition           = &futureTradeError{-6, "get position error"}
	errPlaceOrder            = &futureTradeError{-7, "place order error"}
	errCancelOrderFailed     = &futureTradeError{-8, "cancel order failed"}
	errAssitingIsNotFinished = &futureTradeError{-9, "assist is not finished"}
	errGetKbarFail           = &futureTradeError{-10, "get kbar fail"}
)
