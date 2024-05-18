// Package grpc package grpc
package grpc

import (
	"context"

	"github.com/toc-taiwan/toc-machine-trading/pkg/grpc"

	"github.com/toc-taiwan/toc-trade-protobuf/src/golang/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

// basic -.
type basic struct {
	pool *grpc.ConnPool
}

// NewBasic -.
func NewBasic(client *grpc.ConnPool) BasicgRPCAPI {
	instance := &basic{
		pool: client,
	}
	return instance
}

// CreateLongConnection -.
func (t *basic) CreateLongConnection() error {
	conn := t.pool.Get()
	defer t.pool.Put(conn)
	stream, err := pb.NewBasicDataInterfaceClient(conn).CreateLongConnection(context.Background())
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
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).GetAllStockDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockDetailMessage{}, err
	}
	return r.GetStock(), nil
}

// GetAllFutureDetail -.
func (t *basic) GetAllFutureDetail() ([]*pb.FutureDetailMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).GetAllFutureDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.FutureDetailMessage{}, err
	}
	return r.GetFuture(), nil
}

// GetAllOptionDetail -.
func (t *basic) GetAllOptionDetail() ([]*pb.OptionDetailMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).GetAllOptionDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.OptionDetailMessage{}, err
	}
	return r.GetOption(), nil
}

func (t *basic) CheckUsage() (*pb.ShioajiUsage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).CheckUsage(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.ShioajiUsage{}, err
	}
	return r, nil
}

func (t *basic) Login() error {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	_, err := pb.NewBasicDataInterfaceClient(conn).Login(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}
