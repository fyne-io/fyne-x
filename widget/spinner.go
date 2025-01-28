package widget

import (
	"errors"
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// spinnerButton is the widget used for buttons in Spinners.
type spinnerButton struct {
	widget.Button

	position fyne.Position
	size     fyne.Size
}

// newSpinnerButton creates a spinnerButton for use in Spinner widgets.
//
// Params:
//
//	resource is the resource to be used as the button icon.
//	onTapped is the callback function for button clicks.
func newSpinnerButton(resource fyne.Resource, onTapped func()) *spinnerButton {
	b := &spinnerButton{}
	b.ExtendBaseWidget(b)
	b.setButtonProperties(resource, onTapped)
	return b
}

// MinSize returns the minimum size of the button. Because the minimum size is a constant
// based on the spinner height and theme properties, the minimum size is calculated when
// the button is created.
func (b *spinnerButton) MinSize() fyne.Size {
	return fyne.NewSize(b.size.Height, b.size.Height)
}

// Move moves the button.
func (b *spinnerButton) Move(pos fyne.Position) {
	b.position = pos
	b.BaseWidget.Move(pos)
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
func (b *spinnerButton) containsPoint(pos fyne.Position) bool {
	if pos.X < b.position.X || pos.X > b.position.X+b.size.Width {
		return false
	} else if pos.Y < b.position.Y || pos.Y > b.position.Y+b.size.Height {
		return false
	}
	return true
}

// enableDisable enables or disables the button based on whether the button's
// parent spinner widget is disabled, and whether the spinner's value is at its limit.
//
// Params:
//
//		parentDisabled indicates whether the button's parent spinner is disabled or not.
//		limit indicates whether the spinner's value is at the corresponding limit for this
//	 button. For example, for an up button, the  limit should be true if value == max, and
//	 for a down button, the limit should be true if value == min.
func (b *spinnerButton) enableDisable(parentDisabled, limit bool) {
	if parentDisabled {
		b.Disable()
	} else {
		b.Enable()
		if limit {
			b.Disable()
		}
	}
}

// setButtonProperties sets the button properties.
//
// Params:
//
//	resource is the Resource for the button icon.
//	onTapped is the function to be called when the button is tapped.
func (b *spinnerButton) setButtonProperties(resource fyne.Resource, onTapped func()) {
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

var _ fyne.Disableable = (*Spinner)(nil)
var _ fyne.Focusable = (*Spinner)(nil)
var _ fyne.Tappable = (*Spinner)(nil)
var _ desktop.Mouseable = (*Spinner)(nil)
var _ fyne.Scrollable = (*Spinner)(nil)

// Spinner widget has a minimum, maximum, step and current values along with spinnerButtons
// to increment and decrement the spinner value.
type Spinner struct {
	widget.DisableableWidget

	value  float64
	min    float64
	max    float64
	step   float64
	format string

	upButton   *spinnerButton
	downButton *spinnerButton

	binder      basicBinder
	initialized bool
	hovered     bool
	focused     bool

	OnChanged func(float64) `json:"-"`
}

// NewSpinner creates a new Spinner widget.
//
// Params:
//
//		min is the minimum spinner value. It may be < 0.
//		max is the maximum spinner value. It must be > min.
//		step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//	 	format is the format to display the value in. This format should contain one of the following: "%d", "%+d", or "%.Xf"
//		where X is an unsigned integer.
//		onChanged is the callback function that is called whenever the spinner value changes.
func NewSpinner(min, max, step float64, format string, onChanged func(float64)) *Spinner {
	if min >= max {
		panic(errors.New("spinner max must be greater than min value"))
	}
	if step < 1 {
		panic(errors.New("spinner step must be greater than 0"))
	}
	if step > max-min {
		panic(errors.New("spinner step must be less than or equal to max - min"))
	}
	s := &Spinner{
		min:       min,
		max:       max,
		step:      step,
		format:    format,
		OnChanged: onChanged,
	}
	s.initialized = true
	s.upButton = newSpinnerButton(theme.Icon(theme.IconNameArrowDropUp), s.upButtonClicked)
	s.downButton = newSpinnerButton(theme.Icon(theme.IconNameArrowDropDown), s.downButtonClicked)
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
//
// Params:
//
//	format is the format to display the value in. This format should contain one of the following: "%d", "%+d", or "%.Xf"
//	where X is an unsigned integer.
func NewSpinnerUninitialized(format string) *Spinner {
	s := &Spinner{format: format}
	s.initialized = false
	s.upButton = newSpinnerButton(theme.Icon(theme.IconNameArrowDropUp), s.upButtonClicked)
	s.downButton = newSpinnerButton(theme.Icon(theme.IconNameArrowDropDown), s.downButtonClicked)
	s.Disable()
	return s
}

// NewSpinnerWithData returns a new Spinner widget connected to the specified data source.
//
// Params:
//
//		min is the minimum spinner value. It may be < 0.
//		max is the maximum spinner value. It must be > min.
//		step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//	 	format is the format to display the value in. This format should contain one of the following: "%d", "%+d", or "%.Xf"
//		where X is an unsigned integer.
//		data is the value that is bound to the spinner value.
func NewSpinnerWithData(min, max, step float64, format string, data binding.Float) *Spinner {
	s := NewSpinner(min, max, step, format, nil)
	s.Bind(data)
	s.OnChanged = func(_ float64) {
		s.binder.CallWithData(s.writeData)
	}

	return s
}

// Bind connects the specified data source to this Spinner widget.
// The current value will be displayed and any changes in the data will cause the widget
// to update.
func (s *Spinner) Bind(data binding.Float) {
	s.binder.SetCallback(s.updateFromData)
	s.binder.Bind(data)
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer.
func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	th := s.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	box := canvas.NewRectangle(th.Color(theme.ColorNameBackground, v))
	border := canvas.NewRectangle(color.Transparent)

	value := fmt.Sprintf(s.format, s.value)
	text := canvas.NewText(value, th.Color(theme.ColorNameForeground, v))
	text.Alignment = fyne.TextAlignTrailing

	objects := []fyne.CanvasObject{
		box,
		border,
		text,
		s.upButton,
		s.downButton,
	}
	r := &SpinnerRenderer{
		spinner: s,
		box:     box,
		border:  border,
		text:    text,
		objects: objects,
	}
	return r
}

// Disable disables the Spinner and its buttons.
func (s *Spinner) Disable() {
	s.downButton.Disable()
	s.upButton.Disable()
	s.DisableableWidget.Disable()
	s.Refresh()
}

// Enable enables the Spinner and its buttons as appropriate.
func (s *Spinner) Enable() {
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

// FocusGained is called when the spinner has been given focus.
//
// Implements: fyne.Focusable
func (s *Spinner) FocusGained() {
	s.focused = true
	s.Refresh()
}

// FocusLost is called when the spinner has had focus removed.
//
// Implements: fyne.Focusable
func (s *Spinner) FocusLost() {
	s.focused = false
	s.Refresh()
}

func (s *Spinner) GetFormat() string {
	return s.format
}

// / GetValue retrieves the current Spinner value.
func (s *Spinner) GetValue() float64 {
	return s.value
}

// MinSize returns the minimum size of the Spinner widget. The minimum size is calculated
// based on the maximum width that the value could require based on its format.
func (s *Spinner) MinSize() fyne.Size {
	th := s.Theme()
	padding := fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding) * 2)
	textSize := s.textSize()
	tHeight := textSize.Height + padding.Height
	upButtonHeight := s.upButton.MinSize().Height
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

// MouseUp called on mouse release.
//
// Implements: desktop.Mouseable
func (s *Spinner) MouseUp(m *desktop.MouseEvent) {}

// Scrolled handles mouse scroller events.
//
// Implements fyne.Scrollable
func (s *Spinner) Scrolled(evt *fyne.ScrollEvent) {
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
func (s *Spinner) SetMinMaxStep(min, max, step float64) {
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
func (s *Spinner) SetValue(val float64) {
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
func (s *Spinner) Tapped(evt *fyne.PointEvent) {
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
func (s *Spinner) TypedKey(key *fyne.KeyEvent) {
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
func (s *Spinner) TypedRune(rune rune) {
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
func (s *Spinner) Unbind() {
	s.binder.Unbind()
}

// requestFocus requests that this Spinner receive focus.
func (s *Spinner) requestFocus() {
	if c := fyne.CurrentApp().Driver().CanvasForObject(s); c != nil {
		c.Focus(s)
	}

}

// Calculate the max size of the text that can be displayed for the spinner.
// The size cannot be larger than the larger of the sizes required to display the
// spinner's min and max values.
func (s *Spinner) textSize() fyne.Size {
	var minVal, maxVal string
	if strings.Contains(s.format, "%d") ||
		strings.Contains(s.format, "%+d") {
		minVal = fmt.Sprintf(s.format, int(s.min))
		maxVal = fmt.Sprintf(s.format, int(s.max))
	} else {
		minVal = fmt.Sprintf(s.format, s.min)
		maxVal = fmt.Sprintf(s.format, s.max)
	}
	return maxTextSize(minVal, maxVal)
}

// updateFromData updates the spinner to the value set in the bound data.
func (s *Spinner) updateFromData(data binding.DataItem) {
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
func (s *Spinner) writeData(data binding.DataItem) {
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

// SpinnerRenderer is the renderer for the Spinner widget
type SpinnerRenderer struct {
	spinner *Spinner
	box     *canvas.Rectangle
	border  *canvas.Rectangle
	text    *canvas.Text
	objects []fyne.CanvasObject
}

// Destroy destroys any objects that must be destroyed when the renderer is
// destroyed.
func (r *SpinnerRenderer) Destroy() {}

// Layout positions and sizes all of the objects that make up the Float64Spinner widget.
func (r *SpinnerRenderer) Layout(size fyne.Size) {
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
func (r *SpinnerRenderer) MinSize() fyne.Size {
	return r.spinner.MinSize()
}

// Objects returns the objects associated with the Float64Spinner renderer.
func (r *SpinnerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh refreshes (redisplays) the Float64Spinner widget.
func (r *SpinnerRenderer) Refresh() {
	th := r.spinner.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	fgColor, bgColor, borderColor := spinnerColors(
		r.spinner.Disabled(), r.spinner.focused, r.spinner.hovered)
	r.box.FillColor = th.Color(bgColor, v)
	r.box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	r.border.CornerRadius = r.box.CornerRadius
	r.border.StrokeColor = th.Color(borderColor, v)

	if strings.Contains(r.spinner.format, "%d") ||
		strings.Contains(r.spinner.format, "%+d") {
		r.text.Text = fmt.Sprintf(r.spinner.format, int(r.spinner.value))
	} else {
		r.text.Text = fmt.Sprintf(r.spinner.format, r.spinner.value)
	}
	r.text.Color = th.Color(fgColor, v)
	r.text.Refresh()

	r.spinner.upButton.enableDisable(r.spinner.Disabled(), r.spinner.GetValue() == r.spinner.max)
	r.spinner.downButton.enableDisable(r.spinner.Disabled(), r.spinner.GetValue() == r.spinner.min)
}

// downButtonClicked handles tap events for the Spinner's down button.
func (s *Spinner) downButtonClicked() {
	s.SetValue(s.value - s.step)
}

// / upButtonClicked handles tap events for the Spinner's up button.
func (s *Spinner) upButtonClicked() {
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

// spinnerColors determines display colors for spinners.
func spinnerColors(disabled, focused, hovered bool) (fgColor, bgColor, borderColor fyne.ThemeColorName) {
	fgColor = theme.ColorNameForeground
	bgColor = ""
	borderColor = theme.ColorNameInputBorder
	if disabled {
		fgColor = theme.ColorNameDisabled
		borderColor = theme.ColorNameDisabled
	} else if focused {
		borderColor = theme.ColorNamePrimary
	} else if hovered {
		bgColor = theme.ColorNameHover
	}
	return fgColor, bgColor, borderColor
}
