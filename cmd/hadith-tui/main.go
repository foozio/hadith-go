package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"

    "github.com/nuzlilatief/hadith-go/internal/data"
    "github.com/nuzlilatief/hadith-go/internal/search"
)

// A minimal line-based TUI: type a query, see paginated results, navigate with n/p, open detail with o <index>, q to quit.
func main() {
    log.SetFlags(0)
    root := findBooksRoot()
    store, err := data.NewStore(filepath.Join(root, "books"))
    if err != nil {
        log.Fatalf("load books: %v", err)
    }
    fmt.Println("Hadith TUI — type query + Enter. Commands: :help, q to quit.")
    in := bufio.NewScanner(os.Stdin)
    page := 0
    const pageSize = 10
    var hits []search.Result
    // UI state
    truncWidth := 140
    showFull := false
    colorOn := true
    for {
        fmt.Print("query> ")
        if !in.Scan() {
            return
        }
        line := strings.TrimSpace(in.Text())
        if line == "q" || line == ":q" || line == "quit" {
            return
        }
        if line == "" {
            continue
        }
        // Extended commands: :help, :full, :short, :width N, :color on|off
        if strings.HasPrefix(line, ":") {
            cmd := strings.TrimSpace(strings.TrimPrefix(line, ":"))
            switch {
            case cmd == "help":
                printHelp()
            case cmd == "full":
                showFull = true
                fmt.Println("Full rows enabled (no truncation).")
                if len(hits) > 0 { renderPage(hits, page, pageSize, truncWidth, showFull, colorOn) }
            case cmd == "short":
                showFull = false
                fmt.Printf("Truncation enabled (width=%d).\n", truncWidth)
                if len(hits) > 0 { renderPage(hits, page, pageSize, truncWidth, showFull, colorOn) }
            case strings.HasPrefix(cmd, "width"):
                arg := strings.TrimSpace(strings.TrimPrefix(cmd, "width"))
                if arg == "" {
                    fmt.Printf("Current width: %d\n", truncWidth)
                    break
                }
                n := parseInt(strings.TrimSpace(arg))
                if n <= 0 {
                    fmt.Println("Width must be a positive integer.")
                    break
                }
                truncWidth = n
                showFull = false
                fmt.Printf("Width set to %d (truncation enabled).\n", truncWidth)
                if len(hits) > 0 { renderPage(hits, page, pageSize, truncWidth, showFull, colorOn) }
            case strings.HasPrefix(cmd, "color"):
                arg := strings.TrimSpace(strings.TrimPrefix(cmd, "color"))
                if arg == "on" || arg == "enable" { colorOn = true; fmt.Println("Color enabled.") }
                if arg == "off" || arg == "disable" { colorOn = false; fmt.Println("Color disabled.") }
                if len(hits) > 0 { renderPage(hits, page, pageSize, truncWidth, showFull, colorOn) }
            default:
                fmt.Println("Unknown command. Try :help")
            }
            continue
        }
        // Paging or open commands
        if line == "n" && len(hits) > 0 {
            if (page+1)*pageSize < len(hits) { page++ }
            renderPage(hits, page, pageSize, truncWidth, showFull, colorOn)
            continue
        }
        if line == "p" && len(hits) > 0 {
            if page > 0 { page-- }
            renderPage(hits, page, pageSize, truncWidth, showFull, colorOn)
            continue
        }
        if strings.HasPrefix(line, "o ") && len(hits) > 0 {
            idxStr := strings.TrimSpace(strings.TrimPrefix(line, "o "))
            // convert 1-based index on current page to absolute index
            idx := parseInt(idxStr) - 1
            abs := page*pageSize + idx
            if idx >= 0 && abs >= 0 && abs < len(hits) {
                h := hits[abs].Hadith
                // Full detail view always prints full text
                fmt.Printf("\n%s #%d\nID: %s\nAR: %s\n\n", h.Book, h.Number, h.ID, h.Arab)
            } else {
                fmt.Println("invalid index")
            }
            continue
        }
        // Otherwise treat the line as the new query
        hits = search.SimpleSearch(store.All(), line, 0)
        page = 0
        renderPage(hits, page, pageSize, truncWidth, showFull, colorOn)
    }
}

// ANSI colors
const (
    clrReset  = "\x1b[0m"
    clrDim    = "\x1b[2m"
    clrCyan   = "\x1b[36m"
    clrYellow = "\x1b[33m"
    clrGreen  = "\x1b[32m"
    clrBlue   = "\x1b[34m"
)

func colorize(on bool, color, s string) string {
    if !on { return s }
    return color + s + clrReset
}

func renderPage(hits []search.Result, page, size, truncWidth int, showFull, colorOn bool) {
    total := len(hits)
    if total == 0 {
        fmt.Println("No results. Try another query.")
        return
    }
    start := page * size
    if start >= total { page = 0; start = 0 }
    end := start + size
    if end > total { end = total }

    mode := "short"
    if showFull { mode = "full" }
    fmt.Printf("\nResults %d–%d of %d — (n)ext, (p)rev, (o N) open, (q)uit — mode:%s width:%d\n", start+1, end, total, mode, truncWidth)
    for i, r := range hits[start:end] {
        h := r.Hadith
        width := truncWidth
        if showFull { width = 0 }
        // Separator
        fmt.Printf("%s\n", colorize(colorOn, clrDim, "────────────────────────────────────────────────────────"))
        // Header line with book and number and score
        title := fmt.Sprintf("%2d. %s #%d", i+1, h.Book, h.Number)
        if r.Score > 0 { title += fmt.Sprintf("  score:%d", r.Score) }
        fmt.Println(colorize(colorOn, clrYellow, title))
        // Lines with labels
        fmt.Printf("    %s %s\n", colorize(colorOn, clrGreen, "ID:"), oneLine(h.ID, width))
        fmt.Printf("    %s %s\n", colorize(colorOn, clrCyan, "AR:"), oneLine(h.Arab, width))
    }
}

func oneLine(s string, width int) string {
    s = strings.ReplaceAll(s, "\n", " ")
    s = strings.TrimSpace(s)
    if width > 0 && len(s) > width {
        return s[:width] + "…"
    }
    return s
}

func parseInt(s string) int {
    n := 0
    for _, r := range s {
        if r < '0' || r > '9' {
            return -1
        }
        n = n*10 + int(r-'0')
    }
    return n
}

func printHelp() {
    fmt.Println("Commands:")
    fmt.Println("  :help           Show this help")
    fmt.Println("  :full           Show full text (no truncation)")
    fmt.Println("  :short          Enable truncation mode")
    fmt.Println("  :width N        Set truncation width to N characters")
    fmt.Println("  :color on|off   Toggle ANSI colors")
    fmt.Println("Keys:")
    fmt.Println("  n, p            Next/previous page")
    fmt.Println("  o N             Open full entry N on page")
    fmt.Println("  q               Quit")
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
