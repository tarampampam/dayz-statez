package data_areas

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"math"
	"net/http"

	"server/internal/mod"
	"server/pkg"
)

// New creates a new HTTP handler that serves the areas data as JSON.
func New(dd *mod.DataDir, log *slog.Logger) http.HandlerFunc {
	const fName = mod.AreasFileName

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, cacheHit, dataErr := dd.ReadFile[mod.Areas](fName)
		if dataErr != nil {
			if errors.Is(dataErr, fs.ErrNotExist) {
				w.WriteHeader(http.StatusNoContent)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"failed to read areas data file"}`))

			log.LogAttrs(r.Context(), slog.LevelError, "failed to read areas data file",
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

		var resp pkg.Areas
		fill(data, &resp)

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.LogAttrs(r.Context(), slog.LevelError, "failed to encode response", slog.String("error", err.Error()))
		}
	}
}

// round4 rounds f to at most 4 decimal places, trimming needless precision so the API JSON stays compact.
func round4(f float64) float64 { return math.Round(f*1e4) / 1e4 }

// fill copies the data from mod's data structure to a public data structure, that will be serialized to JSON and
// sent to the client.
func fill(src *mod.Areas, dst *pkg.Areas) {
	if src == nil || dst == nil {
		return
	}

	*dst = make(pkg.Areas, len(*src))

	for i := range *src {
		v, a := &(*src)[i], &(*dst)[i]

		a.Type = v.Type
		a.Pos.X = round4(v.PosX)
		a.Pos.Y = round4(v.PosY)
		a.Radius = round4(v.Radius)
		a.Insiders = v.Insiders
	}
}
