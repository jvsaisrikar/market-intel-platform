# market-intel-platform
A distributed backend platform that ingests market news, enriches it with sentiment, and serves insights via GraphQL and gRPC.

## Run and test locally

From the repository root:

```bash
# Start the Ranker gRPC service on :50052
go run ./services/ranker/cmd/server
```

In a second terminal:

```bash
# Call GetTopStories (expects 2 stories in response)
grpcurl -plaintext \
  -import-path proto/ranker \
  -proto ranker.proto \
  -d '{"ticker":"AAPL","limit":2}' \
  localhost:50052 \
  ranker.v1.RankingService/GetTopStories
```

Optional compile check:

```bash
go build ./services/ranker/cmd/server
```
