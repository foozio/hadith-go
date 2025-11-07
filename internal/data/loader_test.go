package data

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStore(t *testing.T) {
	// Create temporary directory with test data
	tmpDir := t.TempDir()

	// Create test JSON files
	testData := `[
		{"number": 1, "arab": "Arabic text 1", "id": "Indonesian 1"},
		{"number": 2, "arab": "Arabic text 2", "id": "Indonesian 2"}
	]`

	err := os.WriteFile(filepath.Join(tmpDir, "test.json"), []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	store, err := NewStore(tmpDir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	if len(store.Books()) != 1 {
		t.Errorf("Expected 1 book, got %d", len(store.Books()))
	}

	if store.Books()[0] != "test" {
		t.Errorf("Expected book name 'test', got '%s'", store.Books()[0])
	}

	if store.Count() != 2 {
		t.Errorf("Expected 2 hadiths, got %d", store.Count())
	}
}

func TestStore_Get(t *testing.T) {
	tmpDir := t.TempDir()

	testData := `[
		{"number": 1, "arab": "Arabic 1", "id": "Indonesian 1"},
		{"number": 2, "arab": "Arabic 2", "id": "Indonesian 2"}
	]`

	err := os.WriteFile(filepath.Join(tmpDir, "book.json"), []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	store, err := NewStore(tmpDir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	// Test existing hadith
	h, ok := store.Get("book", 1)
	if !ok {
		t.Error("Expected to find hadith 1")
	}
	if h.Number != 1 || h.Book != "book" {
		t.Errorf("Wrong hadith data: %+v", h)
	}

	// Test non-existing hadith
	_, ok = store.Get("book", 99)
	if ok {
		t.Error("Should not find non-existing hadith")
	}

	// Test non-existing book
	_, ok = store.Get("nonexistent", 1)
	if ok {
		t.Error("Should not find hadith in non-existing book")
	}
}

func TestStore_All(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple books
	book1Data := `[{"number": 1, "arab": "A1", "id": "I1"}]`
	book2Data := `[{"number": 1, "arab": "A2", "id": "I2"}]`

	os.WriteFile(filepath.Join(tmpDir, "book1.json"), []byte(book1Data), 0644)
	os.WriteFile(filepath.Join(tmpDir, "book2.json"), []byte(book2Data), 0644)

	store, err := NewStore(tmpDir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	all := store.All()
	if len(all) != 2 {
		t.Errorf("Expected 2 hadiths, got %d", len(all))
	}

	// Check that books are included
	bookNames := make(map[string]bool)
	for _, h := range all {
		bookNames[h.Book] = true
	}

	if len(bookNames) != 2 {
		t.Errorf("Expected hadiths from 2 books, got %d", len(bookNames))
	}
}

func TestStore_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Invalid JSON
	err := os.WriteFile(filepath.Join(tmpDir, "invalid.json"), []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = NewStore(tmpDir)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestStore_NonArrayJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// JSON object instead of array
	err := os.WriteFile(filepath.Join(tmpDir, "object.json"), []byte(`{"key": "value"}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = NewStore(tmpDir)
	if err == nil {
		t.Error("Expected error for non-array JSON")
	}
}

func TestStore_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	store, err := NewStore(tmpDir)
	if err != nil {
		t.Fatalf("NewStore failed on empty directory: %v", err)
	}

	if len(store.Books()) != 0 {
		t.Errorf("Expected no books, got %d", len(store.Books()))
	}

	if store.Count() != 0 {
		t.Errorf("Expected no hadiths, got %d", store.Count())
	}
}
