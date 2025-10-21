.PHONY: help build run test clean docker-up docker-down

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the application
	go build -o asset-pulse-api main.go

run: ## Run the application
	go run main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -f asset-pulse-api
	go clean

tidy: ## Tidy and download dependencies
	go mod tidy
	go mod download

docker-build: ## Build docker image
	docker build -t asset-pulse-api:latest .

docker-up: ## Start docker compose services
	docker-compose up -d

docker-down: ## Stop docker compose services
	docker-compose down

docker-logs: ## View docker compose logs
	docker-compose logs -f

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

