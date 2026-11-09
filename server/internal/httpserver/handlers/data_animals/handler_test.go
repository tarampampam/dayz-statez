package data_animals_test

import (
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"server/internal/httpserver/handlers/data_animals"
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
			fsys:     fstest.MapFS{"animals.json": {Data: []byte(`{"by_type":{"wolf":3,"deer":7}}`)}},
			wantCode: http.StatusOK,
			wantBody: `{"by_type":{"wolf":3,"deer":7}}`,
		},
		"empty payload": {
			fsys:     fstest.MapFS{"animals.json": {Data: []byte(`{}`)}},
			wantCode: http.StatusOK,
			wantBody: `{"by_type":{}}`,
		},
		"missing file": {
			fsys:     fstest.MapFS{},
			wantCode: http.StatusNoContent,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := data_animals.New(mod.OpenDataDir(tt.fsys), slog.New(slog.DiscardHandler))

			rec := httptest.NewRecorder()
			h(rec, httptest.NewRequest(http.MethodGet, "/animals", http.NoBody))

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
