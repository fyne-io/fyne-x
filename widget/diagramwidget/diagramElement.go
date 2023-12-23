package diagramwidget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// DiagramElementProperties are the rendering properties of a DiagramElement
type DiagramElementProperties struct {
	ForegroundColor   color.Color
	BackgroundColor   color.Color
	HandleColor       color.Color
	PadColor          color.Color
	TextSize          float32
	CaptionTextSize   float32
	Padding           float32
	StrokeWidth       float32
	PadStrokeWidth    float32
	HandleStrokeWidth float32
}

// DiagramElement is a widget that can be placed directly in a diagram. The most common
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
	// GetPadColor returns the color for the element's pads
	GetPadColor() color.Color
	// GetProperties returns the properties of the DiagramElement
	GetProperties() DiagramElementProperties
	// handleDragged responds to drag events
	handleDragged(handle *Handle, event *fyne.DragEvent)
	// handleDragEnd responds to the end of a drag
	handleDragEnd(handle *Handle)
	// HideHandles hides the handles on the DiagramElement
	HideHandles()
	// IsLink returns true if the diagram element is a link
	IsLink() bool
	// IsNode returns true of the diagram element is a node
	IsNode() bool
	// Position returns the position of the diagram element
	Position() fyne.Position
	// SetForegroundColor sets the foreground color for the widget
	SetForegroundColor(color.Color)
	// SetBackgroundColor sets the background color for the widget
	SetBackgroundColor(color.Color)
	// SetProperties sets the foreground, background, and handle colors
	SetProperties(DiagramElementProperties)
	// ShowHandles shows the handles on the DiagramElement
	ShowHandles()
	// Size returns the size of the diagram element
	Size() fyne.Size
}

type diagramElement struct {
	widget.BaseWidget
	diagram    *DiagramWidget
	properties DiagramElementProperties
	// foregroundColor color.Color
	// backgroundColor color.Color
	// handleColor     color.Color
	id      string
	handles map[string]*Handle
	pads    map[string]ConnectionPad
}

func (de *diagramElement) GetDiagram() *DiagramWidget {
	return de.diagram
}

func (de *diagramElement) GetDiagramElementID() string {
	return de.id
}

func (de *diagramElement) GetBackgroundColor() color.Color {
	return de.properties.BackgroundColor
}

func (de *diagramElement) GetConnectionPads() map[string]ConnectionPad {
	return de.pads
}

func (de *diagramElement) GetForegroundColor() color.Color {
	return de.properties.ForegroundColor
}

func (de *diagramElement) GetHandle(handleName string) *Handle {
	return de.handles[handleName]
}

func (de *diagramElement) GetHandleColor() color.Color {
	return de.properties.HandleColor
}

func (de *diagramElement) GetPadColor() color.Color {
	return de.properties.PadColor
}

func (de *diagramElement) GetProperties() DiagramElementProperties {
	return de.properties
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
	de.properties = de.diagram.DefaultDiagramElementProperties
	de.pads = make(map[string]ConnectionPad)
}

func (de *diagramElement) SetBackgroundColor(backgroundColor color.Color) {
	de.properties.BackgroundColor = backgroundColor
	de.Refresh()
}

func (de *diagramElement) SetForegroundColor(foregroundColor color.Color) {
	de.properties.ForegroundColor = foregroundColor
	de.Refresh()
}

func (de *diagramElement) SetHandleColor(handleColor color.Color) {
	de.properties.HandleColor = handleColor
	de.Refresh()
}

func (de *diagramElement) SetProperties(properties DiagramElementProperties) {
	de.properties = properties
}

func (de *diagramElement) ShowHandles() {
	for _, handle := range de.handles {
		handle.Show()
		de.Refresh()
	}
}
