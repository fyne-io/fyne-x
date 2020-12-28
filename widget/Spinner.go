package widget // import "fyne.io/x/fyne/widget"

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/driver/mobile"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// SpinnerDefaultPrecision controls the default number of digits shown after the decimal point in
// the Spinner's entry. The special precision -1 uses the smallest number of digits necessary
// such that the value is represented to the best of a float64's ability.
//
// See Also
//
// strconv.ParseFloat
var SpinnerDefaultPrecision = 12

// Declare conformity with fyne.Container interfaces
var _ fyne.Disableable = (*Spinner)(nil)
var _ fyne.Widget = (*Spinner)(nil)

// Declare conformity with widget.Entry interfaces
var _ fyne.Disableable = (*spinnerEntry)(nil)
var _ fyne.Draggable = (*spinnerEntry)(nil)
var _ fyne.Focusable = (*spinnerEntry)(nil)
var _ fyne.Scrollable = (*spinnerEntry)(nil)
var _ fyne.Tappable = (*spinnerEntry)(nil)
var _ fyne.Widget = (*spinnerEntry)(nil)
var _ desktop.Mouseable = (*spinnerEntry)(nil)
var _ desktop.Keyable = (*spinnerEntry)(nil)
var _ mobile.Keyboardable = (*spinnerEntry)(nil)

// Spinner widget has a floating point value shown in a widget.Entry, with widget.Button s to the
// right which increment or decrement the spinner value by a set step value. Also provides an
// integer version, a set minimum and/or maximum, and visual precision for floationg point values.
//
// Note: Since this component uses float64 it should not be considered 100% mathematically
// accurate. For more information, please see: https://0.30000000000000004.com/. If you need this,
// then you would have to re-implement this component using a decimal math library such as:
// github.com/ericlagergren/decimal. A conscious choice not to use this library was made due to
// the significant performance hit involved in using software over hardware floating point math.
type Spinner struct {
	fyne.Container

	buttonDown *widget.Button
	buttonUp   *widget.Button
	entry      *spinnerEntry

	hasMax    bool
	hasMin    bool
	lastValue float64
	max       float64
	min       float64
	precision int
	step      float64
	value     float64

	integer bool
}

// NewSpinner creates a spinner with a floating point initial value and step increment.
func NewSpinner(value, step float64) *Spinner {
	return newSpinner(value, step, false)
}

// NewIntSpinner creates a spinner with an integer initial value and step increment.
func NewIntSpinner(value, step int) *Spinner {
	return newSpinner(float64(value), float64(step), true)
}

func newSpinner(value, step float64, integer bool) *Spinner {
	s := &Spinner{
		buttonDown: widget.NewButtonWithIcon("", theme.MoveDownIcon(), nil),
		buttonUp:   widget.NewButtonWithIcon("", theme.MoveUpIcon(), nil),
		entry:      newSpinnerEntry(),
		precision:  SpinnerDefaultPrecision,
		step:       step,
		value:      value,
		integer:    integer,
	}

	s.buttonDown.OnTapped = s.onDown
	s.buttonUp.OnTapped = s.onUp
	s.entry.spinner = s

	// ! Buttons should be positioned one atop the other vertically... However, this would require
	// ! manual Layout of widgets v.s. using fyne.Container
	buttons := widget.NewHBox(s.buttonUp, s.buttonDown)

	// ! Changing the above would replace this, and the AddObject calls.
	s.Layout = layout.NewBorderLayout(nil, nil, nil, buttons)

	s.Add(s.entry)
	s.Add(buttons)

	// Force the value in the entry to refresh (with the correct decimal count)
	s.updateVal()

	return s
}

// Max returns the maximal value of the spinner.
func (s *Spinner) Max() float64 {
	return s.max
}

// SetMax sets the maximal value of the spinner and updates the current value.
func (s *Spinner) SetMax(max float64) {
	s.hasMax = true
	s.max = max
	s.updateVal()
}

// Min returns the minimal value of the spinner.
func (s *Spinner) Min() float64 {
	return s.min
}

// SetMin sets the minimal value of the spinner and updates the current value.
func (s *Spinner) SetMin(min float64) {
	s.hasMin = true
	s.min = min
	s.updateVal()
}

// Precision returns the number of significant digits that will be displayed in the spinner entry.
//
// Note: This does not change the accuracy of any calculations, only the visual representation.
func (s *Spinner) Precision() int {
	return s.precision
}

// SetPrecision sets the number of significant digits that will be displayed in the spinner entry.
//
// Note: This does not change the accuracy of any calculations, only the visual representation.
func (s *Spinner) SetPrecision(precision int) {
	s.precision = precision
}

// Step returns the increment by which the spinner's value is modified. This should be a positive
// (> 0) value for correct widget behaviour.
func (s *Spinner) Step() float64 {
	return s.step
}

