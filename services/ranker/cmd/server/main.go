package main

import (
	"context"
	"log"
	"net"
	"strings"

	rankerv1 "github.com/saisrikar/market-intel-platform/proto/gen/ranker/v1"
	"google.golang.org/grpc"
)

const (
	defaultLimit = 3
	listenAddr   = ":50052"
)

type rankingServer struct {
	rankerv1.UnimplementedRankingServiceServer
}

func (s *rankingServer) GetTopStories(_ context.Context, req *rankerv1.TopStoriesRequest) (*rankerv1.TopStoriesResponse, error) {
	ticker := strings.ToUpper(strings.TrimSpace(req.GetTicker()))
	stories := mockStories(ticker)

	limit := int(req.GetLimit())
	switch {
	case limit <= 0:
		limit = defaultLimit
	case limit > len(stories):
		limit = len(stories)
	}

	return &rankerv1.TopStoriesResponse{
		Stories: stories[:limit],
	}, nil
}

func mockStories(ticker string) []*rankerv1.Story {
	if ticker == "" {
		ticker = "UNKNOWN"
	}

	return []*rankerv1.Story{
		{
			Id:          "story-1",
			Title:       ticker + " rallies after stronger-than-expected earnings",
			Source:      "Reuters",
			PublishedAt: "2026-03-07T13:30:00Z",
			Sentiment:   0.82,
		},
		{
			Id:          "story-2",
			Title:       ticker + " faces margin pressure concerns from analysts",
			Source:      "Bloomberg",
			PublishedAt: "2026-03-07T12:05:00Z",
			Sentiment:   -0.31,
		},
		{
			Id:          "story-3",
			Title:       "Institutional flows rotate into large-cap tech including " + ticker,
			Source:      "Financial Times",
			PublishedAt: "2026-03-07T10:40:00Z",
			Sentiment:   0.45,
		},
		{
			Id:          "story-4",
			Title:       ticker + " announces expansion of AI infrastructure partnerships",
			Source:      "CNBC",
			PublishedAt: "2026-03-06T21:15:00Z",
			Sentiment:   0.67,
		},
	}
}

func main() {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", listenAddr, err)
	}

	grpcServer := grpc.NewServer()
	rankerv1.RegisterRankingServiceServer(grpcServer, &rankingServer{})

	log.Printf("ranker gRPC server listening on %s", listenAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
