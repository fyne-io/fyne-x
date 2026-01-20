package diagramwidget

import (
	"image/color"

	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
)

// DiagramNode is a rectangular DiagramElement typically containing one or more widgets
type DiagramNode interface {
	DiagramElement
	getBaseDiagramNode() *BaseDiagramNode
	GetEdgePad() ConnectionPad
	R2Center() r2.Vec2
	SetInnerObject(fyne.CanvasObject)
}

// Validate that BaseDiagramNode implements DiagramElement and Tappable
var _ DiagramElement = (*BaseDiagramNode)(nil)
var _ fyne.Tappable = (*BaseDiagramNode)(nil)

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
	typedNode DiagramElement
	// InnerSize stores size that the inner object should have, may not
	// be respected if not large enough for the object.
	InnerSize fyne.Size
	// innerObject is the canvas object that should be drawn inside of
	// the diagram node.
	innerObject fyne.CanvasObject
	// MovedCallback, if present, is invoked when the node is moved
	MovedCallback func()
}

// NewDiagramNode creates a DiagramNode widget and adds it to the DiagramWidget. The user-supplied
// nodeID string must be unique across all of the DiagramElements in the diagram. It can be used
// to retrieve the DiagramNode from the DiagramWidget. It is permissible for the canvas object to
// be nil when this function is called and then add the canvas object later.
func NewDiagramNode(diagram *DiagramWidget, obj fyne.CanvasObject, nodeID string) DiagramNode {
	var diagramNode DiagramNode = &BaseDiagramNode{}
	InitializeBaseDiagramNode(diagramNode, diagram, obj, nodeID)
	return diagramNode
}

// InitializeBaseDiagramNode is used to initailize the BaseDiagramNode. It must be called by any extensions to the BaseDiagramNode
func InitializeBaseDiagramNode(diagramNode DiagramNode, diagram *DiagramWidget, obj fyne.CanvasObject, nodeID string) {
	bdn := diagramNode.getBaseDiagramNode()
	bdn.typedNode = diagramNode
	bdn.InnerSize = fyne.Size{Width: defaultWidth, Height: defaultHeight}
	bdn.innerObject = obj
	bdn.diagramElement.initialize(diagram, nodeID)
	bdn.SetConnectionPad(NewRectanglePad(), "default")
	bdn.pads["default"].Hide()
	for _, handleKey := range []string{"upperLeft", "upperMiddle", "upperRight", "leftMiddle", "rightMiddle", "lowerLeft", "lowerMiddle", "lowerRight"} {
		newHandle := NewHandle(diagramNode)
		bdn.handles[handleKey] = newHandle
		newHandle.Hide()
	}
	bdn.ExtendBaseWidget(diagramNode)
	bdn.diagram.addNode(diagramNode)
	diagramNode.Refresh()
}

// CreateRenderer creates the renderer for the diagram node
func (bdn *BaseDiagramNode) CreateRenderer() fyne.WidgetRenderer {
	dnr := diagramNodeRenderer{
		node: bdn,
		box:  canvas.NewRectangle(bdn.diagram.GetForegroundColor()),
	}

	dnr.box.StrokeWidth = bdn.properties.StrokeWidth
	dnr.box.FillColor = bdn.diagram.GetBackgroundColor()

	(&dnr).Refresh()

	return &dnr
}

// Center reurns the position of the center of the node
func (bdn *BaseDiagramNode) Center() fyne.Position {
	return fyne.Position{X: float32(bdn.R2Center().X), Y: float32(bdn.R2Center().Y)}
}

// Cursor returns the desktop default cursor
func (bdn *BaseDiagramNode) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// DragEnd is presently a no-op
func (bdn *BaseDiagramNode) DragEnd() {
}

// Dragged passes the DragEvent to the diagram for processing
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
	return bdn.pads["default"]
}

// GetTypedElement returns the instantiated type of the node
func (bdn *BaseDiagramNode) GetTypedElement() DiagramElement {
	return bdn.typedNode
}

// HandleDragged modifies the node size when the handle is dragged
func (bdn *BaseDiagramNode) HandleDragged(handle *Handle, event *fyne.DragEvent) {
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
}

// HandleDragEnd determines node behavior when the handle drag ends. By default, it does nothing.
func (bdn *BaseDiagramNode) HandleDragEnd(handle *Handle) {
}

