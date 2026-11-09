package data_events

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"maps"
	"net/http"

	"server/internal/mod"
	"server/pkg"
)

// New creates a new HTTP handler that serves the events data as JSON.
func New(dd *mod.DataDir, log *slog.Logger) http.HandlerFunc {
	const fName = mod.EventsFileName

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, cacheHit, dataErr := dd.ReadFile[mod.Events](fName)
		if dataErr != nil {
			if errors.Is(dataErr, fs.ErrNotExist) {
				w.WriteHeader(http.StatusNoContent)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"failed to read events data file"}`))

			log.LogAttrs(r.Context(), slog.LevelError, "failed to read events data file",
				slog.String("file", string(fName)),
				slog.String("error", dataErr.Error()),
			)

			return
		}

		w.Header().Set("X-Cache-Hit", func() string {
			if cacheHit {
				return "true"
			}

			return "false"
		}())

		var resp pkg.Events
		fill(data, &resp)

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.LogAttrs(r.Context(), slog.LevelError, "failed to encode response", slog.String("error", err.Error()))
		}
	}
}

// fill copies the data from mod's data structure to a public data structure, that will be serialized to JSON and
// sent to the client.
func fill(src *mod.Events, dst *pkg.Events) {
	if src == nil || dst == nil {
		return
	}

	dst.Events = make(map[string]int, len(src.Events))

	maps.Copy(dst.Events, src.Events)
}
