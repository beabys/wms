.PHONY: build-all test-all lint proto-gen docker-up docker-down \
        run-login run-customer run-inbound tidy seed

# Build all Go packages
build-all:
	go build ./...

# Run all unit and integration tests
test-all:
	go test ./... -race -cover -count=1

# Run golangci-lint
lint:
	golangci-lint run ./...

# Generate proto Go code using Buf
proto-gen:
	cd proto/defs && buf generate

# Docker Compose lifecycle
docker-up:
	docker compose -f deployment/docker-compose.yml up -d

docker-down:
	docker compose -f deployment/docker-compose.yml down

# Service run targets (requires docker-up for infra)
.PHONY: run
run:  ## Start all services as modular monolith (one process)
	go run ./cmd/wms

.PHONY: run-login
run-login:
	go run ./login-service/cmd/server

.PHONY: run-customer
run-customer:
	go run ./customer-service/cmd/server

.PHONY: run-inbound
run-inbound:
	go run ./inbound-service/cmd/server

# Go module maintenance
tidy:
	go mod tidy

# =============================================================================
# Database Migrations (golang-migrate/migrate CLI)
# =============================================================================

# Install migrate CLI
.PHONY: install-migrate
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create new migration for a service (usage: make migrate-create SERVICE=login-service)
.PHONY: migrate-create
migrate-create:
	@read -p "Migration name: " name; \
	cd $(SERVICE) && migrate create -ext sql -dir migrations -seq $$name

# Run all pending migrations for a service (usage: make migrate-up SERVICE=login-service DATABASE_URL=...)
.PHONY: migrate-up
migrate-up:
	cd $(SERVICE) && migrate -database "$(DATABASE_URL)" -path migrations up

# Rollback last migration for a service (usage: make migrate-down SERVICE=login-service DATABASE_URL=...)
.PHONY: migrate-down
migrate-down:
	cd $(SERVICE) && migrate -database "$(DATABASE_URL)" -path migrations down 1

# Create a seed user for development (requires running services + PostgreSQL)
.PHONY: seed
seed:
	@echo "Registering seed user..."
	@curl -s -X POST http://localhost:8080/v1/auth/register \
		-H 'Content-Type: application/json' \
		-d '{"email":"admin@wms.com","password":"password123","company_name":"WMS Admin"}' && echo ""
	@echo "Seed user created: admin@wms.com / password123"

# Show help
help:
	@echo "Targets:"
	@echo "  build-all    - go build ./..."
	@echo "  test-all     - go test ./... -race -cover"
	@echo "  lint         - golangci-lint run"
	@echo "  proto-gen    - generate Go code from protos"
	@echo "  docker-up    - start all infra (PostgreSQL, Redis, Kafka)"
	@echo "  docker-down  - stop all infra"
	@echo "  run-login    - go run login-service"
	@echo "  run-customer - go run customer-service"
	@echo "  run-inbound  - go run inbound-service"
	@echo "  tidy         - go mod tidy"
	@echo "  install-migrate - install golang-migrate CLI"
	@echo "  migrate-create  - create new migration (SERVICE=<name>)"
	@echo "  migrate-up      - run pending migrations (SERVICE=... DATABASE_URL=...)"
	@echo "  migrate-down    - rollback last migration (SERVICE=... DATABASE_URL=...)"
