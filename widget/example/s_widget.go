package example

import (
	"sync"

	"fyne.io/fyne/v2"
	wgt "fyne.io/x/fyne/widget"
)

// SampleSWidget implements a table with row based selection and column names.
type SampleSWidget struct {
	wgt.BaseSimpleWidget
	propertyLock sync.RWMutex

	Property int
}

// NewSampleSWidget creates a new SampleSWidget.
func NewSampleSWidget(property int) *SampleSWidget {
	w := &SampleSWidget{Property: property}
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer implements fyne.Widget.
func (s *SampleSWidget) CreateRenderer() fyne.WidgetRenderer {
	rend := &sampleSWidgetRenderer{
		widget: s,
	}
	rend.BaseSimpleRenderer = *wgt.NewBaseSimpleRenderer(s, rend)
	return rend
}

// SetStateSafe sets or changes the state of a widget in a safe way. A Refresh
// is triggered after the state changes have been applied.
func (s *SampleSWidget) SetStateSafe(setState func()) {
	s.BaseSimpleWidget.SetStateSafe(&s.propertyLock, setState)
}

func (s *SampleSWidget) getProperty() int {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	return s.Property
}

type sampleSWidgetRenderer struct {
	wgt.BaseSimpleRenderer

	widget *SampleSWidget
}

// Build renders the SampleSWidget.
func (s *sampleSWidgetRenderer) Build() (objects []fyne.CanvasObject, layout func(size fyne.Size)) {
	// create objects needed for rendering and append them to the objects slice.
	// (executed in Widget.CreateRenderer function)

	return objects, func(size fyne.Size) {
		property := s.widget.getProperty()
		_ = property

		// position, resize or otherwise update objects based on changed properties.
		// (executed in WidgetRenderer.Layout)
	}
}
