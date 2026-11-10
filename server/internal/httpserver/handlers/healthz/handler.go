package healthz

import "net/http"

// Handler is a simple health check endpoint that always returns HTTP 200 OK.
func Handler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte("OK"))
}
