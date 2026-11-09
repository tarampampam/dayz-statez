package data_areas_test

import (
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"server/internal/httpserver/handlers/data_areas"
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
			fsys: fstest.MapFS{"areas.json": {Data: []byte(
				`[{"type":"STATIC","class_name":"ContaminatedArea_Static","pos_x":100.5,"pos_y":200.25,"radius":75,"insiders":2}]`,
			)}},
			wantCode: http.StatusOK,
			wantBody: `[{"type":"STATIC","pos":{"x":100.5,"y":200.25},"radius":75,"insiders":2}]`,
		},
		"empty payload": {
			fsys:     fstest.MapFS{"areas.json": {Data: []byte(`[]`)}},
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

			h := data_areas.New(mod.OpenDataDir(tt.fsys), slog.New(slog.DiscardHandler))

			rec := httptest.NewRecorder()
			h(rec, httptest.NewRequest(http.MethodGet, "/areas", http.NoBody))

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
