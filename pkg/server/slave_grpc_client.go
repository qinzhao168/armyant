package server

import (
	"log"
	"time"

	"github.com/yiqinguo/armyant/pkg/models"
	//pb "github.com/yiqinguo/armyant/pkg/server"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type grpcClient struct {
	client ApiserverClient
}

func NewGrpcClient(addr string) (*grpcClient, error) {

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := NewApiserverClient(conn)

	return &grpcClient{client}, nil
}

func (g *grpcClient) ReportStatus(ctx context.Context, status *models.Status) error {
	for {
		select {
		case <-time.After(time.Second):
			_, err := g.client.ReportStatus(ctx, status)
			log.Printf("report status error: %v", err)
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}

func (g *grpcClient) ReportResult(ctx context.Context, stats *models.Stats) error {
	_, err := g.client.ReportResult(ctx, stats)
	return err
}
