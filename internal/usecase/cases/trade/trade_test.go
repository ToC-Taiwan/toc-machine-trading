// Package trade package trade
package trade_test

import (
	"context"
	"testing"

	"tmt/internal/entity"
	"tmt/internal/usecase/cases/trade"
	"tmt/pb"

	"github.com/golang/mock/gomock"
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

func TestTradeUseCase_BuyFuture(t *testing.T) {
	type args struct {
		order *entity.FutureOrder
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   entity.OrderStatus
		wantErr bool
	}{
		{
			name: "empty code",
			args: args{
				order: &entity.FutureOrder{},
			},
			want:    "",
			want1:   entity.StatusUnknow,
			wantErr: true,
		},
		{
			name: "filled order",
			args: args{
				order: &entity.FutureOrder{
					BaseOrder: entity.BaseOrder{},
					Code:      "TEST",
					Future:    &entity.Future{},
				},
			},
			want:    "order_id",
			want1:   entity.StatusFilled,
			wantErr: false,
		},
	}

	controller := gomock.NewController(t)
	scMock := NewMockTradegRPCAPI(controller)
	scMock.EXPECT().BuyFuture(gomock.Any()).Return(&pb.TradeResult{
		OrderId: "order_id",
		Status:  "Filled",
		Error:   "",
	}, nil)
	uc := &trade.TradeUseCase{
		Sinopac: scMock,
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			id, status, err := uc.BuyFuture(tt.args.order)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
				So(id, ShouldEqual, tt.want)
				So(status, ShouldEqual, tt.want1)
			}
		})
	}
}
