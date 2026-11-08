package mod

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type (
	// DataDir provides cached, sandboxed read access to the mod's JSON output files inside a single data directory.
	//
	// All access goes through [os.Root], so a requested file name can never escape the directory via "..", an
	// absolute path, or a symlink.
	//
	// Decoded files are cached and re-read only when a file's size or mtime changes. The mod rewrites a file only
	// when its content actually changes (write-if-changed), so size+mtime is a reliable "has it changed?" signal.
	DataDir struct {
		root *os.Root

		cacheMu sync.Mutex                         // guards cache
		cache   map[DataFileName]dataDirCacheEntry // last decoded value per file
	}

	// dataDirCacheEntry is one file's last successfully decoded value together
	// with the size and mtime it was decoded from (used to detect staleness).
	dataDirCacheEntry struct {
		mtime, size int64 // file mtime (ms) and size the data was decoded from
		data        any   // decoded value; its concrete type is the T of the producing ReadFile call
	}
)

// NewDataDir returns a DataDir that reads from root, which must be an already opened sandboxed directory
// (see [os.OpenRoot]). The caller owns root's lifetime.
func NewDataDir(root *os.Root) *DataDir {
	return &DataDir{
		root:  root,
		cache: make(map[DataFileName]dataDirCacheEntry),
	}
}

// ReadFile reads, decodes and caches the JSON file fName as type T.
//
// The bool return reports whether the result came from cache (the file was unchanged since the last read) rather
// than a fresh disk read.
//
// The returned value shares backing storage (maps/slices) with the cache and must be treated as READ-ONLY;
// mutating it corrupts the cached copy.
func (dd *DataDir) ReadFile[T any](fName DataFileName) (zero T, _ bool /* from cache? */, _ error) {
	stat, statErr := dd.root.Stat(string(fName))
	if statErr != nil {
		return zero, false, statErr
	}

	size, mtime := stat.Size(), stat.ModTime().UnixMilli()

	{ // fast path - if the file size and mtime are unchanged, return the cached data
		dd.cacheMu.Lock()

		if entry, hit := dd.cache[fName]; hit && entry.size == size && entry.mtime == mtime {
			if data, ok := entry.data.(T); ok {
				dd.cacheMu.Unlock()

				return data, true, nil
			}

			dd.cacheMu.Unlock()

			return zero, false, errors.New("cached data type mismatch for " + string(fName))
		}

		dd.cacheMu.Unlock()
	}

	// slow path - read the file and update the cache
	f, fErr := dd.root.Open(string(fName))
	if fErr != nil {
		return zero, false, fErr
	}
	defer func() { _ = f.Close() }()

	var fresh T
	if err := json.NewDecoder(f).Decode(&fresh); err != nil {
		return zero, false, err
	}

	// commit only after a successful decoding
	dd.cacheMu.Lock()
	dd.cache[fName] = dataDirCacheEntry{mtime: mtime, size: size, data: fresh}
	dd.cacheMu.Unlock()

	return fresh, false, f.Close() // it's better to return the error from Close() than to silently ignore it
}
