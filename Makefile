.PHONY: help build run test test-cover lint clean up down logs migrate-up migrate-down

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Go commands
build: ## Build the application
	go build -o bin/alarm-agent cmd/server/main.go

run: ## Run the application locally
	go run cmd/server/main.go

test: ## Run tests
	go test -v ./...

test-cover: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linting
	golangci-lint run

clean: ## Clean build artifacts
	rm -rf bin/ coverage.out coverage.html

# Docker commands
up: ## Start all services with docker-compose
	docker-compose up -d

down: ## Stop all services
	docker-compose down

logs: ## View application logs
	docker-compose logs -f app

# Database commands
migrate-up: ## Apply database migrations
	docker-compose run --rm migrate -path=/migrations -database=postgres://alarm_user:alarm_pass@postgres:5432/alarm_agent?sslmode=disable up

migrate-down: ## Rollback database migrations
	docker-compose run --rm migrate -path=/migrations -database=postgres://alarm_user:alarm_pass@postgres:5432/alarm_agent?sslmode=disable down

# Development commands
dev: up ## Start development environment
	@echo "Development environment started!"
	@echo "API: http://localhost:8080"
	@echo "Health: http://localhost:8080/health"
	@echo "Metrics: http://localhost:8080/metrics"

# Testing commands
test-webhook: ## Test webhook endpoint with sample data
	curl -X POST http://localhost:8080/webhook/whatsapp \
		-H "Content-Type: application/json" \
		-d '{"results":[{"messageId":"test-123","from":"5511999999999","to":"5511888888888","receivedAt":"2024-01-01T10:00:00Z","message":{"type":"TEXT","text":"Marcar dentista amanhã 14h"}}]}'

# Production deployment
deploy: ## Build and deploy (placeholder)
	@echo "Deploy target - implement according to your deployment strategy"
	docker build -t alarm-agent:latest .

# Database utilities
db-shell: ## Connect to database shell
	docker-compose exec postgres psql -U alarm_user -d alarm_agent

# Monitoring
status: ## Check service status
	docker-compose ps
	@echo ""
	@echo "Health checks:"
	@curl -s http://localhost:8080/health || echo "Service not responding"
	@echo ""
	@curl -s http://localhost:8080/ready || echo "Service not ready"