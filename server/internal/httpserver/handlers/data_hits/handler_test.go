package data_hits_test

import (
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"testing/fstest"
	"time"

	"server/internal/httpserver/handlers/data_hits"
	"server/internal/mod"
	"server/internal/testutil/assert"
)

var testMtime = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

func TestNew(t *testing.T) {
	t.Parallel()

	for name, tt := range map[string]struct {
		giveJSON string // data file content; empty means the file is absent
		wantCode int
		wantBody string
	}{
		"ok": {
			giveJSON: `{"player_hit":{"762x39":3},"creature_hit":{"9x19":4}}`,
			wantCode: http.StatusOK,
			wantBody: `{"player_hit":{"762x39":3},"creature_hit":{"9x19":4}}`,
		},
		"empty payload": {
			giveJSON: `{}`,
			wantCode: http.StatusOK,
			wantBody: `{"player_hit":{},"creature_hit":{}}`,
		},
		"missing file": {
			wantCode: http.StatusNoContent,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var fsys fs.FS = fstest.MapFS{}
			if tt.giveJSON != "" {
				fsys = fstest.MapFS{"hits.json": {Data: []byte(tt.giveJSON), ModTime: testMtime}}
			}

			h := data_hits.New(mod.OpenDataDir(fsys), slog.New(slog.DiscardHandler))

			rec := httptest.NewRecorder()
			h(rec, httptest.NewRequest(http.MethodGet, "/hits", http.NoBody))

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			if tt.wantCode != http.StatusOK {
				return
			}

			assert.JSONEq(t, tt.wantBody, rec.Body.String())

			// freshness + cache headers on a fresh (non-cached) read
			assert.Equal(t, "false", rec.Header().Get("X-Cache-Hit"))
			assert.Equal(t, testMtime.Format(http.TimeFormat), rec.Header().Get("Last-Modified"))

			age, ageErr := strconv.Atoi(rec.Header().Get("X-Data-Age-Seconds"))
			assert.NoError(t, ageErr)
			assert.True(t, age >= 0)
		})
	}
}
