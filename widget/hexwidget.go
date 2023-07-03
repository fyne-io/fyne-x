package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// segmentLookupTable is used by h.Set() - the i-th index into this table
// represents the raw value that should be sent to UpdateSegments to show
// the value i.
var segmentLookupTable []uint8 = []uint8{
	1 << 6,
	(1<<0 | (1 << 1) | (1 << 2) | (1 << 3) | (1 << 6) | (1 << 7)),
	(1<<2 | (1 << 5)),
	(1<<4 | (1 << 5)),
	(1<<0 | (1 << 3) | (1 << 4)),
	(1<<1 | (1 << 4)),
	(1 << 1),
	(1<<3 | (1 << 4) | (1 << 5) | (1 << 6)),
	0,
	(1<<3 | (1 << 4)),
	(1 << 3),
	(1<<0 | (1 << 1)),
	(1<<0 | (1 << 1) | (1 << 2) | (1 << 5)),
	(1<<0 | (1 << 5)),
	(1<<1 | (1 << 2)),
	(1<<1 | (1 << 2) | (1 << 3)),
}

// size of the hex widget
const defaultHexHeight float32 = 58.0
const defaultHexWidth float32 = defaultHexHeight * (8 / 14.0)

// slant angle
const defaultHexOffset float32 = 0.1 * defaultHexWidth

var defaultHexOnColor color.Color = color.RGBA{200, 25, 25, 255}
var defaultHexOffColor color.Color = color.RGBA{25, 15, 15, 64}

type hexRenderer struct {
	hex            *HexWidget
	segmentObjects []fyne.CanvasObject
}

func (h *hexRenderer) MinSize() fyne.Size {
	return fyne.NewSize(
		h.hex.size.Width+h.hex.hexOffset,
		h.hex.size.Height,
	)
}

func (h *hexRenderer) Layout(_ fyne.Size) {
	hexSegmentWidth := 0.2 * h.hex.size.Width

	hexSegmentVLength := (h.hex.size.Height - hexSegmentWidth) / 2
	hexSegmentHLength := h.hex.size.Width - hexSegmentWidth
	pos := fyne.NewPos(h.hex.hexOffset, hexSegmentWidth/2)

	pt0Center := fyne.NewPos(pos.X+h.hex.size.Width/2.0+h.hex.hexOffset, pos.Y)
	pt05 := fyne.NewPos(float32(pt0Center.X)-(hexSegmentHLength/2), pt0Center.Y)
	pt01 := fyne.NewPos(float32(pt0Center.X)+(hexSegmentHLength/2), pt0Center.Y)

	pt6Center := fyne.NewPos(pos.X+h.hex.size.Width/2.0, float32(pt0Center.Y)+hexSegmentVLength)
	pt65 := fyne.NewPos(float32(pt6Center.X)-(hexSegmentHLength/2), pt6Center.Y)
	pt61 := fyne.NewPos(float32(pt6Center.X)+(hexSegmentHLength/2), pt6Center.Y)

	pt3Center := fyne.NewPos(pos.X+h.hex.size.Width/2.0-h.hex.hexOffset, float32(pt0Center.Y)+2*hexSegmentVLength)
	pt34 := fyne.NewPos(float32(pt3Center.X)-(hexSegmentHLength/2), pt3Center.Y)
	pt32 := fyne.NewPos(float32(pt3Center.X)+(hexSegmentHLength/2), pt3Center.Y)

	setLineEndpoints(h.segmentObjects[0].(*canvas.Line), pt05, pt01)
	setLineEndpoints(h.segmentObjects[1].(*canvas.Line), pt01, pt61)
	setLineEndpoints(h.segmentObjects[2].(*canvas.Line), pt61, pt32)
	setLineEndpoints(h.segmentObjects[3].(*canvas.Line), pt32, pt34)
	setLineEndpoints(h.segmentObjects[4].(*canvas.Line), pt34, pt65)
	setLineEndpoints(h.segmentObjects[5].(*canvas.Line), pt65, pt05)
	setLineEndpoints(h.segmentObjects[6].(*canvas.Line), pt65, pt61)
}

