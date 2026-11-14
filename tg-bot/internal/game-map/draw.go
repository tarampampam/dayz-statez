package game_map

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"tg-bot/internal/game-map/maps"
)

type (
	// drawPlayerLocationOptions holds the resolved settings for a DrawPlayer call.
	drawPlayerLocationOptions struct {
		marginXpx, marginYpx int // half-size of the viewport on each side of the player, in pixels

		drawPointsOfInterest bool // draw the in-view points of interest as faint watermarks?
		drawCoordinates      bool // draw the coordinates below the player?

		markerDiameterPx int         // player marker diameter, px
		markerColor      color.Color // player marker fill color

		nickFontPt  float64 // nickname label size, points
		coordFontPt float64 // coordinates label size, points
		poiFontPt   float64 // POI watermark size, points

		nickOffsetPx  int // distance between the player point and the nickname (drawn above it), px
		coordOffsetPx int // distance between the player point and the coordinates (drawn below it), px
	}

	// DrawPlayerOption customizes a DrawPlayer call.
	DrawPlayerOption func(*drawPlayerLocationOptions)
)

// WithDrawPlayerMarginPixels sets the half-size of the viewport kept on each side of the player, in pixels
// (the output is 2*x by 2*y, clamped to the map).
func WithDrawPlayerMarginPixels(x, y int) DrawPlayerOption {
	return func(o *drawPlayerLocationOptions) { o.marginXpx, o.marginYpx = x, y }
}

// WithDrawPlayerWithPointsOfInterest toggles drawing the in-view points of interest as faint watermarks.
func WithDrawPlayerWithPointsOfInterest(draw bool) DrawPlayerOption {
	return func(o *drawPlayerLocationOptions) { o.drawPointsOfInterest = draw }
}

// WithDrawPlayerCoordinates toggles drawing the coordinates below the player.
func WithDrawPlayerCoordinates(draw bool) DrawPlayerOption {
	return func(o *drawPlayerLocationOptions) { o.drawCoordinates = draw }
}

// WithDrawPlayerMarker sets the player marker diameter (px) and fill color.
func WithDrawPlayerMarker(diameterPx int, fill color.Color) DrawPlayerOption {
	return func(o *drawPlayerLocationOptions) { o.markerDiameterPx, o.markerColor = diameterPx, fill }
}

// WithDrawPlayerNicknameFont sets the nickname label size, in points.
func WithDrawPlayerNicknameFont(pt float64) DrawPlayerOption {
	return func(o *drawPlayerLocationOptions) { o.nickFontPt = pt }
}

// WithDrawPlayerCoordinatesFont sets the coordinates label size, in points.
func WithDrawPlayerCoordinatesFont(pt float64) DrawPlayerOption {
	return func(o *drawPlayerLocationOptions) { o.coordFontPt = pt }
}

// WithDrawPlayerPointsOfInterestFont sets the POI watermark size, in points.
func WithDrawPlayerPointsOfInterestFont(pt float64) DrawPlayerOption {
	return func(o *drawPlayerLocationOptions) { o.poiFontPt = pt }
}

// WithDrawPlayerNicknameOffset sets the distance between the player point and the nickname (drawn above it),
// in pixels.
func WithDrawPlayerNicknameOffset(px int) DrawPlayerOption {
	return func(o *drawPlayerLocationOptions) { o.nickOffsetPx = px }
}

// WithDrawPlayerCoordinatesOffset sets the distance between the player point and the coordinates (drawn below
// it), in pixels.
func WithDrawPlayerCoordinatesOffset(px int) DrawPlayerOption {
	return func(o *drawPlayerLocationOptions) { o.coordOffsetPx = px }
}

