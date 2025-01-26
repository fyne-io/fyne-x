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

// MinSize returns the minimum size of the button. Because the minimum size is a constant
// based on the spinner height and theme properties, the minimum size is calculated when
// the button is created.
func (b *baseSpinnerButton) MinSize() fyne.Size {
	return fyne.NewSize(b.size.Height, b.size.Height)
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

	// calculate the minimum button size (really just its height).
	th := b.Theme()
	padding := fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding) * 2)
	text := canvas.NewText("0", color.Black)
	textSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(text.Text,
		text.TextSize, text.TextStyle, text.FontSource)
	tHeight := textSize.Height + padding.Height

	h := tHeight/2 - th.Size(theme.SizeNameInputBorder) - 2
	b.size = fyne.NewSize(h, h)

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
//		s is a pointer to the parent IntSpinner widget that this button is contained within.
//	 resource is the resource for the button icon.
//		onTapped is the function to be called when the button is "tapped".
func newIntSpinnerButton(s *IntSpinner, resource fyne.Resource, onTapped func()) *intSpinnerButton {
	b := &intSpinnerButton{
		spinner: s,
	}
	b.ExtendBaseWidget(b)
	b.setButtonProperties(resource, onTapped)
	return b
}

// / float64SpinnerButton widget is a specialized button for use in the float64Spinner widget
type float64SpinnerButton struct {
	baseSpinnerButton
	spinner *Float64Spinner
}

// newFloat64SpinnerButton creates a float64SpinnerButton widget. It should only be called from
// NewFloat64Spinner or NewFloat64SpinnerUninitialized.
//
// Params:
//
//	s is a pointer to the parent Float64Spinner widget that this button is contained within.
//	resource is the resource for the button icon. For the icon to be properly displayed
//
// for each button status (enabled/disabled, light/dark mode, etc.), the resource should be
// a theme resource.
//
//	onTapped is the function to be called when the button is "tapped".
func newFloat64SpinnerButton(s *Float64Spinner, resource fyne.Resource, onTapped func()) *float64SpinnerButton {
	b := &float64SpinnerButton{
		spinner: s,
	}
	b.ExtendBaseWidget(b)
	b.setButtonProperties(resource, onTapped)
	return b
}

// baseSpinner contains the basic functionality shared by IntSpinner and Float64Spinner.
type baseSpinner struct {
	widget.DisableableWidget
	initialized bool

	binder  basicBinder
	hovered bool
	focused bool
}

// FocusGained is called when the spinner has been given focus.
//
// Implements: fyne.Focusable
func (s *baseSpinner) FocusGained() {
	s.focused = true
	s.Refresh()
}

// FocusLost is called when the spinner has had focus removed.
//
// Implements: fyne.Focusable
func (s *baseSpinner) FocusLost() {
	s.focused = false
	s.Refresh()
}

// MouseIn is called when a desktop pointer enters the widget.
func (s *baseSpinner) MouseIn(evt *desktop.MouseEvent) {
	s.hovered = true
	s.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget.
func (s *baseSpinner) MouseMoved(evt *desktop.MouseEvent) {}

// MouseOut is called when a desktop pointer exits the widget.
func (s *baseSpinner) MouseOut() {
	s.hovered = false
	s.Refresh()
}

var _ fyne.Disableable = (*IntSpinner)(nil)
var _ fyne.Focusable = (*IntSpinner)(nil)
var _ fyne.Tappable = (*IntSpinner)(nil)
var _ desktop.Mouseable = (*IntSpinner)(nil)
var _ fyne.Scrollable = (*IntSpinner)(nil)

// IntSpinner widget has a minimum, maximum, step and current values along with spinnerButtons
// to increment and decrement the spinner value.
type IntSpinner struct {
	baseSpinner

	value int
	min   int
	max   int
	step  int

	upButton   *intSpinnerButton
	downButton *intSpinnerButton

	OnChanged func(int) `json:"-"`
}

// NewIntSpinner creates a new IntSpinner widget.
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
		min:       min,
		max:       max,
		step:      step,
		OnChanged: onChanged,
	}
	s.baseSpinner.initialized = true
	s.upButton = newIntSpinnerButton(s, theme.Icon(theme.IconNameArrowDropUp), s.upButtonClicked)
	s.downButton = newIntSpinnerButton(s, theme.Icon(theme.IconNameArrowDropDown), s.downButtonClicked)
	s.SetValue(s.min)
	return s
}

