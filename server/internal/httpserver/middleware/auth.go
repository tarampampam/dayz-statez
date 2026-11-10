package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewBearerAuth returns a middleware that requires the request to carry the given token, provided either as
// "Authorization: Bearer <token>" or as the "X-API-Key: <token>" header. On a missing/mismatching token it
// responds with 401 Unauthorized and does NOT call the next handler.
//
// The caller must only use this middleware with a non-empty token.
func NewBearerAuth(token string) func(http.Handler) http.Handler {
	unauthorizedBody, err := json.Marshal(map[string]any{
		"error":   "unauthorized",
		"message": "authentication required to access this endpoint",
		"how_to_authenticate": map[string]string{
			"header":             "Authorization",
			"format":             "Bearer <token>",
			"alternative_header": "X-API-Key: <token>",
		},
	})
	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := tokenFromRequest(r); got != token {
				challenge := "Bearer"
				if got != "" { // a token was presented, it just did not match
					challenge = `Bearer error="invalid_token"`
				}

				w.Header().Set("WWW-Authenticate", challenge)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write(unauthorizedBody)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// tokenFromRequest extracts the presented token from the "Authorization: Bearer <token>" header, falling back
// to the "X-API-Key" header. Returns an empty string when neither is present.
func tokenFromRequest(r *http.Request) string {
	if h := r.Header.Get("Authorization"); h != "" {
		if token, ok := strings.CutPrefix(h, "Bearer "); ok {
			return token
		}
	}

	return r.Header.Get("X-Api-Key")
}
