package data_world

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

// New creates a new HTTP handler that serves the world data as JSON.
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

		var resp pkg.World
		fill(data, &resp)

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.LogAttrs(r.Context(), slog.LevelError, "failed to encode response", slog.String("error", err.Error()))
		}
	}
}

// round4 rounds f to at most 4 decimal places, trimming needless precision so the API JSON stays compact.
func round4(f float64) float64 { return math.Round(f*1e4) / 1e4 } //nolint:mnd // 1e4 = 4-decimal rounding scale

// fill copies the data from mod's data structure to a public data structure, that will be serialized to JSON and
// sent to the client.
func fill(src *mod.World, dst *pkg.World) { //nolint:funlen
	if src == nil || dst == nil {
		return
	}

	dst.MapName = src.World
	dst.ServerVersion = src.ServerVersion
	dst.UptimeSeconds = int(src.UptimeSeconds)

	dst.Time.Hour = src.Hour
	dst.Time.Minute = src.Minute

	//nolint:mnd // engine daytime phase (0..3) -> human-readable string
	switch src.TimeOfDay {
	case 0:
		dst.Time.TimeOfDay = "dawn"
	case 1:
		dst.Time.TimeOfDay = "day"
	case 2:
		dst.Time.TimeOfDay = "dusk"
	case 3:
		dst.Time.TimeOfDay = "night"
	default:
		dst.Time.TimeOfDay = "unknown"
	}

	dst.Players.Online = src.PlayersOnline
	dst.Players.Max = src.MaxPlayers
	dst.Players.Spawns = src.PlayersSpawns
	dst.Players.Deaths = src.PlayersDeaths

	dst.Infected.Alive = src.Infected
	dst.Infected.Dead = src.InfectedDead
	dst.Infected.Deaths = src.InfectedDeaths

	dst.Animals.Alive = src.Animals
	dst.Animals.Dead = src.AnimalsDead
	dst.Animals.Deaths = src.AnimalsDeaths

	dst.Vehicles.Cars = src.Cars
	dst.Vehicles.Boats = src.Boats
	dst.Vehicles.Helicopters = src.Helicopters
	dst.Vehicles.Destroyed.Cars = src.CarsDestroyed
	dst.Vehicles.Destroyed.Boats = src.BoatsDestroyed
	dst.Vehicles.Destroyed.Helicopters = src.HelicoptersDestroyed

	dst.Weapons.Count = src.Weapons
	dst.Weapons.Suppressors = src.Suppressors
	dst.Weapons.Optics = src.Optics
	dst.Weapons.AmmoBoxes = src.AmmoBoxes
	dst.Weapons.Magazines = src.Magazines
	dst.Weapons.Ammo = src.Ammo
	dst.Weapons.Explosives = src.Explosives
	dst.Weapons.ExplosivesDetonated = src.ExplosivesDetonated

	dst.Entities.Items = src.Items
	dst.Entities.Bottles = src.Bottles
	dst.Entities.Clothes = src.Clothes
	dst.Entities.Boxes = src.Boxes
	dst.Entities.TerritoryFlags = src.TerritoryFlags
	dst.Entities.BaseBuildings = src.BaseBuildings
	dst.Entities.Tents = src.Tents
	dst.Entities.Stashes = src.Stashes
	dst.Entities.Gardens = src.Gardens
	dst.Entities.Containers = src.Containers
	dst.Entities.CarParts = src.CarParts
	dst.Entities.CarWheels = src.CarWheels

	dst.Activity.EffectAreas = src.EffectAreas
	dst.Activity.Corpses = src.Corpses
	dst.Activity.CrashSites = src.CrashSites
	dst.Activity.ArtilleryBarrages = src.ArtilleryBarrages
	dst.Activity.GiftsUnpacked = src.GiftsUnpacked
	dst.Activity.ChatMessages = src.ChatMessages
	dst.Activity.PathGraphUpdates = src.PathGraphUpdates

	dst.EAI.Survivors = src.EAI
	dst.EAI.NPCs = src.EAINPC
	dst.EAI.Deaths = src.EAIDeaths
	dst.EAI.NPCDeaths = src.EAINPCDeaths

	dst.FPS.Current = round4(src.FPS.Current)
	dst.FPS.WindowMin = round4(src.FPS.WindowMin)
	dst.FPS.WindowMax = round4(src.FPS.WindowMax)
	dst.FPS.WindowAvg = round4(src.FPS.WindowAvg)
	dst.FPS.WindowSamples = src.FPS.WindowSamples
	dst.FPS.Limit = src.FPS.Limit

	dst.Weather.Temperature = round4(src.Weather.Temperature)
	dst.Weather.TemperatureBase = round4(src.Weather.TemperatureBase)
	dst.Weather.WindSpeed = round4(src.Weather.WindSpeed)
	dst.Weather.WindDirection = round4(src.Weather.WindDirection)
	dst.Weather.Rain = round4(src.Weather.Rain)
	dst.Weather.Snow = round4(src.Weather.Snow)
	dst.Weather.Fog = round4(src.Weather.Fog)
	dst.Weather.Clouds = round4(src.Weather.Clouds)

	dst.Food.Total = src.Food.Total
	dst.Food.Other = src.Food.Other
	dst.Food.Fruit = src.Food.Fruit
	dst.Food.Mushroom = src.Food.Mushroom
	dst.Food.Corpse = src.Food.Corpse
	dst.Food.Worm = src.Food.Worm
	dst.Food.Guts = src.Food.Guts
	dst.Food.Meat = src.Food.Meat
	dst.Food.HumanMeat = src.Food.HumanMeat
	dst.Food.Disinfectant = src.Food.Disinfectant
	dst.Food.Pills = src.Food.Pills
	dst.Food.Drink = src.Food.Drink
	dst.Food.Snack = src.Food.Snack
	dst.Food.Candy = src.Food.Candy
	dst.Food.CannedSmall = src.Food.CannedSmall
	dst.Food.CannedMedium = src.Food.CannedMedium
	dst.Food.CannedBig = src.Food.CannedBig
	dst.Food.Jar = src.Food.Jar

	dst.Mining.Bushes = src.Mining.Bushes
	dst.Mining.Rocks = src.Mining.Rocks
	dst.Mining.Trees = src.Mining.Trees

	dst.Fishing.Attempts = src.Fishing.Attempts
	dst.Fishing.Catches = src.Fishing.Catches
}
