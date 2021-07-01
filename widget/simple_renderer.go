package widget

import (
	"fyne.io/fyne/v2"
)

// simpleRenderer is a renderer providing the basic rendering functionality used
// by SimpleWidgetBase.
type simpleRenderer struct {
	widget  fyne.Widget
	objects []fyne.CanvasObject

	layout  func(size fyne.Size)
	destroy func()
	refresh func()
	minSize func() fyne.Size
}

// MinSize returns the minimum size of the widget that is rendered by this renderer.
func (s *simpleRenderer) MinSize() fyne.Size {
	return s.minSize()
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a Layout.
func (s *simpleRenderer) Refresh() {
	s.refresh()
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
func (s *simpleRenderer) Destroy() {
	s.destroy()
}
