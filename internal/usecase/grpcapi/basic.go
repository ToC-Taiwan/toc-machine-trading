// Package grpcapi package grpcapi
package grpcapi

import (
	"context"
	"errors"
	"io"

	"tmt/pb"
	"tmt/pkg/grpc"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

// BasicgRPCAPI -.
type BasicgRPCAPI struct {
	conn     *grpc.ConnPool
	clientID string
}

// NewBasic -.
func NewBasic(client *grpc.ConnPool, devMode bool) *BasicgRPCAPI {
	instance := &BasicgRPCAPI{
		conn:     client,
		clientID: uuid.New().String(),
	}

	if devMode {
		instance.clientID = "debug"
	}

	return instance
}

// Heartbeat Heartbeat
func (t *BasicgRPCAPI) Heartbeat() error {
	conn := t.conn.Get()
	defer t.conn.Put(conn)

	stream, err := pb.NewBasicDataInterfaceClient(conn).Heartbeat(context.Background())
	if err != nil {
		return err
	}

	err = stream.Send(&pb.BeatMessage{Message: t.clientID})
	if err != nil {
		return err
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			if !errors.Is(io.EOF, err) {
				return err
			}
			continue
		}

		if response.GetError() != "" {
			return errors.New(response.GetError())
		}
	}
}

// Terminate -.
func (t *BasicgRPCAPI) Terminate() error {
	conn := t.conn.Get()
	defer t.conn.Put(conn)

	_, err := pb.NewBasicDataInterfaceClient(conn).Terminate(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// GetAllStockDetail GetAllStockDetail
func (t *BasicgRPCAPI) GetAllStockDetail() ([]*pb.StockDetailMessage, error) {
	conn := t.conn.Get()
	defer t.conn.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).GetAllStockDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockDetailMessage{}, err
	}
	return r.GetStock(), nil
}

// GetAllFutureDetail -.
func (t *BasicgRPCAPI) GetAllFutureDetail() ([]*pb.FutureDetailMessage, error) {
	conn := t.conn.Get()
	defer t.conn.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).GetAllFutureDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.FutureDetailMessage{}, err
	}
	return r.GetFuture(), nil
}

// GetAllOptionDetail -.
func (t *BasicgRPCAPI) GetAllOptionDetail() ([]*pb.OptionDetailMessage, error) {
	conn := t.conn.Get()
	defer t.conn.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).GetAllOptionDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.OptionDetailMessage{}, err
	}
	return r.GetOption(), nil
}

func (t *BasicgRPCAPI) CheckUsage() (*pb.ShioajiUsage, error) {
	conn := t.conn.Get()
	defer t.conn.Put(conn)

	r, err := pb.NewBasicDataInterfaceClient(conn).CheckUsage(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.ShioajiUsage{}, err
	}
	return r, nil
}