// DrawPlayer renders a crop of the map centered on the player at the given in-game coordinates (posX east,
// posY north, meters) and overlays the player marker with the nickname above it and the coordinates below.
// It returns an in-memory canvas ready to be encoded into any format.
func DrawPlayer(
	m maps.Map,
	nickname string,
	posX, posY float64,
	opts ...DrawPlayerOption,
) (*image.RGBA, error) {
	opt := drawPlayerLocationOptions{
		marginXpx:        900, //nolint:mnd
		marginYpx:        600, //nolint:mnd
		markerDiameterPx: 14,  //nolint:mnd
		markerColor:      color.RGBA{R: 0xEE, G: 0x22, B: 0x22, A: 0xFF},
		nickFontPt:       18, //nolint:mnd
		coordFontPt:      14, //nolint:mnd
		poiFontPt:        40, //nolint:mnd
		nickOffsetPx:     21, //nolint:mnd
		coordOffsetPx:    19, //nolint:mnd
		drawCoordinates:  true,
	}
	for _, o := range opts {
		o(&opt)
	}

	canvas, cErr := newMapCanvas(m, posX, posY, opt.marginXpx, opt.marginYpx)
	if cErr != nil {
		return nil, cErr
	}

	// draw POI watermarks
	if opt.drawPointsOfInterest {
		for _, l := range m.Locations() {
			if _, inView := canvas.LocToPoint(l.PosX(), l.PosY()); !inView {
				continue // not in the viewport, skip it
			}

			if err := drawBackgroundText(canvas, l.PosX(), l.PosY(), l.Name(), opt.poiFontPt); err != nil {
				return nil, err
			}
		}
	}

	// draw player marker
	drawPoint(canvas, posX, posY, opt.markerDiameterPx, opt.markerColor)

	// nickname goes above the player (negative Y), coordinates below it (positive Y)
	if err := drawForegroundText(canvas, posX, posY, nickname, opt.nickFontPt, withDrawOffsetPixels(0, -opt.nickOffsetPx)); err != nil {
		return nil, err
	}

	if opt.drawCoordinates {
		coords := fmt.Sprintf("[%d, %d]", int(math.Round(posX)), int(math.Round(posY)))
		if err := drawForegroundText(canvas, posX, posY, coords, opt.coordFontPt, withDrawOffsetPixels(0, opt.coordOffsetPx)); err != nil {
			return nil, err
		}
	}

	return canvas.Image(), nil
}

type (
	// drawOptions holds the tunable knobs shared by the draw* helpers.
	drawOptions struct {
		// xOffsetPx, yOffsetPx shift the drawing away from its anchor point, in canvas pixels. Positive x
		// moves it right, positive y moves it down (image Y grows downward); negatives move left/up.
		xOffsetPx, yOffsetPx int
	}

	// drawOption customizes a draw* call.
	drawOption func(*drawOptions)
)

// withDrawOffsetPixels shifts the drawing away from its anchor point by (x, y) canvas pixels (see drawOptions).
func withDrawOffsetPixels(x, y int) drawOption {
	return func(o *drawOptions) { o.xOffsetPx, o.yOffsetPx = x, y }
}

var (
	goBoldFont    = mustFont(gobold.TTF)
	goRegularFont = mustFont(goregular.TTF)
)

func mustFont(b []byte) *opentype.Font {
	f, err := opentype.Parse(b)
	if err != nil {
		panic(err)
	}

	return f
}

// drawBackgroundText stamps a big label centered (both axes) on the in-game point (x east, y north, meters).
// An optional pixel offset (withDrawOffsetPixels) nudges it away from that point.
func drawBackgroundText(c *mapCanvas, x, y float64, text string, fz float64, opts ...drawOption) error {
	if text == "" {
		return nil
	}

	var opt drawOptions
	for _, o := range opts {
		o(&opt)
	}

	// it never returns an error (checked in [opentype.NewFace] source)
	face, fErr := opentype.NewFace(goBoldFont, &opentype.FaceOptions{Size: fz, DPI: 72, Hinting: font.HintingFull})
	if fErr != nil {
		return fErr
	}
	defer func() { _ = face.Close() }()

	shadow := max(1, int(fz/12))
	metrics := face.Metrics()
	lw := font.MeasureString(face, text).Ceil() + shadow
	lh := (metrics.Ascent + metrics.Descent).Ceil() + shadow

	// render the text opaque onto a scratch layer (black shadow first, white body back over it)
	layer := image.NewRGBA(image.Rect(0, 0, lw, lh))
	base := fixed.Point26_6{X: 0, Y: metrics.Ascent}

	drawer := &font.Drawer{Dst: layer, Face: face}

	drawer.Src = image.NewUniform(color.Black)
	drawer.Dot = fixed.Point26_6{X: base.X + fixed.I(shadow), Y: base.Y + fixed.I(shadow)}
	drawer.DrawString(text)

	drawer.Src = image.NewUniform(color.White)
	drawer.Dot = base
	drawer.DrawString(text)

	// place the layer centered on the projected point (plus the optional offset) and blend it in at 20%
	center, _ := c.LocToPoint(x, y) // ignore the in-view flag (clipping handles off-canvas)
	center = center.Add(image.Pt(opt.xOffsetPx, opt.yOffsetPx))
	dstX, dstY := center.X-lw/2, center.Y-lh/2

	draw.DrawMask(
		c.Image(),
		image.Rect(dstX, dstY, dstX+lw, dstY+lh),
		layer,
		image.Point{},
		image.NewUniform(color.Alpha{A: 60}), // `A` is opacity 0-255
		image.Point{},
		draw.Over,
	)

	return nil
}

