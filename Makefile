.PHONY: init generate clean test format lint precommit local-setup local-db-up local-db-down local-migrate local-api local-down

# OpenAPI schema location (managed in this repository)
OPENAPI_FILE := openapi.yaml
DOCKER_COMPOSE := docker compose -f docker/docker-compose.yml
LOCAL_DATABASE_URL := postgres://grumble:grumble@localhost:5432/grumble?sslmode=disable
LOCAL_HTTP_ADDR := :8080

__init_go__:
	@go mod download
	@go mod tidy

__init_oapi_codegen__:
	@if ! command -v oapi-codegen &> /dev/null; then \
		echo "oapi-codegen not found, installing..."; \
		go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest; \
	else \
		echo "oapi-codegen is already installed."; \
	fi

init: __init_go__ __init_oapi_codegen__

# Generate code from OpenAPI file
generate:
	@if [ ! -f "$(OPENAPI_FILE)" ]; then \
		echo "Error: OpenAPI file not found: $(OPENAPI_FILE)"; \
		exit 1; \
	fi
	@echo "Generating from OpenAPI file: $(OPENAPI_FILE)"
	@oapi-codegen -config oapi-codegen/api.yml $(OPENAPI_FILE)

# Clean generated files
clean:
	@rm -f internal/api/api.gen.go

# Run tests
test:
	@go test -v ./...

# Format code
format:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatting complete."

# Lint code
lint:
	@echo "Running linter..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Install it with:"; \
		echo "  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2"; \
		exit 1; \
	fi
	@golangci-lint run ./...
	@echo "Linting complete."

# Pre-commit checks: format, lint, and test
precommit: format lint test
	@echo "✓ All pre-commit checks passed!"

vendor:
	@go mod vendor

build/api: init generate
	@go build -o bin/ ./cmd/api/main.go

build/batch: init generate
	@go build -o bin/ ./cmd/batch/main.go

local-setup:
	@$(MAKE) init
	@$(MAKE) generate
	@$(MAKE) local-db-up
	@$(MAKE) local-migrate
	@echo "Local environment is ready."

local-db-up:
	@$(DOCKER_COMPOSE) up -d db
	@echo "Waiting for database to become ready..."
	@until $(DOCKER_COMPOSE) exec -T db pg_isready -U grumble >/dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "Database is ready."

local-db-down local-down:
	@$(DOCKER_COMPOSE) down

local-migrate:
	@DATABASE_URL="$(LOCAL_DATABASE_URL)" go run ./cmd/migrate

local-api:
	@CURDIR="$(CURDIR)" bash -lc 'set -a; if [ -f .env ]; then source .env; fi; set +a; \
		if [ -z "$$FIREBASE_CREDENTIALS_FILE" ] && [ -f firebase_secrets.json ]; then \
			export FIREBASE_CREDENTIALS_FILE="$$CURDIR/firebase_secrets.json"; \
		fi; \
		export DATABASE_URL="$${DATABASE_URL:-$(LOCAL_DATABASE_URL)}"; \
		export GRUMBLE_HTTP_ADDR="$${GRUMBLE_HTTP_ADDR:-$(LOCAL_HTTP_ADDR)}"; \
		go run ./cmd/api'

local-batch:
	@CURDIR="$(CURDIR)" bash -lc 'set -a; if [ -f .env ]; then source .env; fi; set +a; \
		if [ -z "$$FIREBASE_CREDENTIALS_FILE" ] && [ -f firebase_secrets.json ]; then \
			export FIREBASE_CREDENTIALS_FILE="$$CURDIR/firebase_secrets.json"; \
		fi; \
		export DATABASE_URL="$${DATABASE_URL:-$(LOCAL_DATABASE_URL)}"; \
		go run ./cmd/batch'
