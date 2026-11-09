package data_players

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

// New creates a new HTTP handler that serves the players data as JSON.
func New(dd *mod.DataDir, log *slog.Logger) http.HandlerFunc {
	const fName = mod.PlayersFileName

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, cacheHit, dataErr := dd.ReadFile[mod.Players](fName)
		if dataErr != nil {
			if errors.Is(dataErr, fs.ErrNotExist) {
				w.WriteHeader(http.StatusNoContent)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"failed to read players data file"}`))

			log.LogAttrs(r.Context(), slog.LevelError, "failed to read players data file",
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

		var resp pkg.Players
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
func fill(src *mod.Players, dst *pkg.Players) {
	if src == nil || dst == nil {
		return
	}

	*dst = make(pkg.Players, len(*src))

	for i := range *src {
		v, p := &(*src)[i], &(*dst)[i]

		p.ID.SteamID = v.SteamID
		p.ID.GUID = v.GUID
		p.ID.Name = v.Name

		p.Character.Type = v.CharacterType
		p.Character.Gender = v.Gender
		p.Character.BloodType = v.BloodType

		p.Pos.X = round4(v.PosX)
		p.Pos.Y = round4(v.PosY)
		p.Pos.Alt = round4(v.PosAlt)
		p.Pos.Angle = round4(v.PosAngle)

		p.Vitals.Health = round4(v.Health)
		p.Vitals.Blood = round4(v.Blood)
		p.Vitals.Shock = round4(v.Shock)
		p.Vitals.Energy = round4(v.Energy)
		p.Vitals.Water = round4(v.Water)
		p.Vitals.Toxicity = round4(v.Toxicity)
		p.Vitals.Wetness = round4(v.Wetness)
		p.Vitals.TemperatureCelsius = round4(v.TemperatureCelsius)

		p.State.IsLoaded = bool(v.IsLoaded)
		p.State.IsUnconscious = bool(v.IsUnconscious)
		p.State.IsRestrained = bool(v.IsRestrained)
		p.State.Godmode = bool(v.Godmode)
		p.State.ImmunityBoosted = bool(v.ImmunityBoosted)
		p.State.AgentsActive = v.AgentsActive
		p.State.BleedingSources = v.BleedingSources

		p.Network.Ping.Min = v.PingMin
		p.Network.Ping.Max = v.PingMax
		p.Network.Throttle = round4(v.Throttle)

		p.Stats.PlaytimeSeconds = round4(v.StatPlaytimeSeconds)
		p.Stats.DistanceMeters = round4(v.StatDistanceMeters)
		p.Stats.LongestHitMeters = round4(v.StatLongestHitMeters)
		p.Stats.PlayersKilled = v.StatPlayersKilled
		p.Stats.InfectedKilled = v.StatInfectedKilled

		p.WeightGrams = round4(v.WeightGrams)
		p.LifetimeSeconds = round4(v.LifetimeSeconds)
	}
}
