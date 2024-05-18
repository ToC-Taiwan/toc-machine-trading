// Package grpc package grpc
package grpc

import (
	"context"

	"github.com/toc-taiwan/toc-trade-protobuf/src/golang/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// basic -.
type basic struct {
	conn *grpc.ClientConn
}

// NewBasic -.
func NewBasic(client *grpc.ClientConn) BasicgRPCAPI {
	instance := &basic{
		conn: client,
	}
	return instance
}

// CreateLongConnection -.
func (t *basic) CreateLongConnection() error {
	stream, err := pb.NewBasicDataInterfaceClient(t.conn).CreateLongConnection(context.Background())
	if err != nil {
		panic(err)
	}
	data := &emptypb.Empty{}
	for {
		err := stream.RecvMsg(data)
		if err != nil {
			return err
		}
	}
}

// GetAllStockDetail GetAllStockDetail
func (t *basic) GetAllStockDetail() ([]*pb.StockDetailMessage, error) {
	r, err := pb.NewBasicDataInterfaceClient(t.conn).GetAllStockDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockDetailMessage{}, err
	}
	return r.GetStock(), nil
}

// GetAllFutureDetail -.
func (t *basic) GetAllFutureDetail() ([]*pb.FutureDetailMessage, error) {
	r, err := pb.NewBasicDataInterfaceClient(t.conn).GetAllFutureDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.FutureDetailMessage{}, err
	}
	return r.GetFuture(), nil
}

// GetAllOptionDetail -.
func (t *basic) GetAllOptionDetail() ([]*pb.OptionDetailMessage, error) {
	r, err := pb.NewBasicDataInterfaceClient(t.conn).GetAllOptionDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.OptionDetailMessage{}, err
	}
	return r.GetOption(), nil
}

func (t *basic) CheckUsage() (*pb.ShioajiUsage, error) {
	r, err := pb.NewBasicDataInterfaceClient(t.conn).CheckUsage(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.ShioajiUsage{}, err
	}
	return r, nil
}

func (t *basic) Login() error {
	_, err := pb.NewBasicDataInterfaceClient(t.conn).Login(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}
