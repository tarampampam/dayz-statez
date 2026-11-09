package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"server/internal/httpserver/middleware"
	"server/internal/testutil/assert"
)

func TestNewBearerAuth(t *testing.T) {
	t.Parallel()

	const token = "s3cr3t"

	for name, tt := range map[string]struct {
		giveHeaders   map[string]string
		wantStatus    int
		wantNext      bool
		wantChallenge string // expected WWW-Authenticate header on 401
	}{
		"valid bearer": {
			giveHeaders: map[string]string{"Authorization": "Bearer " + token},
			wantStatus:  http.StatusOK,
			wantNext:    true,
		},
		"valid x-api-key": {
			giveHeaders: map[string]string{"X-API-Key": token},
			wantStatus:  http.StatusOK,
			wantNext:    true,
		},
		"wrong token": {
			giveHeaders:   map[string]string{"Authorization": "Bearer nope"},
			wantStatus:    http.StatusUnauthorized,
			wantNext:      false,
			wantChallenge: `Bearer error="invalid_token"`,
		},
		"missing token": {
			giveHeaders:   map[string]string{},
			wantStatus:    http.StatusUnauthorized,
			wantNext:      false,
			wantChallenge: "Bearer",
		},
		"wrong scheme": {
			giveHeaders:   map[string]string{"Authorization": "Basic " + token},
			wantStatus:    http.StatusUnauthorized,
			wantNext:      false,
			wantChallenge: "Bearer", // no bearer credentials presented -> no error code
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var nextCalled bool

			handler := middleware.NewBearerAuth(token)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				nextCalled = true

				w.WriteHeader(http.StatusOK)
			}))

			r := httptest.NewRequest(http.MethodGet, "/players", http.NoBody)
			for k, v := range tt.giveHeaders {
				r.Header.Set(k, v)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantNext, nextCalled)

			if tt.wantStatus == http.StatusUnauthorized {
				assert.Equal(t, tt.wantChallenge, w.Header().Get("WWW-Authenticate"))
			}
		})
	}
}
