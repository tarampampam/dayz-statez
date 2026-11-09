package data_territories_test

import (
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"server/internal/httpserver/handlers/data_territories"
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
			fsys: fstest.MapFS{"territories.json": {Data: []byte(
				`[{"pos_x":7500.25,"pos_y":5500.5,"lifetime_fraction":0.75,"refresher_radius":100}]`,
			)}},
			wantCode: http.StatusOK,
			wantBody: `[{"pos":{"x":7500.25,"y":5500.5},"refresher":{"lifetime_fraction":0.75,"radius":100}}]`,
		},
		"empty payload": {
			fsys:     fstest.MapFS{"territories.json": {Data: []byte(`[]`)}},
			wantCode: http.StatusOK,
			wantBody: `[]`,
		},
		"missing file": {
			fsys:     fstest.MapFS{},
			wantCode: http.StatusNoContent,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := data_territories.New(mod.OpenDataDir(tt.fsys), slog.New(slog.DiscardHandler))

			rec := httptest.NewRecorder()
			h(rec, httptest.NewRequest(http.MethodGet, "/territories", http.NoBody))

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
