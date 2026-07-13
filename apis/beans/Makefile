.PHONY: generate-protos build clean help

help:
	@echo "Available targets:"
	@echo "  make generate-protos  - Regenerate protobuf code from .proto files"
	@echo "  make build            - Build the application"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make clean            - Remove build artifacts and .tools directory"
	@echo "  make run              - Run the application locally"

generate-protos:
	@cd .. && ./deployment/generate-protos.sh

build: generate-protos
	@cd .. && go build -o beansapi ./beans

docker-build: generate-protos
	@docker build -f Dockerfile -t soumitsr/go-beans-api:latest ..

run: generate-protos
	@cd .. && go run ./beans

clean:
	@rm -rf beansapi
	@cd .. && rm -rf .tools internal/embedding/.proto
	@echo "✓ Cleaned build artifacts"
