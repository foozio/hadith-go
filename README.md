# hadith-go

[![CI](https://github.com/nuzlilatief/hadith-go/actions/workflows/ci.yml/badge.svg)](https://github.com/nuzlilatief/hadith-go/actions/workflows/ci.yml)

A fast, minimal Go project to browse and search hadith collections stored as JSON. It ships a CLI, a TUI, a REST API with a polished web UI, and optional gRPC.

- Data: `books/*.json` arrays with `{ number, arab, id }`
- Storage: in-memory `internal/data.Store`
- Search: case-insensitive substring with simple scoring in `internal/search`
- Interfaces: CLI, TUI, REST (`/books`, `/count`, `/search`, `/hadith/{book}/{number}`), optional gRPC

## Quick Start

- Run the REST API and web UI

```
ADDR=:8080 go run ./cmd/hadith-api
# open http://localhost:8080/
```

- CLI examples

```
go run ./cmd/hadith-cli books

go run ./cmd/hadith-cli count

go run ./cmd/hadith-cli get bukhari 1

go run ./cmd/hadith-cli search -limit 10 niat
```

- TUI

```
go run ./cmd/hadith-tui
```

## REST API

- `GET /healthz` → `ok`
- `GET /books` → `[]string`
- `GET /count` → `{ "count": N }`
- `GET /search`
  - Query params:
    - `q`: search string (optional). If empty and `book` is set, returns all entries in the book (browse mode).
    - `book`: exact book name (filename without `.json`). Optional filter; also enables browse mode when `q` is empty.
    - Pagination (three compatible modes):
      - Offset/limit: `offset` (>=0), `limit` (>0, default 50, max 200). Precedence when `offset` is present.
      - Page-based: `page` (>=1), `page_size` (>0, default 50, max 200).
      - Legacy: `limit` only (applied after search when neither `offset` nor `page/page_size` is provided).
  - Response: JSON array of results, each like `{ hadith, score }` (score omitted in browse mode).
  - Headers (when paginated):
    - Offset/limit: `X-Total-Count`, `X-Offset`, `X-Limit`
    - Page-based: `X-Total-Count`, `X-Page`, `X-Page-Size`

- `GET /hadith/{book}/{number}` → hadith entry or 404

## Web UI

- Served statically from `web/` by the API (at `/`).
- Features: search, book filter, server-side pagination, browse-by-book, responsive layout.
- Typography: local Inter (Latin) and Amiri + Noto Naskh Arabic (Arabic). See `web/fonts/README.md` for bundling instructions. The UI falls back to system fonts if local files are missing.

## Optional gRPC

- Proto at `api/proto/hadith.proto` (Go package `api/gen/go/hadithpb`).
- Generate and build:

```
make proto
make grpc
```

## API Schema

- OpenAPI spec: `api/openapi.yaml`
- View options:
  - Redocly CLI: `npx @redocly/cli preview-docs api/openapi.yaml` (or install globally)
  - Swagger UI (Docker): `docker run -p 8081:8080 -e SWAGGER_JSON=/foo/openapi.yaml -v $(pwd)/api/openapi.yaml:/foo/openapi.yaml swaggerapi/swagger-ui`
  - VS Code: use an OpenAPI extension to preview the file directly

## Development

- Go 1.21+ recommended.
- Useful targets:

```
make run-api
make run-cli
make run-tui
make build
```

- Project layout: see `agents.yml` and `internal/*` packages for data and search internals.

## License

MIT — see `LICENSE`.

## Notes on Data

The hadith JSON data is included under `books/`. Please verify data licensing and provenance according to your intended use. The application code is MIT-licensed; dataset licensing may differ.
