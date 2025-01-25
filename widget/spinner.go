package widget

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// baseSpinnerButton widget is the base widget for intSpinnerButton and float64SpinnerButton.
type baseSpinnerButton struct {
	widget.Button

	position fyne.Position
	size     fyne.Size
}

// Move moves the button.
func (b *baseSpinnerButton) Move(pos fyne.Position) {
	b.position = pos
	b.BaseWidget.Move(pos)
}

// Resize resizes the button.
func (b *baseSpinnerButton) Resize(sz fyne.Size) {
	b.size = sz
	b.BaseWidget.Resize(sz)
}

// containsPoint is a helper method that is called to determine if the point
// is within the button. Returns true if point is within the button and
// false otherwise.
//
// Params:
//
//	pos is the position to check. This point is relative to the upper-left
//
// corner of the containing spinner widget.
func (b *baseSpinnerButton) containsPoint(pos fyne.Position) bool {
	if pos.X < b.position.X || pos.X > b.position.X+b.size.Width {
		return false
	} else if pos.Y < b.position.Y || pos.Y > b.position.Y+b.size.Height {
		return false
	}
	return true
}

// setButtonProperties sets the button properties.
//
// Params:
//
//	resource is the Resource for the button icon.
//	onTapped is the function to be called when the button is tapped.
func (b *baseSpinnerButton) setButtonProperties(resource fyne.Resource, onTapped func()) {
	b.Icon = resource
	b.Text = ""
	b.OnTapped = onTapped
}

// intSpinnerButton widget is a specialized button for use in the IntSpinner widget
type intSpinnerButton struct {
	baseSpinnerButton
	spinner *IntSpinner
}

// newIntSpinnerButton creates an intSpinnerButton widget. It should only be called from
// NewIntSpinner or NewIntSpinnerUninitialized.
//
// Params:
//
//	s is a pointer to the parent IntSpinner widget that this button is contained within.
//	onTapped is the function to be called when the button is "tapped".
func newIntSpinnerButton(s *IntSpinner, resource fyne.Resource, onTapped func()) *intSpinnerButton {
	b := &intSpinnerButton{
		spinner: s,
	}
	b.ExtendBaseWidget(b)
	b.setButtonProperties(resource, onTapped)
	return b
}

// MinSize returns the minimum (actual) size of the spinner button.
func (b *intSpinnerButton) MinSize() fyne.Size {
	th := b.Theme()
	h := b.spinner.MinSize().Height/2 - th.Size(theme.SizeNameInputBorder)
	return fyne.NewSize(h, h)
}

var _ fyne.Disableable = (*IntSpinner)(nil)
var _ fyne.Tappable = (*IntSpinner)(nil)
var _ fyne.Focusable = (*IntSpinner)(nil)
var _ desktop.Mouseable = (*IntSpinner)(nil)
var _ fyne.Scrollable = (*IntSpinner)(nil)

// IntSpinner widget has a minimum, maximum, step and current values along with spinnerButtons
// to increment and decrement the spinner value.
type IntSpinner struct {
	widget.DisableableWidget

	value       int
	min         int
	max         int
	step        int
	initialized bool

	upButton   *intSpinnerButton
	downButton *intSpinnerButton

	binder basicBinder

	hovered bool
	focused bool

	OnChanged func(int) `json:"-"`
}

// NewIntSpinner creates a new Spinner widget.
//
// Params:
//
//	min is the minimum spinner value. It may be < 0.
//	max is the maximum spinner value. It must be > min.
//	step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//	onChanged is the callback function that is called whenever the spinner value changes.
func NewIntSpinner(min, max, step int, onChanged func(int)) *IntSpinner {
	if min >= max {
		panic(errors.New("spinner max must be greater than min value"))
	}
	if step < 1 {
		panic(errors.New("spinner step must be greater than 0"))
	}
	if step > max-min {
		panic(errors.New("spinner step must be less than or equal to max - min"))
	}
	s := &IntSpinner{
		min:         min,
		max:         max,
		step:        step,
		initialized: true,
		OnChanged:   onChanged,
	}
	s.upButton = newIntSpinnerButton(s, theme.Icon(theme.IconNameArrowDropUp), s.upButtonClicked)
	s.downButton = newIntSpinnerButton(s, theme.Icon(theme.IconNameArrowDropDown), s.downButtonClicked)
	s.SetValue(s.min)
	return s
}