// NewIntSpinnerUninitialized returns a new uninitialized IntSpinner widget.
//
// An uninitialized IntSpinner widget is useful when you need to create a Spinner
// but the initial settings are unknown.
// Calling Enable on an uninitialized spinner will not enable the spinner; you
// must first call SetMinMaxStep to initialize spinner values before enabling
// the spinner widget.
func NewIntSpinnerUninitialized() *IntSpinner {
	s := &IntSpinner{}
	s.baseSpinner.initialized = false
	s.upButton = newIntSpinnerButton(s, theme.Icon(theme.IconNameArrowDropUp), s.upButtonClicked)
	s.downButton = newIntSpinnerButton(s, theme.Icon(theme.IconNameArrowDropDown), s.downButtonClicked)
	s.Disable()
	return s
}

// NewIntSpinnerWithData returns a new IntSpinner widget connected to the specified data source.
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

// CreateRenderer is a private method to fyne which links this IntSpinner widget to its
// renderer.
func (s *IntSpinner) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	th := s.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	box := canvas.NewRectangle(th.Color(theme.ColorNameBackground, v))
	border := canvas.NewRectangle(color.Transparent)

	text := canvas.NewText("", color.Black)
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

// Enable enables the IntSpinner and its buttons as appropriate.
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

// GetValue retrieves the current Spinner value.
func (s *IntSpinner) GetValue() int {
	return s.value
}

// MinSize returns the minimum size for the widget.
func (s *IntSpinner) MinSize() fyne.Size {
	th := s.Theme()
	padding := fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding) * 2)
	textSize := s.textSize()
	tHeight := textSize.Height + padding.Height
	upButtonHeight := s.upButton.MinSize().Height
	tWidth := textSize.Width + upButtonHeight + padding.Width
	return fyne.NewSize(tWidth, tHeight)
}

// MouseDown called on mouse click.
// This action causes the IntSpinner to request focus.
//
// Implements: desktop.Mouseable
func (s *IntSpinner) MouseDown(m *desktop.MouseEvent) {
	s.requestFocus()
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

// TypedKey receives key input events when the IntSpinner widget has focus.
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

// TypedRune receives text input events when the INtSpinner widget is focused.
// Increments/decrements the IntSpinner's value when the '+' or '-' key is
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

// Unbind disconnects any configured data source from this IntSpinner.
// The current value will remain at the last value of the data source.
func (s *IntSpinner) Unbind() {
	s.binder.Unbind()
}

// downButtonClicked handles taps on the IntSpinner's down button.
func (s *IntSpinner) downButtonClicked() {
	s.SetValue(s.value - s.step)
}

// requestFocus requests that this IntSpinner receive focus.
func (s *IntSpinner) requestFocus() {
	if c := fyne.CurrentApp().Driver().CanvasForObject(s); c != nil {
		c.Focus(s)
	}
}

// Calculate the max size of the text that can be displayed for the IntSpinner.
// The size cannot be larger than the larger of the sizes for the Spinner
// min and max values.
func (s *IntSpinner) textSize() fyne.Size {
	minText := strconv.Itoa(s.min)
	maxText := strconv.Itoa(s.max)
	return maxTextSize(minText, maxText)
}

// upButtonClicked handles tap events for the IntSpinner's up button.
func (s *IntSpinner) upButtonClicked() {
	s.SetValue(s.value + s.step)
}

// updateFromData updates the IntSpinner to the value set in the bound data.
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

// writeData updates the bound data item as the result of changes in the IntSpinner value.
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

// intSpinnerRenderer is the renderer for the IntSpinner widget
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

// Layout positions and sizes all of the objects that make up the IntSpinner widget.
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
	buttonSize := r.spinner.upButton.MinSize()
	xPos := size.Width - buttonSize.Width - borderSize - padding/2
	yPos := (rMinSize.Height - textSize.Height) / 2
	r.text.Move(fyne.NewPos(xPos, yPos))

	xPos += padding / 4
	yPos -= padding - 2
	r.spinner.upButton.Resize(buttonSize)
	r.spinner.upButton.Move(fyne.NewPos(xPos, yPos))

	yPos = r.spinner.upButton.MinSize().Height + padding/2 - 1
	r.spinner.downButton.Resize(buttonSize)
	r.spinner.downButton.Move(fyne.NewPos(xPos, yPos))
}

