package mod

import (
	"encoding/json"
	"errors"
	"io/fs"
	"sync"
	"time"
)

type (
	// DataDir provides cached read access to the mod's JSON output files inside a single data directory.
	//
	// Decoded files are cached and re-read only when a file's size or mtime changes. The mod rewrites a file only
	// when its content actually changes (write-if-changed), so size+mtime is a reliable "has it changed?" signal.
	DataDir struct {
		fsys fs.FS // file source; prefer a sandboxed os.Root.FS()

		cacheMu sync.Mutex                         // guards cache
		cache   map[DataFileName]dataDirCacheEntry // last decoded value per file
	}

	// dataDirCacheEntry is one file's last successfully decoded value together with the size and mtime it was decoded
	// from (used to detect staleness).
	dataDirCacheEntry struct {
		mtime, size int64 // file mtime (ms) and size the data was decoded from
		data        any   // decoded value; its concrete type is the T of the producing ReadFile call
	}
)

// OpenDataDir returns a DataDir that reads files from fsys. For a sandboxed data directory pass an [os.Root]'s FS,
// e.g. OpenDataDir(root.FS()); the caller owns the underlying root's lifetime.
func OpenDataDir(fsys fs.FS) *DataDir {
	return &DataDir{
		fsys:  fsys,
		cache: make(map[DataFileName]dataDirCacheEntry),
	}
}

// ReadFile reads, decodes and caches the JSON file fName as type T.
//
// The bool return reports whether the result came from cache (the file was unchanged since the last read) rather
// than a fresh disk read. The returned time is the file's modification time (when the mod last wrote it), usable
// as a data-freshness signal.
//
// The returned value shares backing storage (maps/slices) with the cache and must be treated as READ-ONLY;
// mutating it corrupts the cached copy.
func (dd *DataDir) ReadFile[T any](fName DataFileName) (*T, bool /* from cache? */, time.Time /* mod time */, error) {
	info, statErr := fs.Stat(dd.fsys, string(fName))
	if statErr != nil {
		return nil, false, time.Time{}, statErr
	}

	modTime := info.ModTime()
	size, mtime := info.Size(), modTime.UnixMilli()

	{ // fast path - if the file size and mtime are unchanged, return the cached data
		dd.cacheMu.Lock()

		if entry, hit := dd.cache[fName]; hit && entry.size == size && entry.mtime == mtime {
			if data, ok := entry.data.(*T); ok {
				dd.cacheMu.Unlock()

				return data, true, modTime, nil
			}

			dd.cacheMu.Unlock()

			return nil, false, time.Time{}, errors.New("cached data type mismatch for " + string(fName))
		}

		dd.cacheMu.Unlock()
	}

	// slow path - read the file and update the cache
	f, fErr := dd.fsys.Open(string(fName))
	if fErr != nil {
		return nil, false, time.Time{}, fErr
	}
	defer func() { _ = f.Close() }()

	var data T
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, false, time.Time{}, err
	}

	// commit only after a successful decoding
	dd.cacheMu.Lock()
	dd.cache[fName] = dataDirCacheEntry{mtime: mtime, size: size, data: &data}
	dd.cacheMu.Unlock()

	return &data, false, modTime, f.Close() // it's better to return the error from Close() than to silently ignore it
}
