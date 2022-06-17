// Package grpcapi package grpcapi
package grpcapi

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"toc-machine-trading/pb"
	"toc-machine-trading/pkg/sinopac"

	"google.golang.org/protobuf/types/known/emptypb"
)

// BasicgRPCAPI -.
type BasicgRPCAPI struct {
	conn         *sinopac.Connection
	heartbeatMsg string
}

// NewBasic -.
func NewBasic(client *sinopac.Connection) *BasicgRPCAPI {
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
	c := pb.NewHealthCheckClient(conn)
	stream, err := c.Heartbeat(context.Background())
	if err != nil {
		return err
	}

	err = stream.Send(&pb.Beat{Message: t.heartbeatMsg})
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
		err = stream.Send(&pb.Beat{Message: response.GetMessage()})
		if err != nil {
			return err
		}
	}
}

// Terminate -.
func (t *BasicgRPCAPI) Terminate() error {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewHealthCheckClient(conn)
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
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetAllStockDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockDetailMessage{}, err
	}
	return r.GetStock(), nil
}

// GetAllStockSnapshot GetAllStockSnapshot
func (t *BasicgRPCAPI) GetAllStockSnapshot() ([]*pb.StockSnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetAllStockSnapshot(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockSnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotByNumArr GetStockSnapshotByNumArr
func (t *BasicgRPCAPI) GetStockSnapshotByNumArr(stockNumArr []string) ([]*pb.StockSnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockSnapshotByNumArr(context.Background(), &pb.StockNumArr{
		StockNumArr: stockNumArr,
	})
	if err != nil {
		return []*pb.StockSnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotTSE GetStockSnapshotTSE
func (t *BasicgRPCAPI) GetStockSnapshotTSE() ([]*pb.StockSnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockSnapshotTSE(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockSnapshotMessage{}, err
	}
	return r.GetData(), nil
}
