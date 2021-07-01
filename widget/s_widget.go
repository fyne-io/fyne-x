package widget

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// BaseSimpleWidget defines the base for a fyne.Widget implementation.
// To create a new widget base it on BaseSimpleWidget using composition.
// Create a `New` function initialising the widget and make sure to call
// ExtendBaseWidget in the New function. Always use the `New` function to
// create the widget or make sure `ExtendBaseWidget` is called elsewhere.
//
// See ./example/s_wigdet.go for a bootstraped widget implementation.
type BaseSimpleWidget struct {
	widget.BaseWidget

	impl         fyne.Widget
	propertyLock sync.RWMutex
}

// ExtendBaseWidget is used by an extending widget to make use of BaseSimpleWidget functionality.
func (s *BaseSimpleWidget) ExtendBaseWidget(wid fyne.Widget) {
	impl := s.getImpl()
	if impl != nil {
		return
	}

	s.BaseWidget.ExtendBaseWidget(wid)

	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.impl = wid
}

// SetState sets or changes the state of a widget. A Refresh
// is triggered after the state changes have been applied.
func (s *BaseSimpleWidget) SetState(setState func()) {
	setState()
	s.super().Refresh()
}

// SetStateSafe sets or changes the state of a widget in a safe way. A Refresh
// is triggered after the state changes have been applied.
// The provided sync.Locker should be the same you use for read protection of the
// widget properties.
func (s *BaseSimpleWidget) SetStateSafe(m sync.Locker, setState func()) {
	m.Lock()
	setState()
	m.Unlock()

	s.super().Refresh()
}

func (s *BaseSimpleWidget) getImpl() fyne.Widget {
	s.propertyLock.RLock()
	impl := s.impl
	s.propertyLock.RUnlock()

	if impl == nil {
		return nil
	}
	return impl
}

func (s *BaseSimpleWidget) super() fyne.Widget {
	impl := s.getImpl()
	if impl == nil {
		var x interface{} = s
		return x.(fyne.Widget)
	}
	return impl
}
