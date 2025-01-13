package widget

import (
	"fmt"
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

func (b *spinnerBox) Move(pos fyne.Position) {
	b.box.Move(pos)
	b.border.Move(pos)
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
	return fyne.NewSize(20, 30)

}

// Objects returns the objects that make up the spinnerBox
func (r *spinnerBoxRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh redisplays the spinnerBox.
func (r *spinnerBoxRenderer) Refresh() {
	r.b.box.FillColor = color.Gray{Y: 128}
}

var _ fyne.Tappable = (*spinnerButton)(nil)

// spinnerButton is a button used to increment or decrement the value in the Spinner.
type spinnerButton struct {
	widget.DisableableWidget
	spinner *Spinner

	background *canvas.Rectangle

	position fyne.Position
	size     fyne.Size

	OnTapped func() `json:"-"`
}

func newSpinnerButton(s *Spinner, tapped func()) *spinnerButton {
	button := &spinnerButton{spinner: s, OnTapped: tapped}
	button.background = canvas.NewRectangle(color.Gray{192})
	button.ExtendBaseWidget(button)
	return button
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer.
func (b *spinnerButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	//	b.background = canvas.NewRectangle(color.Gray{192})
	objects := []fyne.CanvasObject{b.background}
	c := container.NewWithoutLayout(b.background)
	r := &spinnerButtonRenderer{
		button:    b,
		container: c,
		objects:   objects,
	}
	return r
}

// MinSize returns the minimum size for the spinnerButton
func (b *spinnerButton) MinSize() fyne.Size {
	return fyne.NewSize(12, 12)
}

// Tapped processes click events on the spinnerButton.
func (b *spinnerButton) Tapped(*fyne.PointEvent) {
	if onTapped := b.OnTapped; onTapped != nil {
		onTapped()
	}
}

func (b *spinnerButton) containsPoint(pos fyne.Position) bool {
	if pos.X < b.position.X || pos.X > b.position.X+b.size.Width {
		return false
	} else if pos.Y < b.position.Y || pos.Y > b.position.Y+b.size.Height {
		return false
	}
	return true
}

// Renderer for the spinnerButton
type spinnerButtonRenderer struct {
	button    *spinnerButton
	container *fyne.Container
	objects   []fyne.CanvasObject
}

// Destroy destroys any objects that are created for the spinnerButtonRenderer.
func (r *spinnerButtonRenderer) Destroy() {}

// Layout lays out the components of the spinnerButton.
func (r *spinnerButtonRenderer) Layout(_ fyne.Size) {}

// MinSize returns the minimum size of the spinnerButton.
// While a value is returned here, it is actually overridden
// in the spinnerLayout.
func (r *spinnerButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(50, 50)
}

// Objects returns the CanvasObjects that make up the spinnerButtton.
func (r *spinnerButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh redisplays the s[innerButton.
func (r *spinnerButtonRenderer) Refresh() {
	r.button.background.FillColor = color.Gray{Y: 192}
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

	topLeft.X += l.spinner.box.MinSize().Width
	l.spinner.upButton.background.Resize(l.spinner.upButton.MinSize())
	l.spinner.upButton.Move(topLeft)
	l.spinner.upButton.position = topLeft
	l.spinner.upButton.size = l.spinner.upButton.MinSize()
}

// MinSize returns the minimum size of the Spinner widget
func (l *spinnerLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(20, 30)
}

var _ fyne.Tappable = (*Spinner)(nil)

// Spinner widget has a minimum, maximum, step and current values along with spinnerButtons
// to increment and decrement the spinner value.
type Spinner struct {
	widget.DisableableWidget
	box *spinnerBox

	upButton *spinnerButton

	layout *spinnerLayout
}

// NewSpinner creates a new Spinner widget.
func NewSpinner() *Spinner {
	box := newSpinnerBox()
	s := &Spinner{box: box}
	s.layout = &spinnerLayout{spinner: s}
	s.upButton = newSpinnerButton(s, s.upButtonClicked)
	return s
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer.
func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	c := container.New(s.layout, s.box, s.upButton)
	return widget.NewSimpleRenderer(c)
}

// MinSize returns the minimum size of the Spinner object.
func (s *Spinner) MinSize() fyne.Size {
	return fyne.NewSize(30, 100)
}

// Tapped handles primary button clicks with the cursor over
// the Spinner object.
// If actually over one of the spinnerButtons, processing
// is passed to that button for handling.
func (b *Spinner) Tapped(evt *fyne.PointEvent) {
	fmt.Printf("evt = %v\n", evt)
	if b.upButton.containsPoint(evt.Position) {
		b.upButton.Tapped(evt)
	}
}

func (s *Spinner) upButtonClicked() {
	fmt.Println("upButtonClicked")
}
