package data_animals

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"maps"
	"net/http"
	"strconv"
	"time"

	"server/internal/mod"
	"server/pkg"
)

// New creates a new HTTP handler that serves the animals data as JSON.
func New(dd *mod.DataDir, log *slog.Logger) http.HandlerFunc {
	const fName = mod.AnimalsFileName

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, cacheHit, modTime, dataErr := dd.ReadFile[mod.Animals](fName)
		if dataErr != nil {
			if errors.Is(dataErr, fs.ErrNotExist) {
				w.WriteHeader(http.StatusNoContent)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"failed to read animals data file"}`))

			log.LogAttrs(r.Context(), slog.LevelError, "failed to read animals data file",
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

		w.Header().Set("Last-Modified", modTime.UTC().Format(http.TimeFormat))
		w.Header().Set("X-Data-Age-Seconds", strconv.Itoa(max(0, int(time.Since(modTime).Seconds()))))

		var resp pkg.Animals
		fill(data, &resp)

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.LogAttrs(r.Context(), slog.LevelError, "failed to encode response", slog.String("error", err.Error()))
		}
	}
}

// fill copies the data from mod's data structure to a public data structure, that will be serialized to JSON and
// sent to the client.
func fill(src *mod.Animals, dst *pkg.Animals) {
	if src == nil || dst == nil {
		return
	}

	dst.ByType = make(map[string]int, len(src.ByType))

	maps.Copy(dst.ByType, src.ByType)
}
