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

// Move moves the button.
func (b *spinnerButton) Move(pos fyne.Position) {
	b.position = pos
	b.BaseWidget.Move(pos)
}

// Resize resizes the button.
func (b *spinnerButton) Resize(sz fyne.Size) {
	b.size = sz
	b.BaseWidget.Resize(sz)
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
func (r *spinnerButtonRenderer) Layout(size fyne.Size) {
	//	r.button.background.Move(fyne.NewPos(0, 0))
	r.button.background.Resize(size)
}

// MinSize returns the minimum size of the spinnerButton.
// While a value is returned here, it is actually overridden
// in the spinnerLayout.
func (r *spinnerButtonRenderer) MinSize() fyne.Size {
	h := r.button.spinner.MinSize().Height / 2
	return fyne.NewSize(h, h)
}

// Objects returns the CanvasObjects that make up the spinnerButtton.
func (r *spinnerButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh redisplays the s[innerButton.
func (r *spinnerButtonRenderer) Refresh() {
	r.button.background = canvas.NewRectangle(color.Gray{Y: 192})
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

func (s *Spinner) MinSize() fyne.Size {
	th := s.Theme()
	padding := fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding) * 2)
	textSize := s.textSize()
	tHeight := textSize.Height + padding.Height
	upButtonHeight := tHeight - padding.Height/2
	tWidth := textSize.Width + upButtonHeight + padding.Width
	return fyne.NewSize(tWidth, tHeight)
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

// Calculate the max size of the text that can be displayed for the Spinner.
// The size cannot be larger than the larger of the sizes for the Spinner
// min and max values.
func (s *Spinner) textSize() fyne.Size {
	minText := canvas.NewText(strconv.Itoa(s.min), color.Black)
	maxText := canvas.NewText(strconv.Itoa(s.max), color.Black)
	minTextSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(minText.Text,
		minText.TextSize, minText.TextStyle, minText.FontSource)
	maxTextSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(maxText.Text,
		maxText.TextSize, maxText.TextStyle, maxText.FontSource)
	return fyne.NewSize(max(minTextSize.Width, maxTextSize.Width),
		max(minTextSize.Height, maxTextSize.Height))
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
	r.spinner.Refresh()
	th := r.spinner.Theme()
	borderSize := th.Size(theme.SizeNameInputBorder)
	padding := th.Size(theme.SizeNameInnerPadding)

	// 0.5 is removed so on low DPI it rounds down on the trailing edge
	newSize := fyne.NewSize(size.Width-0.5, size.Height-0.5)
	topLeft := fyne.NewPos(0, 0)
	r.box.Resize(newSize)
	r.box.Move(topLeft)
	r.border.Resize(newSize)
	r.border.StrokeWidth = borderSize
	r.border.Move(topLeft)

	textSize := r.spinner.textSize()
	rMinSize := r.MinSize()
	xPos := borderSize + padding + textSize.Width
	yPos := (rMinSize.Height - textSize.Height) / 2
	r.text.Move(fyne.NewPos(xPos, yPos))

	xPos += padding / 4
	yPos -= theme.Padding()
	r.spinner.upButton.Resize(fyne.NewSize((textSize.Height+padding)/2,
		(textSize.Height+padding)/2))
	r.spinner.upButton.Move(fyne.NewPos(xPos, yPos))
	r.spinner.upButton.Refresh()
	// TODO: add positioning of downButton
}

// MinSize returns the minimum size of the Spinner widget.
func (r *spinnerRenderer) MinSize() fyne.Size {
	return r.spinner.MinSize()
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

// max returns the larger of the two arguments.
// This can/should be replaced by the appropriate go max function
// when the version of go used to build fyne-x is updated to version
// 1.21 or later.
func max(a, b float32) float32 {
	max := a
	if b > a {
		max = b
	}
	return max
}
