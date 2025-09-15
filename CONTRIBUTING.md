# Contributing to hadith-go

Thanks for your interest in contributing! This guide explains how to get set up and how to propose changes.

## Development setup

- Requirements: Go 1.21+.
- Optional (for gRPC): `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`.

Clone and run:

```
go run ./cmd/hadith-api   # starts REST API + serves web UI at /
# or
go run ./cmd/hadith-cli books
go run ./cmd/hadith-tui
```

Make targets:

```
make run-api
make run-cli
make run-tui
make build
make proto   # generate gRPC code
make grpc    # build gRPC server (after make proto)
```

## Project structure

- `internal/data` — JSON loader + in-memory store
- `internal/search` — simple substring search with scoring and ordering
- `cmd/` — CLI, TUI, REST API, optional gRPC
- `web/` — static web UI (served by API)

## Coding conventions

- Keep changes minimal and focused; preserve stable JSON shapes returned by the API.
- Maintain book discovery behavior (walk up from CWD for `books` dir).
- Run `go fmt ./...` before submitting.

## Proposing changes

- Open an issue describing the problem or proposal.
- For PRs:
  - Include a clear description of behavior changes and testing steps.
  - Update docs when adding flags, endpoints, or UI behavior.
  - Avoid unrelated refactors.

## Reporting security issues

Please see `SECURITY.md` for responsible disclosure guidelines.

