.DEFAULT_GOAL := help

all: clean compile ## Clean and compile

clean: ## Remove compiled files
	rm *.pb.go || true

compile: ## Compile schema.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative schema.proto

help: ## Print help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'