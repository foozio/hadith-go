# Changelog

All notable changes to this project will be documented in this file.

## [v1.1.1] - 2026-01-31

### Added
- **Search:** Implemented fuzzy matching using Levenshtein distance to handle typos in search queries (`-fuzzy` flag in CLI, `fuzzy=true` in API).
- **Performance:** Added file-based persistent caching for search results in `.hadith_cache/`.
- **API:** Added `X-Cache` header to search responses to indicate cache hits/misses.

### Documentation
- Updated README with new CLI examples and API parameters.
- Updated OpenAPI specification.
