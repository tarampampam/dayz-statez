package maps_test

import (
	"testing"

	"tg-bot/internal/game-map/maps"
	"tg-bot/internal/game-map/maps/chernarusplus"
	"tg-bot/internal/testutil/assert"
)

func TestMaps(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name          string
		m             maps.Map
		wantMapMeters int
	}{
		{"chernarusplus", chernarusplus.Map, 15360},
		// {"livonia", livonia.Map, 12800},
		// {"nasdara", nasdara.Map, 16384},
		// {"sakhal", sakhal.Map, 15360},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.NotEmpty(t, tt.m.Name())
			assert.Equal(t, 256, tt.m.TileSizePX())
			assert.Equal(t, 8192, tt.m.MapSizePX())
			assert.Equal(t, tt.wantMapMeters, tt.m.MapSize())
			assert.Equal(t, 31, tt.m.MaxXTileIdx())
			assert.Equal(t, 31, tt.m.MaxYTileIdx())

			tile, err := tt.m.Tile(0, 0)
			assert.Nil(t, err)
			assert.NotNil(t, tile)
			assert.Equal(t, 256, tile.Bounds().Dx())
			assert.Equal(t, 256, tile.Bounds().Dy())

			_, err = tt.m.Tile(-1, 0)
			assert.Error(t, err)

			assert.NotEmpty(t, len(tt.m.Locations()))

			for _, l := range tt.m.Locations() {
				assert.NotEmpty(t, l.Category())
				assert.NotEmpty(t, l.PosX())
				assert.NotEmpty(t, l.PosY())
				assert.NotEmpty(t, l.Name())
			}
		})
	}
}
