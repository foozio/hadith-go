package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nuzlilatief/hadith-go/internal/data"
	"github.com/nuzlilatief/hadith-go/internal/search"
)

func TestFileCache(t *testing.T) {
	tmpDir := t.TempDir()
	c := NewFileCache(tmpDir)

	results := []search.Result{
		{Hadith: data.Hadith{Book: "test", Number: 1, Arab: "A", ID: "I"}, Score: 10},
	}
	key := "test-key"

	// Get non-existent
	got, ok := c.Get(key)
	if ok || got != nil {
		t.Error("Expected nil/false for non-existent key")
	}

	// Put
	err := c.Put(key, results)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Get existent
	got, ok = c.Get(key)
	if !ok {
		t.Fatal("Expected to find key")
	}
	if len(got) != 1 || got[0].Hadith.Number != 1 {
		t.Errorf("Got wrong results: %+v", got)
	}

	// Persistency check: new cache instance with same dir
	c2 := NewFileCache(tmpDir)
	got2, ok2 := c2.Get(key)
	if !ok2 || len(got2) != 1 {
		t.Error("Cache should persist to disk")
	}
}

func TestFileCache_RejectsPathTraversalKey(t *testing.T) {
	tmpDir := t.TempDir()
	c := NewFileCache(tmpDir)

	results := []search.Result{
		{Hadith: data.Hadith{Book: "test", Number: 1, Arab: "A", ID: "I"}, Score: 10},
	}
	key := "../evil"

	if err := c.Put(key, results); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Ensure we did not write outside the cache directory.
	outsidePath := filepath.Join(tmpDir, "..", "evil.json")
	if _, err := os.Stat(outsidePath); err == nil {
		t.Fatalf("expected no file outside cache dir, found %s", outsidePath)
	}
}
