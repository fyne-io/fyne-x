package diagramwidget

import (
	"image/color"

	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Validate that DiagramNode implements DiagramElement and Tappable
var _ DiagramElement = (*DiagramNode)(nil)
var _ fyne.Tappable = (*DiagramNode)(nil)

const (
	// default inner size
	defaultWidth  float32 = 50
	defaultHeight float32 = 25
)

// DiagramNode represents a node in the diagram widget. It contains an inner
// widget, and also draws a border, and a "handle" that can be used to drag it
// around.
type DiagramNode struct {
	widget.BaseWidget
	diagramElement
	// InnerSize stores size that the inner object should have, may not
	// be respected if not large enough for the object.
	InnerSize fyne.Size
	// innerObject is the canvas object that should be drawn inside of
	// the diagram node.
	innerObject fyne.CanvasObject
	// Padding is the distance between the inner object's drawing area
	// and the box.
	Padding float32
	// BoxStrokeWidth is the stroke width of the box which delineates the
	// node. Defaults to 1.
	BoxStrokeWidth float32
	// BoxFill is the fill color of the node, the inner object will be
	// drawn on top of this. Defaults to the DiagramTheme's BackgroundColor.
	HandleColor color.Color
	// HandleStrokeWidth is the stroke width of the node handle, defaults
	// to 3.
	HandleStroke float32
	edgePad      *RectanglePad
}

// NewDiagramNode creates a DiagramNode widget and adds it to the DiagramWidget. The user-supplied
// nodeID string must be unique across all of the DiagramElements in the diagram. It can be used
// to retrieve the DiagramNode from the DiagramWidget. It is permissible for the canvas object to
// be nil when this function is called and then add the canvas object later.
func NewDiagramNode(diagram *DiagramWidget, obj fyne.CanvasObject, nodeID string) *DiagramNode {
	dn := &DiagramNode{
		InnerSize:      fyne.Size{Width: defaultWidth, Height: defaultHeight},
		innerObject:    obj,
		Padding:        diagram.DiagramTheme.Size(theme.SizeNamePadding),
		BoxStrokeWidth: 1,
		HandleColor:    diagram.DiagramTheme.Color(theme.ColorNameForeground, diagram.ThemeVariant),
		HandleStroke:   3,
	}
	dn.diagramElement.initialize(diagram, nodeID)
	dn.edgePad = NewRectanglePad(dn)
	dn.edgePad.Hide()
	for _, handleKey := range []string{"upperLeft", "upperMiddle", "upperRight", "leftMiddle", "rightMiddle", "lowerLeft", "lowerMiddle", "lowerRight"} {
		newHandle := NewHandle(dn)
		dn.handles[handleKey] = newHandle
		newHandle.Hide()
	}
	dn.ExtendBaseWidget(dn)
	dn.diagram.AddNode(dn)
	dn.Refresh()
	return dn
}

func (dn *DiagramNode) CreateRenderer() fyne.WidgetRenderer {
	dnr := diagramNodeRenderer{
		node: dn,
		box:  canvas.NewRectangle(dn.diagram.GetForegroundColor()),
	}

	dnr.box.StrokeWidth = dn.BoxStrokeWidth
	dnr.box.FillColor = dn.diagram.GetBackgroundColor()

	(&dnr).Refresh()

	return &dnr
}

func (dn *DiagramNode) Center() fyne.Position {
	return fyne.Position{X: float32(dn.R2Center().X), Y: float32(dn.R2Center().Y)}
}

func (dn *DiagramNode) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (dn *DiagramNode) DragEnd() {
}

func (dn *DiagramNode) Dragged(event *fyne.DragEvent) {
	dn.diagram.DiagramNodeDragged(dn, event)
}

func (dn *DiagramNode) effectiveInnerSize() fyne.Size {
	if dn.innerObject == nil {
		return dn.InnerSize
	}
	return dn.InnerSize.Max(dn.innerObject.MinSize())
}

func (dn *DiagramNode) findKeyForHandle(handle *Handle) string {
	for k, v := range dn.handles {
		if v == handle {
			return k
		}
	}
	return ""
}

func (dn *DiagramNode) GetEdgePad() ConnectionPad {
	return dn.edgePad
}

func (dn *DiagramNode) handleDragged(handle *Handle, event *fyne.DragEvent) {
	// determine which handle it is
	handleKey := dn.findKeyForHandle(handle)
	positionChange := fyne.Position{X: 0, Y: 0}
	sizeChange := fyne.Size{Height: 0, Width: 0}
	switch handleKey {
	case "upperLeft":
		positionChange.X = event.Dragged.DX
		positionChange.Y = event.Dragged.DY
		sizeChange.Height = -event.Dragged.DY
		sizeChange.Width = -event.Dragged.DY
	case "upperMiddle":
		positionChange.Y = event.Dragged.DY
		sizeChange.Height = -event.Dragged.DY
	case "upperRight":
		positionChange.Y = event.Dragged.DY
		sizeChange.Height = -event.Dragged.DY
		sizeChange.Width = event.Dragged.DX
	case "leftMiddle":
		positionChange.X = event.Dragged.DX
		sizeChange.Width = -event.Dragged.DX
	case "rightMiddle":
		sizeChange.Width = event.Dragged.DX
	case "lowerLeft":
		positionChange.X = event.Dragged.DX
		sizeChange.Height = event.Dragged.DY
		sizeChange.Width = -event.Dragged.DX
	case "lowerMiddle":
		sizeChange.Height = event.Dragged.DY
	case "lowerRight":
		sizeChange.Height = event.Dragged.DY
		sizeChange.Width = event.Dragged.DX
	}
	dn.Move(dn.Position().Add(positionChange))
	trialInnerSize := dn.InnerSize.Add(sizeChange)
	dn.InnerSize = dn.innerObject.MinSize().Max(trialInnerSize)
	dn.Resize(dn.Size().Add(sizeChange))
	dn.Refresh()
	dn.GetDiagram().forceRepaint()
}

func (dn *DiagramNode) innerPos() fyne.Position {
	return fyne.Position{
		X: float32(dn.Padding),
		Y: float32(dn.Padding),
	}
}

func (dn *DiagramNode) MouseIn(event *desktop.MouseEvent) {
	dn.HandleColor = dn.diagram.DiagramTheme.Color(theme.ColorNameFocus, dn.diagram.ThemeVariant)
	dn.GetDiagram().forceRepaint()
}

func (dn *DiagramNode) MouseOut() {
	dn.HandleColor = dn.diagram.DiagramTheme.Color(theme.ColorNameForeground, dn.diagram.ThemeVariant)
	dn.GetDiagram().forceRepaint()
}

func (dn *DiagramNode) MouseMoved(event *desktop.MouseEvent) {
}

func (dn *DiagramNode) Move(position fyne.Position) {
	dn.BaseWidget.Move(position)
	dn.Refresh()
}

func (dn *DiagramNode) R2Box() r2.Box {
	inner := dn.effectiveInnerSize()
	s := r2.V2(
		float64(inner.Width+float32(2*dn.Padding)),
		float64(inner.Height+float32(2*dn.Padding)),
	)

	return r2.MakeBox(dn.R2Position(), s)
}

func (dn *DiagramNode) R2Center() r2.Vec2 {
	return dn.R2Box().Center()
}

func (dn *DiagramNode) R2Position() r2.Vec2 {
	return r2.V2(float64(dn.Position().X), float64(dn.Position().Y))
}

func (dn *DiagramNode) SetInnerObject(obj fyne.CanvasObject) {
	dn.innerObject = obj
	dn.Refresh()
	dn.diagram.refreshDependentLinks(dn)
}

func (dn *DiagramNode) Tapped(event *fyne.PointEvent) {
	dn.diagram.DiagramElementTapped(dn, event)
}

// diagramNodeRenderer
type diagramNodeRenderer struct {
	node *DiagramNode
	box  *canvas.Rectangle
}

func (dnr *diagramNodeRenderer) ApplyTheme(size fyne.Size) {
}

func (dnr *diagramNodeRenderer) BackgroundColor() color.Color {
	return dnr.node.diagram.DiagramTheme.Color(theme.ColorNameBackground, dnr.node.diagram.ThemeVariant)
}

func (dnr *diagramNodeRenderer) Destroy() {
}

func (dnr *diagramNodeRenderer) MinSize() fyne.Size {
	// space for the inner widget, plus padding on all sides.
	inner := dnr.node.effectiveInnerSize()
	return fyne.Size{
		Width:  inner.Width + float32(2*dnr.node.Padding),
		Height: inner.Height + float32(2*dnr.node.Padding),
	}
}

func (dnr *diagramNodeRenderer) Layout(size fyne.Size) {
}

func (dnr *diagramNodeRenderer) Objects() []fyne.CanvasObject {
	obj := make([]fyne.CanvasObject, 0)
	obj = append(obj, dnr.box)
	obj = append(obj, dnr.node.edgePad)
	obj = append(obj, dnr.node.innerObject)
	for _, handle := range dnr.node.handles {
		obj = append(obj, handle)
	}
	return obj
}

func (dnr *diagramNodeRenderer) Refresh() {
	nodeSize := dnr.MinSize()
	dnr.node.Resize(nodeSize)
	dnr.node.edgePad.Resize(nodeSize)
	dnr.node.edgePad.Move(fyne.NewPos(0, 0))

	if dnr.node.innerObject != nil {
		dnr.node.innerObject.Move(dnr.node.innerPos())
		dnr.node.innerObject.Resize(dnr.node.effectiveInnerSize())
	}

	dnr.box.Resize(nodeSize)

	// calculate the handle positions
	width := nodeSize.Width
	height := nodeSize.Height
	for key, handle := range dnr.node.handles {
		switch key {
		case "upperLeft":
			handle.Move(fyne.NewPos(0, 0))
		case "upperMiddle":
			handle.Move(fyne.Position{X: width / 2, Y: 0})
		case "upperRight":
			handle.Move(fyne.Position{X: width, Y: 0})
		case "leftMiddle":
			handle.Move(fyne.Position{X: 0, Y: height / 2})
		case "rightMiddle":
			handle.Move(fyne.Position{X: width, Y: height / 2})
		case "lowerLeft":
			handle.Move(fyne.Position{X: 0, Y: height})
		case "lowerMiddle":
			handle.Move(fyne.Position{X: width / 2, Y: height})
		case "lowerRight":
			handle.Move(fyne.Position{X: width, Y: height})
		}
	}

	dnr.box.StrokeWidth = dnr.node.BoxStrokeWidth
	dnr.box.FillColor = color.Transparent
	dnr.box.StrokeColor = dnr.node.diagram.GetForegroundColor()
	dnr.node.edgePad.Refresh()
	dnr.node.diagram.refreshDependentLinks(dnr.node)
	dnr.node.GetDiagram().forceRepaint()
}
