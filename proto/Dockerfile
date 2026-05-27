# Proto builder image — contains buf + protoc-gen-go + protoc-gen-go-grpc
# Build once, use for all proto generation.
# Usage: docker build -t wms-proto-builder:latest -f proto/Dockerfile .

FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

ARG BUF_VERSION=1.69.0
ARG PROTOC_GEN_GO_VERSION=latest
ARG PROTOC_GEN_GO_GRPC_VERSION=latest

# Install protoc-gen-go and protoc-gen-go-grpc via Go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION} \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION}

# Download and install buf
RUN go install github.com/bufbuild/buf/cmd/buf@v${BUF_VERSION}

FROM alpine:3.21

RUN apk add --no-cache libc6-compat

COPY --from=builder /go/bin/buf /usr/local/bin/buf
COPY --from=builder /go/bin/protoc-gen-go /usr/local/bin/protoc-gen-go
COPY --from=builder /go/bin/protoc-gen-go-grpc /usr/local/bin/protoc-gen-go-grpc

ENTRYPOINT ["buf"]
