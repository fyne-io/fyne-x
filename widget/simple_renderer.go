package widget

import (
	"fyne.io/fyne/v2"
)

// newSimpleRenderer creates a new simpleRenderer.
func newSimpleRenderer(widget SimpleWidget, objects []fyne.CanvasObject, layout func(size fyne.Size)) *simpleRenderer {
	return &simpleRenderer{
		widget:  widget,
		objects: objects,
		layout:  layout,
	}
}

// simpleRenderer is a renderer providing the basic rendering functionality used
// by SimpleWidgetBase.
type simpleRenderer struct {
	widget  fyne.Widget
	objects []fyne.CanvasObject

	layout func(size fyne.Size)
}

// MinSize returns the minimum size of the widget that is rendered by this renderer.
func (s *simpleRenderer) MinSize() fyne.Size {
	return fyne.Size{Width: 0, Height: 0}
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a Layout.
func (s *simpleRenderer) Refresh() {
	s.Layout(s.widget.Size())
}

// Layout is a hook that is called if the widget needs to be laid out.
// This should not be overwritten.
func (s *simpleRenderer) Layout(size fyne.Size) {
	s.layout(size)
}

// Objects returns the objects that should be rendered.
//
// Implements: fyne.WidgetRenderer
func (s *simpleRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

// SetObjects updates the objects of the renderer.
func (s *simpleRenderer) SetObjects(objects []fyne.CanvasObject) {
	s.objects = objects
}

// Destroy does nothing in the base implementation.
//
// Implements: fyne.WidgetRenderer
func (s *simpleRenderer) Destroy() {}
