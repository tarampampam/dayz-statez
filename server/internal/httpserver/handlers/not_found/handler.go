package not_found

import "net/http"

// Handler is the HTTP handler for the "not found" route.
// It responds with a 404 status code and a JSON error message.
func Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	_, _ = w.Write([]byte(`{"error": "not found"}`))
}