// NewSpinnerUninitialized returns a new uninitialized Spinner widget.
//
// An uninitialized Spinner widget is useful when you need to create a Spinner
// but the initial settings are unknown.
// Calling Enable on an uninitialized spinner will not enable the spinner; you
// must first call SetMinMaxStep to initialize spinner values before enabling
// the spinner widget.
func NewIntSpinnerUninitialized() *IntSpinner {
	s := &IntSpinner{initialized: false}
	s.upButton = newIntSpinnerButton(s, theme.Icon(theme.IconNameArrowDropUp), s.upButtonClicked)
	s.downButton = newIntSpinnerButton(s, theme.Icon(theme.IconNameArrowDropDown), s.downButtonClicked)
	s.Disable()
	return s
}

// NewIntSpinnerWithData returns a new Spinner widget connected to the specified data source.
//
// Params:
//
//	min is the minimum spinner value. It may be < 0.
//	max is the maximum spinner value. It must be > min.
//	step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//	data is the value that is bound to the spinner value.
func NewIntSpinnerWithData(min, max, step int, data binding.Int) *IntSpinner {
	s := NewIntSpinner(min, max, step, nil)
	s.Bind(data)
	s.OnChanged = func(_ int) {
		s.binder.CallWithData(s.writeData)
	}

	return s
}

