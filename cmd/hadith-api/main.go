package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sort"

    "github.com/nuzlilatief/hadith-go/internal/data"
    "github.com/nuzlilatief/hadith-go/internal/search"
)

func main() {
    log.SetFlags(0)
    root := findBooksRoot()
    store, err := data.NewStore(filepath.Join(root, "books"))
    if err != nil {
        log.Fatalf("load books: %v", err)
    }
    mux := http.NewServeMux()
    // Static web UI (if web/ directory exists at repo root)
    staticDir := filepath.Join(root, "web")
    if st, err := os.Stat(staticDir); err == nil && st.IsDir() {
        // File server registered on "/". Explicit API routes below will take precedence.
        mux.Handle("/", http.FileServer(http.Dir(staticDir)))
    }
    // Serve OpenAPI spec at /openapi.yaml if present
    specPath := filepath.Join(root, "api", "openapi.yaml")
    if st, err := os.Stat(specPath); err == nil && !st.IsDir() {
        mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
            http.ServeFile(w, r, specPath)
        })
    }
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })
    mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, http.StatusOK, store.Books())
    })
    mux.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, http.StatusOK, map[string]int{"count": store.Count()})
    })
    mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
        q := r.URL.Query().Get("q")
        book := r.URL.Query().Get("book")
        // Back-compat: if page/page_size or offset not provided, honor legacy 'limit'.
        limitStr := r.URL.Query().Get("limit")
        pageStr := r.URL.Query().Get("page")
        pageSizeStr := r.URL.Query().Get("page_size")
        offsetStr := r.URL.Query().Get("offset")

        // Choose mode precedence: offset/limit > page/page_size > legacy limit
        useOffset := offsetStr != ""
        usePagination := !useOffset && (pageStr != "" || pageSizeStr != "")

        // Defaults and caps
        const defaultPageSize = 50
        const maxPageSize = 200

        // Build base corpus (optionally filtered by book before search).
        var corpus []data.Hadith
        if book == "" {
            corpus = store.All()
        } else {
            // Filter by book name exact match.
            for _, h := range store.All() {
                if h.Book == book {
                    corpus = append(corpus, h)
                }
            }
        }

        // Build results: if query is empty, browse corpus; else run search.
        var hits []search.Result
        if strings.TrimSpace(q) == "" {
            // Browse mode: list all hadiths in (optionally) selected book
            hits = make([]search.Result, 0, len(corpus))
            for _, h := range corpus {
                hits = append(hits, search.Result{Hadith: h, Score: 0})
            }
            sort.Slice(hits, func(i, j int) bool {
                if hits[i].Hadith.Book != hits[j].Hadith.Book {
                    return hits[i].Hadith.Book < hits[j].Hadith.Book
                }
                return hits[i].Hadith.Number < hits[j].Hadith.Number
            })
        } else {
            // Search without cap to allow pagination afterwards.
            hits = search.SimpleSearch(corpus, q, 0)
        }

        if useOffset {
            // Offset/limit style pagination for compatibility with some clients
            offset := 0
            if n, err := strconv.Atoi(offsetStr); err == nil && n >= 0 {
                offset = n
            }
            lim := defaultPageSize
            if limitStr != "" {
                if n, err := strconv.Atoi(limitStr); err == nil && n > 0 {
                    lim = n
                }
            }
            if lim > maxPageSize { lim = maxPageSize }
            total := len(hits)
            if offset > total { offset = total }
            end := offset + lim
            if end > total { end = total }
            w.Header().Set("X-Total-Count", strconv.Itoa(total))
            w.Header().Set("X-Offset", strconv.Itoa(offset))
            w.Header().Set("X-Limit", strconv.Itoa(lim))
            writeJSON(w, http.StatusOK, hits[offset:end])
            return
        }

        if usePagination {
            // Parse pagination params
            page := 1
            if pageStr != "" {
                if n, err := strconv.Atoi(pageStr); err == nil && n > 0 {
                    page = n
                }
            }
            pageSize := defaultPageSize
            if pageSizeStr != "" {
                if n, err := strconv.Atoi(pageSizeStr); err == nil && n > 0 {
                    pageSize = n
                }
            }
            if pageSize > maxPageSize { pageSize = maxPageSize }
            total := len(hits)
            start := (page - 1) * pageSize
            if start < 0 { start = 0 }
            if start > total { start = total }
            end := start + pageSize
            if end > total { end = total }

            // Pagination headers (body remains array for compatibility)
            w.Header().Set("X-Total-Count", strconv.Itoa(total))
            w.Header().Set("X-Page", strconv.Itoa(page))
            w.Header().Set("X-Page-Size", strconv.Itoa(pageSize))
            writeJSON(w, http.StatusOK, hits[start:end])
            return
        }

        // Legacy limit behavior
        limit := defaultPageSize
        if limitStr != "" {
            if n, err := strconv.Atoi(limitStr); err == nil && n > 0 {
                limit = n
            }
        }
        if limit > maxPageSize { limit = maxPageSize }
        if limit > 0 && len(hits) > limit {
            hits = hits[:limit]
        }
        writeJSON(w, http.StatusOK, hits)
    })
    mux.HandleFunc("/hadith/", func(w http.ResponseWriter, r *http.Request) {
        // /hadith/{book}/{number}
        parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/hadith/"), "/")
        if len(parts) != 2 {
            http.Error(w, "use /hadith/{book}/{number}", http.StatusBadRequest)
            return
        }
        num, err := strconv.Atoi(parts[1])
        if err != nil {
            http.Error(w, "invalid number", http.StatusBadRequest)
            return
        }
        h, ok := store.Get(parts[0], num)
        if !ok {
            http.NotFound(w, r)
            return
        }
        writeJSON(w, http.StatusOK, h)
    })

    addr := envOr("ADDR", ":8080")
    log.Printf("hadith API listening on %s", addr)
    log.Fatal(http.ListenAndServe(addr, cors(mux)))
}

func cors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(status)
    enc := json.NewEncoder(w)
    enc.SetIndent("", "  ")
    _ = enc.Encode(v)
}

func envOr(k, def string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return def
}

// findBooksRoot walks up from CWD to find a directory containing a "books" folder.
func findBooksRoot() string {
    dir, _ := os.Getwd()
    for i := 0; i < 5; i++ { // up to 5 levels
        if st, err := os.Stat(filepath.Join(dir, "books")); err == nil && st.IsDir() {
            return dir
        }
        parent := filepath.Dir(dir)
        if parent == dir {
            break
        }
        dir = parent
    }
    return "."
}
