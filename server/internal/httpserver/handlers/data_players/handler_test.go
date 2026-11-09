package data_players_test

import (
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"testing/fstest"
	"time"

	"server/internal/httpserver/handlers/data_players"
	"server/internal/mod"
	"server/internal/testutil/assert"
)

var testMtime = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

func TestNew(t *testing.T) {
	t.Parallel()

	for name, tt := range map[string]struct {
		giveJSON string // data file content; empty means the file is absent
		wantCode int
		wantBody string
	}{
		"ok": {
			giveJSON: `[{
				"steam_id":"76561198000000001","guid":"abc","name":"tester",
				"blood_type":"A+","gender":"female","character_type":"Judy",
				"is_loaded":1,"is_unconscious":1,"is_restrained":0,"godmode":0,"immunity_boosted":0,
				"health":1.0,"blood":0.89,"shock":1.0,"energy":0.5,"water":0.6,"toxicity":0.1,"wetness":0.2,"temperature_celsius":36.5,
				"weight_grams":7000.5,"lifetime_seconds":123.5,
				"pos_x":13377.5,"pos_y":6292.14,"pos_alt":7.17,"pos_angle":-121.33,
				"ping_min":2,"ping_max":8,"throttle":0.123456,"agents_active":1,"bleeding_sources":0,
				"stat_playtime_seconds":5000,"stat_distance_meters":693.69,"stat_infected_killed":4
			}]`,
			wantCode: http.StatusOK,
			wantBody: `[{
				"id":{"steam_id":"76561198000000001","guid":"abc","name":"tester"},
				"character":{"type":"Judy","gender":"female","blood_type":"A+"},
				"pos":{"x":13377.5,"y":6292.14,"alt":7.17,"angle":-121.33},
				"vitals":{"health":1,"blood":0.89,"shock":1,"energy":0.5,"water":0.6,"toxicity":0.1,"wetness":0.2,"temperature_celsius":36.5},
				"state":{"is_loaded":true,"is_unconscious":true,"is_restrained":false,"godmode":false,"immunity_boosted":false,"agents_active":1,"bleeding_sources":0},
				"network":{"ping":{"min":2,"max":8},"throttle":0.1235},
				"stats":{"playtime_sec":5000,"distance_meters":693.69,"longest_hit_meters":0,"players_killed":0,"infected_killed":4},
				"weight_grams":7000.5,"lifetime_seconds":123.5
			}]`,
		},
		"empty payload": {
			giveJSON: `[]`,
			wantCode: http.StatusOK,
			wantBody: `[]`,
		},
		"missing file": {
			wantCode: http.StatusNoContent,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var fsys fs.FS = fstest.MapFS{}
			if tt.giveJSON != "" {
				fsys = fstest.MapFS{"players.json": {Data: []byte(tt.giveJSON), ModTime: testMtime}}
			}

			h := data_players.New(mod.OpenDataDir(fsys), slog.New(slog.DiscardHandler))

			rec := httptest.NewRecorder()
			h(rec, httptest.NewRequest(http.MethodGet, "/players", http.NoBody))

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			if tt.wantCode != http.StatusOK {
				return
			}

			assert.JSONEq(t, tt.wantBody, rec.Body.String())

			// freshness + cache headers on a fresh (non-cached) read
			assert.Equal(t, "false", rec.Header().Get("X-Cache-Hit"))
			assert.Equal(t, testMtime.Format(http.TimeFormat), rec.Header().Get("Last-Modified"))

			age, ageErr := strconv.Atoi(rec.Header().Get("X-Data-Age-Seconds"))
			assert.NoError(t, ageErr)
			assert.True(t, age >= 0)
		})
	}
}
