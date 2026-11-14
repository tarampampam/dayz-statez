package game_map_test

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"testing"

	gamemap "tg-bot/internal/game-map"
	"tg-bot/internal/game-map/maps"
	"tg-bot/internal/testutil/assert"
)

func TestDrawPlayerDimensions(t *testing.T) {
	t.Parallel()

	for name, tt := range map[string]struct {
		opts         []gamemap.DrawPlayerOption
		wantW, wantH int
	}{
		// defaults are 900x600 half-sizes -> 1800x1200, and the map is big enough not to clamp.
		"default margins": {nil, 1800, 1200},
		// 2*margin on each axis.
		"custom margins": {[]gamemap.DrawPlayerOption{gamemap.WithDrawPlayerMarginPixels(100, 80)}, 200, 160},
		// margins past the map edge clamp to the map size.
		"margins exceed map clamp": {[]gamemap.DrawPlayerOption{gamemap.WithDrawPlayerMarginPixels(5000, 5000)}, 2048, 2048},
		// axes are independent.
		"asymmetric margins": {[]gamemap.DrawPlayerOption{gamemap.WithDrawPlayerMarginPixels(300, 50)}, 600, 100},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			img, err := gamemap.DrawPlayer(newStubMap(), "hero", stubCenter, stubCenter, tt.opts...)
			assert.NoError(t, err)
			assert.NotNil(t, img)

			assert.Equal(t, tt.wantW, img.Bounds().Dx())
			assert.Equal(t, tt.wantH, img.Bounds().Dy())
		})
	}
}

func TestDrawPlayerDrawsMarker(t *testing.T) {
	t.Parallel()

	// With margin (100,80) the centered player lands on pixel (100,80). The marker core is pure fill at its
	// center, so that exact pixel must come out as the fill color.
	fill := color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}

	img, err := gamemap.DrawPlayer(newStubMap(), "", stubCenter, stubCenter,
		gamemap.WithDrawPlayerMarginPixels(100, 80),
		gamemap.WithDrawPlayerMarker(40, fill),
		gamemap.WithDrawPlayerCoordinates(false),
	)
	assert.NoError(t, err)
	assert.NotNil(t, img)

	assert.Equal(t, fill, img.RGBAAt(100, 80))
}

func TestDrawPlayerOptionEffect(t *testing.T) {
	t.Parallel()

	inView := stubLoc{name: "Town", x: 1000, y: 1050} // inside the (100,80) crop around the center
	outOfView := stubLoc{name: "Far", x: 100, y: 100} // nowhere near the crop

	for name, tt := range map[string]struct {
		locs         []maps.Loc
		optsA, optsB []gamemap.DrawPlayerOption
		wantEqual    bool // true => the toggle must NOT change the output
	}{
		"coordinates toggle changes output": {
			nil,
			nil, // coordinates default to on
			[]gamemap.DrawPlayerOption{gamemap.WithDrawPlayerCoordinates(false)},
			false,
		},
		"in-view POI changes output": {
			[]maps.Loc{inView},
			[]gamemap.DrawPlayerOption{gamemap.WithDrawPlayerWithPointsOfInterest(true)},
			[]gamemap.DrawPlayerOption{gamemap.WithDrawPlayerWithPointsOfInterest(false)},
			false,
		},
		"out-of-view POI is culled": {
			[]maps.Loc{outOfView},
			[]gamemap.DrawPlayerOption{gamemap.WithDrawPlayerWithPointsOfInterest(true)},
			[]gamemap.DrawPlayerOption{gamemap.WithDrawPlayerWithPointsOfInterest(false)},
			true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := newStubMap()
			m.locations = tt.locs

			base := gamemap.WithDrawPlayerMarginPixels(100, 80)

			a, err := gamemap.DrawPlayer(m, "hero", stubCenter, stubCenter, append([]gamemap.DrawPlayerOption{base}, tt.optsA...)...)
			assert.NoError(t, err)

			b, err := gamemap.DrawPlayer(m, "hero", stubCenter, stubCenter, append([]gamemap.DrawPlayerOption{base}, tt.optsB...)...)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantEqual, bytes.Equal(a.Pix, b.Pix))
		})
	}
}

func TestDrawPlayerNickname(t *testing.T) {
	t.Parallel()

	base := gamemap.WithDrawPlayerMarginPixels(100, 80)

	// an empty nickname is a no-op for the label, not an error
	blank, err := gamemap.DrawPlayer(newStubMap(), "", stubCenter, stubCenter, base)
	assert.NoError(t, err)
	assert.NotNil(t, blank)

	named, err := gamemap.DrawPlayer(newStubMap(), "hero", stubCenter, stubCenter, base)
	assert.NoError(t, err)

	// a non-empty nickname must leave a visible mark
	assert.False(t, bytes.Equal(blank.Pix, named.Pix))
}

func TestDrawPlayerTileError(t *testing.T) {
	t.Parallel()

	m := newStubMap()
	m.tileErr = errTileBoom

	img, err := gamemap.DrawPlayer(m, "hero", stubCenter, stubCenter, gamemap.WithDrawPlayerMarginPixels(100, 80))

	assert.ErrorIs(t, err, errTileBoom)
	assert.Nil(t, img)
}

// --------------------------------------------------------------------------------------------------------------------

// The stub map is a flat 2048px / 2048m square (scale 1) of uniform gray tiles. Scale 1 keeps the math trivial -
// a player at (1024,1024) with margin (mx,my) and far from every edge projects to pixel (mx,my).
const stubCenter = 1024.0 // player position used across the tests (map center)

var errTileBoom = errors.New("tile boom")

// stubLoc is a fixed point of interest.
type stubLoc struct {
	name string
	x, y float64
}

func (l stubLoc) Category() string { return "city" }
func (l stubLoc) Name() string     { return l.name }
func (l stubLoc) PosX() float64    { return l.x }
func (l stubLoc) PosY() float64    { return l.y }

// stubMap is a test double for maps.Map: a flat square of uniform gray tiles, with optional points of interest
// and an optional tile error for the failure path.
type stubMap struct {
	tileSizePX, mapSizePX, mapSize int
	locations                      []maps.Loc
	tileErr                        error // if set, every Tile call fails with it
}

func newStubMap() *stubMap {
	return &stubMap{tileSizePX: 64, mapSizePX: 2048, mapSize: 2048}
}

func (m *stubMap) Name() string          { return "stub" }
func (m *stubMap) TileSizePX() int       { return m.tileSizePX }
func (m *stubMap) MapSizePX() int        { return m.mapSizePX }
func (m *stubMap) MapSize() int          { return m.mapSize }
func (m *stubMap) MaxXTileIdx() int      { return m.mapSizePX/m.tileSizePX - 1 }
func (m *stubMap) MaxYTileIdx() int      { return m.mapSizePX/m.tileSizePX - 1 }
func (m *stubMap) Locations() []maps.Loc { return m.locations }

func (m *stubMap) Tile(_, _ int) (image.Image, error) {
	if m.tileErr != nil {
		return nil, m.tileErr
	}

	tile := image.NewRGBA(image.Rect(0, 0, m.tileSizePX, m.tileSizePX))
	draw.Draw(tile, tile.Bounds(), image.NewUniform(color.RGBA{R: 40, G: 40, B: 40, A: 0xFF}), image.Point{}, draw.Src)

	return tile, nil
}
