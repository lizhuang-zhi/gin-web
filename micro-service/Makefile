.PHONY: proto
proto: ## Generate protobuf
	./gen-pb.sh

.PHONY: run
run: ## Startup Micro Service
	go run main.go -c=./configs/local/config.yaml

