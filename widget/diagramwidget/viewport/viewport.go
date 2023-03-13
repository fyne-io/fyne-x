package viewport

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type viewportRenderer struct {
	viewport   *ViewportWidget
	statusText *canvas.Text
}

func (r *viewportRenderer) MinSize() fyne.Size {
	return fyne.Size{Width: float32(r.viewport.Width), Height: float32(r.viewport.Height)}
}

func (r *viewportRenderer) Layout(size fyne.Size) {
	r.Refresh()
}

func (r *viewportRenderer) ApplyTheme(size fyne.Size) {
}

func (r *viewportRenderer) Refresh() {
	r.statusText.Move(fyne.Position{X: 0, Y: 0})
	r.statusText.Text = fmt.Sprintf("x=%f y=%f zoom=%f", r.viewport.XOffset, r.viewport.YOffset, r.viewport.Zoom)

	for _, viewportObj := range r.viewport.Objects {
		viewportObj.Refresh(r.viewport)
	}

	// XXX: I think this might be causing Fyne to refresh the whole canvas,
	// since without this the ViewPort widgets don't seem to update
	// themselves ??? Might need Refresh() to also call some kind of
	// Refresh() function of the ViewportObjects.
	r.statusText.Refresh()
}

func (r *viewportRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *viewportRenderer) Destroy() {
}

func (r *viewportRenderer) Objects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0)
	for _, viewportObj := range r.viewport.Objects {
		objects = append(objects, viewportObj.CanvasObjects(r.viewport)...)
	}
	objects = append(objects, r.statusText)
	return objects
}

type ViewportWidget struct {
	widget.BaseWidget
	Width   int
	Height  int
	Zoom    float64
	XOffset float64
	YOffset float64
	Objects []ViewportObject
}

func (w *ViewportWidget) Tapped(ev *fyne.PointEvent) {
}

func (w *ViewportWidget) TappedSecondary(ev *fyne.PointEvent) {
}

func (w *ViewportWidget) CreateRenderer() fyne.WidgetRenderer {
	r := viewportRenderer{
		viewport:   w,
		statusText: canvas.NewText("status", color.RGBA{255, 255, 255, 255}),
	}

	r.Refresh()
	return &r
}

func NewViewportWidget(width, height int) *ViewportWidget {
	vp := &ViewportWidget{
		Width:   width,
		Height:  height,
		Zoom:    1.0,
		XOffset: 0,
		YOffset: 0,
		Objects: make([]ViewportObject, 0),
	}

	vp.ExtendBaseWidget(vp)
	return vp
}

func (w *ViewportWidget) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (w *ViewportWidget) DragEnd() {
	w.Refresh()
}

func (w *ViewportWidget) Dragged(event *fyne.DragEvent) {
	w.XOffset += float64(event.Dragged.DX) / w.Zoom
	w.YOffset += float64(event.Dragged.DY) / w.Zoom
	w.Refresh()
}

func (w *ViewportWidget) MouseIn(event *desktop.MouseEvent) {
}

func (w *ViewportWidget) MouseOut() {
}

func (w *ViewportWidget) MouseMoved(event *desktop.MouseEvent) {
}

func (w *ViewportWidget) Scrolled(ev *fyne.ScrollEvent) {
	if ev.Scrolled.DY > 0 {
		w.Zoom *= 1.15
	} else {
		w.Zoom *= 0.85
	}
	w.Refresh()
}

type ViewportObject interface {
	CanvasObjects(viewport *ViewportWidget) []fyne.CanvasObject
	Refresh(viewport *ViewportWidget)
}

type ViewportLine struct {
	obj         *canvas.Line
	X1          float64
	Y1          float64
	X2          float64
	Y2          float64
	StrokeColor color.Color
	StrokeWidth float64
}

func setLineEndpoints(l *canvas.Line, X1, Y1, X2, Y2 float64) {
	l.Move(fyne.NewPos(float32(X1), float32(Y1)))
	l.Resize(fyne.NewSize(float32(X2)-float32(X1), float32(Y2)-float32(Y1)))
}

func (l *ViewportLine) CanvasObjects(viewport *ViewportWidget) []fyne.CanvasObject {
	if l.obj != nil {
		return []fyne.CanvasObject{l.obj}
	}
	return []fyne.CanvasObject{}
}

func (l *ViewportLine) Refresh(viewport *ViewportWidget) {
	if l.obj == nil {
		l.obj = canvas.NewLine(l.StrokeColor)
	}

	setLineEndpoints(l.obj,
		(l.X1+viewport.XOffset)*viewport.Zoom,
		(l.Y1+viewport.YOffset)*viewport.Zoom,
		(l.X2+viewport.XOffset)*viewport.Zoom,
		(l.Y2+viewport.YOffset)*viewport.Zoom,
	)
	l.obj.StrokeWidth = float32(l.StrokeWidth * viewport.Zoom)
	l.obj.Hidden = false
}
