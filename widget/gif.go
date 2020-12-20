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

// AnimatedGif widget shows a Gif image with many frames.
type AnimatedGif struct {
	widget.BaseWidget

	src      *gif.GIF
	dst      *canvas.Image
	stopping bool
}

// NewAnimatedGif creates a new widget loaded to show the specified image.
// If there is an error loading the image it will be returned in the error value.
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

// Load is used to change the gif file shown.
// It will change the loaded content and prepare the new frames for animation.
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

// CreateRenderer loads the widget renderer for this widget. This is an internal requirement for Fyne.
func (g *AnimatedGif) CreateRenderer() fyne.WidgetRenderer {
	return &gifRenderer{gif: g}
}

// Start begins the animation. The speed of the transition is controlled by the loaded gif file.
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

// Stop will request that the animation stops running, the last frame will remain visible
func (g *AnimatedGif) Stop() {
	g.stopping = true
}

type gifRenderer struct {
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
