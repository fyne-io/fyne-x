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

const (
	// default inner size
	defaultWidth  float32 = 50
	defaultHeight float32 = 50

	// default padding around the inner object in a node
	defaultPadding float32 = 10
)

type diagramNodeRenderer struct {
	node   *DiagramNode
	handle *canvas.Line
	box    *canvas.Rectangle
}

// DiagramNode represents a node in the diagram widget. It contains an inner
// widget, and also draws a border, and a "handle" that can be used to drag it
// around.
type DiagramNode struct {
	widget.BaseWidget

	Diagram *DiagramWidget

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
	// drawn on top of this. Defaults to the theme.BackgroundColor().
	BoxFillColor color.Color

	// BoxStrokeColor is the stroke color of the node rectangle. Defaults
	// to theme.TextColor().
	BoxStrokeColor color.Color

	// HandleColor is the color of node handle.
	HandleColor color.Color

	// HandleStrokeWidth is the stroke width of the node handle, defaults
	// to 3.
	HandleStroke float32
}

func (r *diagramNodeRenderer) MinSize() fyne.Size {
	// space for the inner widget, plus padding on all sides.
	inner := r.node.effectiveInnerSize()
	return fyne.Size{
		Width:  inner.Width + float32(2*r.node.Padding),
		Height: inner.Height + float32(2*r.node.Padding),
	}
}

func (r *diagramNodeRenderer) Layout(size fyne.Size) {
	r.node.Resize(r.MinSize())

	r.node.InnerObject.Move(r.node.innerPos())
	r.node.InnerObject.Resize(r.node.effectiveInnerSize())

	r.box.Resize(r.MinSize())

	canvas.Refresh(r.node.InnerObject)
}

func (r *diagramNodeRenderer) ApplyTheme(size fyne.Size) {
}

func (r *diagramNodeRenderer) Refresh() {
	// move and size the inner object appropriately
	r.node.InnerObject.Move(r.node.innerPos())
	r.node.InnerObject.Resize(r.node.effectiveInnerSize())

	// move the box and update it's colors
	r.box.StrokeWidth = r.node.BoxStrokeWidth
	r.box.FillColor = r.node.BoxFillColor
	r.box.StrokeColor = r.node.BoxStrokeColor
	r.box.Resize(r.MinSize())

	// calculate the handle positions
	r.handle.Position1 = fyne.Position{
		X: float32(r.node.Padding),
		Y: float32(r.node.Padding / 2),
	}

	r.handle.Position2 = fyne.Position{
		X: r.node.effectiveInnerSize().Width + float32(r.node.Padding),
		Y: float32(r.node.Padding / 2),
	}

	r.handle.StrokeWidth = r.node.HandleStroke
	r.handle.StrokeColor = r.node.HandleColor

	for _, e := range r.node.Diagram.GetEdges(r.node) {
		e.Refresh()
	}

	canvas.Refresh(r.box)
	canvas.Refresh(r.handle)
	canvas.Refresh(r.node.InnerObject)
}

func (r *diagramNodeRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *diagramNodeRenderer) Destroy() {
}

func (r *diagramNodeRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.box,
		r.handle,
		r.node.InnerObject,
	}
}

func (n *DiagramNode) CreateRenderer() fyne.WidgetRenderer {
	r := diagramNodeRenderer{
		node:   n,
		handle: canvas.NewLine(n.HandleColor),
		box:    canvas.NewRectangle(n.BoxStrokeColor),
	}

	r.handle.StrokeWidth = n.HandleStroke
	r.box.StrokeWidth = n.BoxStrokeWidth
	r.box.FillColor = n.BoxFillColor

	(&r).Refresh()

	return &r
}

func NewDiagramNode(d *DiagramWidget, obj fyne.CanvasObject) *DiagramNode {
	w := &DiagramNode{
		Diagram:        d,
		InnerSize:      fyne.Size{Width: defaultWidth, Height: defaultHeight},
		InnerObject:    obj,
		Padding:        defaultPadding,
		BoxStrokeWidth: 1,
		BoxFillColor:   theme.BackgroundColor(),
		BoxStrokeColor: theme.TextColor(),
		HandleColor:    theme.TextColor(),
		HandleStroke:   3,
	}

	w.ExtendBaseWidget(w)

	return w
}

func (n *DiagramNode) innerPos() fyne.Position {
	return fyne.Position{
		X: float32(n.Padding),
		Y: float32(n.Padding),
	}
}

func (n *DiagramNode) effectiveInnerSize() fyne.Size {
	return n.InnerSize.Max(n.InnerObject.MinSize())
}

func (n *DiagramNode) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (n *DiagramNode) DragEnd() {
	n.Refresh()
}

func (n *DiagramNode) Dragged(event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	n.Displace(delta)
	n.Refresh()
}

func (n *DiagramNode) MouseIn(event *desktop.MouseEvent) {
	n.HandleColor = theme.FocusColor()
	n.Refresh()
}

func (n *DiagramNode) MouseOut() {
	n.HandleColor = theme.TextColor()
	n.Refresh()
}

func (n *DiagramNode) MouseMoved(event *desktop.MouseEvent) {
}

func (n *DiagramNode) Displace(delta fyne.Position) {
	n.Move(n.Position().Add(delta))
}

func (n *DiagramNode) R2Position() r2.Vec2 {
	return r2.V2(float64(n.Position().X), float64(n.Position().Y))
}

func (n *DiagramNode) R2Box() r2.Box {
	inner := n.effectiveInnerSize()
	s := r2.V2(
		float64(inner.Width+float32(2*n.Padding)),
		float64(inner.Height+float32(2*n.Padding)),
	)

	return r2.MakeBox(n.R2Position(), s)
}

func (n *DiagramNode) R2Center() r2.Vec2 {
	return n.R2Box().Center()
}

func (n *DiagramNode) Center() fyne.Position {
	return fyne.Position{X: float32(n.R2Center().X), Y: float32(n.R2Center().Y)}
}
