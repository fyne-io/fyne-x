package diagramwidget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// Validate implementation of Draggable
var _ fyne.Draggable = (*Handle)(nil)

var defaultHandleSize float32 = 10.0

type Handle struct {
	widget.BaseWidget
	handleSize float32
	de         DiagramElement
}

func NewHandle(diagramElement DiagramElement) *Handle {
	handle := &Handle{
		de:         diagramElement,
		handleSize: defaultHandleSize,
	}
	handle.BaseWidget.ExtendBaseWidget(handle)
	return handle
}

func (h *Handle) CreateRenderer() fyne.WidgetRenderer {
	hr := &handleRenderer{
		handle: h,
		rect:   canvas.NewRectangle(h.getStrokeColor()),
	}
	hr.rect.FillColor = color.Transparent
	hr.Refresh()
	return hr
}

func (h *Handle) Dragged(event *fyne.DragEvent) {
	h.de.handleDragged(h, event)
}

func (h *Handle) DragEnd() {
	h.de.handleDragEnd(h)
}

func (h *Handle) getStrokeColor() color.Color {
	return h.de.GetDiagram().GetForegroundColor()
}

func (h *Handle) getStrokeWidth() float32 {
	return 1.0
}

func (h *Handle) Move(position fyne.Position) {
	delta := fyne.Position{X: -h.handleSize / 2, Y: -h.handleSize / 2}
	h.BaseWidget.Move(position.Add(delta))
}

// handleRenderer
type handleRenderer struct {
	handle *Handle
	rect   *canvas.Rectangle
}

func (hr *handleRenderer) Destroy() {

}

func (hr *handleRenderer) Layout(size fyne.Size) {
	hr.rect.Resize(hr.MinSize())
	hr.handle.Resize(hr.MinSize())
}

func (hr *handleRenderer) MinSize() fyne.Size {
	return fyne.Size{Height: hr.handle.handleSize, Width: hr.handle.handleSize}
}

func (hr *handleRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{
		hr.rect,
	}
	return obj
}

func (hr *handleRenderer) Refresh() {
	hr.rect.StrokeColor = hr.handle.getStrokeColor()
	hr.rect.FillColor = color.Transparent
	hr.rect.StrokeWidth = hr.handle.getStrokeWidth()
	hr.rect.Refresh()
}
