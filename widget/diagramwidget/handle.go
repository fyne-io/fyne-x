package diagramwidget

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// draggable is used to verify that objects implement the Draggable interface. An assignment to
// draggable will fail at compile time if the interface is not fully implemented.
var draggable fyne.Draggable
var hoverable desktop.Hoverable

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
	// This assignment verifies that handle fully implements the Draggable interface
	draggable = handle
	hoverable = handle
	handle.BaseWidget.ExtendBaseWidget(handle)
	return handle
}

func (h *Handle) CreateRenderer() fyne.WidgetRenderer {
	hr := &handleRenderer{
		handle: h,
		rect:   canvas.NewRectangle(h.getStrokeColor()),
	}
	hr.Refresh()
	ForceRefresh()
	return hr
}

func (h *Handle) Dragged(event *fyne.DragEvent) {
	h.de.handleDragged(h, event)
	log.Print("Handle dragged")
}

func (h *Handle) DragEnd() {

}

func (h *Handle) getStrokeColor() color.Color {
	variant := h.de.GetDiagram().ThemeVariant
	return h.de.GetDiagram().DiagramTheme.Color(theme.ColorNameForeground, variant)
}

func (h *Handle) getStrokeWidth() float32 {
	return 1.0
}

func (h *Handle) MouseIn(*desktop.MouseEvent) {
	log.Print("Mouse in handle")
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (h *Handle) MouseMoved(*desktop.MouseEvent) {
	log.Print("Mouse moved in handle")
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (h *Handle) MouseOut() {
	log.Print("Mouse out of handle")
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

func (hr *handleRenderer) MinSize() fyne.Size {
	return fyne.Size{Height: hr.handle.handleSize, Width: hr.handle.handleSize}
}

func (hr *handleRenderer) Layout(size fyne.Size) {
	hr.rect.Resize(hr.MinSize())
	hr.handle.Resize(hr.MinSize())
}

func (hr *handleRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{
		hr.rect,
	}
	return obj
}

func (hr *handleRenderer) Refresh() {
	hr.rect.StrokeColor = hr.handle.getStrokeColor()
	hr.rect.StrokeWidth = hr.handle.getStrokeWidth()
	ForceRefresh()
}
