package widget

import (
	"image"
	"image/draw"
	"image/gif"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// AnimatedGif widget shows a Gif image with many frames.
type AnimatedGif struct {
	widget.BaseWidget
	min fyne.Size

	src               *gif.GIF
	dst               *canvas.Image
	remaining         int
	stopping, running bool
	runLock           sync.RWMutex
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

// CreateRenderer loads the widget renderer for this widget. This is an internal requirement for Fyne.
func (g *AnimatedGif) CreateRenderer() fyne.WidgetRenderer {
	return &gifRenderer{gif: g}
}

// Load is used to change the gif file shown.
// It will change the loaded content and prepare the new frames for animation.
func (g *AnimatedGif) Load(u fyne.URI) error {
	g.dst.Image = nil
	g.dst.Refresh()

	read, err := storage.Reader(u)
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

// MinSize returns the minimum size that this GIF can occupy.
// Because gif images are measured in pixels we cannot use the dimensions, so this defaults to 0x0.
// You can set a minimum size if required using SetMinSize.
func (g *AnimatedGif) MinSize() fyne.Size {
	return g.min
}

// SetMinSize sets the smallest possible size that this AnimatedGif should be drawn at.
// Be careful not to set this based on pixel sizes as that will vary based on output device.
func (g *AnimatedGif) SetMinSize(min fyne.Size) {
	g.min = min
}

// Start begins the animation. The speed of the transition is controlled by the loaded gif file.
func (g *AnimatedGif) Start() {
	if g.isRunning() {
		return
	}
	g.runLock.Lock()
	defer g.runLock.Unlock()
	g.running = true

	buffer := image.NewNRGBA(g.dst.Image.Bounds())
	draw.Draw(buffer, g.dst.Image.Bounds(), g.src.Image[0], image.Point{}, draw.Over)
	g.dst.Image = buffer
	g.dst.Refresh()

	go func() {
		switch g.src.LoopCount {
		case -1: // don't loop
			g.remaining = 1
		case 0: // loop forever
			g.remaining = -1
		default:
			g.remaining = g.src.LoopCount + 1
		}

		for g.remaining != 0 {
			for c, srcImg := range g.src.Image {
				if g.isStopping() {
					break
				}
				draw.Draw(buffer, g.dst.Image.Bounds(), srcImg, image.Point{}, draw.Over)
				g.dst.Refresh()

				time.Sleep(time.Millisecond * time.Duration(g.src.Delay[c]) * 10)
			}
			g.remaining--
		}

		g.running = false
	}()
}

// Stop will request that the animation stops running, the last frame will remain visible
func (g *AnimatedGif) Stop() {
	g.runLock.Lock()
	g.stopping = true
	g.runLock.Unlock()
}

func (g *AnimatedGif) isStopping() bool {
	g.runLock.RLock()
	defer g.runLock.RUnlock()
	return g.stopping
}

func (g *AnimatedGif) isRunning() bool {
	g.runLock.RLock()
	defer g.runLock.RUnlock()
	return g.running
}

type gifRenderer struct {
	gif *AnimatedGif
}

func (g *gifRenderer) Destroy() {
	g.gif.Stop()
}

func (g *gifRenderer) Layout(size fyne.Size) {
	g.gif.dst.Resize(size)
}

func (g *gifRenderer) MinSize() fyne.Size {
	return g.gif.MinSize()
}

func (g *gifRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{g.gif.dst}
}

func (g *gifRenderer) Refresh() {
	g.gif.dst.Refresh()
}
