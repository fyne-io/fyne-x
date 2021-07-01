package widget

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// SimpleWidget defines an interface for fyne widgets that are simple
// to implement when based on SimpleWidgetBase.
// For more info on how to implement, see the documentation of the
// SimpleWidgetBase.
type SimpleWidget interface {
	fyne.Widget

	Build() ([]fyne.CanvasObject, func(size fyne.Size))
}

// SimpleWidgetBase defines the base for a SimpleWidget implementation.
// To create a new widget base it on SimpleWidgetBase using composition.
// Create a `New` function initialising the widget and make sure to call
// ExtendBaseWidget in the New function. Always use the `New` function to
// create the widget or make sure `ExtendBaseWidget` is called elsewhere.
//
// Overwrite the `Build() (objects []fyne.CanvasObject, layout func(size fyne.Size))`
// function. It returns the objects needed to render the widgets content,
// as well as a function `layout` responsible for positioning and resizing the
// different objects.
// Try not to define new objects in the `layout` function as they would be
// recreated every time the widget is refreshed.
//
// Other methods defined by the fyne.Widget interface can be overwritten
// and will be used by the SimpleWidgetBase if overwritten.
// Overwrite `Destroy()` to implement custom cleanup for the widget.
// Implement `Refresh()` to optimize refresh performance. By default `Refresh`
// calls the `Layout` method of the renderer, which is not necessary in most cases.
//
// See ./example/simple_wigdet.go for a bootstraped widget implementation.
type SimpleWidgetBase struct {
	widget.BaseWidget

	propertyLock sync.RWMutex
	impl         SimpleWidget
	layout       func(fyne.Size)
}

// Build must be overwritten in a widget. It returns a slice of child objects
// and a layout function.
// The layout function should be used to position and size the objects (widgets
// and canvas objects). New objects should be created in the Build function body,
// so they are not re-created every time the widget gets refreshed.
func (s *SimpleWidgetBase) Build() (objects []fyne.CanvasObject, layout func(size fyne.Size)) {
	return nil, func(fyne.Size) {}
}

// CreateRenderer implements the Widget interface. It creates a simpleRenderer
// and returns it. No renderer needs to be implemented. If the simpleRenderer
// doesn't do it, SimpleWidget is probably not suitable for the use case.
// Usually this should not be overwritten or called manually.
func (s *SimpleWidgetBase) CreateRenderer() fyne.WidgetRenderer {
	wdgt := s.super()
	objs, layout := wdgt.Build()

	s.propertyLock.Lock()
	s.layout = layout
	s.propertyLock.Unlock()

	renderer := &simpleRenderer{
		widget:  wdgt,
		objects: objs,
		layout:  layout,
		destroy: s.Destroy,
		refresh: s.Refresh,
		minSize: s.MinSize,
	}
	return renderer
}

// SetState sets or changes the state of a widget. A Refresh
// is triggered after the state changes have been applied.
func (s *SimpleWidgetBase) SetState(setState func()) {
	setState()
	s.super().Refresh()
}

// SetStateSafe sets or changes the state of a widget in a safe way. A Refresh
// is triggered after the state changes have been applied.
// The provided sync.Locker should be the same you use for read protection of the
// widget properties.
func (s *SimpleWidgetBase) SetStateSafe(m sync.Locker, setState func()) {
	m.Lock()
	setState()
	m.Unlock()

	s.super().Refresh()
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (s *SimpleWidgetBase) ExtendBaseWidget(wid SimpleWidget) {
	impl := s.getImpl()
	if impl != nil {
		return
	}

	s.BaseWidget.ExtendBaseWidget(wid)

	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.impl = wid
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// By default it triggers a call to Layout. Overwrite to optimize for performance if
// a call to Layout is not needed.
func (s *SimpleWidgetBase) Refresh() {
	s.propertyLock.RLock()
	layout := s.layout
	s.propertyLock.RUnlock()

	if layout == nil {
		return
	}
	layout(s.Size())
}

// MinSize for the widget - it should never be resized below this value.
// By default this returns (0, 0). Overwrite this to return a different
// minimum size for the widget.
func (s *SimpleWidgetBase) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

// Destroy can be overwritten if there is extra cleanup needed.
func (s *SimpleWidgetBase) Destroy() {}

func (s *SimpleWidgetBase) super() SimpleWidget {
	impl := s.getImpl()
	if impl == nil {
		var x interface{} = s
		return x.(SimpleWidget)
	}
	return impl
}

func (s *SimpleWidgetBase) getImpl() SimpleWidget {
	s.propertyLock.RLock()
	impl := s.impl
	s.propertyLock.RUnlock()

	if impl == nil {
		return nil
	}
	return impl
}
