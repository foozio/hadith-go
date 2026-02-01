SHELL := /bin/bash
ROOT := $(shell pwd)
PROTO_DIR := api/proto
GEN_DIR := api/gen/go
DIST_DIR := dist

.PHONY: run-cli run-tui run-api build all proto grpc release clean-dist

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

clean-dist:
	rm -rf $(DIST_DIR)

release: clean-dist
	mkdir -p $(DIST_DIR)
	# Darwin ARM64
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/hadith-cli-darwin-arm64 ./cmd/hadith-cli
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/hadith-api-darwin-arm64 ./cmd/hadith-api
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/hadith-tui-darwin-arm64 ./cmd/hadith-tui
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/hadith-cli-linux-amd64 ./cmd/hadith-cli
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/hadith-api-linux-amd64 ./cmd/hadith-api
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/hadith-cli-windows-amd64.exe ./cmd/hadith-cli
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/hadith-api-windows-amd64.exe ./cmd/hadith-api
	@echo "Build complete. Artifacts in $(DIST_DIR)/"