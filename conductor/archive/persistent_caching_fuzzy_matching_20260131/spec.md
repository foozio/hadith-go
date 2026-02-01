# Specification: Persistent Caching and Fuzzy Matching

## Overview
This track aims to enhance the `hadith-go` search engine by introducing fuzzy matching to handle typos and persistent caching to speed up recurring search queries.

## Requirements
- **Fuzzy Matching:** Integrate a fuzzy string matching algorithm (e.g., Levenshtein distance) into `internal/search`.
- **Persistent Caching:** Implement a file-based or lightweight DB-based cache for search results.
- **Performance:** Ensure that cache lookups are faster than re-running searches on the full in-memory store.
- **Configurability:** Allow enabling/disabling fuzzy matching and caching via configuration or flags.

## Architecture Changes
- `internal/search`: Add a new `FuzzySearch` function or enhance existing `Search`.
- `internal/data`: Add a caching layer between the search interfaces and the data store.