// SetStep sets the increment by which the spinner's value is modified. This should be a positive
// (> 0) value for correct widget behaviour.
func (s *Spinner) SetStep(step float64) {
	s.step = step
}

// Value returns the spinner's current value.
func (s *Spinner) Value() float64 {
	return s.value
}

// SetValue sets the spinner's current value to a given floating point value, respecting set min/max.
func (s *Spinner) SetValue(value float64) {
	s.value = value
	s.updateVal()
}

// ValueInt returns the spinner's integer value (rounding if necessary).
func (s *Spinner) ValueInt() int {
	return int(math.Round(s.value))
}

// SetValueInt sets the spinner's current value to a given integer, respecting set min/max.
func (s *Spinner) SetValueInt(value int) {
	s.value = float64(value)
	s.updateVal()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
//
// Implements: fyne.Widget
func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	return &spinnerRenderer{
		spinner: s,
	}
}

// Disable this widget so that it cannot be interacted with, updating any style appropriately.
//
// Implements: fyne.Disableable
func (s *Spinner) Disable() {
	s.buttonDown.Disable()
	s.buttonUp.Disable()
	s.entry.Disable()
}

// Disabled returns whether the entry is disabled.
//
// Implements: fyne.Disableable
func (s *Spinner) Disabled() bool {
	return s.entry.Disabled()
}

// Enable this widget, updating any style or features appropriately.
//
// Implements: fyne.Disableable
func (s *Spinner) Enable() {
	s.entry.Enable()

	// updateVal enables the down/up buttons based on min/max
	s.updateVal()
}

func (s *Spinner) onDown() {
	s.value -= s.step
	s.updateVal()
}

func (s *Spinner) onUp() {
	s.value += s.step
	s.updateVal()
}

func (s *Spinner) onEnter() {
	if f, err := strconv.ParseFloat(s.entry.Text, 64); err == nil {
		s.value = f
	} else {
		s.value = s.lastValue
	}

	s.updateVal()
}

func (s *Spinner) onScrolled(e *fyne.ScrollEvent) {
	// ! Fyne does not yet interpret/provide a setting for the scroll direction
	if f := float64(e.DeltaY); f > 0 {
		s.onUp()
	} else if f < 0 {
		s.onDown()
	}
}

func (s *Spinner) updateVal() {
	if s.hasMin && s.hasMax {
		s.value = math.Min(math.Max(s.min, s.value), s.max)
	} else if s.hasMax {
		s.value = math.Min(s.value, s.max)
	} else if s.hasMin {
		s.value = math.Max(s.value, s.min)
	}

	s.lastValue = s.value

	if s.integer {
		s.value = math.Round(s.value)
		s.entry.SetText(fmt.Sprintf("%d", int(s.value)))
	} else {
		// ! May produce funky visuals such as `5.000000000001` because floats.
		s.entry.SetText(strconv.FormatFloat(s.value, 'f', s.precision, 64))
	}

	// Disable the down button if we're at the minimum value
	if s.hasMin && s.value <= s.min {
		s.buttonDown.Disable()
	} else {
		s.buttonDown.Enable()
	}

	// Disable the up button if we're at the maximum value
	if s.hasMax && s.value >= s.max {
		s.buttonUp.Disable()
	} else {
		s.buttonUp.Enable()
	}
}

// ---

type spinnerEntry struct {
	widget.Entry
	spinner *Spinner
}

func newSpinnerEntry() *spinnerEntry {
	e := &spinnerEntry{}

	e.ExtendBaseWidget(e)

	return e
}

// Keyboard implements the Keyboardable interface.
//
// Implements: mobile.Keyboardable
func (e *spinnerEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

// KeyDown handler for keypress events - used to store shift modifier state for text selection.
//
// Implements: desktop.Keyable
func (e *spinnerEntry) KeyDown(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn:
		e.spinner.onEnter()
	case fyne.KeyDown:
		e.spinner.onDown()
	case fyne.KeyUp:
		e.spinner.onUp()
	default:
		e.Entry.KeyDown(key)
	}
}

// Scrolled is called when an input device triggers a scroll event, such as a desktop mouse wheel.
//
// Implements: fyne.Scrollable
func (e *spinnerEntry) Scrolled(s *fyne.ScrollEvent) {
	e.spinner.onScrolled(s)
}

// ---

type spinnerRenderer struct {
	spinner *Spinner
}

func (r *spinnerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *spinnerRenderer) Destroy() {}

func (r *spinnerRenderer) Layout(size fyne.Size) {
	r.spinner.Container.Layout.Layout(r.spinner.Container.Objects, size)
}

func (r *spinnerRenderer) MinSize() fyne.Size {
	return r.spinner.Container.MinSize()
}

func (r *spinnerRenderer) Objects() []fyne.CanvasObject {
	return r.spinner.Container.Objects
}

func (r *spinnerRenderer) Refresh() {
	r.spinner.Container.Refresh()
}
