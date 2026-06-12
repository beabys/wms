SHELL:=/bin/bash

PROJECT_PATH := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

.PHONY: run
run:
	go mod tidy && go run cmd/server/main.go

.PHONY: build
build:
	go build -o $(PROJECT_PATH)/bin/server cmd/server/main.go

.PHONY: gen-api
gen-api:
	mkdir -p $(PROJECT_PATH)/internal/api/v1
	cd $(PROJECT_PATH) && oapi-codegen --package=v1 --generate='types,client,spec,chi-server,skip-prune' ./openapi.yaml > ./internal/api/v1/v1.go

.PHONY: mockery
mockery:
	cd $(PROJECT_PATH) && mockery && go mod tidy

.PHONY: test
test:
	go test ./auth-service-bff/... -v -count=1

.PHONY: test-unit
test-unit:
	go test ./auth-service-bff/internal/... -v -count=1

.PHONY: vet
vet:
	go vet ./auth-service-bff/...

.PHONY: build-all
build-all:
	go build ./auth-service-bff/...
