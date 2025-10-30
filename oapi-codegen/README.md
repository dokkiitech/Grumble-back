# oapi-codegen Configuration

This directory contains the configuration for generating Go server code from the OpenAPI specification.

## OpenAPI Specification Source

**Important**: This backend does NOT maintain its own OpenAPI specification. Instead, it uses the OpenAPI spec from the frontend repository:

- **Repository**: `git@github.com:dokkiitech/Grumble.git`
- **File**: `openapi.yaml` (in the root of frontend repo)
- **GitHub URL**: https://raw.githubusercontent.com/dokkiitech/Grumble/main/openapi.yaml

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
1. Look for `../Grumble/openapi.yaml` (if frontend repo is cloned as a sibling directory)
2. If not found, download the OpenAPI spec from GitHub
3. Generate Go server code using the configuration in `api.yml`

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

1. Frontend team updates `openapi.yaml` in the Grumble repository
2. Backend team runs `make generate` to regenerate server code
3. Implement the server interface methods in `internal/controller/`

## Important Notes

- **DO NOT edit** `internal/api/api.gen.go` manually - it will be overwritten
- The OpenAPI spec is the single source of truth, maintained by the frontend team
- Always regenerate code after OpenAPI spec changes
- If the frontend repo structure changes, update `OPENAPI_FILE` in Makefile

## Troubleshooting

### "Local OpenAPI file not found"
- Ensure the Grumble frontend repository is cloned as a sibling directory
- Or the command will automatically download from GitHub

### "oapi-codegen: command not found"
```bash
make init
```

### Generation fails
- Check that the OpenAPI spec is valid YAML
- Verify the GitHub URL is accessible
- Ensure you're using oapi-codegen v2
