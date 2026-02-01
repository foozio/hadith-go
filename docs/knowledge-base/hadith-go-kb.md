# Knowledge Base: hadith-go

## 1. Project Overview
**hadith-go** is a high-performance, minimal Go application for browsing and searching Islamic Hadith collections. It is designed to be fast, portable, and easy to deploy, with no external runtime dependencies for its core functionality.

## 2. System Architecture
The application follows a **Monolithic, In-Memory** architecture.
*   **Data Loading:** Upon startup, all JSON files in `books/` are parsed and loaded into RAM (`internal/data.Store`). This ensures microsecond-level access times.
*   **Concurrency:** Search operations utilize Go routines to parallelize processing across CPU cores (`internal/search.ConcurrentSearch`), critical for maintaining performance as the dataset grows.
*   **Interfaces:** The core logic is shared across multiple interfaces (CLI, TUI, API, gRPC), ensuring consistent behavior.

## 3. Core Components

### 3.1 Data Layer (`internal/data`)
*   **Store**: The central in-memory repository. It is thread-safe (`sync.RWMutex`) for concurrent reads.
*   **Loader**: `NewStore` scans the `books/` directory. It uses `json.Decoder` to stream-parse files, reducing peak memory usage during load compared to `ioutil.ReadAll`.
*   **Models**:
    *   `Hadith`: Contains `Book` (name), `Number` (int), `Arab` (text), `ID` (Indonesian text).
    *   `Book`: A collection of Hadiths.

### 3.2 Search Layer (`internal/search`)
*   **Scoring Logic**: Results are ranked by relevance.
    *   **Weight 3**: Match in Indonesian text (primary language).
    *   **Weight 2**: Match in Arabic text.
    *   **Weight 1**: Match in Book name.
    *   **Tie-Breakers**: Book Name (A-Z) -> Hadith Number (Ascending).
*   **Algorithms**:
    *   **Simple Search**: Case-insensitive substring matching (`strings.Contains`). Optimized with concurrency for large datasets (>1000 items).
    *   **Fuzzy Search**: Implements **Levenshtein Distance**. It tokenizes text and matches the query against individual words. It allows 1 edit distance for short words (<=4 chars) and 2 for longer words.
    *   **Security**: Fuzzy search limits query length to 128 chars to prevent DoS attacks.

### 3.3 Cache Layer (`internal/cache`)
*   **FileCache**: Caches search results to disk (`.hadith_cache/`) to speed up repeated queries.
*   **Keying**: Uses SHA-256 hash of the query parameters (`q`, `book`, `fuzzy`) to generate safe filenames, preventing path traversal attacks.

## 4. Interfaces (`cmd/`)

### 4.1 REST API (`cmd/hadith-api`)
*   **Endpoints**:
    *   `GET /books`: List collections.
    *   `GET /count`: Total hadiths.
    *   `GET /search`: Search with query params (`q`, `book`, `fuzzy`, `limit`, `offset`, `page`).
    *   `GET /hadith/{book}/{number}`: Fetch single entry.
*   **Features**:
    *   Serves the static Web UI (`web/`) at root `/`.
    *   Implements middleware for Logging, Security Headers, and CORS.
    *   Supports OpenAPI spec serving (`/openapi.yaml`).

### 4.2 CLI (`cmd/hadith-cli`)
*   A traditional command-line tool.
*   **Commands**: `books`, `count`, `get <book> <num>`, `search <query>`.
*   Useful for scripting and quick server-side checks.

### 4.3 TUI (`cmd/hadith-tui`)
*   A terminal-based user interface for interactive browsing.
*   Features paginated search results and detail views without leaving the terminal.

### 4.4 gRPC (`cmd/hadith-grpc`)
*   **Status**: Optional build (requires `-tags grpc`).
*   **Definition**: `api/proto/hadith.proto`.
*   **Services**: `ListBooks`, `GetHadith`, `Search`.
*   **Use Case**: For typed, high-performance inter-service communication.

## 5. Development & Operations

### 5.1 Build System
The `Makefile` is the central control point:
*   `make build`: Compiles all binaries.
*   `make proto`: Generates Go code from `.proto` files (requires `protoc`).
*   `make grpc`: Builds the gRPC server.
*   `make release`: Cross-compiles binaries for Linux, Windows, and macOS (ARM64/AMD64).

### 5.2 Conductor
The `conductor/` directory serves as the project management hub, containing:
*   `tracks.md`: Registry of active development tracks.
*   `product.md`: Product definition and goals.
*   `tech-stack.md`: Technology decisions.

### 5.3 Dependencies
*   **Runtime**: Pure Standard Library (except gRPC).
*   **gRPC**: `google.golang.org/grpc`, `google.golang.org/protobuf`.

## 6. Data Flow
1.  **Startup**: `main()` calls `data.NewStore()`.
2.  **Ingest**: `Store` reads `books/*.json` -> parses JSON -> populates `byBook` map.
3.  **Request**: User invokes search (via HTTP/CLI).
4.  **Process**:
    *   System checks `FileCache`.
    *   If miss, calls `search.Search` (Concurrent) or `search.FuzzySearch`.
    *   Results are scored and sorted.
5.  **Response**: Results are serialized (JSON for API, Text for CLI/TUI) and returned.
