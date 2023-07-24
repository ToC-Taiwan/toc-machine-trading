// Package grpcapi package grpcapi
package grpcapi

import (
	"context"

	"tmt/pb"
	"tmt/pkg/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

// BasicgRPCAPI -.
type BasicgRPCAPI struct {
	pool *grpc.ConnPool
}

// NewBasic -.
func NewBasic(client *grpc.ConnPool) *BasicgRPCAPI {
	instance := &BasicgRPCAPI{
		pool: client,
	}
	return instance
}

// CreateLongConnection -.
func (t *BasicgRPCAPI) CreateLongConnection() error {
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

// Terminate -.
func (t *BasicgRPCAPI) Terminate() error {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	_, err := pb.NewBasicDataInterfaceClient(conn).Terminate(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// GetAllStockDetail GetAllStockDetail
func (t *BasicgRPCAPI) GetAllStockDetail() ([]*pb.StockDetailMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).GetAllStockDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockDetailMessage{}, err
	}
	return r.GetStock(), nil
}

// GetAllFutureDetail -.
func (t *BasicgRPCAPI) GetAllFutureDetail() ([]*pb.FutureDetailMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).GetAllFutureDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.FutureDetailMessage{}, err
	}
	return r.GetFuture(), nil
}

// GetAllOptionDetail -.
func (t *BasicgRPCAPI) GetAllOptionDetail() ([]*pb.OptionDetailMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).GetAllOptionDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.OptionDetailMessage{}, err
	}
	return r.GetOption(), nil
}

func (t *BasicgRPCAPI) CheckUsage() (*pb.ShioajiUsage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).CheckUsage(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.ShioajiUsage{}, err
	}
	return r, nil
}

func (t *BasicgRPCAPI) Login() error {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	_, err := pb.NewBasicDataInterfaceClient(conn).Login(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}
