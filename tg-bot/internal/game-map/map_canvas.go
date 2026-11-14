package game_map

import (
	"image"
	"image/draw"
	"math"

	"tg-bot/internal/game-map/maps"
)

// mapCanvas is the rendered crop of a map plus the knowledge of how in-game coordinates land on it.
// It owns the pixels and the game<->pixel projection.
//
// The whole math, start to finish:
//   - Game coords are in METERS (origin at the map's SW corner, Y grows north / up). The image is in PIXELS (origin
//     top-left, Y grows down; different units, and Y points the opposite way)
//   - scale = pixels per meter (MapSizePX / MapSize). It is the unit bridge - multiply any meter value by scale to
//     get its pixel offset on the FULL map. Without it meters and pixels can't be compared
//   - The canvas is only a CROP of the full map, so a point's canvas pixel is its full-map pixel minus the crop's
//     corner. originX and topY are that corner, pre-computed in full-map pixels
//   - X needs no flip (both axes grow east), so it is a plain offset - `canvasX = x*scale - originX`
//     (originX = full-map pixel of the crop's LEFT edge)
//   - Y is flipped (game up vs image down), so instead of subtracting an origin we count DOWN from the top edge -
//     `canvasY = topY - y*scale` (topY = full-map pixel of the crop's TOP edge)
//
// That asymmetry is exactly why the two anchors are named differently: originX is a start we subtract
// (an offset along X), topY is a baseline we subtract FROM (the Y-flip). Both are integer pixels, so
// LocToPoint is one multiply per axis and no division.
type mapCanvas struct {
	img *image.RGBA // the rendered crop

	scale   float64 // pixels per in-game meter (MapSizePX / MapSize); the meter->pixel bridge
	originX int     // full-map pixel of the crop's left edge; canvasX = x*scale - originX
	topY    int     // full-map pixel of the crop's top edge (highest game Y); canvasY = topY - y*scale
}

// newMapCanvas renders the crop of provided map, centered on the in-game point (locX east, locY north, meters) and
// returns it ready to be drawn on.
//
// The crop is a window of 2*marginXpx by 2*marginYpx pixels, near a map edge it slides fully inside the map instead
// of padding with empty space. If a margin made the window exceed the map, that axis is clamped to the map size, so
// the output can come out smaller than requested.
func newMapCanvas(m maps.Map, locX, locY float64, marginXpx, marginYpx int) (*mapCanvas, error) {
	mapPX := m.MapSizePX()
	scale := float64(mapPX) / float64(m.MapSize())

	// viewport is the margin on each side, clamped so it never exceeds the map
	vw, vh := min(2*marginXpx, mapPX), min(2*marginYpx, mapPX)

	// center in whole-map pixels (Y flipped - the top pixel row is the highest game Y), then slide the window fully
	// inside the map so we never address tiles outside the grid
	cx := locX * scale
	cy := float64(mapPX) - locY*scale
	ox := int(math.Round(max(0, min(cx-float64(vw)/2, float64(mapPX-vw)))))
	oy := int(math.Round(max(0, min(cy-float64(vh)/2, float64(mapPX-vh)))))

	img := image.NewRGBA(image.Rect(0, 0, vw, vh))

	// blit every tile overlapping the crop
	ts := m.TileSizePX()
	for ty := oy / ts; ty <= (oy+vh-1)/ts; ty++ {
		for tx := ox / ts; tx <= (ox+vw-1)/ts; tx++ {
			tile, err := m.Tile(tx, ty)
			if err != nil {
				return nil, err
			}

			dx, dy := tx*ts-ox, ty*ts-oy
			draw.Draw(img, image.Rect(dx, dy, dx+ts, dy+ts), tile, image.Point{}, draw.Src)
		}
	}

	return &mapCanvas{
		img:     img,
		scale:   scale,
		originX: ox,
		topY:    mapPX - oy, // whole-map pixel of the crop's top row == highest visible game Y
	}, nil
}

// Image returns the underlying canvas to draw on, encode, etc.
func (c *mapCanvas) Image() *image.RGBA { return c.img }

// LocToPoint projects an in-game position (x east, y north, meters) onto a pixel of the canvas. The point is always
// returned, even when it falls outside the image.
func (c *mapCanvas) LocToPoint(x, y float64) (image.Point, bool) {
	pt := image.Pt(
		int(math.Round(x*c.scale))-c.originX,
		c.topY-int(math.Round(y*c.scale)),
	)

	return pt, pt.In(c.img.Bounds())
}
