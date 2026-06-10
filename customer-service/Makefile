SHELL:=/bin/bash

PROJECT_PATH := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

.PHONY: run
run:
	go mod tidy && go run cmd/server/main.go

.PHONY: build
build:
	go build -o $(PROJECT_PATH)/bin/server cmd/server/main.go

.PHONY: proto
proto:
	$(MAKE) -C $(PROJECT_PATH)/.. proto

.PHONY: mockery
mockery:
	cd $(PROJECT_PATH) && mockery && go mod tidy

.PHONY: migrate-up
migrate-up:
	migrate -path $(PROJECT_PATH)/migrations -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path $(PROJECT_PATH)/migrations -database "$(DATABASE_URL)" down 1

.PHONY: test
test:
	go test ./internal/... -v -count=1
