package diagramwidget

import (
	"image/color"
	"log"

	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

type DiagramNode interface {
	DiagramElement
	getBaseDiagramNode() *BaseDiagramNode
	MouseIn(event *desktop.MouseEvent)
	MouseOut()
	MouseMoved(event *desktop.MouseEvent)
	R2Center() r2.Vec2
}

// Validate that BaseDiagramNode implements DiagramElement and Tappable
var _ DiagramElement = (*BaseDiagramNode)(nil)
var _ fyne.Tappable = (*BaseDiagramNode)(nil)
var _ desktop.Hoverable = (*BaseDiagramNode)(nil)
var _ desktop.Hoverable = (DiagramNode)(nil)
var _ fyne.Widget = (*BaseDiagramNode)(nil)
var _ fyne.Widget = (DiagramNode)(nil)

const (
	// default inner size
	defaultWidth  float32 = 50
	defaultHeight float32 = 25
)

// BaseDiagramNode represents a node in the diagram widget. It contains an inner
// widget, and also draws a border, and a "handle" that can be used to drag it
// around.
type BaseDiagramNode struct {
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
	// HandleStrokeWidth is the stroke width of the node handle, defaults
	// to 3.
	HandleStroke float32
	edgePad      *RectanglePad
	// MovedCallback, if present, is invoked when the node is moved
	MovedCallback func()
}

// NewDiagramNode creates a DiagramNode widget and adds it to the DiagramWidget. The user-supplied
// nodeID string must be unique across all of the DiagramElements in the diagram. It can be used
// to retrieve the DiagramNode from the DiagramWidget. It is permissible for the canvas object to
// be nil when this function is called and then add the canvas object later.
func NewDiagramNode(diagram *DiagramWidget, obj fyne.CanvasObject, nodeID string) *BaseDiagramNode {
	bdn := &BaseDiagramNode{}
	InitializeBaseDiagramNode(bdn, diagram, obj, nodeID)
	return bdn
}

// InitializeBaseDiagramNode is used to initailize the BaseDiagramNode. It must be called by any extensions to the BaseDiagramNode
func InitializeBaseDiagramNode(diagramNode DiagramNode, diagram *DiagramWidget, obj fyne.CanvasObject, nodeID string) {
	bdn := diagramNode.getBaseDiagramNode()
	bdn.InnerSize = fyne.Size{Width: defaultWidth, Height: defaultHeight}
	bdn.innerObject = obj
	bdn.Padding = diagram.DiagramTheme.Size(theme.SizeNamePadding)
	bdn.BoxStrokeWidth = 1
	bdn.HandleStroke = 3
	bdn.diagramElement.initialize(diagram, nodeID)
	bdn.edgePad = NewRectanglePad(bdn)
	bdn.edgePad.Hide()
	for _, handleKey := range []string{"upperLeft", "upperMiddle", "upperRight", "leftMiddle", "rightMiddle", "lowerLeft", "lowerMiddle", "lowerRight"} {
		newHandle := NewHandle(bdn)
		bdn.handles[handleKey] = newHandle
		newHandle.Hide()
	}
	bdn.ExtendBaseWidget(diagramNode)
	bdn.diagram.addNode(diagramNode)
	diagramNode.Refresh()
}

func (bdn *BaseDiagramNode) CreateRenderer() fyne.WidgetRenderer {
	dnr := diagramNodeRenderer{
		node: bdn,
		box:  canvas.NewRectangle(bdn.diagram.GetForegroundColor()),
	}

	dnr.box.StrokeWidth = bdn.BoxStrokeWidth
	dnr.box.FillColor = bdn.diagram.GetBackgroundColor()

	(&dnr).Refresh()

	return &dnr
}

func (bdn *BaseDiagramNode) Center() fyne.Position {
	return fyne.Position{X: float32(bdn.R2Center().X), Y: float32(bdn.R2Center().Y)}
}

func (bdn *BaseDiagramNode) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (bdn *BaseDiagramNode) DragEnd() {
}

func (bdn *BaseDiagramNode) Dragged(event *fyne.DragEvent) {
	bdn.diagram.DiagramNodeDragged(bdn, event)
}

func (bdn *BaseDiagramNode) effectiveInnerSize() fyne.Size {
	if bdn.innerObject == nil {
		return bdn.InnerSize
	}
	return bdn.InnerSize.Max(bdn.innerObject.MinSize())
}

func (bdn *BaseDiagramNode) findKeyForHandle(handle *Handle) string {
	for k, v := range bdn.handles {
		if v == handle {
			return k
		}
	}
	return ""
}

func (bdn *BaseDiagramNode) getBaseDiagramNode() *BaseDiagramNode {
	return bdn
}

// GetDefaultConnectionPad returns the edge pad for the node
func (bdn *BaseDiagramNode) GetDefaultConnectionPad() ConnectionPad {
	return bdn.GetEdgePad()
}

// GetEdgePad returns the edge pad for the node
func (bdn *BaseDiagramNode) GetEdgePad() ConnectionPad {
	return bdn.edgePad
}

func (bdn *BaseDiagramNode) handleDragged(handle *Handle, event *fyne.DragEvent) {
	// determine which handle it is
	currentInnerSize := bdn.effectiveInnerSize()
	handleKey := bdn.findKeyForHandle(handle)
	positionChange := fyne.Position{X: 0, Y: 0}
	sizeChange := fyne.Size{Height: 0, Width: 0}
	switch handleKey {
	case "upperLeft":
		positionChange.X = event.Dragged.DX
		sizeChange.Width = -event.Dragged.DX
		positionChange.Y = event.Dragged.DY
		sizeChange.Height = -event.Dragged.DY
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
		sizeChange.Width = -event.Dragged.DX
		sizeChange.Height = event.Dragged.DY
	case "lowerMiddle":
		sizeChange.Height = event.Dragged.DY
	case "lowerRight":
		sizeChange.Height = event.Dragged.DY
		sizeChange.Width = event.Dragged.DX
	}
	trialInnerSize := bdn.InnerSize.Add(sizeChange)
	bdn.InnerSize = bdn.innerObject.MinSize().Max(trialInnerSize)
	if trialInnerSize.Height < bdn.InnerSize.Height {
		sizeChange.Height = bdn.InnerSize.Height - currentInnerSize.Height
		if positionChange.Y != 0 {
			positionChange.Y = -sizeChange.Height
		}
	}
	if trialInnerSize.Width < bdn.InnerSize.Width {
		sizeChange.Width = bdn.InnerSize.Width - currentInnerSize.Width
		if positionChange.X != 0 {
			positionChange.X = -sizeChange.Width
		}
	}
	bdn.Resize(bdn.Size().Add(sizeChange))
	bdn.Move(bdn.Position().Add(positionChange))
	bdn.Refresh()
	bdn.GetDiagram().ForceRepaint()
}

func (bdn *BaseDiagramNode) innerPos() fyne.Position {
	return fyne.Position{
		X: float32(bdn.Padding),
		Y: float32(bdn.Padding),
	}
}

func (bdn *BaseDiagramNode) MouseIn(event *desktop.MouseEvent) {
	log.Print("Node MouseIn")
	bdn.handleColor = bdn.diagram.DiagramTheme.Color(theme.ColorNameFocus, bdn.diagram.ThemeVariant)
	bdn.Refresh()
	bdn.GetDiagram().ForceRepaint()
}

func (bdn *BaseDiagramNode) MouseOut() {
	log.Print("Node MouseIn")
	bdn.handleColor = bdn.foregroundColor
	bdn.Refresh()
	bdn.GetDiagram().ForceRepaint()
}

func (bdn *BaseDiagramNode) MouseMoved(event *desktop.MouseEvent) {
}

func (bdn *BaseDiagramNode) Move(position fyne.Position) {
	bdn.BaseWidget.Move(position)
	if bdn.MovedCallback != nil {
		bdn.MovedCallback()
	}
	bdn.Refresh()
	bdn.diagram.ForceRepaint()
}

func (bdn *BaseDiagramNode) R2Box() r2.Box {
	inner := bdn.effectiveInnerSize()
	s := r2.V2(
		float64(inner.Width+float32(2*bdn.Padding)),
		float64(inner.Height+float32(2*bdn.Padding)),
	)

	return r2.MakeBox(bdn.R2Position(), s)
}

func (bdn *BaseDiagramNode) R2Center() r2.Vec2 {
	return bdn.R2Box().Center()
}

func (bdn *BaseDiagramNode) R2Position() r2.Vec2 {
	return r2.V2(float64(bdn.Position().X), float64(bdn.Position().Y))
}

func (bdn *BaseDiagramNode) SetInnerObject(obj fyne.CanvasObject) {
	bdn.innerObject = obj
	bdn.Refresh()
	bdn.diagram.refreshDependentLinks(bdn)
}

func (bdn *BaseDiagramNode) Tapped(event *fyne.PointEvent) {
	bdn.diagram.DiagramElementTapped(bdn, event)
}

// diagramNodeRenderer
type diagramNodeRenderer struct {
	node *BaseDiagramNode
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
	dnr.node.edgePad.Refresh()

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
		handle.Resize(fyne.NewSize(handle.handleSize, handle.handleSize))
		handle.Refresh()
	}

	dnr.box.StrokeWidth = dnr.node.BoxStrokeWidth
	dnr.box.FillColor = color.Transparent
	dnr.box.StrokeColor = dnr.node.diagram.GetForegroundColor()

	dnr.node.edgePad.Refresh()
	dnr.node.diagram.refreshDependentLinks(dnr.node)
	dnr.node.GetDiagram().ForceRepaint()
}
