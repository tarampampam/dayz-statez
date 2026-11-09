package middleware

import "net/http"

// CORS returns a middleware that enables cross-origin access to the API from browsers.
//
// It sets "Access-Control-Allow-Origin: *" and "Access-Control-Expose-Headers: *" on every response, so a
// cross-origin Fetch call can both read the body and inspect any response header (including our freshness/cache
// headers), without having to enumerate them. Browser CORS preflight requests (OPTIONS carrying the
// Access-Control-Request-Method header, which the browser adds only to its own automatically generated
// preflights - never to user-initiated OPTIONS) are answered directly with 204 and the matching headers,
// without forwarding to the next handler.
//
// The wildcards are safe here because the API authenticates with a bearer token (a header the client sets
// explicitly), not with cookies/credentials - so "*" is never combined with Access-Control-Allow-Credentials.
// Note: the "*" expose wildcard is ignored by browsers for credentialed (cookie) requests, which this API does
// not use.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdr := w.Header()
		hdr.Set("Access-Control-Allow-Origin", "*")
		hdr.Set("Access-Control-Expose-Headers", "*")

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			hdr.Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))

			// mirror the requested headers (covers Authorization / X-API-Key, not matched by the "*" wildcard)
			if reqHeaders := r.Header.Get("Access-Control-Request-Headers"); reqHeaders != "" {
				hdr.Set("Access-Control-Allow-Headers", reqHeaders)
			}

			hdr.Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)

			return
		}

		next.ServeHTTP(w, r)
	})
}
