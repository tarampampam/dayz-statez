package data_vehicles

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

// New creates a new HTTP handler that serves the vehicles data as JSON.
func New(dd *mod.DataDir, log *slog.Logger) http.HandlerFunc {
	const fName = mod.VehiclesFileName

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, cacheHit, modTime, dataErr := dd.ReadFile[mod.Vehicles](fName)
		if dataErr != nil {
			if errors.Is(dataErr, fs.ErrNotExist) {
				w.WriteHeader(http.StatusNoContent)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"failed to read vehicles data file"}`))

			log.LogAttrs(r.Context(), slog.LevelError, "failed to read vehicles data file",
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

		var resp pkg.Vehicles
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
func fill(src *mod.Vehicles, dst *pkg.Vehicles) {
	if src == nil || dst == nil {
		return
	}

	*dst = make(pkg.Vehicles, len(*src))

	for i := range *src {
		v, c := &(*src)[i], &(*dst)[i]

		c.Hash = v.Hash

		c.Class.Name = v.ClassName
		c.Class.Type = v.VehicleType

		c.Pos.X = round4(v.PosX)
		c.Pos.Y = round4(v.PosY)
		c.Pos.Alt = round4(v.PosAlt)
		c.Pos.Angle = round4(v.PosAngle)

		c.State.Health = round4(v.Health)
		c.State.SpeedMS = round4(v.SpeedMS)
		c.State.Engine.On = bool(v.EngineOn)
		c.State.Engine.RPM = v.RPM
		c.State.Engine.Gear = v.Gear
		c.State.LightsOn = bool(v.LightsOn)
		c.State.Crew.Occupied = v.CrewOccupied
		c.State.Crew.Capacity = v.CrewCapacity
		c.State.FuelFraction = round4(v.FuelFraction)
		c.State.CoolantFraction = round4(v.CoolantFraction)
		c.State.CargoItems = v.CargoItems

		c.Parts = make([]pkg.VehiclePart, len(v.Parts))
		for j := range v.Parts {
			c.Parts[j].Name = v.Parts[j].Name
			c.Parts[j].Health = round4(v.Parts[j].Health)
		}
	}
}
