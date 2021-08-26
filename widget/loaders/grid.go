package loaders

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type GridLoader struct {
	widget.DisableableWidget

	mux      sync.Mutex
	renderer *gridRenderer
}

func NewGridLoader() *GridLoader {
	t := &GridLoader{}
	t.ExtendBaseWidget(t)
	return t
}

func (t *GridLoader) Start() {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.renderer != nil {
		t.renderer.Start()
	}
}

func (t *GridLoader) Stop() {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.renderer != nil {
		t.renderer.Stop()
	}
}

func (t *GridLoader) CreateRenderer() fyne.WidgetRenderer {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.ExtendBaseWidget(t)
	t.renderer = &gridRenderer{
		circles: [9]*canvas.Circle{
			canvas.NewCircle(theme.ForegroundColor()),
			canvas.NewCircle(theme.ForegroundColor()),
			canvas.NewCircle(theme.ForegroundColor()),
			canvas.NewCircle(theme.DisabledColor()),
			canvas.NewCircle(theme.ForegroundColor()),
			canvas.NewCircle(theme.ForegroundColor()),
			canvas.NewCircle(theme.ForegroundColor()),
			canvas.NewCircle(theme.ForegroundColor()),
			canvas.NewCircle(theme.ForegroundColor()),
		},
		anims: [9]*fyne.Animation{},
	}
	return t.renderer
}

type gridRenderer struct {
	circles [9]*canvas.Circle
	anims   [9]*fyne.Animation

	rendermux    sync.Mutex
	animstartmux sync.Mutex
}

func (t *gridRenderer) objects() []fyne.CanvasObject {
	out := make([]fyne.CanvasObject, 9)
	for i, c := range t.circles {
		out[i] = c
	}
	return out
}

func (t *gridRenderer) updateCircleColor(idx int, c color.Color) {
	t.circles[idx].FillColor = c
	t.circles[idx].Refresh()
}

func (t *gridRenderer) getCirclePos(size fyne.Size, idx int) fyne.Position {
	var delta fyne.Delta

	switch idx {
	// Top row
	case 0:
		delta = fyne.NewDelta(0, 0)

	case 1:
		delta = fyne.NewDelta(size.Width, 0)

	case 2:
		delta = fyne.NewDelta(size.Width*2, 0)

	// Middle row
	case 3:
		delta = fyne.NewDelta(0, size.Height+theme.Padding())

	case 4:
		delta = fyne.NewDelta(size.Width, size.Height+theme.Padding())

	case 5:
		delta = fyne.NewDelta(size.Width*2, size.Height+theme.Padding())

	// Bottom row
	case 6:
		delta = fyne.NewDelta(0, size.Height*2+theme.Padding()*2)

	case 7:
		delta = fyne.NewDelta(size.Width, size.Height*2+theme.Padding()*2)

	case 8:
		delta = fyne.NewDelta(size.Width*2, size.Height*2+theme.Padding()*2)

	}

	return fyne.NewPos(0, 0).Add(delta)
}

func (t *gridRenderer) Objects() []fyne.CanvasObject { return t.objects() }

func (t *gridRenderer) Layout(size fyne.Size) {
	t.rendermux.Lock()
	defer t.rendermux.Unlock()

	t.Stop()

	circleSize := fyne.NewSize(size.Width/8.4, size.Height/8.4)
	for i, c := range t.circles {
		c.Resize(circleSize)
		c.Move(t.getCirclePos(circleSize, i))
		this := i
		t.anims[i] = canvas.NewColorRGBAAnimation(theme.ForegroundColor(), theme.DisabledColor(), time.Millisecond*800, func(c color.Color) {
			t.updateCircleColor(this, c)
		})
		t.anims[i].Curve = fyne.AnimationLinear
		t.anims[i].AutoReverse = true
		t.anims[i].RepeatCount = fyne.AnimationRepeatForever
	}

	go func() {
		t.animstartmux.Lock()
		defer t.animstartmux.Unlock()

		a := []int{0, 3, 8, 1, 5, 7, 4, 6, 2}

		var i int
		for range time.NewTicker(time.Millisecond * 100).C {
			t.anims[a[i]].Start()
			if i == 8 {
				return
			}
			i++
		}
	}()

	t.Refresh()
}

func (t *gridRenderer) MinSize() fyne.Size { return fyne.NewSize(110, 70) }

func (t *gridRenderer) Refresh() {
	for _, c := range t.circles {
		c.Refresh()
	}
}

func (t *gridRenderer) Start() {
	for _, anim := range t.anims {
		if anim != nil {
			anim.Start()
		}
	}
}

func (t *gridRenderer) Stop() {
	for _, anim := range t.anims {
		if anim != nil {
			anim.Stop()
		}
	}
}

func (t *gridRenderer) Destroy() { t.Stop() }
