.PHONY: client
client: ## Startup Client
	go run ./client/main.go

.PHONY: proto
proto: ## Generate protobuf
	sh ./scripts/gen-pb.sh

.PHONY: run
run: ## Startup Micro Service
	go run ./cmd/server/main.go -c=./configs/local/config.yaml

