# Variables
DEV_CONFIG_FILE := config/config.yaml
DEV_COMPOSE_FILE := docker-compose.yaml

# Extract values from config.dev.yml using yq (you need to install yq first)
DEV_DB_PASSWORD := $(shell yq -e '.db.password' $(DEV_CONFIG_FILE))
DEV_DB_NAME := $(shell yq -e '.db.dbname' $(DEV_CONFIG_FILE))
DEV_DB_USER := $(shell yq -e '.db.user' $(DEV_CONFIG_FILE))

# Common environment variables
define set_dev_env_vars
	DB_PASSWORD=$(DEV_DB_PASSWORD) \
	DB_NAME=$(DEV_DB_NAME) \
	DB_USER=$(DEV_DB_USER) 
endef

# default commands (TODO: change later)
.PHONY: fmt vet build clean
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vetc
	@ go build

clean: 
	@ go clean


# development commands
.PHONY: dev-build
dev-build:
	docker compose -f $(DEV_COMPOSE_FILE) build

.PHONY: dev-up
dev-up:
	@echo "Starting development environment..."
	$(call set_dev_env_vars) docker compose -f $(DEV_COMPOSE_FILE) up -d

.PHONY: dev-build-up
dev-build-up: dev-build dev-up

.PHONY: dev-down
dev-down:
	docker compose -f $(DEV_COMPOSE_FILE) down

.PHONY: dev-logs
dev-logs:
	docker compose -f $(DEV_COMPOSE_FILE) logs -f

.PHONY: dev-ps
dev-ps:
	docker compose -f $(DEV_COMPOSE_FILE) ps

.PHONY: dev-db-create
dev-db-create:
	@echo "Starting database creation..."
	$(call set_dev_env_vars) docker compose -f $(DEV_COMPOSE_FILE) up -d
	$(call set_dev_env_vars) docker compose -f $(DEV_COMPOSE_FILE) exec main ./main pg_create
	$(call set_dev_env_vars) docker compose -f $(DEV_COMPOSE_FILE) down

.PHONY: dev-db-migrate
dev-db-migrate:
	@echo "Starting migration..."
	$(call set_dev_env_vars) docker compose -f $(DEV_COMPOSE_FILE) up -d
	$(call set_dev_env_vars) docker compose -f $(DEV_COMPOSE_FILE) exec main ./main pg_migrate
	$(call set_dev_env_vars) docker compose -f $(DEV_COMPOSE_FILE) down

.PHONY: dev-clear
dev-clear:
	echo "Starting clearing..."
	echo "Stopping containers..."
	docker compose -f $(DEV_COMPOSE_FILE) down
	docker compose -f $(DEV_COMPOSE_FILE) down --remove-orphans
	echo "Removing containers and volumes..."
	docker compose -f $(DEV_COMPOSE_FILE) rm -v -s
	echo "Removing named volumes..."
	docker compose -f $(DEV_COMPOSE_FILE) down -v
	
.PHONY: help
help:
	@echo ""
	@echo "Available commands:"
	@echo "  dev-build     	- Build the application container"
	@echo "  dev-up        	- Start the development environment"
	@echo "  dev-build-up  	- Build and start the development environment"
	@echo "  dev-down      	- Stop the development environment"
	@echo "  dev-logs      	- Show logs from all containers"
	@echo "  dev-ps        	- List running containers"
	@echo "  dev-db-create 	- Create and initialize databases"
	@echo "  dev-db-migrate	- Run database migrations"
	@echo "  dev-clear     	- Remove all containers and volumes"
	@echo ""


.DEFAULT_GOAL := help