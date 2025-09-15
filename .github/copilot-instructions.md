# hadith-go — AI Assistant Instructions

## Overview
- Purpose: Load hadith collections from `books/*.json`, search them, and expose CLI, TUI, REST API, and optional gRPC.
- Data flow: JSON files → `internal/data.Store` (in-memory) → search via `internal/search.SimpleSearch` → output via CLI/TUI/HTTP/gRPC.

## Repo Layout
- `cmd/hadith-cli`: Lists books, counts, gets by book+number, substring search.
- `cmd/hadith-tui`: Minimal line-based TUI with paging and commands (`:help`, `:full`, `:short`, `:width N`, `:color on|off`).
- `cmd/hadith-api`: REST API with CORS; endpoints for health, books, count, search, and hadith detail.
- `cmd/hadith-grpc`: gRPC server behind build tag `grpc` (see Makefile and proto section).
- `internal/data`: JSON loader and in-memory store with RWMutex and stable book ordering.
- `internal/search`: Case-insensitive substring search with simple scoring and deterministic sort.
- `api/proto/hadith.proto`: Proto definitions; generated Go lives under `api/gen/go/hadithpb`.

## Data Model and Loading
- JSON format: each `books/<name>.json` is an array of `{ number:int, arab:string, id:string }`.
- Loader (`internal/data/loader.go`):
  - Uses filename (sans `.json`) as `Hadith.Book`.
  - Builds `Store.byBook` and `Store.books` (sorted) at init; store is read-only afterwards.
  - Key APIs: `Books()`, `Count()`, `Get(book, number)`, `All()`.
- Book discovery: binaries search upwards from CWD for a directory containing `books` (see `findBooksRoot()` in each `cmd/*`).

## Search Behavior
- `internal/search.SimpleSearch(all, query, limit)`:
  - Case-insensitive `contains` on `ID` (+3), `Arab` (+2), `Book` (+1).
  - Sorts by `Score` desc, then `Book` asc, then `Number` asc.
  - Applies `limit` after sorting; `limit<=0` means no cap.

## REST API (`cmd/hadith-api`)
- Env: `ADDR` (default `:8080`). CORS: `*` with `GET, OPTIONS`.
- Endpoints:
  - `GET /healthz` → `ok`.
  - `GET /books` → `[]string`.
  - `GET /count` → `{ "count": N }`.
  - `GET /search` → array of results `{ hadith, score }`.
    - Query params:
      - `q`: search term; when empty with `book` set, returns all entries in that book (browse mode).
      - `book`: exact book name.
      - Pagination modes (precedence):
        1) `offset` + `limit` (default limit 50, max 200) with headers `X-Total-Count`, `X-Offset`, `X-Limit`
        2) `page` + `page_size` (default 50, max 200) with headers `X-Total-Count`, `X-Page`, `X-Page-Size`
        3) legacy `limit` only.
  - `GET /hadith/{book}/{number}` → `Hadith` or 404.
- JSON is pretty-printed for readability.

## CLI and TUI
- CLI (`cmd/hadith-cli`):
  - `books | count | get <book> <number> | search [-limit N] <query>`.
  - `get` prints indented JSON; `search` prints readable, truncated lines with score.
- TUI (`cmd/hadith-tui`):
  - Type query to search; `n/p` to page; `o N` to open; `q` to quit; see `:help`.
  - Browse by book is supported in the web UI; for CLI/TUI, use search with a book name included in the query as a workaround or extend as needed.

## gRPC (optional)
- Proto: `api/proto/hadith.proto` with service `HadithService`.
- Generate: `make proto` (requires `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`).
- Build server: `make grpc` (uses `-tags grpc`). Without tag, `cmd/hadith-grpc/main_stub.go` runs.

## Dev Workflows
- Run CLI: `go run ./cmd/hadith-cli books`.
- Run TUI: `go run ./cmd/hadith-tui`.
- Run API: `ADDR=:8080 go run ./cmd/hadith-api`.
- Build all: `go build ./...`.
 - Serve web UI from `web/` at `/`. Fonts are local (`web/fonts/*`).

## Conventions and Gotchas
- Book names must match filenames (sans `.json`) for `Get` and URLs.
- Store is in-memory only; no DB or persistent index; all search is linear over `Store.All()`.
- Keep REST types stable: API returns `search.Result` shape from `internal/search`.
- When adding features, preserve upward `books` discovery and stable book ordering.
