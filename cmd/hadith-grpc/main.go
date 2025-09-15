//go:build grpc

package main

import (
    "context"
    "log"
    "net"
    "os"
    "path/filepath"

    "google.golang.org/grpc"

    "github.com/nuzlilatief/hadith-go/api/gen/go/hadithpb"
    "github.com/nuzlilatief/hadith-go/internal/data"
)

type server struct{
    hadithpb.UnimplementedHadithServiceServer
    store *data.Store
}

func (s *server) ListBooks(ctx context.Context, _ *hadithpb.ListBooksRequest) (*hadithpb.ListBooksResponse, error) {
    return &hadithpb.ListBooksResponse{Books: s.store.Books()}, nil
}

func (s *server) GetHadith(ctx context.Context, req *hadithpb.GetHadithRequest) (*hadithpb.GetHadithResponse, error) {
    h, ok := s.store.Get(req.GetBook(), int(req.GetNumber()))
    if !ok {
        return &hadithpb.GetHadithResponse{}, nil
    }
    return &hadithpb.GetHadithResponse{Hadith: &hadithpb.Hadith{Book: h.Book, Number: int32(h.Number), Arab: h.Arab, Id: h.ID}}, nil
}

func (s *server) Search(ctx context.Context, req *hadithpb.SearchRequest) (*hadithpb.SearchResponse, error) {
    // Minimal inline search to avoid import cycle.
    var out []*hadithpb.Hadith
    q := req.GetQuery()
    for _, h := range s.store.All() {
        if q == "" || containsFold(h.ID, q) || containsFold(h.Arab, q) || containsFold(h.Book, q) {
            out = append(out, &hadithpb.Hadith{Book: h.Book, Number: int32(h.Number), Arab: h.Arab, Id: h.ID})
            if req.GetLimit() > 0 && int32(len(out)) >= req.GetLimit() {
                break
            }
        }
    }
    return &hadithpb.SearchResponse{Results: out}, nil
}

func containsFold(hay, needle string) bool {
    // simple case-insensitive containment
    // avoid importing strings to keep this file minimal under build tag
    H := []rune(hay)
    N := []rune(needle)
    for i := range H {
        if i+len(N) > len(H) { return false }
        ok := true
        for j := range N {
            r1 := H[i+j]
            r2 := N[j]
            if r1 >= 'A' && r1 <= 'Z' { r1 = r1 - 'A' + 'a' }
            if r2 >= 'A' && r2 <= 'Z' { r2 = r2 - 'A' + 'a' }
            if r1 != r2 { ok = false; break }
        }
        if ok { return true }
    }
    return false
}

func main() {
    root := findBooksRoot()
    store, err := data.NewStore(filepath.Join(root, "books"))
    if err != nil {
        log.Fatalf("load books: %v", err)
    }
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("listen: %v", err)
    }
    s := grpc.NewServer()
    hadithpb.RegisterHadithServiceServer(s, &server{store: store})
    log.Printf("hadith gRPC listening on %s", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatal(err)
    }
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

