package data_territories

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"

	"server/internal/mod"
	"server/pkg"
)

// New creates a new HTTP handler that serves the territories data as JSON.
func New(dd *mod.DataDir, log *slog.Logger) http.HandlerFunc {
	const fName = mod.TerritoriesFileName

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, cacheHit, modTime, dataErr := dd.ReadFile[mod.Territories](fName)
		if dataErr != nil {
			if errors.Is(dataErr, fs.ErrNotExist) {
				w.WriteHeader(http.StatusNoContent)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"failed to read territories data file"}`))

			log.LogAttrs(r.Context(), slog.LevelError, "failed to read territories data file",
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

		var resp pkg.Territories
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
func fill(src *mod.Territories, dst *pkg.Territories) {
	if src == nil || dst == nil {
		return
	}

	*dst = make(pkg.Territories, len(*src))

	for i := range *src {
		v, t := &(*src)[i], &(*dst)[i]

		t.Pos.X = round4(v.PosX)
		t.Pos.Y = round4(v.PosY)

		t.Refresher.LifetimeFraction = round4(v.LifetimeFraction)
		t.Refresher.Radius = round4(v.RefresherRadius)
	}
}
