SHELL := /bin/bash
ROOT := $(shell pwd)
PROTO_DIR := api/proto
GEN_DIR := api/gen/go

.PHONY: run-cli run-tui run-api build all proto grpc

run-cli:
	go run ./cmd/hadith-cli --help || true

run-tui:
	go run ./cmd/hadith-tui

run-api:
	ADDR=:8080 go run ./cmd/hadith-api

build:
	go build ./...

# Requires protoc and protoc-gen-go installed and on PATH.
proto:
	@mkdir -p $(GEN_DIR)
	protoc -I=$(PROTO_DIR) --go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/hadith.proto

# Build gRPC server after generating code and fetching dependencies.
grpc:
	GOFLAGS="-tags=grpc" go build ./cmd/hadith-grpc

