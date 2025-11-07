# Entity Relationship Diagram (ERD) for hadith-go

## Data Model Overview

The hadith-go application uses a simple, flat data model optimized for in-memory storage and fast search operations. All data is loaded from JSON files and stored in memory for performance.

## Entities

### Hadith
**Primary Entity**: Represents a single Hadith entry

Attributes:
- `book` (string): Book/collection name (e.g., "bukhari", "muslim")
- `number` (int): Sequential number within the book
- `arab` (string): Arabic text of the Hadith
- `id` (string): Indonesian translation

**Key**: Composite key of (book, number)

### Book
**Derived Entity**: Represents a collection of Hadiths

Attributes:
- `name` (string): Book name (derived from filename)
- `hadiths` ([]Hadith): Array of Hadith entries

**Relationship**: One Book contains many Hadiths (1:N)

## Storage Structure

### In-Memory Store (`internal/data.Store`)
```
Store {
  byBook: map[string][]Hadith  // book name -> array of hadiths
  books: []string              // sorted list of book names
  rootDir: string              // source directory path
}
```

### JSON File Format
Each book is stored as a separate JSON file:
```json
[
  {
    "number": 1,
    "arab": "Arabic text...",
    "id": "Indonesian translation..."
  },
  {
    "number": 2,
    "arab": "Another Arabic text...",
    "id": "Another Indonesian translation..."
  }
]
```

## Relationships

```
Book (1) ────contains──── (N) Hadith
  │                              │
  └── name                       ├── book (FK)
  └── hadiths[]                  ├── number (PK)
                                 ├── arab
                                 └── id
```

## Search Results

### SearchResult
**Runtime Entity**: Represents a search match

Attributes:
- `hadith` (Hadith): The matching Hadith
- `score` (int): Relevance score (higher = better match)

## Data Flow

1. **Load Phase**:
   - Read JSON files from `books/` directory
   - Parse JSON arrays into Hadith structs
   - Group by book name in `byBook` map
   - Maintain sorted book list

2. **Query Phase**:
   - For book-specific queries: Filter `byBook[book]`
   - For global queries: Concatenate all Hadiths
   - Apply search logic with scoring
   - Sort results by score, then book, then number

3. **Response Phase**:
   - Paginate results based on request parameters
   - Return JSON/structured data to client

## Constraints

- **Uniqueness**: Each (book, number) combination must be unique
- **Referential Integrity**: All Hadiths must belong to a valid book
- **Data Consistency**: JSON files must contain valid UTF-8 encoded arrays
- **Performance**: All data fits in memory for fast access

## Indexing Strategy

Currently uses simple in-memory arrays with linear search. Future enhancements could include:
- Inverted index for text search
- Trie structures for prefix matching
- Bloom filters for fast existence checks
- Sorted arrays for binary search by number

## Schema Evolution

The current schema is stable and minimal. Future versions might add:
- Metadata fields (authenticity grade, narrator chains)
- Cross-references between Hadiths
- Timestamp/version information
- User annotations (requires persistence layer)</content>
<parameter name="filePath">ERD.md