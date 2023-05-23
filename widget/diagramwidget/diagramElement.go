package diagramwidget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// A DiagramElement is a widget that can be placed directly in a diagram. The most common
// elements are Node and Link widgets.
type DiagramElement interface {
	fyne.Widget
	// GetBackgroundColor returns the background color for the widget
	GetBackgroundColor() color.Color
	// GetConnectionPads() returns all of the connection pads on the element
	GetConnectionPads() map[string]ConnectionPad
	// GetForegroundColor returns the foreground color for the widget
	GetForegroundColor() color.Color
	// GetDefaultConnectionPad returns the default pad for the DiagramElement
	GetDefaultConnectionPad() ConnectionPad
	// GetDiagram returns the DiagramWidget to which the DiagramElement belongs
	GetDiagram() *DiagramWidget
	// GetDiagramElementID returns the string identifier provided at the time the DiagramElement was created
	GetDiagramElementID() string
	// GetHandle returns the handle with the indicated index name
	GetHandle(string) *Handle
	// GetHandleColor returns the color for the element's handles
	GetHandleColor() color.Color
	// handleDragged responds to drag events
	handleDragged(handle *Handle, event *fyne.DragEvent)
	// handleDragEnd responds to the end of a drag
	handleDragEnd(handle *Handle)
	// HideHandles hides the handles on the DiagramElement
	HideHandles()
	// ShowHandles shows the handles on the DiagramElement
	ShowHandles()
	// SetForegroundColor sets the foreground color for the widget
	SetForegroundColor(color.Color)
	// SetBackgroundColor sets the background color for the widget
	SetBackgroundColor(color.Color)
}

type diagramElement struct {
	widget.BaseWidget
	diagram         *DiagramWidget
	foregroundColor color.Color
	backgroundColor color.Color
	handleColor     color.Color
	id              string
	handles         map[string]*Handle
	pads            map[string]ConnectionPad
}

func (de *diagramElement) GetDiagram() *DiagramWidget {
	return de.diagram
}

func (de *diagramElement) GetDiagramElementID() string {
	return de.id
}

func (de *diagramElement) GetBackgroundColor() color.Color {
	return de.backgroundColor
}

func (de *diagramElement) GetConnectionPads() map[string]ConnectionPad {
	return de.pads
}

func (de *diagramElement) GetForegroundColor() color.Color {
	return de.foregroundColor
}

func (de *diagramElement) GetHandle(handleName string) *Handle {
	return de.handles[handleName]
}

func (de *diagramElement) GetHandleColor() color.Color {
	return de.handleColor
}

func (de *diagramElement) HideHandles() {
	for _, handle := range de.handles {
		handle.Hide()
	}
}

func (de *diagramElement) initialize(diagram *DiagramWidget, id string) {
	de.diagram = diagram
	de.id = id
	de.handles = make(map[string]*Handle)
	de.handleColor = de.diagram.DiagramTheme.Color(theme.ColorNameForeground, de.diagram.ThemeVariant)
	de.pads = make(map[string]ConnectionPad)
}

func (de *diagramElement) SetBackgroundColor(backgroundColor color.Color) {
	de.backgroundColor = backgroundColor
	de.Refresh()
}

func (de *diagramElement) SetForegroundColor(foregroundColor color.Color) {
	de.foregroundColor = foregroundColor
	de.Refresh()
}

func (de *diagramElement) SetHandleColor(handleColor color.Color) {
	de.handleColor = handleColor
	de.Refresh()
}

func (de *diagramElement) ShowHandles() {
	for _, handle := range de.handles {
		handle.Show()
	}
}
