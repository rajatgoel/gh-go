# Contributing

This project uses `just` to orchestrate development tasks and relies on generated code (Protobuf and SQL). The sections below outline common workflows, tooling, architecture, and testing.

## Build Commands

```bash
# Generate code (Protocol Buffers and SQL)
just gen

# Run linting (depends on gen)
just lint

# Run tests (depends on lint)
just test

# Build Docker image
just docker
```

### Individual Commands

```bash
# Generate Protocol Buffer code
go tool buf generate proto

# Generate SQL code
find . -name sqlc.yaml | xargs go tool sqlc generate -f

# Format imports (groups: std, current package, 3rd-party)
go tool goimports -local github.com/dynoinc/gh-go -w .

# Update dependencies
go mod tidy

# Lint Protocol Buffers
go tool buf lint proto

# Check for breaking changes in Protocol Buffers
go tool buf breaking proto --against '.git#branch=main,subdir=proto'

# Apply automatic fixes
go fix ./...

# Run Go's built-in static analysis
go vet ./...

# Run additional static analysis
go tool staticcheck ./...

# Run security vulnerability check
go tool govulncheck ./...

# Run comprehensive linting
go tool golangci-lint run ./...

# Run tests
go test ./...

# Run tests with race detection and shuffle
go test -v -count=1 -race -shuffle=on ./...

# Run single test
go test ./itest -run TestBasic

# Check for dead code
go tool deadcode ./...
```

### Tool Management

All development tools are managed via `go get -tool` and listed in the `tool` section of `go.mod`. To add a new tool:

```bash
go get -tool <tool-package>
```

Current tools:
- `buf` - Protocol Buffer tooling
- `golangci-lint` - Go linting
- `sqlc` - SQL code generation
- `goimports` - Import formatting and organization
- `govulncheck` - Vulnerability scanning
- `staticcheck` - Static analysis

## Project Architecture

### Overview

gh-go is a key-value store service with a gRPC API. It follows a clean architecture pattern with separation between API definition, service implementation, and storage.

### Key Components

1. Frontend Service (`internal/frontend/handler.go`)
   - Implements the gRPC service defined in Protocol Buffers
   - Adapter between client-facing API and backend storage
   - Methods: `Put` (store key-value) and `Get` (retrieve by key)
   - Returns `NotFound` for missing keys

2. Backend Storage (`internal/sqlbackend/`)
   - In-memory SQLite database
   - `Backend` interface with `Put` and `Get`
   - `sqliteBackend` uses `sqlc`-generated queries
   - Migrations via `golang-migrate`, embedded with `go:embed`

3. API Definition (`proto/frontend/v1/service.proto`)
   - Message formats and RPC methods for gRPC

4. Entry Point (`cmd/frontend/main.go`)
   - Bootstrap, gRPC server, OTEL instrumentation, graceful shutdown

5. Configuration (`internal/config/config.go`)
   - Env-based config with `.env` support

6. Client Library (`client/client.go`)
   - Type-safe gRPC client with functional options
   - Includes OTEL instrumentation

### Data Flow

1. Client makes a gRPC request
2. Handler calls backend
3. Backend executes operation in SQLite
4. Result returns to the client

## Testing

Integration tests in `itest/frontend_test.go` cover end-to-end functionality and error handling. Unit tests in `internal/frontend/server_test.go` validate logging behavior of the interceptor.

## Code Generation

- Protobuf: configured via `buf.gen.yaml` and `proto/buf.yaml` (v2); code goes to `proto/`
- SQL: configured in `internal/sqlbackend/sqlc.yaml`; generated code in `internal/sqlbackend/sqlgen/`

Regenerate after changes to proto or SQL:

```bash
just gen
```

## Docker

Build container image:

```bash
just docker
```

Multi-stage build produces a small, distroless image.

## Database Migrations

- Migrations in `internal/sqlbackend/migrations/`
- Naming: `NNNNNN_description.up.sql`
- Embedded with `go:embed`
- Applied automatically on startup

Add a migration:
1. Create `internal/sqlbackend/migrations/000002_*.up.sql`
2. Write SQL changes
3. Rebuild or restart the service

