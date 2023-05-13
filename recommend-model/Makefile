.DEFAULT_GOAL := help

all: clean compile ## Clean and compile

clean: ## Remove compiled files
	rm schema_pb2.py schema_pb2_grpc.py || true

compile: ## Compile schema.proto
	python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. schema.proto

help: ## Print help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'