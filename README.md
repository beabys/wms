# WMS — Modular Monorepo

Warehouse Management System built with Go gRPC services, Kafka event-driven architecture, and Vue 3 frontend.

## Architecture

- **Monorepo** — single `go.mod` at root, all services share proto contracts
- **gRPC** — synchronous internal communication between services
- **Kafka** — asynchronous event-driven workflows
- **BFF** — Backend-for-Frontend pattern (REST/OpenAPI per domain)
- **PostgreSQL** — schema-per-domain isolation
- **DDD + Hexagonal** — per-service structure based on [beabys/go-template](https://github.com/beabys/go-template)

## Prerequisites

- Go 1.26+
- Docker & Docker Compose
- Buf CLI (`brew install buf`)
- golangci-lint (`brew install golangci-lint`)

## Quick Start

```bash
# Start infrastructure (PostgreSQL, Redis, Kafka)
make docker-up

# Install dependencies
make tidy

# Generate proto Go code
make proto-gen

# Build all services
make build-all

# Run all tests
make test-all
```

## Project Structure

```
├── proto/defs/              # Proto contract definitions
├── proto/gen/go/            # Generated Go code
├── pkg/                     # Shared libraries
│   ├── authinterceptor/     # gRPC auth (JWT + API key)
│   ├── logger/              # Structured logging (zap)
│   └── postgres/            # PostgreSQL connection pool
├── deployment/              # Docker Compose + Kafka scripts
├── <domain>-service/        # Domain services (future)
├── <domain>-service-bff/    # BFF services (future)
└── FE/                      # Vue 3 frontend
```

## Make Targets

| Target | Description |
|--------|-------------|
| `build-all` | Build all Go packages |
| `test-all` | Run all tests with race detection |
| `lint` | Run golangci-lint |
| `proto-gen` | Generate Go code from protos |
| `docker-up` | Start all infrastructure |
| `docker-down` | Stop all infrastructure |
| `tidy` | Go mod tidy |
