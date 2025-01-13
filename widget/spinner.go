package widget

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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

var _ fyne.Tappable = (*Spinner)(nil)

// Spinner widget has a minimum, maximum, step and current values along with spinnerButtons
// to increment and decrement the spinner value.
type Spinner struct {
	widget.DisableableWidget

	value int
	min   int
	max   int
	step  int

	upButton *spinnerButton

	// layout *spinnerLayout
}

// NewSpinner creates a new Spinner widget.
func NewSpinner(min, max, step int, tapped func()) *Spinner {
	s := &Spinner{
		min:   min,
		max:   max,
		step:  step,
		value: min,
	}
	//	s.layout = &spinnerLayout{spinner: s}
	s.upButton = newSpinnerButton(s, s.upButtonClicked)
	return s
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer.
func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	th := s.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	box := canvas.NewRectangle(th.Color(theme.ColorNameInputBackground, v))
	border := canvas.NewRectangle(color.Transparent)

	text := canvas.NewText(strconv.Itoa(s.value), th.Color(theme.ColorNameForeground, v))

	objects := []fyne.CanvasObject{
		box,
		border,
		text,
		s.upButton,
		// TODO: add downButton
	}
	r := &spinnerRenderer{
		spinner: s,
		box:     box,
		border:  border,
		text:    text,
		objects: objects,
	}
	return r
}

// MinSize returns the minimum size of the Spinner object.
func (s *Spinner) MinSize() fyne.Size {
	return fyne.NewSize(30, 20)
}

// Tapped handles+ primary button clicks with the cursor over
// the Spinner object.
// If actually over one of the spinnerButtons, processing
// is passed to that button for handling.
func (s *Spinner) Tapped(evt *fyne.PointEvent) {
	fmt.Printf("evt = %v\n", evt)
	if s.upButton.containsPoint(evt.Position) {
		s.upButton.Tapped(evt)
	}
}

func (s *Spinner) textSize() fyne.Size {
	// TODO: calculate size based on spinner value
	return fyne.NewSize(25, 15)
}

// spinnerRenderer is the renderer for the Spinner widget
type spinnerRenderer struct {
	spinner *Spinner
	box     *canvas.Rectangle
	border  *canvas.Rectangle
	text    *canvas.Text
	objects []fyne.CanvasObject
}

// Destroy destroys any objects that must be destroyed when the renderer is
// destroyed.
func (r *spinnerRenderer) Destroy() {}

// Layout positions and sizes all of the objects that make up the Spinner widget.
func (r *spinnerRenderer) Layout(size fyne.Size) {
	th := r.spinner.Theme()
	borderSize := th.Size(theme.SizeNameInputBorder)
	padding := th.Size(theme.SizeNameInnerPadding)

	// 0.5 is removed so on low DPI it rounds down on the trailing edge
	r.border.Resize(fyne.NewSize(size.Width-borderSize-0.5,
		size.Height-borderSize-0.5))
	r.border.StrokeWidth = borderSize
	r.border.Move(fyne.NewSquareOffsetPos(borderSize / 2))
	r.box.Resize(size.Subtract(fyne.NewSquareSize(borderSize * 2)))
	r.box.Move(fyne.NewSquareOffsetPos(borderSize))

	textSize := r.spinner.textSize()
	rMinSize := r.MinSize()
	xPos := borderSize + padding + textSize.Width
	yPos := (rMinSize.Height - textSize.Height) / 2
	r.text.Move(fyne.NewPos(xPos, yPos))

	xPos += padding
	yPos -= theme.Padding()
	r.spinner.upButton.Move(fyne.NewPos(xPos, yPos))
	// TODO: add positioning of downButton
	r.spinner.Refresh()
}

// MinSize returns the minimum size of the Spinner widget.
func (r *spinnerRenderer) MinSize() fyne.Size {
	th := r.spinner.Theme()
	padding := fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding) * 2)
	textSize := r.spinner.textSize()
	tWidth := textSize.Width + r.spinner.upButton.MinSize().Width + padding.Width
	tHeight := textSize.Height + padding.Height
	return fyne.NewSize(tWidth, tHeight)
}

// Objects returns the objects associated with the Spinner renderer.
func (r *spinnerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh refreshes (redisplays) the Spinner widget.
func (r *spinnerRenderer) Refresh() {
	th := r.spinner.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.box.FillColor = th.Color(theme.ColorNameInputBackground, v)
	r.box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	r.border.CornerRadius = r.box.CornerRadius

	r.border.StrokeColor = th.Color(theme.ColorNameInputBorder, v)
	r.text.Text = strconv.Itoa(r.spinner.value)
	r.text.Alignment = fyne.TextAlignTrailing
}

func (s *Spinner) upButtonClicked() {
	fmt.Println("upButtonClicked")
}