// MinSize returns the minimum size of the IntSpinner widget.
func (r *intSpinnerRenderer) MinSize() fyne.Size {
	return r.spinner.MinSize()
}

// Objects returns the objects associated with the IntSpinner renderer.
func (r *intSpinnerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh refreshes (redisplays) the IntSpinner widget.
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

var _ fyne.Disableable = (*Float64Spinner)(nil)
var _ fyne.Focusable = (*Float64Spinner)(nil)
var _ fyne.Tappable = (*Float64Spinner)(nil)
var _ desktop.Mouseable = (*Float64Spinner)(nil)
var _ fyne.Scrollable = (*Float64Spinner)(nil)

// Float64Spinner widget has a minimum, maximum, step and current values along with spinnerButtons
// to increment and decrement the spinner value.
type Float64Spinner struct {
	baseSpinner

	value     float64
	min       float64
	max       float64
	step      float64
	precision uint

	upButton   *float64SpinnerButton
	downButton *float64SpinnerButton

	OnChanged func(float64) `json:"-"`
}

// NewFloat64Spinner creates a new Float64Spinner widget.
//
// Params:
//
//		min is the minimum spinner value. It may be < 0.
//		max is the maximum spinner value. It must be > min.
//		step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//	 precision is the number of digits to display after the decimal point.
//		onChanged is the callback function that is called whenever the spinner value changes.
func NewFloat64Spinner(min, max, step float64, precision uint, onChanged func(float64)) *Float64Spinner {
	if min >= max {
		panic(errors.New("spinner max must be greater than min value"))
	}
	if step < 1 {
		panic(errors.New("spinner step must be greater than 0"))
	}
	if step > max-min {
		panic(errors.New("spinner step must be less than or equal to max - min"))
	}
	s := &Float64Spinner{
		min:       min,
		max:       max,
		step:      step,
		precision: precision,
		OnChanged: onChanged,
	}
	s.baseSpinner.initialized = true
	s.upButton = newFloat64SpinnerButton(s, theme.Icon(theme.IconNameArrowDropUp), s.upButtonClicked)
	s.downButton = newFloat64SpinnerButton(s, theme.Icon(theme.IconNameArrowDropDown), s.downButtonClicked)
	s.SetValue(s.min)
	return s
}

// NewFloat64SpinnerUninitialized returns a new uninitialized Float64Spinner widget.
//
// An uninitialized Float64Spinner widget is useful when you need to create a Float64Spinner
// but the initial settings are unknown.
// Calling Enable on an uninitialized spinner will not enable the spinner; you
// must first call SetMinMaxStep to initialize spinner values before enabling
// the spinner widget.
func NewFloat64SpinnerUninitialized() *Float64Spinner {
	s := &Float64Spinner{}
	s.baseSpinner.initialized = false
	s.upButton = newFloat64SpinnerButton(s, theme.Icon(theme.IconNameArrowDropUp), s.upButtonClicked)
	s.downButton = newFloat64SpinnerButton(s, theme.Icon(theme.IconNameArrowDropDown), s.downButtonClicked)
	s.Disable()
	return s
}

// NewFloat64SpinnerWithData returns a new Flaot64Spinner widget connected to the specified data source.
//
// Params:
//
//		min is the minimum spinner value. It may be < 0.
//		max is the maximum spinner value. It must be > min.
//		step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//	 precision is the number of digits to display after the decimal point.
//		data is the value that is bound to the spinner value.
func NewFloat64SpinnerWithData(min, max, step float64, precision uint, data binding.Float) *Float64Spinner {
	s := NewFloat64Spinner(min, max, step, precision, nil)
	s.Bind(data)
	s.OnChanged = func(_ float64) {
		s.binder.CallWithData(s.writeData)
	}

	return s
}

// Bind connects the specified data source to this Float64Spinner widget.
// The current value will be displayed and any changes in the data will cause the widget
// to update.
func (s *Float64Spinner) Bind(data binding.Float) {
	s.binder.SetCallback(s.updateFromData)
	s.binder.Bind(data)
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer.
func (s *Float64Spinner) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	th := s.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	box := canvas.NewRectangle(th.Color(theme.ColorNameBackground, v))
	border := canvas.NewRectangle(color.Transparent)

	format := fmt.Sprintf("%%.%df", s.precision)
	value := fmt.Sprintf(format, s.value)
	text := canvas.NewText(value, th.Color(theme.ColorNameForeground, v))

	objects := []fyne.CanvasObject{
		box,
		border,
		text,
		s.upButton,
		s.downButton,
	}
	r := &float64SpinnerRenderer{
		spinner: s,
		box:     box,
		border:  border,
		text:    text,
		objects: objects,
	}
	return r
}

// Disable disables the Float64Spinner and its buttons.
func (s *Float64Spinner) Disable() {
	s.downButton.Disable()
	s.upButton.Disable()
	s.DisableableWidget.Disable()
	s.Refresh()
}

// Enable enables the Float64Spinner and its buttons as appropriate.
func (s *Float64Spinner) Enable() {
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

// / GetValue retrieves the current Float64Spinner value.
func (s *Float64Spinner) GetValue() float64 {
	return s.value
}

func (s *Float64Spinner) MinSize() fyne.Size {
	th := s.Theme()
	padding := fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding) * 2)
	textSize := s.textSize()
	tHeight := textSize.Height + padding.Height
	upButtonHeight := s.upButton.MinSize().Height
	tWidth := textSize.Width + upButtonHeight + padding.Width
	return fyne.NewSize(tWidth, tHeight)
}

