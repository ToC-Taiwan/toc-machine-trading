// Package trade package trade
package trade_test

import (
	"context"
	"testing"

	"tmt/internal/usecase/cases/trade"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTradeUseCase_UpdateTradeBalanceByTradeDay(t *testing.T) {
	type args struct {
		ctx  context.Context
		date string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty date",
			args: args{
				ctx:  nil,
				date: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			uc := &trade.TradeUseCase{}
			err := uc.UpdateTradeBalanceByTradeDay(tt.args.ctx, tt.args.date)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
			}
		})
	}
}
