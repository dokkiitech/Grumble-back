.PHONY: init generate clean test

# OpenAPI schema location (external repository)
OPENAPI_REPO := git@github.com:dokkiitech/Grumble.git
OPENAPI_FILE := ../Grumble/openapi.yaml
OPENAPI_REMOTE_URL := https://raw.githubusercontent.com/dokkiitech/Grumble/main/openapi.yaml

__init_go__:
	@go mod download
	@go mod tidy

__init_oapi_codegen__:
	@if ! command -v oapi-codegen &> /dev/null; then \
		echo "oapi-codegen not found, installing..."; \
		go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest; \
	else \
		echo "oapi-codegen is already installed."; \
	fi

init: __init_go__ __init_oapi_codegen__

# Generate code from local OpenAPI file (if Grumble repo is cloned as sibling)
generate:
	@if [ -f "$(OPENAPI_FILE)" ]; then \
		echo "Generating from local OpenAPI file: $(OPENAPI_FILE)"; \
		oapi-codegen -config oapi-codegen/api.yml $(OPENAPI_FILE); \
	else \
		echo "Local OpenAPI file not found. Downloading from GitHub..."; \
		curl -sSL $(OPENAPI_REMOTE_URL) -o /tmp/openapi.yaml; \
		oapi-codegen -config oapi-codegen/api.yml /tmp/openapi.yaml; \
	fi

# Clean generated files
clean:
	@rm -f internal/api/api.gen.go

# Run tests
test:
	@go test -v ./...

vendor:
	@go mod vendor
