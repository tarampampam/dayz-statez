package healthz

import "net/http"

func Handler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte("OK"))
}
