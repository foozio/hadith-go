package search

import (
	"sort"
	"strings"

	"github.com/nuzlilatief/hadith-go/internal/data"
)

// levenshtein calculates the Levenshtein distance between two strings.
func levenshtein(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	n, m := len(r1), len(r2)

	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}

	row := make([]int, m+1)
	for i := 0; i <= m; i++ {
		row[i] = i
	}

	for i := 1; i <= n; i++ {
		prev := i
		var val int
		for j := 1; j <= m; j++ {
			if r1[i-1] == r2[j-1] {
				val = row[j-1]
			} else {
				val = min(row[j-1]+1, min(prev+1, row[j]+1))
			}
			row[j-1] = prev
			prev = val
		}
		row[m] = prev
	}
	return row[m]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// FuzzySearch performs a search using Levenshtein distance on tokens.
// It matches the query against individual words in the hadith text.
func FuzzySearch(all []data.Hadith, query string, limit int) []Result {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil
	}
	ql := strings.ToLower(q)
	
	// Threshold: allow 1 error for short words (len <= 4), 2 for longer.
	threshold := 2
	if len(ql) <= 4 {
		threshold = 1
	}

	results := make([]Result, 0, limit)

	for _, h := range all {
		score := 0
		
		// Helper to check tokens
		checkTokens := func(text string, weight int) int {
			localScore := 0
			tokens := strings.Fields(strings.ToLower(text))
			for _, token := range tokens {
				// Clean token (remove punctuation) - basic version
				cleanToken := strings.Trim(token, ".,:;\"'()[]{}!?")
				
				dist := levenshtein(ql, cleanToken)
				if dist == 0 {
					localScore += weight * 3 // Exact match bonus
				} else if dist <= threshold {
					// Score inversely proportional to distance
					localScore += weight * (3 - dist) 
				}
			}
			return localScore
		}

		score += checkTokens(h.ID, 3)
		score += checkTokens(h.Arab, 2)
		score += checkTokens(h.Book, 1)

		if score > 0 {
			results = append(results, Result{Hadith: h, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		if results[i].Hadith.Book != results[j].Hadith.Book {
			return results[i].Hadith.Book < results[j].Hadith.Book
		}
		return results[i].Hadith.Number < results[j].Hadith.Number
	})

	if limit > 0 && len(results) > limit {
		return results[:limit]
	}
	return results
}