// Bind connects the specified data source to this Spinner widget.
// The current value will be displayed and any changes in the data will cause the widget
// to update.
func (s *IntSpinner) Bind(data binding.Int) {
	s.binder.SetCallback(s.updateFromData)
	s.binder.Bind(data)
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer.
func (s *IntSpinner) CreateRenderer() fyne.WidgetRenderer {
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
	r := &intSpinnerRenderer{
		spinner: s,
		box:     box,
		border:  border,
		text:    text,
		objects: objects,
	}
	return r
}

// Disable disables the IntSpinner and its buttons.
func (s *IntSpinner) Disable() {
	s.downButton.Disable()
	s.upButton.Disable()
	s.DisableableWidget.Disable()
	s.Refresh()
}

// Enable enables the Spinner and its buttons as appropriate.
func (s *IntSpinner) Enable() {
	if !s.initialized {
		return
	}
	if s.GetValue() < s.max {
		s.upButton.Enable()
	}
	if s.GetValue() > s.min {
		s.downButton.Enable()
	}
	s.DisableableWidget.Enable()
	s.SetValue(s.value)
	s.Refresh()
}

// FocusGained is called when the Spinner has been given focus.
//
// Implements: fyne.Focusable
func (s *IntSpinner) FocusGained() {
	s.focused = true
	s.Refresh()
}

// FocusLost is called when the Spinner has had focus removed.
//
// Implements: fyne.Focusable
func (s *IntSpinner) FocusLost() {
	s.focused = false
	s.Refresh()
}

// GetValue retrieves the current Spinner value.
func (s *IntSpinner) GetValue() int {
	return s.value
}

func (s *IntSpinner) MinSize() fyne.Size {
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
func (s *IntSpinner) MouseDown(m *desktop.MouseEvent) {
	s.requestFocus()
	s.Refresh()
}

// MouseIn is called when a desktop pointer enters the widget.
func (s *IntSpinner) MouseIn(evt *desktop.MouseEvent) {
	s.hovered = true
	s.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget.
func (s *IntSpinner) MouseMoved(evt *desktop.MouseEvent) {}

// MouseOut is called when a desktop pointer exits the widget.
func (s *IntSpinner) MouseOut() {
	s.hovered = false
	s.Refresh()
}

// MouseUp called on mouse release.
//
// Implements: desktop.Mouseable
func (s *IntSpinner) MouseUp(m *desktop.MouseEvent) {}

// Scrolled handles mouse scroller events.
//
// Implements fyne.Scrollable
func (s *IntSpinner) Scrolled(evt *fyne.ScrollEvent) {
	if s.Disabled() || !s.focused {
		return
	}
	if evt.Scrolled.DY > 0 {
		s.SetValue(s.value + s.step)
	} else if evt.Scrolled.DY < 0 {
		s.SetValue(s.value - s.step)
	}
}

// SetMinMaxStep sets the widget's minimum, maximum, and step values.
//
// Params:
//
//	min is the minimum spinner value. It may be < 0.
//	max is the maximum spinner value. It must be > min.
//	step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//
// If the previously set value is less than min, then the value is set to min.
// If the previously set value is greater than max, then the value is set to max.
func (s *IntSpinner) SetMinMaxStep(min, max, step int) {
	if max <= min {
		panic(errors.New("spinner max must be greater than min value"))
	}
	if step < 1 {
		panic(errors.New("spinner step must be greater than 0"))
	}
	if step > max-min {
		panic(errors.New("spinner step must be less than or equal to max - min"))
	}
	s.min = min
	s.max = max
	s.step = step
	s.initialized = true

	if s.value < s.min {
		s.SetValue(s.min)
	} else if s.value > s.max {
		s.SetValue(s.max)
	}
}

// SetValue sets the spinner value. It ensures that the value is always >= min and
// <= max.
func (s *IntSpinner) SetValue(val int) {
	if s.Disabled() {
		return
	}
	s.value = val
	if s.value >= s.max {
		s.value = s.max
		s.upButton.Disable()
	} else {
		s.upButton.Enable()
	}
	if s.value <= s.min {
		s.value = s.min
		s.downButton.Disable()
	} else {
		s.downButton.Enable()
	}
	s.Refresh()
	if s.OnChanged != nil {
		s.OnChanged(s.value)
	}
}

// Tapped handles primary button clicks with the cursor over
// the Spinner object.
// If actually over one of the spinnerButtons, processing
// is passed to that button for handling.
func (s *IntSpinner) Tapped(evt *fyne.PointEvent) {
	if s.Disabled() {
		return
	}
	if s.upButton.containsPoint(evt.Position) {
		s.upButton.Tapped(evt)
	} else if s.downButton.containsPoint(evt.Position) {
		s.downButton.Tapped(evt)
	}
}

// TypedKey receives key input events when the Spinner widget has focus.
// Increments/decrements the Spinner's value when the up or down key is
// pressed.
//
// Implements: fyne.Focusable
func (s *IntSpinner) TypedKey(key *fyne.KeyEvent) {
	if s.Disabled() || !s.focused {
		return
	}
	switch key.Name {
	case fyne.KeyUp:
		s.SetValue(s.value + s.step)
	case fyne.KeyDown:
		s.SetValue(s.value - s.step)
	default:
		return
	}
}

// TypedRune receives text input events when the Spinner widget is focused.
// Increments/decrements the Spinner's value when the '+' or '-' key is
// pressed.
//
// Implements: fyne.Focusable
func (s *IntSpinner) TypedRune(rune rune) {
	if s.Disabled() || !s.focused {
		return
	}
	switch rune {
	case '+':
		s.SetValue(s.value + s.step)
	case '-':
		s.SetValue(s.value - s.step)
	default:
		return
	}
}

// Unbind disconnects any configured data source from this Spinner.
// The current value will remain at the last value of the data source.
func (s *IntSpinner) Unbind() {
	s.binder.Unbind()
}

func (s *IntSpinner) requestFocus() {
	if c := fyne.CurrentApp().Driver().CanvasForObject(s); c != nil {
		c.Focus(s)
	}

}

// Calculate the max size of the text that can be displayed for the Spinner.
// The size cannot be larger than the larger of the sizes for the Spinner
// min and max values.
func (s *IntSpinner) textSize() fyne.Size {
	minText := canvas.NewText(strconv.Itoa(s.min), color.Black)
	maxText := canvas.NewText(strconv.Itoa(s.max), color.Black)
	minTextSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(minText.Text,
		minText.TextSize, minText.TextStyle, minText.FontSource)
	maxTextSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(maxText.Text,
		maxText.TextSize, maxText.TextStyle, maxText.FontSource)
	return fyne.NewSize(max(minTextSize.Width, maxTextSize.Width),
		max(minTextSize.Height, maxTextSize.Height))
}

// updateFromData updates the Spinner to the value set in the bound data.
func (s *IntSpinner) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	textSource, ok := data.(binding.Int)
	if !ok {
		return
	}
	val, err := textSource.Get()
	if err != nil {
		fyne.LogError("Error getting current data value", err)
		return
	}
	s.SetValue(val)
}

// writeData updates the bound data item as the result of changes in the Spinner value.
func (s *IntSpinner) writeData(data binding.DataItem) {
	if data == nil {
		return
	}
	intTarget, ok := data.(binding.Int)
	if !ok {
		return
	}
	currentValue, err := intTarget.Get()
	if err != nil {
		return
	}
	if currentValue != s.GetValue() {
		err := intTarget.Set(s.GetValue())
		if err != nil {
			fyne.LogError(fmt.Sprintf("Failed to set binding value to %d", s.GetValue()), err)
		}
	}
}

// spinnerRenderer is the renderer for the Spinner widget
type intSpinnerRenderer struct {
	spinner *IntSpinner
	box     *canvas.Rectangle
	border  *canvas.Rectangle
	text    *canvas.Text
	objects []fyne.CanvasObject
}

// Destroy destroys any objects that must be destroyed when the renderer is
// destroyed.
func (r *intSpinnerRenderer) Destroy() {}

// Layout positions and sizes all of the objects that make up the Spinner widget.
func (r *intSpinnerRenderer) Layout(size fyne.Size) {
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
	buttonSize := fyne.NewSize((textSize.Height+padding)/2-1,
		(textSize.Height+padding)/2-1)
	xPos := size.Width - buttonSize.Width - borderSize - padding/2
	yPos := (rMinSize.Height - textSize.Height) / 2
	r.text.Move(fyne.NewPos(xPos, yPos))

	xPos += padding / 4
	yPos -= theme.Padding()
	r.spinner.upButton.Resize(buttonSize)
	r.spinner.upButton.Move(fyne.NewPos(xPos, yPos))

	yPos = r.spinner.upButton.MinSize().Height + padding/4 - 1
	r.spinner.downButton.Resize(buttonSize)
	r.spinner.downButton.Move(fyne.NewPos(xPos, yPos))
}

// MinSize returns the minimum size of the Spinner widget.
func (r *intSpinnerRenderer) MinSize() fyne.Size {
	return r.spinner.MinSize()
}

// Objects returns the objects associated with the Spinner renderer.
func (r *intSpinnerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh refreshes (redisplays) the Spinner widget.
func (r *intSpinnerRenderer) Refresh() {
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

	if r.spinner.Disabled() {
		r.spinner.upButton.Disable()
		r.spinner.downButton.Disable()
	} else {
		r.spinner.upButton.Enable()
		if r.spinner.GetValue() == r.spinner.max {
			r.spinner.upButton.Disable()
		}
		r.spinner.downButton.Enable()
		if r.spinner.GetValue() == r.spinner.min {
			r.spinner.downButton.Disable()
		}
	}
	r.spinner.upButton.Refresh()
	r.spinner.downButton.Refresh()
}

// spinnerColors returns the colors to display the button in.
// This is based on spinnerButtonRenderer.buttonColorNames, above.
func (r *intSpinnerRenderer) spinnerColors() (fgColor, bgColor, borderColor fyne.ThemeColorName) {
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

func (s *IntSpinner) upButtonClicked() {
	s.SetValue(s.value + s.step)
}

func (s *IntSpinner) downButtonClicked() {
	s.SetValue(s.value - s.step)
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
