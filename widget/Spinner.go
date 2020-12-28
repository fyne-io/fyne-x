package widget // import "fyne.io/x/fyne/widget"

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var SpinnerDefaultPrecision = 12

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

func NewSpinner(value, step float64) *Spinner {
	return newSpinner(value, step, false)
}

func NewIntSpinner(value, step int) *Spinner {
	return newSpinner(float64(value), float64(step), true)
}

func newSpinner(value, step float64, integer bool) *Spinner {
	s := &Spinner{
		buttonDown: widget.NewButtonWithIcon("", theme.MoveDownIcon(), nil),
		buttonUp:   widget.NewButtonWithIcon("", theme.MoveUpIcon(), nil),
		entry:      &spinnerEntry{},
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

	s.AddObject(s.entry)
	s.AddObject(buttons)

	// Force the value in the entry to refresh (with the correct decimal count)
	s.updateVal()

	return s
}

func (s *Spinner) Max() float64 {
	return s.max
}

func (s *Spinner) SetMax(max float64) {
	s.hasMax = true
	s.max = max
	s.updateVal()
}

func (s *Spinner) Min() float64 {
	return s.min
}

func (s *Spinner) SetMin(min float64) {
	s.hasMin = true
	s.min = min
	s.updateVal()
}

func (s *Spinner) Precision() int {
	return s.precision
}

func (s *Spinner) SetPrecision(precision int) {
	s.precision = precision
}

func (s *Spinner) Step() float64 {
	return s.step
}

func (s *Spinner) SetStep(step float64) {
	s.step = step
}

func (s *Spinner) Value() float64 {
	return s.value
}

func (s *Spinner) SetValue(value float64) {
	s.value = value
	s.updateVal()
}

func (s *Spinner) ValueInt() int {
	return int(math.Round(s.value))
}

func (s *Spinner) SetValueInt(value int) {
	s.value = float64(value)
	s.updateVal()
}

func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	return &spinnerRenderer{
		spinner: s,
	}
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
