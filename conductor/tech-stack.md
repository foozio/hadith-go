# Technology Stack

## Core Technologies
- **Language:** Go 1.21+
- **Frontend:** Vanilla JavaScript, HTML, CSS
- **Data Format:** JSON

## Backend & APIs
- **REST API:** Go `net/http` standard library.
- **gRPC:** `google.golang.org/grpc` for high-performance RPC.
- **Documentation:** 
    - OpenAPI 3.0 for REST.
    - Protocol Buffers (v3) for gRPC.

## Storage & Search
- **Storage:** In-memory data store using Go maps and slices, loaded from `books/*.json`.
- **Search Engine:** Custom substring search logic with simple scoring implemented in `internal/search`.

## Infrastructure & Tools
- **Build System:** `Makefile` for managing builds, tests, and proto generation.
- **Deployment:** Go binary with static asset embedding/serving.
- **Architecture:** Monolithic core with multiple interface-specific entry points (CLI, TUI, API).
