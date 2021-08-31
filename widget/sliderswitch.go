package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Tappable = (*SliderSwitch)(nil)

// SliderSwitch is an extension of a slider that switches
// between two values when clicked.
type SliderSwitch struct {
	widget.Slider

	// OnToggle is called AFTER the switch has changed state.
	// Bool is true if the switch is in the "ON" state.
	OnToggle func(bool)
	// If set, the slider will be moved to an intermediate state,
	// the function called, and then moved to its final state after
	// the function returns. Bool is true if the switch is transitioning
	// to the "ON" state.
	OnTransition func(bool)
}

// NewSliderSwitch returns a new SliderSwitch.
func NewSliderSwitch() *SliderSwitch {
	sw := &SliderSwitch{
		Slider: widget.Slider{
			Value:       0,
			Min:         0,
			Max:         2,
			Step:        1, // Make intermediate step - maybe spinner
			Orientation: widget.Horizontal,
		},
	}
	sw.ExtendBaseWidget(sw)
	return sw
}

// Tapped responds to tapped events on the switch.
func (s *SliderSwitch) Tapped(*fyne.PointEvent) {
	originalValue := s.Value
	if s.OnTransition != nil {
		value := s.Value
		s.SetValue(1)
		s.OnTransition(value == 0)

	}
	switch originalValue {
	case 0:
		s.SetValue(2)
	case 2:
		s.SetValue(0)
	}
	if s.OnToggle != nil {
		s.OnToggle(s.Value == 2)
	}
}
