package data_world_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"server/internal/httpserver/handlers/data_world"
	"server/internal/mod"
	"server/internal/testutil/assert"
)

func TestNew_ok(t *testing.T) {
	t.Parallel()

	// time_of_day 1 -> "day"; uptime 205.5 -> int 205; temperature 9.481267 -> round4 9.4813.
	fsys := fstest.MapFS{"world.json": {Data: []byte(
		`{"world":"chernarusplus","uptime_seconds":205.5,"hour":8,"minute":33,"time_of_day":1,` +
			`"weather":{"temperature":9.481267}}`,
	)}}

	h := data_world.New(mod.OpenDataDir(fsys), slog.New(slog.DiscardHandler))

	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/world", http.NoBody))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var got struct {
		MapName       string `json:"map_name"`
		UptimeSeconds int    `json:"uptime_seconds"`
		Time          struct {
			Hour      int    `json:"hour"`
			Minute    int    `json:"minute"`
			TimeOfDay string `json:"time_of_day"`
		} `json:"time"`
		Weather struct {
			Temperature float64 `json:"temperature"`
		} `json:"weather"`
	}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

	assert.Equal(t, "chernarusplus", got.MapName)
	assert.Equal(t, 205, got.UptimeSeconds)
	assert.Equal(t, 8, got.Time.Hour)
	assert.Equal(t, 33, got.Time.Minute)
	assert.Equal(t, "day", got.Time.TimeOfDay)
	assert.Equal(t, 9.4813, got.Weather.Temperature)
}

func TestNew_missingFile(t *testing.T) {
	t.Parallel()

	h := data_world.New(mod.OpenDataDir(fstest.MapFS{}), slog.New(slog.DiscardHandler))

	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/world", http.NoBody))

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
