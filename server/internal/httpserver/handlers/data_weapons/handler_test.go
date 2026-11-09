package data_weapons_test

import (
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"server/internal/httpserver/handlers/data_weapons"
	"server/internal/mod"
	"server/internal/testutil/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	for name, tt := range map[string]struct {
		fsys     fs.FS
		wantCode int
		wantBody string // expected JSON body; empty = not checked
	}{
		"ok": {
			fsys: fstest.MapFS{"weapons.json": {Data: []byte(
				`{"shots":{"akm":120},"world_count":{"mosin9130":65},"player_kills":{"akm":2},"creature_kills":{"zombie":3}}`,
			)}},
			wantCode: http.StatusOK,
			wantBody: `{"shots":{"akm":120},"world_count":{"mosin9130":65},"kills":{"players":{"akm":2},"creatures":{"zombie":3}}}`,
		},
		"empty payload": {
			fsys:     fstest.MapFS{"weapons.json": {Data: []byte(`{}`)}},
			wantCode: http.StatusOK,
			wantBody: `{"shots":{},"world_count":{},"kills":{"players":{},"creatures":{}}}`,
		},
		"missing file": {
			fsys:     fstest.MapFS{},
			wantCode: http.StatusNoContent,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := data_weapons.New(mod.OpenDataDir(tt.fsys), slog.New(slog.DiscardHandler))

			rec := httptest.NewRecorder()
			h(rec, httptest.NewRequest(http.MethodGet, "/weapons", http.NoBody))

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
