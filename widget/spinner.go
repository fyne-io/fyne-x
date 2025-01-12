package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// spinnerBox creates the Spinner background and border.
type spinnerBox struct {
	widget.BaseWidget
	box    *canvas.Rectangle
	border *canvas.Rectangle
}

func newSpinnerBox() *spinnerBox {
	box := canvas.NewRectangle(color.Black)
	border := canvas.NewRectangle(color.Transparent)
	b := &spinnerBox{
		box:    box,
		border: border,
	}
	b.ExtendBaseWidget(b)
	return b
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer.
func (b *spinnerBox) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	objects := []fyne.CanvasObject{b.box, b.border}
	c := container.NewWithoutLayout(b.box, b.border)
	r := &spinnerBoxRenderer{b: b, container: c, objects: objects}
	return r
}

type spinnerBoxRenderer struct {
	b         *spinnerBox
	container *fyne.Container
	objects   []fyne.CanvasObject
}

// Destroy destroys any objects created for the renderer.
func (r *spinnerBoxRenderer) Destroy() {}

// Layout lays out the components of the spinnerBox.
// Because the size and layout depend on other widgets in the
// Spinner, this functionality is handled by the spinnerLayout.
func (r *spinnerBoxRenderer) Layout(_ fyne.Size) {}

// MinSize returns the minimum size of the spinnerBox.
// While this method is called, its value is overridden by a call
// to spinnerLayout.MinSize().
func (r *spinnerBoxRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)

}

// Objects returns the objects that make up the spinnerBox
func (r *spinnerBoxRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh redisplays the spinnerBox.
func (r *spinnerBoxRenderer) Refresh() {
	r.b.box.FillColor = color.Gray{Y: 128}
}

// spinnerLayout is the layout for the Spinner widget.
type spinnerLayout struct {
	spinner *Spinner
}

// Layout resizes and positions each of the elements of the Spinner widget.
func (l *spinnerLayout) Layout(_ []fyne.CanvasObject, _ fyne.Size) {
	topLeft := fyne.NewPos(0, 0)

	size := l.MinSize([]fyne.CanvasObject{})
	l.spinner.box.box.Resize(size)
	l.spinner.box.box.Move(topLeft)
	l.spinner.box.border.Resize(size)
	l.spinner.box.border.Move(topLeft)
}

// MinSize returns the minimum size of the Spinner widget
func (l *spinnerLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(20, 30)
}

// Spinner widget has a minimum, maximum, step and current values along with spinnerButtons
// to increment and decrement the spinner value.
type Spinner struct {
	widget.DisableableWidget
	box *spinnerBox

	layout *spinnerLayout
}

// NewSpinner creates a new Spinner widget.
func NewSpinner() *Spinner {
	box := newSpinnerBox()
	s := &Spinner{box: box}
	s.layout = &spinnerLayout{spinner: s}
	return s
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer.
func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	c := container.New(s.layout, s.box)
	return widget.NewSimpleRenderer(c)
}
