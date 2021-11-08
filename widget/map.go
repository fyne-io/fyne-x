package widget

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"net/http"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/draw"
)

const tileSize = 256

type Map struct {
	widget.BaseWidget

	pixels     *image.NRGBA
	w, h       int
	zoom, x, y int
}

func NewMap() *Map {
	m := &Map{}
	m.ExtendBaseWidget(m)
	return m
}

func (m *Map) MinSize() fyne.Size {
	return fyne.NewSize(64, 64)
}

func (m *Map) CreateRenderer() fyne.WidgetRenderer {
	license, _ := url.Parse("https://openstreetmap.org")
	copyright := widget.NewHyperlink("OpenStreetMap", license)
	copyright.Alignment = fyne.TextAlignTrailing
	zoom := container.NewVBox(
		widget.NewButtonWithIcon("", theme.ZoomInIcon(), func() {
			if m.zoom >= 19 {
				return
			}
			m.zoom++
			m.x *= 2
			m.y *= 2
			m.Refresh()
		}),
		widget.NewButtonWithIcon("", theme.ZoomOutIcon(), func() {
			if m.zoom <= 0 {
				return
			}
			m.zoom--
			m.x /= 2
			m.y /= 2
			m.Refresh()
		}))

	move := container.NewGridWithColumns(3, layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
			m.y--
			m.Refresh()
		}), layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
			m.x--
			m.Refresh()
		}), layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
			m.x++
			m.Refresh()
		}), layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() {
			m.y++
			m.Refresh()
		}), layout.NewSpacer())

	overlay := container.NewBorder(nil, copyright, container.NewVBox(move), zoom)

	c := container.NewMax(canvas.NewRaster(m.draw), overlay)
	return widget.NewSimpleRenderer(c)
}

func (m *Map) draw(w, h int) image.Image {
	if m.w != w || m.h != h {
		m.pixels = image.NewNRGBA(image.Rect(0, 0, w, h))
	}

	midTileX := (w - tileSize*2) / 2
	midTileY := (h - tileSize*2) / 2
	if m.zoom == 0 {
		midTileX += tileSize / 2
		midTileY += tileSize / 2
	}

	count := math.Pow(2, float64(m.zoom))
	mx := m.x + int(count/2-0.5)
	my := m.y + int(count/2-0.5)
	firstTileX := mx - int(math.Ceil(float64(midTileX)/float64(tileSize)))
	firstTileY := my - int(math.Ceil(float64(midTileY)/float64(tileSize)))

	cl := &http.Client{}
	for x := firstTileX; (x-firstTileX)*tileSize <= w+tileSize/2; x++ {
		for y := firstTileY; (y-firstTileY)*tileSize <= h+tileSize/2; y++ {
			if x < 0 || y < 0 || x >= int(count) || y >= int(count) {
				continue
			}

			u := fmt.Sprintf("https://tile.openstreetmap.org/%d/%d/%d.png", m.zoom, x, y)
			req, err := http.NewRequest("GET", u, nil)
			if err != nil {
				fyne.LogError("download error", err)
				continue
			}

			req.Header.Set("User-Agent", "Fyne-X Map Widget/0.1")
			res, err := cl.Do(req)
			if err != nil {
				fyne.LogError("decode error", err)
				continue
			}

			src, err := png.Decode(res.Body)
			_ = res.Body.Close()
			if err != nil {
				fyne.LogError("decode error", err)
				continue
			}

			pos := image.Pt(midTileX+(x-mx)*tileSize,
				midTileY+(y-my)*tileSize)
			draw.Copy(m.pixels, pos, src, src.Bounds(), draw.Over, nil)
		}
	}

	return m.pixels
}