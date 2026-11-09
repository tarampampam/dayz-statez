package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"server/internal/httpserver/middleware"
	"server/internal/testutil/assert"
)

func TestCORS(t *testing.T) {
	t.Parallel()

	for name, tt := range map[string]struct {
		giveMethod  string
		giveHeaders map[string]string
		wantStatus  int
		wantHeaders map[string]string
		wantNext    bool
	}{
		"browser preflight - method only": {
			giveMethod:  http.MethodOptions,
			giveHeaders: map[string]string{"Access-Control-Request-Method": "GET"},
			wantStatus:  http.StatusNoContent,
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":   "*",
				"Access-Control-Allow-Methods":  "GET",
				"Access-Control-Max-Age":        "86400",
				"Access-Control-Expose-Headers": "*",
			},
			wantNext: false,
		},
		"browser preflight - method and headers": {
			giveMethod: http.MethodOptions,
			giveHeaders: map[string]string{
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "Authorization",
			},
			wantStatus: http.StatusNoContent,
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET",
				"Access-Control-Allow-Headers": "Authorization",
				"Access-Control-Max-Age":       "86400",
			},
			wantNext: false,
		},
		"user-initiated OPTIONS without ACRM passes through": {
			giveMethod:  http.MethodOptions,
			giveHeaders: map[string]string{},
			wantStatus:  http.StatusOK,
			wantHeaders: map[string]string{"Access-Control-Allow-Origin": "*"},
			wantNext:    true,
		},
		"actual GET carries CORS headers and passes through": {
			giveMethod:  http.MethodGet,
			giveHeaders: map[string]string{},
			wantStatus:  http.StatusOK,
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":   "*",
				"Access-Control-Expose-Headers": "*",
			},
			wantNext: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var nextCalled bool

			handler := middleware.CORS(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				nextCalled = true

				w.WriteHeader(http.StatusOK)
			}))

			r := httptest.NewRequest(tt.giveMethod, "/", http.NoBody)
			for k, v := range tt.giveHeaders {
				r.Header.Set(k, v)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantNext, nextCalled)

			for k, v := range tt.wantHeaders {
				assert.Equal(t, v, w.Header().Get(k))
			}
		})
	}
}
