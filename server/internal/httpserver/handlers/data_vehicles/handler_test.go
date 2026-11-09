package data_vehicles_test

import (
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"server/internal/httpserver/handlers/data_vehicles"
	"server/internal/mod"
	"server/internal/testutil/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	for name, tt := range map[string]struct {
		fsys     fs.FS
		wantCode int
		wantBody string // expected JSON body; empty = not checked
	}{
		"ok": {
			fsys: fstest.MapFS{"vehicles.json": {Data: []byte(`[{
				"class_name":"Hatchback_02_Blue","vehicle_type":"Car","hash":384863173,
				"health":0.948,"crew_occupied":0,"crew_capacity":4,"speed_ms":12.34567,"rpm":0,"gear":1,
				"engine_on":"","lights_on":0,"fuel_fraction":0.0377,"coolant_fraction":-1.0,"cargo_items":3,
				"pos_x":13376.52,"pos_y":6292.08,"pos_alt":5.96,"pos_angle":86.73,
				"parts":[{"name":"CarBattery","health":0.9551}]
			}]`)}},
			wantCode: http.StatusOK,
			wantBody: `[{
				"hash":384863173,
				"class":{"name":"Hatchback_02_Blue","type":"Car"},
				"pos":{"x":13376.52,"y":6292.08,"alt":5.96,"angle":86.73},
				"state":{
					"health":0.948,"speed_ms":12.3457,
					"engine":{"on":false,"rpm":0,"gear":1},
					"lights_on":false,
					"crew":{"occupied":0,"capacity":4},
					"fuel_fraction":0.0377,"coolant_fraction":-1,"cargo_items":3
				},
				"parts":[{"name":"CarBattery","health":0.9551}]
			}]`,
		},
		"empty payload": {
			fsys:     fstest.MapFS{"vehicles.json": {Data: []byte(`[]`)}},
			wantCode: http.StatusOK,
			wantBody: `[]`,
		},
		"missing file": {
			fsys:     fstest.MapFS{},
			wantCode: http.StatusNoContent,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := data_vehicles.New(mod.OpenDataDir(tt.fsys), slog.New(slog.DiscardHandler))

			rec := httptest.NewRecorder()
			h(rec, httptest.NewRequest(http.MethodGet, "/vehicles", http.NoBody))

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
