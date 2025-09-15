package search

import (
    "sort"
    "strings"

    "github.com/nuzlilatief/hadith-go/internal/data"
)

// Result is a search hit with a simple score heuristic.
type Result struct {
    Hadith data.Hadith
    Score  int
}

// SimpleSearch performs a case-insensitive substring search across arab and id texts and book name.
// It returns up to limit results sorted by score then by book+number.
func SimpleSearch(all []data.Hadith, query string, limit int) []Result {
    q := strings.TrimSpace(query)
    if q == "" {
        return nil
    }
    ql := strings.ToLower(q)
    results := make([]Result, 0, limit)
    for _, h := range all {
        score := 0
        if strings.Contains(strings.ToLower(h.ID), ql) {
            score += 3
        }
        if strings.Contains(strings.ToLower(h.Arab), ql) {
            score += 2
        }
        if strings.Contains(strings.ToLower(h.Book), ql) {
            score += 1
        }
        if score > 0 {
            results = append(results, Result{Hadith: h, Score: score})
            if limit > 0 && len(results) >= limit {
                // do not break; we want to collect all to sort properly
            }
        }
    }
    sort.Slice(results, func(i, j int) bool {
        if results[i].Score != results[j].Score {
            return results[i].Score > results[j].Score
        }
        // tie-breaker: book then number
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

