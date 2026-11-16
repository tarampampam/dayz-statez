package status

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"

	"gh.tarampamp.am/dayz-statez/server/internal/mod"
	"gh.tarampamp.am/dayz-statez/server/pkg/api"
)

// New creates a new public HTTP handler that serves a short summary of the game world/server state as JSON.
func New(dd *mod.DataDir, log *slog.Logger) http.HandlerFunc {
	const fName = mod.WorldFileName

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, cacheHit, modTime, dataErr := dd.ReadFile[mod.World](fName)
		if dataErr != nil {
			if errors.Is(dataErr, fs.ErrNotExist) {
				w.WriteHeader(http.StatusNoContent)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"failed to read world data file"}`))

			log.LogAttrs(r.Context(), slog.LevelError, "failed to read world data file",
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

		var resp api.Status
		fill(data, &resp)

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.LogAttrs(r.Context(), slog.LevelError, "failed to encode response", slog.String("error", err.Error()))
		}
	}
}

// round4 rounds f to at most 4 decimal places, trimming needless precision so the API JSON stays compact.
func round4(f float64) float64 { return math.Round(f*1e4) / 1e4 } //nolint:mnd // 1e4 = 4-decimal rounding scale

// fill copies the needed fields from the mod's world data into the public status structure, that will be serialized
// to JSON and sent to the client.
func fill(src *mod.World, dst *api.Status) {
	if src == nil || dst == nil {
		return
	}

	dst.Version = src.ServerVersion
	dst.UptimeSeconds = int(src.UptimeSeconds)

	dst.Players.Online = src.PlayersOnline
	dst.Players.Max = src.MaxPlayers

	dst.Time.Hours = src.Hour
	dst.Time.Minutes = src.Minute

	dst.Weather.Temperature = round4(src.Weather.Temperature)
	dst.Weather.Rain = round4(src.Weather.Rain)
	dst.Weather.Snow = round4(src.Weather.Snow)
}
