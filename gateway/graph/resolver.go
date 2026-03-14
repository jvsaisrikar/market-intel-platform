package graph

import "github.com/saisrikar/market-intel-platform/gateway/internal/ranker"

type Resolver struct {
	RankerClient ranker.Client
}
