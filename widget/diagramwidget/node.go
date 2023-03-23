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

// Validate that DiagramNode is a DiagramElement
var _ DiagramElement = (*DiagramNode)(nil)

const (
	// default inner size
	defaultWidth  float32 = 50
	defaultHeight float32 = 25

	// default padding around the inner object in a node
	defaultPadding float32 = 10
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
	// InnerObject is the canvas object that should be drawn inside of
	// the diagram node.
	InnerObject fyne.CanvasObject
	// Padding is the distance between the inner object's drawing area
	// and the box.
	Padding float32
	// BoxStrokeWidth is the stroke width of the box which delineates the
	// node. Defaults to 1.
	BoxStrokeWidth float32
	// BoxFill is the fill color of the node, the inner object will be
	// drawn on top of this. Defaults to the DiagramTheme's BackgroundColor.
	BoxFillColor color.Color
	// BoxStrokeColor is the stroke color of the node rectangle. Defaults
	// to DiagramTheme's ForegroundColor
	BoxStrokeColor color.Color
	// HandleColor is the color of node handle.
	HandleColor color.Color
	// HandleStrokeWidth is the stroke width of the node handle, defaults
	// to 3.
	HandleStroke float32
	handles      map[string]*Handle
}

func NewDiagramNode(diagram *DiagramWidget, obj fyne.CanvasObject) *DiagramNode {
	dn := &DiagramNode{
		InnerSize:      fyne.Size{Width: defaultWidth, Height: defaultHeight},
		InnerObject:    obj,
		Padding:        diagram.DiagramTheme.Size(theme.SizeNamePadding),
		BoxStrokeWidth: 1,
		BoxFillColor:   diagram.DiagramTheme.Color(theme.ColorNameBackground, diagram.ThemeVariant),
		BoxStrokeColor: diagram.DiagramTheme.Color(theme.ColorNameForeground, diagram.ThemeVariant),
		HandleColor:    diagram.DiagramTheme.Color(theme.ColorNameForeground, diagram.ThemeVariant),
		HandleStroke:   3,
		handles:        make(map[string]*Handle),
	}
	dn.diagramElement.diagram = diagram
	for _, handleKey := range []string{"upperLeft", "upperMiddle", "upperRight", "leftMiddle", "rightMiddle", "lowerLeft", "lowerMiddle", "lowerRight"} {
		newHandle := NewHandle(dn)
		dn.handles[handleKey] = newHandle
	}
	dn.ExtendBaseWidget(dn)
	return dn
}

func (dn *DiagramNode) CreateRenderer() fyne.WidgetRenderer {
	dnr := diagramNodeRenderer{
		node: dn,
		box:  canvas.NewRectangle(dn.BoxStrokeColor),
	}

	dnr.box.StrokeWidth = dn.BoxStrokeWidth
	dnr.box.FillColor = dn.BoxFillColor

	(&dnr).Refresh()

	return &dnr
}

func (dn *DiagramNode) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (dn *DiagramNode) Displace(delta fyne.Position) {
	dn.Move(dn.Position().Add(delta))
}

func (dn *DiagramNode) DragEnd() {
}

func (dn *DiagramNode) Dragged(event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	dn.Displace(delta)
	ForceRefresh()
}

func (dn *DiagramNode) effectiveInnerSize() fyne.Size {
	return dn.InnerSize.Max(dn.InnerObject.MinSize())
}

func (dn *DiagramNode) findKeyForHandle(handle *Handle) string {
	for k, v := range dn.handles {
		if v == handle {
			return k
		}
	}
	return ""
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
	dn.InnerSize = dn.InnerObject.MinSize().Max(trialInnerSize)
	dn.Resize(dn.Size().Add(sizeChange))
	dn.Refresh()
	ForceRefresh()
}

func (dn *DiagramNode) innerPos() fyne.Position {
	return fyne.Position{
		X: float32(dn.Padding),
		Y: float32(dn.Padding),
	}
}

func (dn *DiagramNode) MouseIn(event *desktop.MouseEvent) {
	dn.HandleColor = dn.diagram.DiagramTheme.Color(theme.ColorNameFocus, dn.diagram.ThemeVariant)
	ForceRefresh()
}

func (dn *DiagramNode) MouseOut() {
	dn.HandleColor = dn.diagram.DiagramTheme.Color(theme.ColorNameForeground, dn.diagram.ThemeVariant)
	ForceRefresh()
}

func (dn *DiagramNode) MouseMoved(event *desktop.MouseEvent) {
}

func (dn *DiagramNode) R2Position() r2.Vec2 {
	return r2.V2(float64(dn.Position().X), float64(dn.Position().Y))
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

func (dn *DiagramNode) Center() fyne.Position {
	return fyne.Position{X: float32(dn.R2Center().X), Y: float32(dn.R2Center().Y)}
}

// diagramNodeRenderer
type diagramNodeRenderer struct {
	node *DiagramNode
	box  *canvas.Rectangle
}

func (r *diagramNodeRenderer) ApplyTheme(size fyne.Size) {
}

func (r *diagramNodeRenderer) BackgroundColor() color.Color {
	return r.node.diagram.DiagramTheme.Color(theme.ColorNameBackground, r.node.diagram.ThemeVariant)
}

func (r *diagramNodeRenderer) Destroy() {
}

func (r *diagramNodeRenderer) MinSize() fyne.Size {
	// space for the inner widget, plus padding on all sides.
	inner := r.node.effectiveInnerSize()
	return fyne.Size{
		Width:  inner.Width + float32(2*r.node.Padding),
		Height: inner.Height + float32(2*r.node.Padding),
	}
}

func (dnr *diagramNodeRenderer) Layout(size fyne.Size) {
	nodeSize := dnr.MinSize().Max(size)
	dnr.node.Resize(nodeSize)

	dnr.node.InnerObject.Move(dnr.node.innerPos())
	dnr.node.InnerObject.Resize(dnr.node.effectiveInnerSize())

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
}

func (dnr *diagramNodeRenderer) Objects() []fyne.CanvasObject {
	obj := make([]fyne.CanvasObject, 0)
	obj = append(obj, dnr.box)
	obj = append(obj, dnr.node.InnerObject)
	for _, handle := range dnr.node.handles {
		obj = append(obj, handle)
	}
	return obj
}

func (r *diagramNodeRenderer) Refresh() {
	r.box.StrokeWidth = r.node.BoxStrokeWidth
	r.box.FillColor = r.node.BoxFillColor
	r.box.StrokeColor = r.node.BoxStrokeColor
	for _, e := range r.node.diagram.GetEdges(r.node) {
		e.Refresh()
	}
}
