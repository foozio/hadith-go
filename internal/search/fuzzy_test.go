package search

import (
	"strings"
	"testing"

	"github.com/nuzlilatief/hadith-go/internal/data"
)

func TestFuzzySearch(t *testing.T) {
	hadiths := []data.Hadith{
		{Book: "test", Number: 1, Arab: "muhammad", ID: "prophet"},
		{Book: "test", Number: 2, Arab: "ahmad", ID: "messenger"},
		{Book: "test", Number: 3, Arab: "irrelevant", ID: "text"},
	}

	tests := []struct {
		query    string
		expected []int // Expected Hadith Numbers
		desc     string
	}{
		{
			query:    "muhammd", // Typo
			expected: []int{1},
			desc:     "fuzzy match with typo",
		},
		{
			query:    "prophit", // Typo
			expected: []int{1},
			desc:     "fuzzy match in ID",
		},
		{
			query:    "irrelevnt",
			expected: []int{3},
			desc:     "fuzzy match another word",
		},
		{
			query:    "",
			expected: nil,
			desc:     "empty query",
		},
		{
			query:    "xyzabc",
			expected: nil,
			desc:     "no matches",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			results := FuzzySearch(hadiths, tt.query, 10)

			if len(results) != len(tt.expected) {
				t.Fatalf("Expected %d results, got %d", len(tt.expected), len(results))
			}

			// Sort by Number for checking presence in expected list (if we had multiple expected)
			// But FuzzySearch returns sorted by score.
			// Let's just check if the expected IDs are present.
			// For single result expectations, simple check is fine.
			
			if len(results) > 0 && len(tt.expected) > 0 {
				if results[0].Hadith.Number != tt.expected[0] {
					t.Errorf("Expected top result %d, got %d", tt.expected[0], results[0].Hadith.Number)
				}
			}
		})
	}
}

func TestFuzzySearch_Sort(t *testing.T) {
	hadiths := []data.Hadith{
		{Book: "test", Number: 1, ID: "word match"},    // dist 0 (match), score 3*3=9
		{Book: "test", Number: 2, ID: "word mtch"},     // dist 1 (mtch), score 3*2=6
		{Book: "test", Number: 3, ID: "wrd mtch"},      // dist 1 (wrd) + dist 1 (mtch) = 3*2 + 3*2 = 12? No wait.
		// Query: "match"
		// 1. "match" -> dist 0 -> score 9
		// 2. "mtch" -> dist 1 -> score 3*(3-1)=6
		// 3. "wrd" (3), "mtch" (1) -> score 0 + 6 = 6. Tie with #2. 
		// Tie breaker: number. #2 < #3.
	}
	
	results := FuzzySearch(hadiths, "match", 10)
	
	if len(results) < 2 {
		t.Fatal("Expected results")
	}
	
	if results[0].Hadith.Number != 1 {
		t.Error("Expected exact match to be first")
	}
}

func TestLevenshtein_EdgeCases(t *testing.T) {
	if levenshtein("", "") != 0 {
		t.Error("Empty strings should have 0 distance")
	}
	if levenshtein("abc", "") != 3 {
		t.Error("Distance to empty string should be length")
	}
	if levenshtein("", "abc") != 3 {
		t.Error("Distance from empty string should be length")
	}
}

func TestFuzzySearch_RejectsLongQuery(t *testing.T) {
	hadiths := []data.Hadith{
		{Book: "test", Number: 1, Arab: "muhammad", ID: "prophet"},
	}

	longQuery := strings.Repeat("a", 129)
	results := FuzzySearch(hadiths, longQuery, 10)
	if results != nil && len(results) != 0 {
		t.Fatalf("expected no results for long query, got %d", len(results))
	}
}
