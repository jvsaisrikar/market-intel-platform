# market-intel-platform
A distributed backend platform that ingests market news, enriches it with sentiment, and serves insights via GraphQL and gRPC.

## Running commands

Run from repo root:

```bash
# 1) Install/update deps
go mod tidy

# 2) Build everything
go build ./...
```

Terminal 1:

```bash
# 3) Start Ranker gRPC service (:50052)
go run ./services/ranker/cmd/server
```

Terminal 2:

```bash
# 4) Start GraphQL gateway (:8080)
go run ./gateway/cmd/server
```

Terminal 3 (gRPC direct test):

```bash
grpcurl -plaintext \
  -import-path proto/ranker \
  -proto ranker.proto \
  -d '{"ticker":"AAPL","limit":2}' \
  localhost:50052 \
  ranker.v1.RankingService/GetTopStories
```

Terminal 3 (GraphQL -> gRPC end-to-end test):

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"query { topStories(ticker: \"AAPL\", limit: 2) { title sentiment source } }"}'
```
