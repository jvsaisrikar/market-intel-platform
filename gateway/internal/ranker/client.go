package ranker

import (
	"context"
	"time"

	rankerv1 "github.com/saisrikar/market-intel-platform/proto/gen/ranker/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	GetTopStories(ctx context.Context, ticker string, limit int32) ([]*rankerv1.Story, error)
	Close() error
}

type grpcClient struct {
	conn   *grpc.ClientConn
	client rankerv1.RankingServiceClient
}

func NewClient(addr string) (Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &grpcClient{
		conn:   conn,
		client: rankerv1.NewRankingServiceClient(conn),
	}, nil
}

func (c *grpcClient) GetTopStories(ctx context.Context, ticker string, limit int32) ([]*rankerv1.Story, error) {
	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.GetTopStories(callCtx, &rankerv1.TopStoriesRequest{
		Ticker: ticker,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetStories(), nil
}

func (c *grpcClient) Close() error {
	return c.conn.Close()
}
