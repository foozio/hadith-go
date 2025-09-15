package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"

    "github.com/nuzlilatief/hadith-go/internal/data"
    "github.com/nuzlilatief/hadith-go/internal/search"
)

func usage() {
    fmt.Fprintf(os.Stderr, "hadith-cli usage:\n")
    fmt.Fprintf(os.Stderr, "  hadith-cli books\n")
    fmt.Fprintf(os.Stderr, "  hadith-cli count\n")
    fmt.Fprintf(os.Stderr, "  hadith-cli get <book> <number>\n")
    fmt.Fprintf(os.Stderr, "  hadith-cli search [-limit N] <query>\n")
}

func main() {
    log.SetFlags(0)
    if len(os.Args) < 2 {
        usage()
        os.Exit(2)
    }
    cmd := os.Args[1]
    // locate books directory relative to working directory by default
    root := findBooksRoot()
    store, err := data.NewStore(filepath.Join(root, "books"))
    if err != nil {
        log.Fatalf("load books: %v", err)
    }
    switch cmd {
    case "books":
        for _, b := range store.Books() {
            fmt.Println(b)
        }
    case "count":
        fmt.Println(store.Count())
    case "get":
        if len(os.Args) < 4 {
            usage()
            os.Exit(2)
        }
        book := os.Args[2]
        n, err := strconv.Atoi(os.Args[3])
        if err != nil {
            log.Fatalf("invalid number: %v", err)
        }
        h, ok := store.Get(book, n)
        if !ok {
            log.Fatalf("not found: %s #%d", book, n)
        }
        enc := json.NewEncoder(os.Stdout)
        enc.SetIndent("", "  ")
        _ = enc.Encode(h)
    case "search":
        fs := flag.NewFlagSet("search", flag.ExitOnError)
        limit := fs.Int("limit", 20, "max results")
        _ = fs.Parse(os.Args[2:])
        q := strings.Join(fs.Args(), " ")
        results := search.SimpleSearch(store.All(), q, *limit)
        for _, r := range results {
            fmt.Printf("%s #%d [score %d]\n", r.Hadith.Book, r.Hadith.Number, r.Score)
            // print Indonesian translation first for readability
            fmt.Printf("ID: %s\n", oneLine(r.Hadith.ID))
            fmt.Printf("AR: %s\n\n", oneLine(r.Hadith.Arab))
        }
    default:
        usage()
        os.Exit(2)
    }
}

func oneLine(s string) string {
    s = strings.ReplaceAll(s, "\n", " ")
    s = strings.TrimSpace(s)
    if len(s) > 240 {
        return s[:240] + "…"
    }
    return s
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

