package example

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/x/fyne/widget"
)

// SampleWidget is a sample widget demonstrating the base structure of a
// widget implementing SimpleWidget using SimpleWidgetBase.
type SampleWidget struct {
	widget.SimpleWidgetBase
	propertyLock sync.RWMutex

	Property int
}

// NewSampleWidget creates a new sample widget.
func NewSampleWidget(property int) *SampleWidget {
	wgt := &SampleWidget{Property: property}
	wgt.ExtendBaseWidget(wgt)
	return wgt
}

// Render renders the SampleWidget.
func (s *SampleWidget) Render() (objects []fyne.CanvasObject, layout func(size fyne.Size)) {
	// create objects needed for rendering and append them to the objects slice.
	// (executed in Widget.CreateRenderer function)

	return objects, func(size fyne.Size) {
		property := s.getProperty()
		_ = property

		// position, resize or otherwise update objects based on changed properties.
		// (executed in WidgetRenderer.Layout)
	}
}

// SetStateSafe sets or changes the state of a widget in a safe way. A Refresh
// is triggered after the state changes have been applied.
func (s *SampleWidget) SetStateSafe(setState func()) {
	s.SimpleWidgetBase.SetStateSafe(&s.propertyLock, setState)
}

func (s *SampleWidget) getProperty() int {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	return s.Property
}

// All other methods of a fyne.Widget can be overwritten as usual. E.g. MinSize, Resize, Refresh, etc.
