SHELL:=/bin/bash
PROJECT_PATH := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
-include .env
export

# Default credentials (match docker-compose.infra.yml)
AUTH_DB_USER   ?= wms
AUTH_DB_PASSWORD ?= wms
AUTH_DB_HOST   ?= localhost
AUTH_DB_PORT   ?= 5432
AUTH_DB_NAME   ?= wms_auth
CUSTOMER_DB_USER    ?= wms
CUSTOMER_DB_PASSWORD ?= wms
CUSTOMER_DB_HOST    ?= localhost
CUSTOMER_DB_PORT ?= 5432
CUSTOMER_DB_NAME    ?= wms_customer

# ─── Proto ───────────────────────────────────────────
.PHONY: proto-gen
proto-gen:
	docker run -w /proto/defs --rm -v $(CURDIR)/proto:/proto --platform linux/amd64 ealves/buf:0.1.0-rc4 generate

# ─── Go ──────────────────────────────────────────────
.PHONY: tidy build-all test lint
tidy:
	go mod tidy
build-all:
	go build ./...
test:
	go test ./... -race -coverprofile .testCoverage.txt -v
lint:
	staticcheck ./... || true

# ─── Mocks ────────────────────────────────────────────
MOCK_DIRS := auth-service auth-service-bff customer-service customer-service-bff
.PHONY: mock-gen
mock-gen:
	@for dir in $(MOCK_DIRS); do \
		echo "Generating mocks in $$dir..."; \
		cd $$dir && mockery && cd ..; \
	done
	go mod tidy

# ─── Docker Network ──────────────────────────────────
.PHONY: network-create network-rm
network-create:
	docker network inspect wms-net >/dev/null 2>&1 || docker network create wms-net
network-rm:
	docker network rm wms-net 2>/dev/null || true

# ─── Infrastructure ──────────────────────────────────
.PHONY: start-infra stop-infra
start-infra:
	docker compose -f deployment/docker-compose.infra.yml up -d
stop-infra:
	docker compose -f deployment/docker-compose.infra.yml down

# ─── Backend (individual) ────────────────────────────
.PHONY: start-auth stop-auth start-customer stop-customer
start-auth: start-infra
	docker compose -f deployment/docker-compose.auth.yml up -d --build
stop-auth:
	docker compose -f deployment/docker-compose.auth.yml down
start-customer: start-infra
	docker compose -f deployment/docker-compose.customer.yml up -d --build
stop-customer:
	docker compose -f deployment/docker-compose.customer.yml down

# ─── Backend (all) ────────────────────────────────────
.PHONY: start-be stop-be be-logs be-ps
start-be: network-create
	docker compose -f deployment/docker-compose.yml up -d --build
stop-be:
	docker compose -f deployment/docker-compose.yml down
be-logs:
	docker compose -f deployment/docker-compose.yml logs -f
be-ps:
	docker compose -f deployment/docker-compose.yml ps

# ─── Migrations ──────────────────────────────────────
.PHONY: migrate-up migrate-down migrate-create
migrate-up:
	@read -p "Service (auth-service/customer-service): " svc; \
	url=""; \
	if [ "$$svc" = "auth-service" ]; then \
		url="postgres://$(AUTH_DB_USER):$(AUTH_DB_PASSWORD)@$(AUTH_DB_HOST):$(AUTH_DB_PORT)/$(AUTH_DB_NAME)?sslmode=disable"; \
	elif [ "$$svc" = "customer-service" ]; then \
		url="postgres://$(CUSTOMER_DB_USER):$(CUSTOMER_DB_PASSWORD)@$(CUSTOMER_DB_HOST):$(CUSTOMER_DB_PORT)/$(CUSTOMER_DB_NAME)?sslmode=disable"; \
	fi; \
	migrate -path $$svc/migrations -database "$$url" up

migrate-down:
	@read -p "Service: " svc; \
	read -p "Steps (default 1): " steps; \
	url="postgres://$(AUTH_DB_USER):$(AUTH_DB_PASSWORD)@$(AUTH_DB_HOST):$(AUTH_DB_PORT)/$(AUTH_DB_NAME)?sslmode=disable"; \
	migrate -path $$svc/migrations -database "$$url" down $${steps:-1}

migrate-create:
	@read -p "Service: " svc; \
	read -p "Description: " desc; \
	migrate create -ext sql -dir $$svc/migrations -seq $$desc

# ─── Frontend (individual) ────────────────────────────
.PHONY: start-admin-ui start-customer-ui
start-admin-ui:
	cd admin-ui && npm run dev
start-customer-ui:
	cd customer-ui && npm run dev

# ─── Frontend (all) ──────────────────────────────────
.PHONY: start-fe stop-fe
start-fe:
	@echo "Starting admin-ui on :5173 and customer-ui on :5174..."
	cd admin-ui && npm run dev & \
	cd customer-ui && npm run dev & \
	wait
stop-fe:
	-pkill -f "vite.*admin-ui" 2>/dev/null || true
	-pkill -f "vite.*customer-ui" 2>/dev/null || true

# ─── Dev (everything) ────────────────────────────────
.PHONY: dev
dev: network-create
	docker compose -f deployment/docker-compose.infra.yml up -d
	@echo "Waiting for databases..."
	@sleep 3
	$(MAKE) start-auth
	$(MAKE) start-customer
	@echo "Backend services started. Run 'make start-fe' for frontend."

# ─── Utility ──────────────────────────────────────────
.PHONY: unit unit-coverage
unit:
	go test $$(go list ./pkg/... ./auth-service/... ./customer-service/... 2>/dev/null | grep -v /mocks) -race -coverprofile .testCoverage.txt -v
unit-coverage: unit
	go tool cover -html=.testCoverage.txt -o unit.html
