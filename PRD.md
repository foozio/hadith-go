# Product Requirements Document (PRD) for hadith-go

## Overview
hadith-go is a minimal, fast Go application for browsing and searching Islamic Hadith collections. It provides multiple interfaces (REST API, CLI, TUI, gRPC) to access Hadith data stored in JSON format.

## Product Vision
To create an accessible, performant tool for scholars, students, and developers to explore Hadith collections with modern interfaces while maintaining simplicity and speed.

## Target Users
- Islamic scholars and students
- Developers building Islamic applications
- General users interested in Hadith study
- API consumers needing programmatic access

## Core Features

### Data Management
- Load Hadith collections from JSON files
- Support multiple book collections (Bukhari, Muslim, etc.)
- In-memory storage for fast access
- Thread-safe data access

### Search Functionality
- Case-insensitive substring search across:
  - Indonesian translations (ID)
  - Arabic text (Arab)
  - Book names
- Scoring system prioritizing Indonesian > Arabic > Book name matches
- Concurrent search for large datasets
- Browse mode (list all Hadiths in a book)

### Interfaces

#### REST API
- `/healthz` - Health check
- `/books` - List available books
- `/count` - Total Hadith count
- `/search` - Search with pagination support
- `/hadith/{book}/{number}` - Get specific Hadith
- Multiple pagination modes (offset/limit, page-based, legacy)
- CORS support for web clients
- OpenAPI 3.0 specification

#### CLI
- `books` - List books
- `count` - Show total count
- `get <book> <number>` - Get specific Hadith
- `search [-limit N] <query>` - Search with results

#### TUI (Terminal User Interface)
- Interactive search interface
- Paginated results display
- Color support
- Configurable truncation width
- Navigation commands (next/prev page, open details)

#### gRPC (Optional)
- Protocol buffer definitions
- ListBooks, GetHadith, Search RPC methods
- Build tag controlled compilation

### Web UI
- Responsive single-page application
- Search form with book filtering
- Server-side pagination
- Clean, accessible design
- Local font support (Inter, Amiri, Noto Naskh)

## Technical Requirements

### Performance
- Fast startup (load data into memory)
- Concurrent search processing
- Efficient pagination
- Low memory footprint

### Compatibility
- Go 1.21+ support
- Cross-platform (Linux, macOS, Windows)
- Standard library dependencies only (except gRPC)

### Data Format
- JSON array format: `[{number, arab, id}]`
- Book names derived from filenames
- UTF-8 encoding support

## Non-Functional Requirements

### Security
- Input validation and sanitization
- Safe file operations
- No external dependencies for core functionality
- CORS configuration for web access

### Reliability
- Graceful error handling
- Thread-safe operations
- Data integrity checks

### Maintainability
- Clean, idiomatic Go code
- Modular architecture
- Comprehensive documentation
- Build system with Makefile

## Future Enhancements
- Advanced search features (regex, fuzzy matching)
- Full-text indexing
- Caching layer
- Authentication/authorization
- Export functionality
- Mobile app interfaces
- Multi-language support

## Success Metrics
- Fast search response times (<100ms for typical queries)
- Accurate search results with good scoring
- Intuitive interfaces across all access methods
- Reliable data loading and serving
- Active community adoption and contributions</content>
<parameter name="filePath">PRD.md