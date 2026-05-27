SHELL:=/bin/bash

PROJECT_PATH := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

.PHONY: gen-api
gen-api:
	docker container run --rm -v $(PWD):/app golang:1.26 sh -c \
	"go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest && \
	cd /app && mkdir -p ./internal/api/v1 && \
	oapi-codegen --package=v1 --generate=\"types,client,spec,chi-server,skip-prune\" ./openapi.yaml > /tmp/v1.go.raw && \
	spec_line=$$(awk '/^var swaggerSpec = /{print NR; exit}' /tmp/v1.go.raw) && \
	head -n $$((spec_line - 1)) /tmp/v1.go.raw | sed 's/V1/v1/g' | sed 's/Id\$$/ID/g' | sed 's/Guid/GUID/g' | sed 's/Sku/SKU/g' | sed 's/Qoh/QOH/g' | sed 's/float32/float64/g' | sed 's/Url/URL/g' > /tmp/v1.go.code && \
	tail -n +$${spec_line} /tmp/v1.go.raw >> /tmp/v1.go.code && \
	mv /tmp/v1.go.code ./internal/api/v1/v1.go"

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test ./... -v

.PHONY: run
run:
	go run cmd/server/main.go
