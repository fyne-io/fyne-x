package diagramwidget

import (
	"image/color"

	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	pointPadSize float32 = 10
	padLineWidth float32 = 3
)

// ConnectionPad is an interface to a connection area on a DiagramElement.
type ConnectionPad interface {
	fyne.Widget
	desktop.Hoverable
	GetPadOwner() DiagramElement
	GetCenter() fyne.Position
	GetConnectionPoint(referencePoint fyne.Position) fyne.Position
}

type connectionPad struct {
	padOwner DiagramElement
}

func (cp *connectionPad) GetPadOwner() DiagramElement {
	return cp.padOwner
}

/******************************
	PointPad
*******************************/

// Validate that PointPad implements ConnectionPad
var _ ConnectionPad = (*PointPad)(nil)

type PointPad struct {
	widget.BaseWidget
	connectionPad
}

func NewPointPad(padOwner DiagramElement) *PointPad {
	pp := &PointPad{}
	pp.connectionPad.padOwner = padOwner
	pp.BaseWidget.ExtendBaseWidget(pp)
	return pp
}

func (pp *PointPad) CreateRenderer() fyne.WidgetRenderer {
	ppr := &pointPadRenderer{
		pp: pp,
		l1: canvas.NewLine(pp.padOwner.GetDiagram().GetHoverColor()),
		l2: canvas.NewLine(pp.padOwner.GetDiagram().GetHoverColor()),
	}
	ppr.l1.StrokeWidth = padLineWidth
	ppr.l2.StrokeWidth = padLineWidth
	return ppr
}

// GetCenter returns the position in diagram coordinates
func (pp *PointPad) GetCenter() fyne.Position {
	return pp.padOwner.Position().Add(pp.Position())
}

func (pp *PointPad) GetConnectionPoint(referencePoint fyne.Position) fyne.Position {
	return pp.GetCenter()
}

func (pp *PointPad) MouseIn(event *desktop.MouseEvent) {

}

func (pp *PointPad) MouseMoved(event *desktop.MouseEvent) {

}

func (pp *PointPad) MouseOut() {

}

// pointPadRenderer
type pointPadRenderer struct {
	pp *PointPad
	l1 *canvas.Line
	l2 *canvas.Line
}

func (ppr *pointPadRenderer) Destroy() {

}

func (ppr *pointPadRenderer) Layout(size fyne.Size) {
	ppr.l1.Position1 = fyne.NewPos(0, 0)
	ppr.l1.Position2 = fyne.NewPos(pointPadSize, pointPadSize)
	ppr.l2.Position1 = fyne.NewPos(pointPadSize, 0)
	ppr.l2.Position2 = fyne.NewPos(0, pointPadSize)
}

func (ppr *pointPadRenderer) MinSize() fyne.Size {
	return fyne.Size{Height: pointPadSize, Width: pointPadSize}
}

func (ppr *pointPadRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{
		ppr.l1,
		ppr.l2,
	}
	return obj
}

func (ppr *pointPadRenderer) Refresh() {
	ppr.l1.StrokeColor = ppr.pp.padOwner.GetDiagram().GetHoverColor()
	ppr.l1.StrokeWidth = padLineWidth
	ppr.l2.StrokeColor = ppr.pp.padOwner.GetDiagram().GetHoverColor()
	ppr.l2.StrokeWidth = padLineWidth
}

/***********************************
	RectanglePad
*************************************/

// Validate that RectanglePad implements ConnectionPad
var _ ConnectionPad = (*RectanglePad)(nil)

type RectanglePad struct {
	widget.BaseWidget
	connectionPad
	// box is the pad shape in diagram coordinates
	// box r2.Box
}

func NewRectanglePad(padOwner DiagramElement) *RectanglePad {
	rp := &RectanglePad{}
	rp.connectionPad.padOwner = padOwner
	rp.BaseWidget.ExtendBaseWidget(rp)
	return rp
}

func (rp *RectanglePad) CreateRenderer() fyne.WidgetRenderer {
	rpr := &rectanglePadRenderer{
		rp:   rp,
		rect: *canvas.NewRectangle(rp.padOwner.GetDiagram().GetForegroundColor()),
	}
	// rpr.rect.FillColor = color.Transparent
	rpr.rect.StrokeWidth = padLineWidth
	return rpr
}

// GetCenter() returns the center of the pad in the diagram's coordinate system
func (rp *RectanglePad) GetCenter() fyne.Position {
	box := rp.makeBox()
	r2Center := box.Center()
	return fyne.NewPos(float32(r2Center.X), float32(r2Center.Y))
}

func (rp *RectanglePad) GetConnectionPoint(referencePoint fyne.Position) fyne.Position {
	box := rp.makeBox()
	r2ReferencePoint := r2.MakeVec2(float64(referencePoint.X), float64(referencePoint.Y))
	linkLine := r2.MakeLineFromEndpoints(box.Center(), r2ReferencePoint)
	r2Intersection, _ := box.Intersect(linkLine)
	return fyne.NewPos(float32(r2Intersection.X), float32(r2Intersection.Y))
}

// makeBox returns an r2 box representing the rectangle pad's position and size in the
// diagram's coorinate system
func (rp *RectanglePad) makeBox() r2.Box {
	diagramCoordinatePosition := rp.padOwner.Position().Add(rp.Position())
	r2Position := r2.V2(float64(diagramCoordinatePosition.X), float64(diagramCoordinatePosition.Y))
	s := r2.V2(
		float64(rp.Size().Width),
		float64(rp.Size().Height),
	)
	return r2.MakeBox(r2Position, s)
}

func (rp *RectanglePad) MouseIn(event *desktop.MouseEvent) {

}

func (rp *RectanglePad) MouseMoved(event *desktop.MouseEvent) {

}

func (rp *RectanglePad) MouseOut() {

}

// rectanglePadRenderer
type rectanglePadRenderer struct {
	rp   *RectanglePad
	rect canvas.Rectangle
}

func (rpr *rectanglePadRenderer) Destroy() {

}

func (rpr *rectanglePadRenderer) Layout(size fyne.Size) {
	padOwnerSize := rpr.rp.padOwner.Size()
	rpr.rp.Resize(padOwnerSize)
	rpr.rect.Resize(padOwnerSize)
}

func (rpr *rectanglePadRenderer) MinSize() fyne.Size {
	return rpr.rp.padOwner.Size()
}

func (rpr *rectanglePadRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{
		&rpr.rect,
	}
	return obj
}

func (rpr *rectanglePadRenderer) Refresh() {
	rpr.rect.StrokeColor = rpr.rp.padOwner.GetDiagram().GetForegroundColor()
	rpr.rect.FillColor = color.Transparent
	rpr.rect.StrokeWidth = padLineWidth
	ForceRepaint()
}
