// Package grpcapi package grpcapi
package grpcapi

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"tmt/pb"
	"tmt/pkg/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

// BasicgRPCAPI -.
type BasicgRPCAPI struct {
	conn         *grpc.Connection
	heartbeatMsg string
}

// NewBasic -.
func NewBasic(client *grpc.Connection) *BasicgRPCAPI {
	instance := &BasicgRPCAPI{
		conn:         client,
		heartbeatMsg: "beat",
	}

	mode, ok := os.LookupEnv("DEPLOYMENT")
	if !ok || mode != "prod" {
		instance.heartbeatMsg = "debug"
	}
	return instance
}

// Heartbeat Heartbeat
func (t *BasicgRPCAPI) Heartbeat() error {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewBasicDataInterfaceClient(conn)
	stream, err := c.Heartbeat(context.Background())
	if err != nil {
		return err
	}

	err = stream.Send(&pb.BeatMessage{Message: t.heartbeatMsg})
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
		time.Sleep(3 * time.Second)
		err = stream.Send(&pb.BeatMessage{Message: response.GetMessage()})
		if err != nil {
			return err
		}
	}
}

// Terminate -.
func (t *BasicgRPCAPI) Terminate() error {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewBasicDataInterfaceClient(conn)
	_, err := c.Terminate(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// GetAllStockDetail GetAllStockDetail
func (t *BasicgRPCAPI) GetAllStockDetail() ([]*pb.StockDetailMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewBasicDataInterfaceClient(conn)
	r, err := c.GetAllStockDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockDetailMessage{}, err
	}
	return r.GetStock(), nil
}

// GetAllFutureDetail -.
func (t *BasicgRPCAPI) GetAllFutureDetail() ([]*pb.FutureDetailMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewBasicDataInterfaceClient(conn)
	r, err := c.GetAllFutureDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.FutureDetailMessage{}, err
	}
	return r.GetFuture(), nil
}
