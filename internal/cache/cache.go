package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nuzlilatief/hadith-go/internal/search"
)

func cacheFileName(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:]) + ".json"
}

// FileCache implements a simple file-based cache for search results.
type FileCache struct {
	dir string
}

// NewFileCache creates a new FileCache using the specified directory.
func NewFileCache(dir string) *FileCache {
	_ = os.MkdirAll(dir, 0755)
	return &FileCache{dir: dir}
}

// Get retrieves results from the cache for a given key.
func (c *FileCache) Get(key string) ([]search.Result, bool) {
	path := filepath.Join(c.dir, cacheFileName(key))
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var results []search.Result
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, false
	}

	return results, true
}

// Put stores results in the cache for a given key.
func (c *FileCache) Put(key string, results []search.Result) error {
	path := filepath.Join(c.dir, cacheFileName(key))
	data, err := json.Marshal(results)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
