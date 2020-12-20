package widget

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
)

type AnimatedGif struct {
	widget.BaseWidget

	src      *gif.GIF
	dst      *canvas.Image
	stopping bool
}

func NewAnimatedGif(u fyne.URI) (*AnimatedGif, error) {
	ret := &AnimatedGif{}
	ret.ExtendBaseWidget(ret)
	ret.dst = &canvas.Image{}
	ret.dst.FillMode = canvas.ImageFillContain

	if u == nil {
		return ret, nil
	}

	return ret, ret.Load(u)
}

func (g *AnimatedGif) Load(u fyne.URI) error {
	g.dst.Image = nil
	g.dst.Refresh()

	read, err := storage.OpenFileFromURI(u)
	if err != nil {
		return err
	}
	pix, err := gif.DecodeAll(read)
	if err != nil {
		return err
	}
	g.src = pix
	g.dst.Image = pix.Image[0]

	return nil
}

func (g *AnimatedGif) CreateRenderer() fyne.WidgetRenderer {
	return &gifRenderer{gif: g}
}

func (g *AnimatedGif) Stop() {
	g.stopping = true
}

func (g *AnimatedGif) Start() {
	buffer := image.NewNRGBA(g.dst.Image.Bounds())
	draw.Draw(buffer, g.dst.Image.Bounds(), g.src.Image[0], image.Point{}, draw.Over)
	g.dst.Image = buffer
	g.dst.Refresh()

	go func() {
		g.stopping = false

		for !g.stopping {
			for c, srcImg := range g.src.Image {
				if g.stopping {
					break
				}
				draw.Draw(buffer, g.dst.Image.Bounds(), srcImg, image.Point{}, draw.Over)
				g.dst.Refresh()

				time.Sleep(time.Millisecond * time.Duration(g.src.Delay[c]) * 10)
			}
		}

		g.dst.Image = nil
	}()
}

type gifRenderer struct{
	gif *AnimatedGif
}

func (g *gifRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (g *gifRenderer) Destroy() {
	g.gif.Stop()
}

func (g *gifRenderer) Layout(size fyne.Size) {
	g.gif.dst.Resize(size)
}

func (g *gifRenderer) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

func (g *gifRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{g.gif.dst}
}

func (g *gifRenderer) Refresh() {
	g.gif.dst.Refresh()
}
