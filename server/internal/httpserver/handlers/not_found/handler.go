package not_found

import "net/http"

func Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	_, _ = w.Write([]byte(`{"error": "not found"}`))
}
