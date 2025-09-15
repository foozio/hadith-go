package data

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "sync"
)

// Store loads and holds all hadiths from the books directory.
type Store struct {
    mu      sync.RWMutex
    byBook  map[string][]Hadith
    books   []string
    rootDir string
}

// NewStore loads JSON files from booksDir. Filenames (without .json) are used as book names.
func NewStore(booksDir string) (*Store, error) {
    st := &Store{byBook: map[string][]Hadith{}, rootDir: booksDir}
    entries, err := os.ReadDir(booksDir)
    if err != nil {
        return nil, fmt.Errorf("read books dir: %w", err)
    }
    for _, e := range entries {
        if e.IsDir() {
            continue
        }
        name := e.Name()
        if !strings.HasSuffix(strings.ToLower(name), ".json") {
            continue
        }
        book := strings.TrimSuffix(name, filepath.Ext(name))
        if err := st.loadBook(filepath.Join(booksDir, name), book); err != nil {
            return nil, err
        }
    }
    // stable ordering of book names
    for b := range st.byBook {
        st.books = append(st.books, b)
    }
    sort.Strings(st.books)
    return st, nil
}

func (s *Store) loadBook(path, book string) error {
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("open %s: %w", path, err)
    }
    defer f.Close()
    dec := json.NewDecoder(f)
    // File is an array of objects with fields number, arab, id
    tok, err := dec.Token()
    if err != nil {
        return fmt.Errorf("decode %s: %w", path, err)
    }
    if delim, ok := tok.(json.Delim); !ok || delim != '[' {
        return fmt.Errorf("%s: expected JSON array", path)
    }
    var hadiths []Hadith
    for dec.More() {
        var raw struct {
            Number int    `json:"number"`
            Arab   string `json:"arab"`
            ID     string `json:"id"`
        }
        if err := dec.Decode(&raw); err != nil {
            if errors.Is(err, io.EOF) {
                break
            }
            return fmt.Errorf("decode hadith in %s: %w", path, err)
        }
        hadiths = append(hadiths, Hadith{
            Book:   book,
            Number: raw.Number,
            Arab:   raw.Arab,
            ID:     raw.ID,
        })
    }
    s.byBook[book] = hadiths
    return nil
}

// Books returns book names in stable order.
func (s *Store) Books() []string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    out := make([]string, len(s.books))
    copy(out, s.books)
    return out
}

// Count returns total hadith count across all books.
func (s *Store) Count() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    n := 0
    for _, v := range s.byBook {
        n += len(v)
    }
    return n
}

// Get returns a hadith by book and number. The number refers to the "number" field in the JSON.
func (s *Store) Get(book string, number int) (Hadith, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    list := s.byBook[book]
    for _, h := range list {
        if h.Number == number {
            return h, true
        }
    }
    return Hadith{}, false
}

// All returns all hadiths across all books.
func (s *Store) All() []Hadith {
    s.mu.RLock()
    defer s.mu.RUnlock()
    var out []Hadith
    for _, b := range s.books {
        out = append(out, s.byBook[b]...)
    }
    return out
}

