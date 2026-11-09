package data_index

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"server/internal/httpserver/handlers"
	"server/internal/mod"
)

func New(dd *mod.DataDir, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			animals     *mod.Animals
			areas       *mod.Areas
			events      *mod.Events
			hits        *mod.Hits
			players     *mod.Players
			territories *mod.Territories
			vehicles    *mod.Vehicles
			weapons     *mod.Weapons
			world       *mod.World
			zombies     *mod.Zombies
		)

		var (
			ctx = r.Context()
			wg  sync.WaitGroup
		)

		wg.Go(func() { animals = read[mod.Animals](ctx, log, dd, mod.AnimalsFileName) })
		wg.Go(func() { areas = read[mod.Areas](ctx, log, dd, mod.AreasFileName) })
		wg.Go(func() { events = read[mod.Events](ctx, log, dd, mod.EventsFileName) })
		wg.Go(func() { hits = read[mod.Hits](ctx, log, dd, mod.HitsFileName) })
		wg.Go(func() { players = read[mod.Players](ctx, log, dd, mod.PlayersFileName) })
		wg.Go(func() { territories = read[mod.Territories](ctx, log, dd, mod.TerritoriesFileName) })
		wg.Go(func() { vehicles = read[mod.Vehicles](ctx, log, dd, mod.VehiclesFileName) })
		wg.Go(func() { weapons = read[mod.Weapons](ctx, log, dd, mod.WeaponsFileName) })
		wg.Go(func() { world = read[mod.World](ctx, log, dd, mod.WorldFileName) })
		wg.Go(func() { zombies = read[mod.Zombies](ctx, log, dd, mod.ZombiesFileName) })

		wg.Wait()

		w.Header().Set("Content-Type", "application/json")

		if everyIsNull(animals, areas, events, hits, players, territories, vehicles, weapons, world, zombies) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"failed to read any data files"}`))

			return
		}

		w.WriteHeader(http.StatusOK)

		response := struct {
			Animals     handlers.Animals     `json:"animals"`
			Areas       handlers.Areas       `json:"areas"`
			Events      handlers.Events      `json:"events"`
			Hits        handlers.Hits        `json:"hits"`
			Players     handlers.Players     `json:"players"`
			Territories handlers.Territories `json:"territories"`
			Vehicles    handlers.Vehicles    `json:"vehicles"`
			Weapons     handlers.Weapons     `json:"weapons"`
			World       handlers.World       `json:"world"`
			Zombies     handlers.Zombies     `json:"zombies"`
		}{}

		response.Animals.FillFrom(animals)
		response.Areas.FillFrom(areas)
		response.Events.FillFrom(events)
		response.Hits.FillFrom(hits)
		response.Players.FillFrom(players)
		response.Territories.FillFrom(territories)
		response.Vehicles.FillFrom(vehicles)
		response.Weapons.FillFrom(weapons)
		response.World.FillFrom(world)
		response.Zombies.FillFrom(zombies)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.LogAttrs(ctx, slog.LevelError, "failed to encode response", slog.String("error", err.Error()))
		}
	}
}

// read decodes name as T. On error it logs a warning and returns nil.
func read[T any](ctx context.Context, log *slog.Logger, dd *mod.DataDir, name mod.DataFileName) *T {
	v, _, err := dd.ReadFile[T](name)
	if err != nil {
		log.LogAttrs(
			ctx,
			slog.LevelWarn,
			"failed to read data file",
			slog.String("file", string(name)),
			slog.String("error", err.Error()),
		)

		return nil
	}

	return &v
}

func everyIsNull(v ...any) bool {
	for _, p := range v {
		if p != nil {
			return false
		}
	}

	return true
}
