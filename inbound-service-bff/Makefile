SHELL:=/bin/bash

.PHONY: gen-api
gen-api:
	docker container run --rm -v $(PWD):/app golang:1.26 sh -c \
	"go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest && \
	cd /app && mkdir -p ./internal/api/v1 && \
	oapi-codegen --package=v1 --generate=\"types,client,spec,chi-server,skip-prune\" ./openapi.yaml | \
	sed 's/V1/v1/g' | sed 's/Id\$$/ID/g' | sed 's/Guid/GUID/g' | sed 's/Sku/SKU/g' | \
	sed 's/Qoh/QOH/g' | sed 's/float32/float64/g' | sed 's/Url/URL/g' > ./internal/api/v1/v1.go"

.PHONY: run
run:
	go mod tidy && go run ./inbound-service-bff/cmd/server/main.go
