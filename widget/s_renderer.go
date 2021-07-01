package widget

import (
	"sync"

	"fyne.io/fyne/v2"
)

type SimpleWidgetRenderer interface {
	fyne.WidgetRenderer

	Build() (objects []fyne.CanvasObject, layout func(size fyne.Size))
}

func NewBaseSimpleRenderer(wgt fyne.Widget, renderer SimpleWidgetRenderer) *BaseSimpleRenderer {
	rend := &BaseSimpleRenderer{
		widget: wgt,
		impl:   renderer,
	}
	objs, layout := rend.super().Build()
	rend.objects = objs
	rend.layout = layout

	return rend
}

type BaseSimpleRenderer struct {
	objects      []fyne.CanvasObject
	propertyLock sync.RWMutex
	widget       fyne.Widget
	impl         SimpleWidgetRenderer

	layout func(fyne.Size)
}

func (s *BaseSimpleRenderer) Build() (objects []fyne.CanvasObject, layout func(size fyne.Size)) {
	return nil, func(fyne.Size) {}
}

// Destroy can be overwritten if there is extra cleanup needed.
func (s *BaseSimpleRenderer) Destroy() {}

// Layout is a hook that is called if the widget needs to be laid out.
// This should not be overwritten.
func (s *BaseSimpleRenderer) Layout(size fyne.Size) {
	s.layout(size)
}

func (s *BaseSimpleRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (s *BaseSimpleRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// By default it triggers a call to Layout. Overwrite to optimize for performance if
// a call to Layout is not needed.
func (s *BaseSimpleRenderer) Refresh() {
	s.propertyLock.RLock()
	layout := s.layout
	s.propertyLock.RUnlock()

	if layout == nil {
		return
	}
	layout(s.widget.Size())
}

func (s *BaseSimpleRenderer) getImpl() SimpleWidgetRenderer {
	s.propertyLock.RLock()
	impl := s.impl
	s.propertyLock.RUnlock()

	if impl == nil {
		return nil
	}
	return impl
}

func (s *BaseSimpleRenderer) super() SimpleWidgetRenderer {
	impl := s.getImpl()
	if impl == nil {
		var x interface{} = s
		return x.(SimpleWidgetRenderer)
	}
	return impl
}
