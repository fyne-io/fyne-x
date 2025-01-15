package widget

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Tappable = (*spinnerButton)(nil)

// spinnerButton is a button used to increment or decrement the value in the Spinner.
type spinnerButton struct {
	widget.DisableableWidget
	spinner *Spinner
	up      bool

	position fyne.Position
	size     fyne.Size

	OnTapped func() `json:"-"`
}

func newSpinnerButton(s *Spinner, up bool, onTapped func()) *spinnerButton {
	button := &spinnerButton{
		spinner:  s,
		up:       up,
		OnTapped: onTapped,
	}
	button.ExtendBaseWidget(button)
	return button
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer.
func (b *spinnerButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)

	r := &spinnerButtonRenderer{
		button: b,
	}
	th := b.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	r.background = canvas.NewRectangle(th.Color(theme.ColorNameButton, v))
	r.line1 = canvas.NewLine(th.Color(theme.ColorNameForeground, v))
	r.line2 = canvas.NewLine(th.Color(theme.ColorNameForeground, v))

	r.objects = []fyne.CanvasObject{
		r.background,
		r.line1,
		r.line2,
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
	if b.Disabled() {
		return
	}
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
	button  *spinnerButton
	objects []fyne.CanvasObject

	background *canvas.Rectangle
	line1      *canvas.Line
	line2      *canvas.Line
}

// Destroy destroys any objects that are created for the spinnerButtonRenderer.
func (r *spinnerButtonRenderer) Destroy() {}

// Layout lays out the components of the spinnerButton.
func (r *spinnerButtonRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	if r.button.up {
		r.line1.Position1 = fyne.NewPos(0.2*size.Width, 0.75*size.Height)
		r.line1.Position2 = fyne.NewPos(0.5*size.Width, 0.25*size.Height)
		r.line2.Position1 = fyne.NewPos(0.5*size.Width, 0.25*size.Height)
		r.line2.Position2 = fyne.NewPos(0.8*size.Width, 0.75*size.Height)
	} else {
		r.line1.Position1 = fyne.NewPos(0.2*size.Width, 0.25*size.Height)
		r.line1.Position2 = fyne.NewPos(0.5*size.Width, 0.75*size.Height)
		r.line2.Position1 = fyne.NewPos(0.5*size.Width, 0.75*size.Height)
		r.line2.Position2 = fyne.NewPos(0.8*size.Width, 0.25*size.Height)
	}
}

// MinSize returns the minimum (actual) size of the spinnerButton.
func (r *spinnerButtonRenderer) MinSize() fyne.Size {
	th := r.button.spinner.Theme()
	h := r.button.spinner.MinSize().Height/2 - th.Size(theme.SizeNameInputBorder)
	return fyne.NewSize(h, h)
}

// Objects returns the CanvasObjects that make up the spinnerButtton.
func (r *spinnerButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh redisplays the s[innerButton.
func (r *spinnerButtonRenderer) Refresh() {
	th := r.button.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	fgColor, bgColor, bgBlend := r.buttonColorNames()
	backColor := th.Color(bgColor, v)
	if bgBlend != "" {
		backColor = blendColor(backColor, th.Color(bgBlend, v))
	}
	r.background.FillColor = backColor
	r.line1.StrokeColor = th.Color(fgColor, v)
	r.line2.StrokeColor = th.Color(fgColor, v)

	r.background.Refresh()
	r.line1.Refresh()
	r.line2.Refresh()
}

// buttonColorNames returns the colors to display the button in.
// This is a copy of widget.Button.buttonColorNames in fyne v2.5.3.
func (r *spinnerButtonRenderer) buttonColorNames() (
	fgColor, bgColor, bgBlend fyne.ThemeColorName) {
	fgColor = theme.ColorNameForeground
	bgColor = theme.ColorNameButton
	if r.button.Disabled() {
		bgColor = theme.ColorNameDisabledButton
		fgColor = theme.ColorNameDisabled
	}
	return fgColor, bgColor, bgBlend
}

var _ fyne.Tappable = (*Spinner)(nil)
var _ fyne.Focusable = (*Spinner)(nil)
var _ desktop.Mouseable = (*Spinner)(nil)

// Spinner widget has a minimum, maximum, step and current values along with spinnerButtons
// to increment and decrement the spinner value.
type Spinner struct {
	widget.DisableableWidget

	value int
	min   int
	max   int
	step  int

	upButton   *spinnerButton
	downButton *spinnerButton

	hovered bool
	focused bool
}

// NewSpinner creates a new Spinner widget.
func NewSpinner(min, max, step int, tapped func()) *Spinner {
	s := &Spinner{
		min:   min,
		max:   max,
		step:  step,
		value: min,
	}
	s.upButton = newSpinnerButton(s, true, s.upButtonClicked)
	s.downButton = newSpinnerButton(s, false, s.downButtonClicked)
	return s
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer.
func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	th := s.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	box := canvas.NewRectangle(th.Color(theme.ColorNameBackground, v))
	border := canvas.NewRectangle(color.Transparent)

	text := canvas.NewText(strconv.Itoa(s.value), th.Color(theme.ColorNameForeground, v))

	objects := []fyne.CanvasObject{
		box,
		border,
		text,
		s.upButton,
		s.downButton,
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

// FocusGained is called when the Spinner has been given focus.
//
// Implements: fyne.Focusable
func (s *Spinner) FocusGained() {
	s.focused = true
	s.Refresh()
}

// FocusLost is called when the Spinner has had focus removed.
//
// Implements: fyne.Focusable
func (s *Spinner) FocusLost() {
	s.focused = false
	s.Refresh()
}

func (s *Spinner) MinSize() fyne.Size {
	th := s.Theme()
	padding := fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding) * 2)
	textSize := s.textSize()
	tHeight := textSize.Height + padding.Height
	upButtonHeight := (tHeight - padding.Height/2) / 2
	tWidth := textSize.Width + upButtonHeight + padding.Width
	return fyne.NewSize(tWidth, tHeight)
}

// MouseDown called on mouse click.
// This action causes the Spinner to request focus.
//
// Implements: desktop.Mouseable
func (s *Spinner) MouseDown(m *desktop.MouseEvent) {
	s.requestFocus()
	s.Refresh()
}

// MouseIn is called when a desktop pointer enters the widget.
func (s *Spinner) MouseIn(evt *desktop.MouseEvent) {
	s.hovered = true
	s.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget.
func (s *Spinner) MouseMoved(evt *desktop.MouseEvent) {

}

// MouseOut is called when a desktop pointer exits the widget.
func (s *Spinner) MouseOut() {
	s.hovered = false
	s.Refresh()
}

// MouseUp called on mouse release.
//
// Implements: desktop.Mouseable
func (s *Spinner) MouseUp(m *desktop.MouseEvent) {}

// SetValue sets the spinner value. It ensures that the value is always >= min and
// <= max.
func (s *Spinner) SetValue(val int) {
	if s.Disabled() {
		return
	}
	s.value = val
	if s.value >= s.max {
		s.value = s.max
	} else if s.value < s.min {
		s.value = s.min
	}
	s.Refresh()
}

// Tapped handles primary button clicks with the cursor over
// the Spinner object.
// If actually over one of the spinnerButtons, processing
// is passed to that button for handling.
func (s *Spinner) Tapped(evt *fyne.PointEvent) {
	if s.Disabled() {
		return
	}
	fmt.Printf("evt = %v\n", evt)
	if s.upButton.containsPoint(evt.Position) {
		s.upButton.Tapped(evt)
	} else if s.downButton.containsPoint(evt.Position) {
		s.downButton.Tapped(evt)
	}
}

// TypedKey receives key input events when the Spinner widget has focus.
//
// Implements: fyne.Focusable
func (s *Spinner) TypedKey(evt *fyne.KeyEvent) {
	// don't do anything yet.
}

// TypedRune receives text input events when the Spinner widget is focused.
//
// Implements: fyne.Focusable
func (s *Spinner) TypedRune(rune rune) {
	// don't do anything yet.
}

func (s *Spinner) requestFocus() {
	if c := fyne.CurrentApp().Driver().CanvasForObject(s); c != nil {
		c.Focus(s)
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
	buttonSize := fyne.NewSize((textSize.Height+padding)/2-1,
		(textSize.Height+padding)/2-1)
	r.spinner.upButton.Resize(buttonSize)
	r.spinner.upButton.Move(fyne.NewPos(xPos, yPos))

	yPos = r.spinner.upButton.MinSize().Height + padding/4 - 1
	r.spinner.downButton.Resize(buttonSize)
	r.spinner.downButton.Move(fyne.NewPos(xPos, yPos))
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

	fgColor, bgColor, borderColor := r.spinnerColors()
	r.box.FillColor = th.Color(bgColor, v)
	r.box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	r.border.CornerRadius = r.box.CornerRadius

	r.border.StrokeColor = th.Color(borderColor, v)
	r.text.Text = strconv.Itoa(r.spinner.value)
	r.text.Color = th.Color(fgColor, v)
	r.text.Refresh()
	r.text.Alignment = fyne.TextAlignTrailing

	r.spinner.upButton.Enable()
	r.spinner.downButton.Enable()
	if r.spinner.Disabled() {
		r.spinner.upButton.Disable()
		r.spinner.downButton.Disable()
	} else if r.spinner.value == r.spinner.min {
		r.spinner.downButton.Disable()
	} else if r.spinner.value == r.spinner.max {
		r.spinner.upButton.Disable()
	}
}

// spinnerColors returns the colors to display the button in.
// This is based on spinnerButtonRenderer.buttonColorNames, above.
func (r *spinnerRenderer) spinnerColors() (fgColor, bgColor, borderColor fyne.ThemeColorName) {
	fgColor = theme.ColorNameForeground
	bgColor = ""
	borderColor = theme.ColorNameInputBorder
	if r.spinner.Disabled() {
		fgColor = theme.ColorNameDisabled
		borderColor = theme.ColorNameDisabled
	} else if r.spinner.focused {
		borderColor = theme.ColorNamePrimary
	} else if r.spinner.hovered {
		bgColor = theme.ColorNameHover
	}
	return fgColor, bgColor, borderColor
}

func (s *Spinner) upButtonClicked() {
	s.SetValue(s.value + s.step)
}

func (s *Spinner) downButtonClicked() {
	s.SetValue(s.value - s.step)
}

// blendColor blends two colors together.
// This is a copy of blendcolor in button.go in fyne v2.5.3.
func blendColor(under, over color.Color) color.Color {
	// This alpha blends with the over operator, and accounts for RGBA() returning
	// alpha-premultiplied values
	dstR, dstG, dstB, dstA := under.RGBA()
	srcR, srcG, srcB, srcA := over.RGBA()

	srcAlpha := float32(srcA) // 0xFFFF
	dstAlpha := float32(dstA) // 0XFFFF

	outAlpha := srcAlpha * dstAlpha * (1 - srcAlpha)
	outR := srcR + uint32(float32(dstR)*(1-srcAlpha))
	outG := srcG + uint32(float32(dstG)*(1-srcAlpha))
	outB := srcB + uint32(float32(dstB)*(1-srcAlpha))
	// We create an RGBA64 here because the color components are already alpha-
	// premultiplied 16-bit values (they're just stored as uint32s).
	return color.RGBA64{
		R: uint16(outR),
		G: uint16(outG),
		B: uint16(outB),
		A: uint16(outAlpha * 0xFFFF),
	}
}

// max returns the larger of the two arguments.
// This can/should be replaced by the appropriate go max function
// when the version of go used to build fyne-x is updated to version
// 1.21 or later.
func max(a, b float32) float32 {
	max := a
	if a < b {
		max = b
	}
	return max
}
