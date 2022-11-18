package wrapper

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*tappableObject)(nil)
var _ fyne.Tappable = (*tappableObject)(nil)

// handles the tap event.
type tappableObject struct {
	widget.BaseWidget
	object   fyne.CanvasObject
	onTapped func(*fyne.PointEvent)
}

// Content returns the encapsulated widget.
func (t *tappableObject) Content() fyne.CanvasObject {
	return t.object
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
//
// Implements: fyne.Widget
func (t *tappableObject) CreateRenderer() fyne.WidgetRenderer {
	if t.object == nil {
		return nil
	}
	if o, ok := t.object.(fyne.Widget); ok {
		return o.CreateRenderer()
	}
	return widget.NewSimpleRenderer(t.object)
}

// Tapped reacts on tap (or click) events.
//
// Implements: fyne.Tappable
func (t *tappableObject) Tapped(e *fyne.PointEvent) {
	if t.object == nil {
		return
	}

	if o, ok := t.object.(fyne.Tappable); ok {
		o.Tapped(e)
	}
	t.onTapped(e)
}

// SetTappable set the object to be tappable.
func SetTappable(object fyne.CanvasObject, ontapped func(*fyne.PointEvent)) fyne.CanvasObject {
	tappable := &tappableObject{object: object, onTapped: ontapped}
	tappable.ExtendBaseWidget(tappable)
	return tappable
}
