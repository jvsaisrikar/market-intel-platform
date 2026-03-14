package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/saisrikar/market-intel-platform/gateway/graph"
	"github.com/saisrikar/market-intel-platform/gateway/internal/ranker"
)

const (
	defaultHTTPAddr   = ":8080"
	defaultRankerAddr = "localhost:50052"
)

func main() {
	rankerAddr := envOrDefault("RANKER_ADDR", defaultRankerAddr)
	httpAddr := envOrDefault("GATEWAY_ADDR", defaultHTTPAddr)

	rankerClient, err := ranker.NewClient(rankerAddr)
	if err != nil {
		log.Fatalf("failed to create ranker client for %s: %v", rankerAddr, err)
	}
	defer func() {
		if closeErr := rankerClient.Close(); closeErr != nil {
			log.Printf("failed to close ranker client: %v", closeErr)
		}
	}()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			RankerClient: rankerClient,
		},
	}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("GraphQL gateway listening on %s (ranker: %s)", httpAddr, rankerAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatalf("failed to start gateway HTTP server: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
