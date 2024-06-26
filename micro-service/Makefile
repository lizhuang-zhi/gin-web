.PHONY: client
client: ## Startup Client
	go run ./client/main.go

.PHONY: proto
proto: ## Generate protobuf
	./scripts/gen-pb.sh

.PHONY: run
run: ## Startup Micro Service
	go run ./cluster/main.go -c=./configs/local/config.yaml

