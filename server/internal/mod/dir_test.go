package mod_test

import (
	"io/fs"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"server/internal/mod"
	"server/internal/testutil/assert"
)

// mapFile builds an in-memory file with an explicit mtime.
func mapFile(data string, mtime time.Time) *fstest.MapFile {
	return &fstest.MapFile{Data: []byte(data), ModTime: mtime}
}

func TestDataDir_ReadFile(t *testing.T) {
	t.Parallel()

	type (
		sample struct {
			N    int            `json:"n"`
			Tags []string       `json:"tags"`
			Set  map[string]int `json:"set"`
		}

		other struct {
			X string `json:"x"`
		}
	)

	const file mod.DataFileName = "sample.json"

	t.Run("decode", func(t *testing.T) {
		t.Parallel()

		for name, tc := range map[string]struct {
			json string
			want sample
		}{
			"full": {
				`{"n":7,"tags":["a","b"],"set":{"x":1}}`,
				sample{N: 7, Tags: []string{"a", "b"}, Set: map[string]int{"x": 1}},
			},
			"empty object": {`{}`, sample{}},
			"nulls":        {`{"n":0,"tags":null,"set":null}`, sample{}},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				dd := mod.OpenDataDir(fstest.MapFS{string(file): mapFile(tc.json, time.Unix(1, 0))})

				got, cached, err := dd.ReadFile[sample](file)
				assert.NoError(t, err)
				assert.False(t, cached) // first read is always fresh
				assert.DeepEqual(t, tc.want, *got)

				got, cached, err = dd.ReadFile[sample](file)
				assert.NoError(t, err)
				assert.True(t, cached) // second read is cached
				assert.DeepEqual(t, tc.want, *got)
			})
		}
	})

	t.Run("re-reads when mtime changes but size stays", func(t *testing.T) {
		t.Parallel()

		mfs := fstest.MapFS{string(file): mapFile(`{"n":1}`, time.Unix(1, 0))}
		dd := mod.OpenDataDir(mfs)

		_, c1, _ := dd.ReadFile[sample](file)
		assert.False(t, c1)

		mfs[string(file)].ModTime = time.Unix(2, 0) // same bytes, newer mtime

		_, c2, err := dd.ReadFile[sample](file)
		assert.NoError(t, err)
		assert.False(t, c2) // mtime differs -> fresh read, not cache
	})

	t.Run("re-reads on content change without stale map keys", func(t *testing.T) {
		t.Parallel()

		mfs := fstest.MapFS{string(file): mapFile(`{"set":{"a":1,"b":2}}`, time.Unix(1, 0))}
		dd := mod.OpenDataDir(mfs)

		first, _, err1 := dd.ReadFile[sample](file)
		assert.NoError(t, err1)
		assert.DeepEqual(t, map[string]int{"a": 1, "b": 2}, first.Set)

		mfs[string(file)] = mapFile(`{"set":{"a":5}}`, time.Unix(2, 0)) // "b" removed

		second, cached, err2 := dd.ReadFile[sample](file)
		assert.NoError(t, err2)
		assert.False(t, cached)
		assert.DeepEqual(t, map[string]int{"a": 5}, second.Set) // decoded fresh: no stale "b"
	})

	t.Run("missing file", func(t *testing.T) {
		t.Parallel()

		dd := mod.OpenDataDir(fstest.MapFS{})

		data, cached, err := dd.ReadFile[sample](file)
		assert.Nil(t, data)
		assert.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrNotExist)
		assert.False(t, cached)
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		dd := mod.OpenDataDir(fstest.MapFS{string(file): mapFile(`{not json`, time.Unix(1, 0))})

		data, cached, err := dd.ReadFile[sample](file)
		assert.Nil(t, data)
		assert.Error(t, err)
		assert.False(t, cached)
	})

	t.Run("type mismatch against a cached file", func(t *testing.T) {
		t.Parallel()

		dd := mod.OpenDataDir(fstest.MapFS{string(file): mapFile(`{"n":1}`, time.Unix(1, 0))})

		_, _, err1 := dd.ReadFile[sample](file) // caches as sample
		assert.NoError(t, err1)

		_, cached, err2 := dd.ReadFile[other](file) // same file, different T
		assert.False(t, cached)
		assert.ErrorContains(t, err2, "type mismatch")
	})

	t.Run("concurrent reads are safe", func(t *testing.T) {
		t.Parallel()

		dd := mod.OpenDataDir(fstest.MapFS{string(file): mapFile(`{"n":9}`, time.Unix(1, 0))})

		const n = 50

		var (
			wg   sync.WaitGroup
			errs = make([]error, n)
		)

		wg.Add(n)

		for i := range n {
			go func(i int) {
				defer wg.Done()

				_, _, errs[i] = dd.ReadFile[sample](file)
			}(i)
		}

		wg.Wait()

		for _, err := range errs {
			assert.NoError(t, err)
		}
	})
}
