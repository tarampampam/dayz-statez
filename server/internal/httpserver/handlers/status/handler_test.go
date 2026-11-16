package status_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"testing/fstest"
	"time"

	"gh.tarampamp.am/dayz-statez/server/internal/httpserver/handlers/status"
	"gh.tarampamp.am/dayz-statez/server/internal/mod"
	"gh.tarampamp.am/dayz-statez/server/internal/testutil/assert"
)

var testMtime = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

func TestNew_ok(t *testing.T) {
	t.Parallel()

	// uptime 205.5 -> int 205; temperature 9.481267 -> round4 9.4813.
	fsys := fstest.MapFS{"world.json": {Data: []byte(
		`{"server_version":"1.30.164289.70","uptime_seconds":205.5,"hour":8,"minute":33,` +
			`"players_online":2,"max_players":60,"weather":{"temperature":9.481267,"rain":0.5,"snow":0}}`,
	), ModTime: testMtime}}

	h := status.New(mod.OpenDataDir(fsys), slog.New(slog.DiscardHandler))

	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/status", http.NoBody))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var got struct {
		Version       string `json:"version"`
		UptimeSeconds int    `json:"uptime_seconds"`
		Players       struct {
			Online int `json:"online"`
			Max    int `json:"max"`
		} `json:"players"`
		Time struct {
			Hours   int `json:"hours"`
			Minutes int `json:"minutes"`
		} `json:"time"`
		Weather struct {
			Temperature float64 `json:"temperature"`
			Rain        float64 `json:"rain"`
			Snow        float64 `json:"snow"`
		} `json:"weather"`
	}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

	assert.Equal(t, "1.30.164289.70", got.Version)
	assert.Equal(t, 205, got.UptimeSeconds)
	assert.Equal(t, 2, got.Players.Online)
	assert.Equal(t, 60, got.Players.Max)
	assert.Equal(t, 8, got.Time.Hours)
	assert.Equal(t, 33, got.Time.Minutes)
	assert.Equal(t, 9.4813, got.Weather.Temperature)
	assert.Equal(t, 0.5, got.Weather.Rain)
	assert.Equal(t, float64(0), got.Weather.Snow)

	// freshness + cache headers on a fresh (non-cached) read
	assert.Equal(t, "false", rec.Header().Get("X-Cache-Hit"))
	assert.Equal(t, testMtime.Format(http.TimeFormat), rec.Header().Get("Last-Modified"))

	age, ageErr := strconv.Atoi(rec.Header().Get("X-Data-Age-Seconds"))
	assert.NoError(t, ageErr)
	assert.True(t, age >= 0)
}

func TestNew_missingFile(t *testing.T) {
	t.Parallel()

	h := status.New(mod.OpenDataDir(fstest.MapFS{}), slog.New(slog.DiscardHandler))

	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/status", http.NoBody))

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
