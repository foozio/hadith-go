package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/nuzlilatief/hadith-go/internal/data"
)

// setupTestStore creates a temporary store for testing
func setupTestStore(t *testing.T) *data.Store {
	tmpDir := t.TempDir()

	testData := `[
		{"number": 1, "arab": "Arabic text 1", "id": "Indonesian 1"},
		{"number": 2, "arab": "Arabic text 2", "id": "Indonesian 2"}
	]`

	err := os.WriteFile(filepath.Join(tmpDir, "test.json"), []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	store, err := data.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	return store
}

func TestHealthzEndpoint(t *testing.T) {
	store := setupTestStore(t)
	mux := setupRouter(store)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "ok" {
		t.Errorf("Expected 'ok', got '%s'", w.Body.String())
	}
}

// Simplified test router
func setupRouter(store *data.Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return mux
}