func (h *hexRenderer) Refresh() {
	hexSegmentWidth := 0.2 * h.hex.size.Width

	h.segmentObjects[0].(*canvas.Line).StrokeWidth = float32(hexSegmentWidth / 2)
	h.segmentObjects[1].(*canvas.Line).StrokeWidth = float32(hexSegmentWidth / 2)
	h.segmentObjects[2].(*canvas.Line).StrokeWidth = float32(hexSegmentWidth / 2)
	h.segmentObjects[3].(*canvas.Line).StrokeWidth = float32(hexSegmentWidth / 2)
	h.segmentObjects[4].(*canvas.Line).StrokeWidth = float32(hexSegmentWidth / 2)
	h.segmentObjects[5].(*canvas.Line).StrokeWidth = float32(hexSegmentWidth / 2)
	h.segmentObjects[6].(*canvas.Line).StrokeWidth = float32(hexSegmentWidth / 2)

	for i, v := range h.segmentObjects {
		v.(*canvas.Line).StrokeColor = h.hex.getSegmentColor(i)
		canvas.Refresh(v)
	}
}

func (h *hexRenderer) Destroy() {
}

func (h *hexRenderer) Objects() []fyne.CanvasObject {
	return h.segmentObjects
}

// HexWidget represents a 7-segment hexadecimal display. The segments
// of the display mapped active-low onto 7 state bits, with segment 0 in
// the least significant bit.
//
//	     0
//	   -----
//	  |     |
//	5 |     | 1
//	  |  6  |
//	   -----
//	  |     |
//	4 |     | 2
//	  |  3  |
//	   -----
type HexWidget struct {
	widget.BaseWidget
	segments uint8

	// size of the hex widget
	size fyne.Size

	// slant angle
	hexOffset float32

	// color when the hex is on
	hexOnColor color.Color

	// color when the hex is off
	hexOffColor color.Color
}

// SetOnColor changes the color that segments are shown as when they are
// active/on.
func (h *HexWidget) SetOnColor(c color.Color) {
	h.hexOnColor = c
	h.Refresh()
}

// SetOffColor changes the color that segments are shown as when they are
// inactive/off.
func (h *HexWidget) SetOffColor(c color.Color) {
	h.hexOffColor = c
	h.Refresh()
}

// SetSize changes the size of the hex widget.
func (h *HexWidget) SetSize(s fyne.Size) {
	h.size = s
	h.Refresh()
}

// SetSlant changes the amount of "slant" in the hex widgets. The topmost
// segment is offset by slant many units to the right. A value of 0 means no
// slant at all. For example, setting the slant equal to the height should
// result in a 45 degree angle.
func (h *HexWidget) SetSlant(s float32) {
	h.hexOffset = s
	h.Refresh()
}

func (h *HexWidget) getSegmentColor(segno int) color.Color {
	if (h.segments & (1 << uint(segno))) == 0 {
		return h.hexOnColor
	}

	return h.hexOffColor
}

// CreateRenderer implements fyne.Widget
func (h *HexWidget) CreateRenderer() fyne.WidgetRenderer {

	seg0 := canvas.NewLine(h.hexOffColor)
	seg1 := canvas.NewLine(h.hexOffColor)
	seg2 := canvas.NewLine(h.hexOffColor)
	seg3 := canvas.NewLine(h.hexOffColor)
	seg4 := canvas.NewLine(h.hexOffColor)
	seg5 := canvas.NewLine(h.hexOffColor)
	seg6 := canvas.NewLine(h.hexOffColor)

	r := &hexRenderer{
		hex:            h,
		segmentObjects: []fyne.CanvasObject{seg0, seg1, seg2, seg3, seg4, seg5, seg6},
	}

	r.Refresh()

	return r
}

// NewHexWidget instantiates a new widget instance, with all of the segments
// disabled.
func NewHexWidget() *HexWidget {
	h := &HexWidget{
		segments:    0xff,
		size:        fyne.NewSize(defaultHexWidth, defaultHexHeight),
		hexOffset:   defaultHexOffset,
		hexOnColor:  defaultHexOnColor,
		hexOffColor: defaultHexOffColor,
	}

	h.ExtendBaseWidget(h)
	return h
}

// UpdateSegments changes the state of the segments and causes the widget to
// refresh so the changes are visible to the user. Segments values are packed
// into the 8-bit segments integer, see the documentation for HexWidget for
// more information on the appropriate packing.
func (h *HexWidget) UpdateSegments(segments uint8) {
	h.segments = segments
	h.Refresh()
}

// Set updates the hex widget to show a specific number between 0 and 15, which
// will be rendered in hexadecimal in 0...f. If the number is greater than 15,
// it will be modulo-ed by 16.
func (h *HexWidget) Set(val uint) {
	val = val % 16
	h.UpdateSegments(segmentLookupTable[val])
}

func setLineEndpoints(l *canvas.Line, pt1, pt2 fyne.Position) {
	l.Position1 = pt1
	l.Position2 = pt2
}
