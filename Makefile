# Makefile for Social Login Application

# Variables
DOCKER_COMPOSE := docker-compose
BACKEND_DIR := backend
FRONTEND_DIR := frontend
INFRA_DIR := infra

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  setup      - Initial setup (environment variables + dependencies)"
	@echo "  setup-env  - Setup environment variables (.env file)"
	@echo "  build      - Build all services"
	@echo "  up         - Start all services"
	@echo "  down       - Stop all services"
	@echo "  restart    - Restart all services"
	@echo "  logs       - Show logs"
	@echo "  clean      - Clean up containers and volumes"
	@echo "  backend    - Run backend in development mode"
	@echo "  frontend   - Run frontend in development mode"
	@echo "  infra      - Deploy infrastructure"
	@echo "  infra-destroy - Destroy infrastructure"

# Setup
.PHONY: setup
setup: setup-env setup-backend setup-frontend

.PHONY: setup-env
setup-env:
	@echo "Setting up environment variables..."
	./scripts/setup-env.sh

.PHONY: setup-backend
setup-backend:
	@echo "Setting up backend..."
	cd $(BACKEND_DIR) && go mod download
	cd $(BACKEND_DIR) && go mod tidy

.PHONY: setup-frontend
setup-frontend:
	@echo "Setting up frontend..."
	cd $(FRONTEND_DIR) && npm install

# Docker operations
.PHONY: build
build:
	@echo "Building all services..."
	$(DOCKER_COMPOSE) build

.PHONY: up
up:
	@echo "Starting all services..."
	$(DOCKER_COMPOSE) up -d

.PHONY: down
down:
	@echo "Stopping all services..."
	$(DOCKER_COMPOSE) down

.PHONY: restart
restart: down up

.PHONY: logs
logs:
	$(DOCKER_COMPOSE) logs -f

.PHONY: clean
clean:
	@echo "Cleaning up containers and volumes..."
	$(DOCKER_COMPOSE) down -v --rmi all --remove-orphans

# Development
.PHONY: backend
backend:
	@echo "Running backend in development mode..."
	cd $(BACKEND_DIR) && air

.PHONY: frontend
frontend:
	@echo "Running frontend in development mode..."
	cd $(FRONTEND_DIR) && npm run dev

.PHONY: dev
dev:
	@echo "Starting development environment..."
	$(DOCKER_COMPOSE) up -d postgres
	@echo "Waiting for database to be ready..."
	sleep 5
	@echo "Run 'make backend' and 'make frontend' in separate terminals"

# Infrastructure
.PHONY: infra
infra:
	@echo "Deploying infrastructure..."
	cd $(INFRA_DIR) && terraform init
	cd $(INFRA_DIR) && terraform plan
	cd $(INFRA_DIR) && terraform apply

.PHONY: infra-plan
infra-plan:
	@echo "Planning infrastructure changes..."
	cd $(INFRA_DIR) && terraform plan

.PHONY: infra-destroy
infra-destroy:
	@echo "Destroying infrastructure..."
	cd $(INFRA_DIR) && terraform destroy

# Database operations
.PHONY: db-up
db-up:
	@echo "Starting database..."
	$(DOCKER_COMPOSE) up -d postgres

.PHONY: db-down
db-down:
	@echo "Stopping database..."
	$(DOCKER_COMPOSE) stop postgres

.PHONY: db-reset
db-reset:
	@echo "Resetting database..."
	$(DOCKER_COMPOSE) down postgres
	$(DOCKER_COMPOSE) up -d postgres

# Linting and formatting
.PHONY: lint
lint: lint-backend lint-frontend

.PHONY: lint-backend
lint-backend:
	@echo "Linting backend..."
	cd $(BACKEND_DIR) && golangci-lint run

.PHONY: lint-frontend
lint-frontend:
	@echo "Linting frontend..."
	cd $(FRONTEND_DIR) && npm run lint

.PHONY: format
format: format-backend format-frontend

.PHONY: format-backend
format-backend:
	@echo "Formatting backend..."
	cd $(BACKEND_DIR) && go fmt ./...

.PHONY: format-frontend
format-frontend:
	@echo "Formatting frontend..."
	cd $(FRONTEND_DIR) && npm run format

# Status check
.PHONY: status
status:
	@echo "Checking service status..."
	$(DOCKER_COMPOSE) ps

# Install Air for hot reload
.PHONY: install-air
install-air:
	@echo "Installing Air for hot reload..."
	go install github.com/air-verse/air@latest 