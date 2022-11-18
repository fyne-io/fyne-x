package wrapper

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*mouseableObject)(nil)
var _ desktop.Hoverable = (*mouseableObject)(nil)

// handles the canvas object and mouse events
type mouseableObject struct {
	widget.BaseWidget
	object     fyne.CanvasObject
	mouseIn    func(*desktop.MouseEvent)
	mouseMoved func(*desktop.MouseEvent)
	mouseOut   func()
}

// Content returns the originl object that was set to be hoverable.
func (m *mouseableObject) Content() fyne.CanvasObject {
	return m.object
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
//
// Implements: fyne.Widget
func (m *mouseableObject) CreateRenderer() fyne.WidgetRenderer {
	if m.object == nil {
		return nil
	}
	if o, ok := m.object.(fyne.Widget); ok {
		return o.CreateRenderer()
	}
	return widget.NewSimpleRenderer(m.object)
}

// MouseIn is called when the mouse enters the widget.
//
// Implements: desktop.Hoverable
func (m *mouseableObject) MouseIn(e *desktop.MouseEvent) {
	if m.mouseIn == nil {
		return
	}

	if o, ok := m.object.(desktop.Hoverable); ok {
		o.MouseIn(e)
	}

	m.mouseIn(e)
}

// MouseMoved is called when the mouse moves over the widget.
//
// Implements: desktop.Hoverable
func (m *mouseableObject) MouseMoved(e *desktop.MouseEvent) {
	if m.mouseMoved == nil {
		return
	}

	if o, ok := m.object.(desktop.Hoverable); ok {
		o.MouseMoved(e)
	}

	m.mouseMoved(e)
}

// MouseOut is called when the mouse exits the widget.
//
// Implements: desktop.Hoverable
func (m *mouseableObject) MouseOut() {
	if m.mouseOut == nil {
		return
	}

	if o, ok := m.object.(desktop.Hoverable); ok {
		o.MouseOut()
	}

	m.mouseOut()
}

// SetHoverable sets the object to be hoverable.
func SetHoverable(object fyne.CanvasObject, mouseIn, mouseMoved func(*desktop.MouseEvent), mouseout func()) fyne.CanvasObject {
	m := &mouseableObject{
		object:     object,
		mouseIn:    mouseIn,
		mouseMoved: mouseMoved,
		mouseOut:   mouseout,
	}
	m.ExtendBaseWidget(m)
	return m
}
