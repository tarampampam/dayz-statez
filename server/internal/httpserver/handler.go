package httpserver

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"server/internal/httpserver/handlers/data_animals"
	"server/internal/httpserver/handlers/data_areas"
	"server/internal/httpserver/handlers/data_events"
	"server/internal/httpserver/handlers/data_hits"
	"server/internal/httpserver/handlers/data_players"
	"server/internal/httpserver/handlers/data_territories"
	"server/internal/httpserver/handlers/data_vehicles"
	"server/internal/httpserver/handlers/data_weapons"
	"server/internal/httpserver/handlers/data_world"
	"server/internal/httpserver/handlers/data_zombies"
	"server/internal/httpserver/handlers/healthz"
	"server/internal/httpserver/handlers/not_found"
	"server/internal/httpserver/middleware"
	"server/internal/mod"
)

// NewHandler creates the main HTTP handler for the server.
func NewHandler(log *slog.Logger, dd *mod.DataDir, apiToken string) http.Handler { //nolint:funlen
	mux := http.NewServeMux()

	auth := func() func(h http.HandlerFunc) http.HandlerFunc {
		if apiToken != "" {
			mw := middleware.NewBearerAuth(apiToken)

			return func(h http.HandlerFunc) http.HandlerFunc { return mw(h).ServeHTTP }
		}

		return func(h http.HandlerFunc) http.HandlerFunc { return h } // auth disabled
	}()

	const (
		animalsRoute     = "/animals"
		areasRoute       = "/areas"
		eventsRoute      = "/events"
		hitsRoute        = "/hits"
		playersRoute     = "/players"
		territoriesRoute = "/territories"
		vehiclesRoute    = "/vehicles"
		weaponsRoute     = "/weapons"
		worldRoute       = "/world"
		zombiesRoute     = "/zombies"

		healthRoute = "/healthz"
	)

	helpMessage, jErr := json.Marshal(map[string]string{
		"$help$": "StateZ read-only HTTP API. Each endpoint returns the latest snapshot the DayZ server mod wrote " +
			"for its data domain, cached in memory (see the X-Cache-Hit response header)",
		animalsRoute: "Live wildlife census - count of alive animals per canonical type. Use for " +
			"population/spawn-balance monitoring and wildlife density dashboards",
		areasRoute: "Active effect areas - static preset gas zones and dynamic toxic/artillery events - each with its " +
			"type, center position, radius and the number of players currently inside. Use for a live hazard map and " +
			"player-in-zone alerts",
		eventsRoute: "Cumulative counters (since server start) of engine events keyed by name - session lifecycle, " +
			"client connect/disconnect, deaths, chat, VON, etc. Use for activity analytics and diagnosing " +
			"connection/session issues",
		hitsRoute: "Cumulative qualifying hit counts since server start, split into player_hit and creature_hit and " +
			"broken down by ammo type. Use for combat/weapon-usage analytics",
		playersRoute: "Snapshot of every online player - identity (steam id, guid, name), character, position, vitals " +
			"(health/blood/shock/energy/water/toxicity/temperature), state flags (unconscious, restrained, godmode, " +
			"bleeding, etc.), network and lifetime stats (playtime, distance, kills). Use for a live player map",
		territoriesRoute: "Base territory flags - position, remaining refresh-timer fraction and refresh radius. " +
			"Use to map bases and track which are kept alive vs decaying toward cleanup",
		vehiclesRoute: "Every tracked vehicle - class/type, position, health, crew occupied/capacity, engine, lights, " +
			"fuel & radiator-coolant levels, cargo item count and per-part health. Use for a live vehicle map, fleet " +
			"condition monitoring and spotting free vs player-claimed cars (empty cargo + full hp + still at a spawn " +
			"point = unclaimed)",
		weaponsRoute: "Cumulative shots fired per weapon type, current live weapon-entity counts, and kills grouped by " +
			"source (players vs creatures). Use for weapon-usage/balance analytics and loot-economy insight",
		worldRoute: "Single world snapshot - map name, uptime, in-game clock (hour, minute, time-of-day phase), " +
			"population counts (players/infected/animals plus cumulative deaths), vehicle counts, loot & entity " +
			"breakdowns, server FPS, weather, and food/mining/fishing stats. Use for a server status dashboard and " +
			"performance/weather/economy monitoring",
		zombiesRoute: "Live infected census - counts per canonical type and per AI mind-state (calm, disturbed, " +
			"alerted, chase, fight, dead). Use for horde-density and AI-activity monitoring and correlating it with " +
			"server FPS",
		healthRoute: "Liveness probe - always returns 200",
	})
	if jErr != nil {
		panic(jErr)
	}

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(helpMessage)
	})

	mux.HandleFunc("GET "+animalsRoute, auth(data_animals.New(dd, log)))
	mux.HandleFunc("GET "+areasRoute, auth(data_areas.New(dd, log)))
	mux.HandleFunc("GET "+eventsRoute, auth(data_events.New(dd, log)))
	mux.HandleFunc("GET "+hitsRoute, auth(data_hits.New(dd, log)))
	mux.HandleFunc("GET "+playersRoute, auth(data_players.New(dd, log)))
	mux.HandleFunc("GET "+territoriesRoute, auth(data_territories.New(dd, log)))
	mux.HandleFunc("GET "+vehiclesRoute, auth(data_vehicles.New(dd, log)))
	mux.HandleFunc("GET "+weaponsRoute, auth(data_weapons.New(dd, log)))
	mux.HandleFunc("GET "+worldRoute, auth(data_world.New(dd, log)))
	mux.HandleFunc("GET "+zombiesRoute, auth(data_zombies.New(dd, log)))

	mux.HandleFunc("GET "+healthRoute, healthz.Handler)
	mux.HandleFunc("/", not_found.Handler)

	return middleware.Apply(mux,
		middleware.CORS,
		middleware.NewAccessLog(log, slog.LevelInfo, func() func(*http.Request) bool {
			// the following routes will NOT be logged
			skipMap := map[string]struct{}{
				healthRoute:             {},
				"/apple-touch-icon.png": {},
				"/favicon.ico":          {},
				"/robots.txt":           {},
			}

			return func(r *http.Request) bool {
				_, skip := skipMap[strings.ToLower(r.URL.Path)]

				return skip
			}
		}()),
	)
}
