package search

import (
    "runtime"
    "sort"
    "strings"
    "sync"

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

// ConcurrentSearch performs the same search as SimpleSearch but uses concurrency
// to parallelize the search across chunks of hadith data for better performance
// on large datasets.
func ConcurrentSearch(all []data.Hadith, query string, limit int) []Result {
    q := strings.TrimSpace(query)
    if q == "" {
        return nil
    }
    
    if len(all) == 0 {
        return nil
    }
    
    // For small datasets, use the simple version to avoid goroutine overhead
    if len(all) < 1000 {
        return SimpleSearch(all, query, limit)
    }
    
    ql := strings.ToLower(q)
    numWorkers := runtime.NumCPU()
    if numWorkers > len(all) {
        numWorkers = len(all)
    }
    
    chunkSize := len(all) / numWorkers
    if chunkSize == 0 {
        chunkSize = 1
    }
    
    // Channel to collect results from workers
    resultsChan := make(chan []Result, numWorkers)
    var wg sync.WaitGroup
    
    // Launch workers
    for i := 0; i < numWorkers; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if i == numWorkers-1 {
            end = len(all) // ensure last worker gets remaining items
        }
        if start >= len(all) {
            break
        }
        
        wg.Add(1)
        go func(chunk []data.Hadith) {
            defer wg.Done()
            chunkResults := searchChunk(chunk, ql)
            resultsChan <- chunkResults
        }(all[start:end])
    }
    
    // Close channel when all workers are done
    go func() {
        wg.Wait()
        close(resultsChan)
    }()
    
    // Collect all results
    var allResults []Result
    for chunkResults := range resultsChan {
        allResults = append(allResults, chunkResults...)
    }
    
    // Sort all results
    sort.Slice(allResults, func(i, j int) bool {
        if allResults[i].Score != allResults[j].Score {
            return allResults[i].Score > allResults[j].Score
        }
        // tie-breaker: book then number
        if allResults[i].Hadith.Book != allResults[j].Hadith.Book {
            return allResults[i].Hadith.Book < allResults[j].Hadith.Book
        }
        return allResults[i].Hadith.Number < allResults[j].Hadith.Number
    })
    
    if limit > 0 && len(allResults) > limit {
        return allResults[:limit]
    }
    return allResults
}

// searchChunk performs search on a chunk of hadith data
// This is the core search logic extracted for use by workers
func searchChunk(chunk []data.Hadith, queryLower string) []Result {
    var results []Result
    for _, h := range chunk {
        score := 0
        if strings.Contains(strings.ToLower(h.ID), queryLower) {
            score += 3
        }
        if strings.Contains(strings.ToLower(h.Arab), queryLower) {
            score += 2
        }
        if strings.Contains(strings.ToLower(h.Book), queryLower) {
            score += 1
        }
        if score > 0 {
            results = append(results, Result{Hadith: h, Score: score})
        }
    }
    return results
}

// Search is the recommended search function that automatically chooses
// the optimal search method based on dataset size and other factors.
// It provides the best performance while maintaining the same interface.
func Search(all []data.Hadith, query string, limit int) []Result {
    return ConcurrentSearch(all, query, limit)
}