// MouseDown called on mouse click.
// This action causes the Float64Spinner to request focus.
//
// Implements: desktop.Mouseable
func (s *Float64Spinner) MouseDown(m *desktop.MouseEvent) {
	s.requestFocus()
	s.Refresh()
}

// MouseUp called on mouse release.
//
// Implements: desktop.Mouseable
func (s *Float64Spinner) MouseUp(m *desktop.MouseEvent) {}

// Scrolled handles mouse scroller events.
//
// Implements fyne.Scrollable
func (s *Float64Spinner) Scrolled(evt *fyne.ScrollEvent) {
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
func (s *Float64Spinner) SetMinMaxStep(min, max, step float64) {
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
func (s *Float64Spinner) SetValue(val float64) {
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
// the spinner object.
// If actually over one of the spinnerButtons, processing
// is passed to that button for handling.
func (s *Float64Spinner) Tapped(evt *fyne.PointEvent) {
	if s.Disabled() {
		return
	}
	if s.upButton.containsPoint(evt.Position) {
		s.upButton.Tapped(evt)
	} else if s.downButton.containsPoint(evt.Position) {
		s.downButton.Tapped(evt)
	}
}

// TypedKey receives key input events when the spinner widget has focus.
// Increments/decrements the spinner's value when the up or down key is
// pressed.
//
// Implements: fyne.Focusable
func (s *Float64Spinner) TypedKey(key *fyne.KeyEvent) {
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

// TypedRune receives text input events when the spinner widget is focused.
// Increments/decrements the spinner's value when the '+' or '-' key is
// pressed.
//
// Implements: fyne.Focusable
func (s *Float64Spinner) TypedRune(rune rune) {
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

// Unbind disconnects any configured data source from this spinner.
// The current value will remain at the last value of the data source.
func (s *Float64Spinner) Unbind() {
	s.binder.Unbind()
}

// requestFocus requests that this IntSpinner receive focus.
func (s *Float64Spinner) requestFocus() {
	if c := fyne.CurrentApp().Driver().CanvasForObject(s); c != nil {
		c.Focus(s)
	}

}

// Calculate the max size of the text that can be displayed for the spinner.
// The size cannot be larger than the larger of the sizes for the spinner
// min and max values.
func (s *Float64Spinner) textSize() fyne.Size {
	format := fmt.Sprintf("%%.%df", s.precision)
	minVal := fmt.Sprintf(format, s.min)
	maxVal := fmt.Sprintf(format, s.max)
	return maxTextSize(minVal, maxVal)
}

// updateFromData updates the spinner to the value set in the bound data.
func (s *Float64Spinner) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	textSource, ok := data.(binding.Float)
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

// writeData updates the bound data item as the result of changes in the spinner value.
func (s *Float64Spinner) writeData(data binding.DataItem) {
	if data == nil {
		return
	}
	floatTarget, ok := data.(binding.Float)
	if !ok {
		return
	}
	currentValue, err := floatTarget.Get()
	if err != nil {
		return
	}
	if currentValue != s.GetValue() {
		err := floatTarget.Set(s.GetValue())
		if err != nil {
			fyne.LogError(fmt.Sprintf("Failed to set binding value to %f", s.GetValue()), err)
		}
	}
}

// float64SpinnerRenderer is the renderer for the Float64Spinner widget
type float64SpinnerRenderer struct {
	spinner *Float64Spinner
	box     *canvas.Rectangle
	border  *canvas.Rectangle
	text    *canvas.Text
	objects []fyne.CanvasObject
}

// Destroy destroys any objects that must be destroyed when the renderer is
// destroyed.
func (r *float64SpinnerRenderer) Destroy() {}

// Layout positions and sizes all of the objects that make up the Float64Spinner widget.
func (r *float64SpinnerRenderer) Layout(size fyne.Size) {
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
	buttonSize := r.spinner.upButton.MinSize()
	xPos := size.Width - buttonSize.Width - borderSize - padding/2
	yPos := (rMinSize.Height - textSize.Height) / 2
	r.text.Move(fyne.NewPos(xPos, yPos))

	xPos += padding / 4
	yPos -= padding - 2
	r.spinner.upButton.Resize(buttonSize)
	r.spinner.upButton.Move(fyne.NewPos(xPos, yPos))

	yPos = r.spinner.upButton.MinSize().Height + padding/2 - 1
	r.spinner.downButton.Resize(buttonSize)
	r.spinner.downButton.Move(fyne.NewPos(xPos, yPos))
}

// MinSize returns the minimum size of the Flaot64Spinner widget.
func (r *float64SpinnerRenderer) MinSize() fyne.Size {
	return r.spinner.MinSize()
}

// Objects returns the objects associated with the Float64Spinner renderer.
func (r *float64SpinnerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh refreshes (redisplays) the Float64Spinner widget.
func (r *float64SpinnerRenderer) Refresh() {
	th := r.spinner.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	fgColor, bgColor, borderColor := r.spinnerColors()
	r.box.FillColor = th.Color(bgColor, v)
	r.box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	r.border.CornerRadius = r.box.CornerRadius

	r.border.StrokeColor = th.Color(borderColor, v)
	format := fmt.Sprintf("%%.%df", r.spinner.precision)
	r.text.Text = fmt.Sprintf(format, r.spinner.value)
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

// downButtonClicked handles tap events for the Float64Spinner's down button.
func (s *Float64Spinner) downButtonClicked() {
	s.SetValue(s.value - s.step)
}

// spinnerColors returns the colors to display the button in.
func (r *float64SpinnerRenderer) spinnerColors() (fgColor, bgColor, borderColor fyne.ThemeColorName) {
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

// upButtonClicked handles tap events for the Float64Spinner's up button.
func (s *Float64Spinner) upButtonClicked() {
	s.SetValue(s.value + s.step)
}

// maxTextSize calculates the larger of the canvas.Text sizes for the two string params
func maxTextSize(minText, maxText string) fyne.Size {
	// color does not affect the text size, so use Black.
	minT := canvas.NewText(minText, color.Black)
	maxT := canvas.NewText(maxText, color.Black)
	minTextSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(minT.Text,
		minT.TextSize, minT.TextStyle, minT.FontSource)
	maxTextSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(maxT.Text,
		maxT.TextSize, maxT.TextStyle, maxT.FontSource)
	return fyne.NewSize(max(minTextSize.Width, maxTextSize.Width),
		max(minTextSize.Height, maxTextSize.Height))
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
