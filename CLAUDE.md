# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

The project uses the `just` command runner for key development tasks:

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
go tool goimports -local github.com/rajatgoel/gh-go -w .

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

1. **Frontend Service** (`internal/frontend/handler.go`)
   - Implements the gRPC service defined in Protocol Buffers
   - Acts as an adapter between client-facing API and backend storage
   - Handles two RPC methods: `Put` (store key-value) and `Get` (retrieve by key)
   - Uses standard gRPC error handling with status codes

2. **Backend Storage** (`internal/sqlbackend/`)
   - Handles data persistence using SQLite as an in-memory database
   - Defines a `Backend` interface with `Put` and `Get` operations
   - Current implementation (`sqliteBackend`) uses generated SQL code
   - Uses golang-migrate/migrate/v4 for database schema migrations
   - Migrations are embedded in the binary using go:embed

3. **API Definition** (`proto/frontend/v1/service.proto`)
   - Defines service contract using Protocol Buffers
   - Specifies message formats and RPC methods

4. **Entry Point** (`cmd/frontend/main.go`)
   - Application bootstrap
   - Sets up native gRPC server with TCP listener
   - Integrates OpenTelemetry instrumentation by default
   - Implements graceful shutdown with signal handling
   - Uses errgroup for coordinating goroutines

5. **Configuration** (`internal/config/config.go`)
   - Environment-based configuration using envconfig
   - Supports .env files via godotenv
   - Configurable service name, environment, and port

6. **Client Library** (`client/client.go`)
   - Type-safe gRPC client for the frontend service
   - Provides functional options pattern for configuration (WithTarget, WithDialer, WithInsecure)
   - Includes OpenTelemetry instrumentation for client-side tracing
   - Used in integration tests with custom dialers for in-memory testing

### Data Flow

1. Client makes a gRPC request to the server
2. Request is handled by the gRPC framework and routed to the appropriate handler
3. Handler calls the corresponding method on the backend
4. Backend executes the operation on the SQLite database
5. Results flow back to the client

### Testing

Integration tests in `itest/frontend_test.go` demonstrate end-to-end functionality:
- Setting up in-memory SQLite database
- Creating test gRPC server using `bufconn` for in-memory connections
- Using gRPC client to make requests
- Verifying that values can be stored and retrieved
- Testing error handling with both real and mock backends

Key tests:
- `TestBasic` - validates the complete Put/Get workflow including error handling for non-existent keys
- `TestErrorHandling` - tests error propagation using a mock backend that always fails

The tests use a shared `setupTestServer` helper function to reduce code duplication.

## Code Generation

The project relies on generated code:

1. **Protocol Buffers**: API definitions compiled to Go code using the opaque API
   - Generates both client and server interfaces for gRPC
   - Uses builder pattern for message construction (e.g., `GetResponse_builder{Value: value}.Build()`)
   - Configuration in `buf.gen.yaml` with paths=source_relative and default_api_level=API_OPAQUE
2. **SQL**: Queries generated from SQL files using sqlc
   - Configuration in `internal/sqlbackend/sqlc.yaml`
   - Generates type-safe Go code from SQL queries and schema
   - Uses SQLite engine with embedded schema

Always regenerate code after changes to proto files or SQL files:
```bash
just gen
```

## Docker

The project includes a Dockerfile for containerization. Build the image with:
```bash
just docker
```

This creates a multi-stage build optimized for production deployment.

## Database Migrations

The project uses golang-migrate/migrate/v4 for database schema management:

### Migration System
- Migrations are stored in `internal/sqlbackend/migrations/`
- Migration files follow the naming pattern: `NNNNNN_description.up.sql`
- Only up migrations are used (no down migrations)
- Migrations are embedded in the binary using `go:embed`
- All pending migrations are automatically applied when the database is opened

### Adding New Migrations
1. Create a new file in `internal/sqlbackend/migrations/`
2. Use the next sequential number (e.g., `000002_add_index.up.sql`)
3. Write your SQL schema changes
4. Rebuild and restart the service

### Current Migrations
- `000001_init_schema.up.sql` - Initial keyvalue table creation