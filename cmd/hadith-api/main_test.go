package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/nuzlilatief/hadith-go/internal/data"
	"github.com/nuzlilatief/hadith-go/internal/search"
)

// setupTestStore creates a temporary store for testing
func setupTestStore(t *testing.T) (*data.Store, string) {
	tmpDir := t.TempDir()

	// create books dir
	booksDir := filepath.Join(tmpDir, "books")
	if err := os.Mkdir(booksDir, 0755); err != nil {
		t.Fatalf("Failed to create books dir: %v", err)
	}

	testData := `[
		{"number": 1, "arab": "muhammad", "id": "prophet"},
		{"number": 2, "arab": "ahmad", "id": "messenger"}
	]`

	err := os.WriteFile(filepath.Join(booksDir, "test.json"), []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	store, err := data.NewStore(booksDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	return store, tmpDir
}

func TestHealthzEndpoint(t *testing.T) {
	store, root := setupTestStore(t)
	mux := setupRouter(store, root)

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

func TestSearchFuzzy(t *testing.T) {
	store, root := setupTestStore(t)
	mux := setupRouter(store, root)

	// Query with typo "muhammd" -> should match "muhammad" if fuzzy=true
	req := httptest.NewRequest("GET", "/search?q=muhammd&fuzzy=true", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var results []search.Result
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected results for fuzzy search, got none")
	} else {
		if results[0].Hadith.Number != 1 {
			t.Errorf("Expected hadith #1, got #%d", results[0].Hadith.Number)
		}
	}
}

func TestSearchNoFuzzy(t *testing.T) {
	store, root := setupTestStore(t)
	mux := setupRouter(store, root)

	// Query with typo "muhammd" -> should NOT match if fuzzy is default (false)
	req := httptest.NewRequest("GET", "/search?q=muhammd", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var results []search.Result
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

		if len(results) != 0 {

			t.Errorf("Expected 0 results for exact search with typo, got %d", len(results))

		}

	}

	

	func TestGetBookAndHadith(t *testing.T) {

		store, root := setupTestStore(t)

		mux := setupRouter(store, root)

	

		// Test /books

		req := httptest.NewRequest("GET", "/books", nil)

		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {

			t.Errorf("GET /books: expected 200, got %d", w.Code)

		}

	

		// Test /hadith/test/1

		req = httptest.NewRequest("GET", "/hadith/test/1", nil)

		w = httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {

			t.Errorf("GET /hadith/test/1: expected 200, got %d", w.Code)

		}

		

		// Test /hadith/test/999 (not found)

		req = httptest.NewRequest("GET", "/hadith/test/999", nil)

		w = httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {

			t.Errorf("GET /hadith/test/999: expected 404, got %d", w.Code)

		}

	}

	