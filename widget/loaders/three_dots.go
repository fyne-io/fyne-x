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

type ThreeDotsLoader struct {
	widget.DisableableWidget

	mux      sync.Mutex
	renderer *threeDotsRenderer
}

func NewThreeDotsLoader() *ThreeDotsLoader {
	t := &ThreeDotsLoader{}
	t.ExtendBaseWidget(t)
	return t
}

func (t *ThreeDotsLoader) Start() {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.renderer != nil {
		t.renderer.Start()
	}
}

func (t *ThreeDotsLoader) Stop() {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.renderer != nil {
		t.renderer.Stop()
	}
}

func (t *ThreeDotsLoader) CreateRenderer() fyne.WidgetRenderer {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.ExtendBaseWidget(t)
	t.renderer = &threeDotsRenderer{
		circle1: canvas.NewCircle(theme.ForegroundColor()),
		circle2: canvas.NewCircle(theme.DisabledColor()),
		circle3: canvas.NewCircle(theme.ForegroundColor()),
	}
	return t.renderer
}

type threeDotsRenderer struct {
	maxSize, minSize fyne.Size

	circle1, circle2, circle3          *canvas.Circle
	circle1pos, circle2pos, circle3pos fyne.Position

	circle13Color, circle2Color *fyne.Animation
	circle13Size, circle2Size   *fyne.Animation

	mux sync.Mutex
}

func (t *threeDotsRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{t.circle1, t.circle2, t.circle3}
}

func (t *threeDotsRenderer) Layout(size fyne.Size) {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.Stop()

	t.maxSize = fyne.NewSize(size.Width/4, size.Height/4)
	t.minSize = fyne.NewSize(size.Width/6.66, size.Height/6.66)

	t.circle1.Resize(t.maxSize)
	t.circle2.Resize(t.minSize)
	t.circle3.Resize(t.maxSize)

	delta := fyne.NewDelta(t.minSize.Width+theme.Padding(), 0)

	t.circle1pos = fyne.NewPos(theme.Padding(), theme.Padding())
	t.circle2pos = t.circle1pos.Add(delta)
	t.circle3pos = t.circle2pos.Add(delta)

	t.circle1.Move(t.circle1pos)
	t.circle2.Move(t.circle2pos)
	t.circle3.Move(t.circle3pos)

	t.circle13Color = canvas.NewColorRGBAAnimation(theme.ForegroundColor(), theme.DisabledColor(), time.Millisecond*800, t.updateCircle13Color)
	t.circle13Color.Curve = fyne.AnimationLinear
	t.circle13Color.AutoReverse = true
	t.circle13Color.RepeatCount = fyne.AnimationRepeatForever

	t.circle13Size = canvas.NewSizeAnimation(t.maxSize, t.minSize, time.Millisecond*800, t.updateCircle13Size)
	t.circle13Size.Curve = fyne.AnimationLinear
	t.circle13Size.AutoReverse = true
	t.circle13Size.RepeatCount = fyne.AnimationRepeatForever

	t.circle2Color = canvas.NewColorRGBAAnimation(theme.DisabledColor(), theme.ForegroundColor(), time.Millisecond*800, t.updateCircle2Color)
	t.circle2Color.Curve = fyne.AnimationLinear
	t.circle2Color.AutoReverse = true
	t.circle2Color.RepeatCount = fyne.AnimationRepeatForever

	t.circle2Size = canvas.NewSizeAnimation(t.minSize, t.maxSize, time.Millisecond*800, t.updateCircle2Size)
	t.circle2Size.Curve = fyne.AnimationLinear
	t.circle2Size.AutoReverse = true
	t.circle2Size.RepeatCount = fyne.AnimationRepeatForever

	t.Start()
}

func (t *threeDotsRenderer) updateCircle13Color(c color.Color) {
	t.circle1.FillColor = c
	t.circle3.FillColor = c

	t.circle1.Refresh()
	t.circle3.Refresh()
}

func (t *threeDotsRenderer) getOffset(s fyne.Size) fyne.Vector2 {
	offset := t.maxSize.Subtract(s)
	offset = fyne.NewSize(offset.Width/2, offset.Height/2)
	return offset
}

func (t *threeDotsRenderer) updateCircle13Size(s fyne.Size) {
	offset := t.getOffset(s)

	t.circle1.Resize(s)
	t.circle1.Position1 = t.circle1pos.Add(offset)

	t.circle3.Resize(s)
	t.circle3.Position1 = t.circle3pos.Add(offset)

	t.circle1.Refresh()
	t.circle3.Refresh()
}

func (t *threeDotsRenderer) updateCircle2Color(c color.Color) {
	t.circle2.FillColor = c
	t.circle2.Refresh()
}

func (t *threeDotsRenderer) updateCircle2Size(s fyne.Size) {
	t.circle2.Resize(s)
	t.circle2.Position1 = t.circle2pos.Add(t.getOffset(s))
	t.circle2.Refresh()
}

func (t *threeDotsRenderer) MinSize() fyne.Size { return fyne.NewSize(100, 20) }

func (t *threeDotsRenderer) Refresh() {
	for _, circle := range []*canvas.Circle{t.circle1, t.circle2, t.circle3} {
		circle.Refresh()
	}
}

func (t *threeDotsRenderer) Start() {
	for _, anim := range []*fyne.Animation{t.circle13Color, t.circle2Color, t.circle13Size, t.circle2Size} {
		if anim != nil {
			anim.Start()
		}
	}
}

func (t *threeDotsRenderer) Stop() {
	for _, anim := range []*fyne.Animation{t.circle13Color, t.circle2Color, t.circle13Size, t.circle2Size} {
		if anim != nil {
			anim.Stop()
		}
	}
}

func (t *threeDotsRenderer) Destroy() { t.Stop() }