// drawForegroundText stamps a label centered (both axes) on the in-game point (x east, y north, meters).
// An optional pixel offset (withDrawOffsetPixels) nudges it away from that point.
func drawForegroundText(c *mapCanvas, x, y float64, text string, fz float64, opts ...drawOption) error {
	if text == "" {
		return nil
	}

	var opt drawOptions
	for _, o := range opts {
		o(&opt)
	}

	face, fErr := opentype.NewFace(goRegularFont, &opentype.FaceOptions{Size: fz, DPI: 72, Hinting: font.HintingFull})
	if fErr != nil {
		return fErr
	}
	defer func() { _ = face.Close() }()

	shadow := max(1, int(fz/16))
	metrics := face.Metrics()
	width := font.MeasureString(face, text)

	// center the glyph box on the projected point (plus the optional offset)
	center, _ := c.LocToPoint(x, y) // ignore the in-view flag (clipping handles off-canvas)
	center = center.Add(image.Pt(opt.xOffsetPx, opt.yOffsetPx))
	base := fixed.Point26_6{
		X: fixed.I(center.X) - width/2,
		Y: fixed.I(center.Y) - (metrics.Ascent+metrics.Descent)/2 + metrics.Ascent,
	}

	// opaque text draws straight onto the canvas: black shadow first, white body back over it
	drawer := &font.Drawer{Dst: c.Image(), Face: face}

	drawer.Src = image.NewUniform(color.Black)
	drawer.Dot = fixed.Point26_6{X: base.X + fixed.I(shadow), Y: base.Y + fixed.I(shadow)}
	drawer.DrawString(text)

	drawer.Src = image.NewUniform(color.White)
	drawer.Dot = base
	drawer.DrawString(text)

	return nil
}

// drawPoint stamps a filled dot of the given diameter (pixels) and fill color on the in-game point (x east, y north,
// meters), ringed by a white outline. The outline sits outside the dot and its width is 10% of the diameter (at least
// 1px), so the total size is diameter + 2*outline.
func drawPoint(c *mapCanvas, x, y float64, diameter int, fill color.Color, opts ...drawOption) {
	var opt drawOptions
	for _, o := range opts {
		o(&opt)
	}

	center, _ := c.LocToPoint(x, y) // ignore the in-view flag (clipping handles off-canvas)
	center = center.Add(image.Pt(opt.xOffsetPx, opt.yOffsetPx))

	var (
		core  = float64(diameter) / 2
		outer = core + math.Max(1, float64(diameter)/10) // + white ring, 10% of the diameter (min 1px)
		img   = c.Image()
		bound = int(outer) + 1
	)

	fr, fg, fb, _ := fill.RGBA()

	// blend a 16-bit channel from white toward the fill by t, returning an 8-bit value
	toFill := func(ch uint32, t float64) float64 { return (float64(ch)*t + 0xFFFF*(1-t)) / 257 }

	for dy := -bound; dy <= bound; dy++ {
		for dx := -bound; dx <= bound; dx++ {
			d := math.Sqrt(float64(dx*dx + dy*dy))

			shape := math.Min(1, math.Max(0, outer+0.5-d)) // anti-aliased outer edge (dot vs map)
			if shape == 0 {
				continue
			}
			fillT := math.Min(1, math.Max(0, core+0.5-d)) // anti-aliased inner edge (core vs white ring)

			// color = white ring leaped to the fill across the core edge, composited over the map by shape
			px, py := center.X+dx, center.Y+dy
			bg := img.RGBAAt(px, py)
			img.SetRGBA(px, py, color.RGBA{
				R: uint8(toFill(fr, fillT)*shape + float64(bg.R)*(1-shape)),
				G: uint8(toFill(fg, fillT)*shape + float64(bg.G)*(1-shape)),
				B: uint8(toFill(fb, fillT)*shape + float64(bg.B)*(1-shape)),
				A: 0xFF,
			})
		}
	}
}

// drawBitmap composites a bitmap (e.g. a decoded PNG icon with an alpha channel) centered on the in-game
// point (x east, y north, meters). Alpha is respected (draw.Over) and the bitmap is clipped at the canvas
// edge. An optional pixel offset (withDrawOffsetPixels) nudges it away from that point.
func drawBitmap(c *mapCanvas, x, y float64, bitmap image.Image, opts ...drawOption) {
	var opt drawOptions
	for _, o := range opts {
		o(&opt)
	}

	center, _ := c.LocToPoint(x, y) // ignore the in-view flag (clipping handles off-canvas)
	center = center.Add(image.Pt(opt.xOffsetPx, opt.yOffsetPx))

	b := bitmap.Bounds()
	m := center.Sub(image.Pt(b.Dx()/2, b.Dy()/2)) // top-left so the bitmap is centered on the point

	draw.Draw(c.Image(), image.Rectangle{Min: m, Max: m.Add(b.Size())}, bitmap, b.Min, draw.Over)
}
