# oapi-codegen Configuration

This directory contains the configuration for generating Go server code from the OpenAPI specification.

## OpenAPI Specification Source

The OpenAPI specification is managed in this repository:

- **File**: `openapi.yaml` (in the root of this repository)
- The spec was originally from the frontend repository but is now maintained here

## Code Generation

### Prerequisites

```bash
make init
```

This will:
- Install Go dependencies
- Install oapi-codegen v2

### Generate Server Code

```bash
make generate
```

This command will:
1. Check if `openapi.yaml` exists in the root directory
2. Generate Go server code using the configuration in `api.yml`

### Generated Output

- **File**: `internal/api/api.gen.go`
- **Package**: `api`
- **Framework**: Gin
- **Contains**:
  - Model types (structs for request/response)
  - Server interface
  - Gin route registration functions

## Configuration Details

The `api.yml` file contains:

- **Package**: `api`
- **Framework**: `gin-server` (Gin web framework)
- **Models**: Generated from OpenAPI schemas
- **Strict Server**: Enabled for better type safety
- **Output**: `internal/api/api.gen.go`

## Directory Structure

```
Grumble-back/
├── oapi-codegen/
│   ├── api.yml           # oapi-codegen configuration
│   └── README.md         # This file
├── internal/
│   └── api/
│       └── api.gen.go    # Generated code (do not edit manually)
└── Makefile              # Build commands
```

## Development Workflow

1. Update `openapi.yaml` in the root directory
2. Run `make generate` to regenerate server code
3. Implement the server interface methods in `internal/controller/`

## Important Notes

- **DO NOT edit** `internal/api/api.gen.go` manually - it will be overwritten
- The OpenAPI spec is the single source of truth
- Always regenerate code after OpenAPI spec changes

## Troubleshooting

### "OpenAPI file not found"
- Ensure `openapi.yaml` exists in the root directory of this repository

### "oapi-codegen: command not found"
```bash
make init
```

### Generation fails
- Check that the OpenAPI spec is valid YAML
- Ensure you're using oapi-codegen v2