func (bdn *BaseDiagramNode) innerPos() fyne.Position {
	return fyne.Position{
		X: float32(bdn.properties.Padding),
		Y: float32(bdn.properties.Padding),
	}
}

// IsLink returns false because this is a node
func (bdn *BaseDiagramNode) IsLink() bool {
	return false
}

// IsNode returns true because this is a node
func (bdn *BaseDiagramNode) IsNode() bool {
	return true
}

// Move moves the node and invokes the callback if present.
func (bdn *BaseDiagramNode) Move(position fyne.Position) {
	bdn.BaseWidget.Move(position)
	if bdn.MovedCallback != nil {
		bdn.MovedCallback()
	}
	bdn.Refresh()
}

// R2Box returns the bounding box in r2 coordinates
func (bdn *BaseDiagramNode) R2Box() r2.Box {
	inner := bdn.effectiveInnerSize()
	s := r2.V2(
		float64(inner.Width+float32(2*bdn.properties.Padding)),
		float64(inner.Height+float32(2*bdn.properties.Padding)),
	)

	return r2.MakeBox(bdn.R2Position(), s)
}

// R2Center returns the r2 vector for the center of the bounding box
func (bdn *BaseDiagramNode) R2Center() r2.Vec2 {
	return bdn.R2Box().Center()
}

// R2Position returns the position of the node as an r2 vector
func (bdn *BaseDiagramNode) R2Position() r2.Vec2 {
	return r2.V2(float64(bdn.Position().X), float64(bdn.Position().Y))
}

// SetConnectionPad sets the connection pad for the indicated key.
func (bdn *BaseDiagramNode) SetConnectionPad(pad ConnectionPad, key string) {
	if pad != nil {
		pad.SetLineWidth(bdn.GetProperties().PadStrokeWidth)
		pad.setPadOwner(bdn)
	}
	bdn.pads[key] = pad
}

// SetInnerObject makes the skupplied canvas object the center of the node
func (bdn *BaseDiagramNode) SetInnerObject(obj fyne.CanvasObject) {
	bdn.innerObject = obj
	bdn.Refresh()
	bdn.diagram.refreshDependentLinks(bdn)
}

// Tapped passes the tapped event on to the Diagram
func (bdn *BaseDiagramNode) Tapped(event *fyne.PointEvent) {
	bdn.diagram.DiagramElementTapped(bdn)
}

// diagramNodeRenderer
type diagramNodeRenderer struct {
	node *BaseDiagramNode
	box  *canvas.Rectangle
}

func (dnr *diagramNodeRenderer) ApplyTheme(size fyne.Size) {
}

func (dnr *diagramNodeRenderer) BackgroundColor() color.Color {
	return dnr.node.properties.BackgroundColor
}

func (dnr *diagramNodeRenderer) Destroy() {
}

func (dnr *diagramNodeRenderer) MinSize() fyne.Size {
	// space for the inner widget, plus padding on all sides.
	inner := dnr.node.effectiveInnerSize()
	return fyne.Size{
		Width:  inner.Width + float32(2*dnr.node.properties.Padding),
		Height: inner.Height + float32(2*dnr.node.properties.Padding),
	}
}

func (dnr *diagramNodeRenderer) Layout(size fyne.Size) {
}

func (dnr *diagramNodeRenderer) Objects() []fyne.CanvasObject {
	obj := make([]fyne.CanvasObject, 0)
	obj = append(obj, dnr.box)
	obj = append(obj, dnr.node.innerObject)
	for _, pad := range dnr.node.pads {
		if pad != nil {
			obj = append(obj, pad)
		}
	}
	for _, handle := range dnr.node.handles {
		obj = append(obj, handle)
	}
	return obj
}

func (dnr *diagramNodeRenderer) Refresh() {
	nodeSize := dnr.MinSize()
	dnr.node.Resize(nodeSize)
	if dnr.node.pads["default"] != nil {
		dnr.node.pads["default"].Resize(nodeSize)
		dnr.node.pads["default"].Refresh()
	}

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

	dnr.box.StrokeWidth = dnr.node.properties.StrokeWidth
	dnr.box.FillColor = dnr.node.properties.BackgroundColor
	dnr.box.StrokeColor = dnr.node.properties.ForegroundColor
	dnr.box.Refresh()

	for _, pad := range dnr.node.pads {
		if pad != nil {
			pad.Refresh()
		}
	}
	dnr.node.diagram.refreshDependentLinks(dnr.node)
}
