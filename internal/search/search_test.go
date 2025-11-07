package search

import (
	"sort"
	"testing"

	"github.com/nuzlilatief/hadith-go/internal/data"
)

func TestSimpleSearch(t *testing.T) {
	hadiths := []data.Hadith{
		{Book: "test", Number: 1, Arab: "arabic one", ID: "indonesian satu"},
		{Book: "test", Number: 2, Arab: "arabic two", ID: "indonesian dua"},
		{Book: "test", Number: 3, Arab: "arabic three", ID: "indonesian tiga"},
	}

	tests := []struct {
		query    string
		expected []Result
		desc     string
	}{
		{
			query: "indonesian",
			expected: []Result{
				{Hadith: hadiths[0], Score: 3},
				{Hadith: hadiths[1], Score: 3},
				{Hadith: hadiths[2], Score: 3},
			},
			desc: "search in ID field",
		},
		{
			query: "arabic",
			expected: []Result{
				{Hadith: hadiths[0], Score: 2},
				{Hadith: hadiths[1], Score: 2},
				{Hadith: hadiths[2], Score: 2},
			},
			desc: "search in Arab field",
		},
		{
			query: "test",
			expected: []Result{
				{Hadith: hadiths[0], Score: 1},
				{Hadith: hadiths[1], Score: 1},
				{Hadith: hadiths[2], Score: 1},
			},
			desc: "search in Book field",
		},
		{
			query: "satu",
			expected: []Result{
				{Hadith: hadiths[0], Score: 3},
			},
			desc: "search specific word",
		},
		{
			query:    "",
			expected: nil,
			desc:     "empty query",
		},
		{
			query:    "nonexistent",
			expected: nil,
			desc:     "no matches",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			results := SimpleSearch(hadiths, tt.query, 0)

			// Sort results for consistent comparison
			sort.Slice(results, func(i, j int) bool {
				if results[i].Score != results[j].Score {
					return results[i].Score > results[j].Score
				}
				return results[i].Hadith.Number < results[j].Hadith.Number
			})

			if len(results) != len(tt.expected) {
				t.Errorf("Expected %d results, got %d", len(tt.expected), len(results))
				return
			}

			for i, expected := range tt.expected {
				if results[i].Hadith.Number != expected.Hadith.Number ||
					results[i].Score != expected.Score {
					t.Errorf("Result %d: expected %+v, got %+v", i, expected, results[i])
				}
			}
		})
	}
}

func TestSimpleSearch_Limit(t *testing.T) {
	hadiths := []data.Hadith{
		{Book: "test", Number: 1, Arab: "match", ID: "match"},
		{Book: "test", Number: 2, Arab: "match", ID: "match"},
		{Book: "test", Number: 3, Arab: "match", ID: "match"},
	}

	results := SimpleSearch(hadiths, "match", 2)
	if len(results) != 2 {
		t.Errorf("Expected 2 results with limit, got %d", len(results))
	}
}

func TestConcurrentSearch(t *testing.T) {
	hadiths := []data.Hadith{
		{Book: "test", Number: 1, Arab: "arabic one", ID: "indonesian satu"},
		{Book: "test", Number: 2, Arab: "arabic two", ID: "indonesian dua"},
		{Book: "test", Number: 3, Arab: "arabic three", ID: "indonesian tiga"},
	}

	// Test with small dataset (should use simple search)
	results := ConcurrentSearch(hadiths, "indonesian", 0)
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Verify results are sorted by score
	for i := 0; i < len(results)-1; i++ {
		if results[i].Score < results[i+1].Score {
			t.Error("Results not sorted by score descending")
		}
	}
}

func TestSearch(t *testing.T) {
	hadiths := []data.Hadith{
		{Book: "test", Number: 1, Arab: "arabic", ID: "indonesian"},
	}

	// Test that Search uses ConcurrentSearch
	results := Search(hadiths, "indonesian", 0)
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0].Hadith.Number != 1 {
		t.Errorf("Wrong hadith returned")
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	hadiths := []data.Hadith{
		{Book: "test", Number: 1, Arab: "text", ID: "text"},
	}

	results := Search(hadiths, "", 0)
	if results != nil {
		t.Errorf("Expected nil for empty query, got %v", results)
	}
}

func TestSearch_EmptyHadiths(t *testing.T) {
	results := Search([]data.Hadith{}, "query", 0)
	if results != nil {
		t.Errorf("Expected nil for empty hadiths, got %v", results)
	}
}
