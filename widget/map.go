package widget

import (
	"image"
	"math"
	"net/http"
	"net/url"

	"github.com/nfnt/resize"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"golang.org/x/image/draw"
)

const tileSize = 256

// Map widget renders an interactive map using OpenStreetMap tile data.
type Map struct {
	widget.BaseWidget

	pixels     *image.NRGBA
	w, h       int
	zoom, x, y int

	cl *http.Client

	tileSource      string // url to download xyz tiles (example: "https://tile.openstreetmap.org/%d/%d/%d.png")
	hideAttribution bool   // enable copyright disclaimer
	disclaimerLabel string // label for disclaimer (example: "OpenStreetMap")
	disclaimerUrl   string // url for disclaimer (example: "https://openstreetmap.org")
	hideZoomButtons bool   // enable zoom buttons
	hideMoveButtons bool   // enable move map buttons
}

type MapOption func(*Map)

func WithOsmTiles() MapOption {
	return func(m *Map) {
		m.tileSource = "https://tile.openstreetmap.org/%d/%d/%d.png"
		m.disclaimerLabel = "OpenStreetMap"
		m.disclaimerUrl = "https://openstreetmap.org"
		m.hideAttribution = false
	}
}

func WithTileSource(tileSource string) MapOption {
	return func(m *Map) {
		m.tileSource = tileSource
	}
}

func WithAttribution(enable bool, label, url string) MapOption {
	return func(m *Map) {
		m.hideAttribution = !enable
		m.disclaimerLabel = label
		m.disclaimerUrl = url
	}
}

func WithZoomButtons(enable bool) MapOption {
	return func(m *Map) {
		m.hideZoomButtons = !enable
	}
}

func WithScrollButtons(enable bool) MapOption {
	return func(m *Map) {
		m.hideMoveButtons = !enable
	}
}

// NewMap creates a new instance of the map widget.
func NewMap() *Map {
	m := &Map{cl: &http.Client{}}
	WithOsmTiles()(m)
	m.ExtendBaseWidget(m)
	return m
}

// NewMap creates a new instance of the map widget with provided map options.
func NewMapWithOptions(opts ...MapOption) *Map {
	m := NewMap()
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// MinSize returns the smallest possible size for a widget.
// For our map this is a constant size representing a single tile on a device with
// the highest known DPI (4x).
func (m *Map) MinSize() fyne.Size {
	return fyne.NewSize(64, 64)
}

// CreateRenderer returns the renderer for this widget.
// A map renderer is simply the map Raster with user interface elements overlaid.
func (m *Map) CreateRenderer() fyne.WidgetRenderer {
	var zoom fyne.CanvasObject
	if !m.hideZoomButtons {
		zoom = container.NewVBox(
			newMapButton(theme.ZoomInIcon(), func() {
				if m.zoom >= 19 {
					return
				}
				m.zoom++
				m.x *= 2
				m.y *= 2
				m.Refresh()
			}),
			newMapButton(theme.ZoomOutIcon(), func() {
				if m.zoom <= 0 {
					return
				}
				m.zoom--
				m.x /= 2
				m.y /= 2
				m.Refresh()
			}))
	}

	var move fyne.CanvasObject
	if !m.hideMoveButtons {
		buttonLayout := container.NewGridWithColumns(3, layout.NewSpacer(),
			newMapButton(theme.MoveUpIcon(), func() {
				m.y--
				m.Refresh()
			}), layout.NewSpacer(),
			newMapButton(theme.NavigateBackIcon(), func() {
				m.x--
				m.Refresh()
			}), layout.NewSpacer(),
			newMapButton(theme.NavigateNextIcon(), func() {
				m.x++
				m.Refresh()
			}), layout.NewSpacer(),
			newMapButton(theme.MoveDownIcon(), func() {
				m.y++
				m.Refresh()
			}), layout.NewSpacer())
		move = container.NewVBox(buttonLayout)
	}

	var copyright fyne.CanvasObject
	if !m.hideAttribution {
		license, _ := url.Parse(m.disclaimerUrl)
		view := widget.NewHyperlink(m.disclaimerLabel, license)
		view.Alignment = fyne.TextAlignTrailing
		copyright = view
	}

	overlay := container.NewBorder(nil, copyright, move, zoom)

	c := container.NewMax(canvas.NewRaster(m.draw), overlay)
	return widget.NewSimpleRenderer(c)
}

func (m *Map) draw(w, h int) image.Image {
	scale := 1
	tileSize := tileSize
	// TODO use retina tiles once OSM supports it in their server (text scaling issues)...
	if c := fyne.CurrentApp().Driver().CanvasForObject(m); c != nil {
		scale = int(c.Scale())
		if scale < 1 {
			scale = 1
		}
		tileSize = tileSize * scale
	}

	if m.w != w || m.h != h {
		m.pixels = image.NewNRGBA(image.Rect(0, 0, w, h))
	}

	midTileX := (w - tileSize*2) / 2
	midTileY := (h - tileSize*2) / 2
	if m.zoom == 0 {
		midTileX += tileSize / 2
		midTileY += tileSize / 2
	}

	count := 1 << m.zoom
	mx := m.x + int(float32(count)/2-0.5)
	my := m.y + int(float32(count)/2-0.5)
	firstTileX := mx - int(math.Ceil(float64(midTileX)/float64(tileSize)))
	firstTileY := my - int(math.Ceil(float64(midTileY)/float64(tileSize)))

	for x := firstTileX; (x-firstTileX)*tileSize <= w+tileSize; x++ {
		for y := firstTileY; (y-firstTileY)*tileSize <= h+tileSize; y++ {
			if x < 0 || y < 0 || x >= int(count) || y >= int(count) {
				continue
			}

			src, err := getTile(m.tileSource, x, y, m.zoom, m.cl)
			if err != nil {
				fyne.LogError("tile fetch error", err)
				continue
			}

			pos := image.Pt(midTileX+(x-mx)*tileSize,
				midTileY+(y-my)*tileSize)
			scaled := src
			if scale > 1 {
				scaled = resize.Resize(uint(tileSize), uint(tileSize), src, resize.Lanczos2)
			}
			draw.Copy(m.pixels, pos, scaled, image.Rect(0, 0, tileSize, tileSize), draw.Over, nil)
		}
	}

	return m.pixels
}
