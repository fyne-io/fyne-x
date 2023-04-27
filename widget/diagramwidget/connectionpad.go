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

// ConnectionPad is an interface to a connection shape on a DiagramElement.
type ConnectionPad interface {
	fyne.Widget
	desktop.Hoverable
	GetPadOwner() DiagramElement
	GetCenter() fyne.Position
	getConnectionPoint(referencePoint fyne.Position) fyne.Position
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

// PointPad is a ConnectionPad consisting of a single point (the Position of the PointPad)
type PointPad struct {
	widget.BaseWidget
	connectionPad
}

// NewPointPad creates a PointPad and associates it with the DiagramElement. Note that, by default,
// the position of the PointPad will be (0,0), i.e. the origin of the DiagramElement.
func NewPointPad(padOwner DiagramElement) *PointPad {
	pp := &PointPad{}
	pp.connectionPad.padOwner = padOwner
	pp.BaseWidget.ExtendBaseWidget(pp)
	return pp
}

// CreateRenderer creates the WidgetRenderer for a PointPad
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

// getConnectionPoint returns the point on the pad to which a connection will be made from the referencePoint.
// For a point pad, this is always the center.
func (pp *PointPad) getConnectionPoint(referencePoint fyne.Position) fyne.Position {
	return pp.GetCenter()
}

// MouseIn responds to mouse movements within the pointPadSize distance of the center
func (pp *PointPad) MouseIn(event *desktop.MouseEvent) {
	// TODO implement this
}

// MouseMoved responds to mouse movements within the pointPadSize distance of the center
func (pp *PointPad) MouseMoved(event *desktop.MouseEvent) {
	// TODO implement this
}

// MouseOut responds to mouse movements within the pointPadSize distance of the center
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

// RectanglePad provides a ConnectionPad corresponding to the perimeter of the DiagramElement owning the pad.
type RectanglePad struct {
	widget.BaseWidget
	connectionPad
}

// NewRectanglePad creates a RectanglePad and associates it with the DiagramElement. The size of the
// pad becomes the size of the padOwner.
func NewRectanglePad(padOwner DiagramElement) *RectanglePad {
	rp := &RectanglePad{}
	rp.connectionPad.padOwner = padOwner
	rp.BaseWidget.ExtendBaseWidget(rp)
	return rp
}

// CreateRenderer creates the WidgetRenderer for the RectanglePad
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

// getConnectionPoint returns the point at which the connection should be made from a reference point.
// The reference point is in diagram coordinates and the returned point is also in diagram coordinates.
// For a RectanglePad this point is the intersection of a line segment from the reference point to the center
// of the rectangle pad and the rectangle bounding the pad.
func (rp *RectanglePad) getConnectionPoint(referencePoint fyne.Position) fyne.Position {
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

// MouseIn responds to the mouse entering the bounds of the RectanglePad
func (rp *RectanglePad) MouseIn(event *desktop.MouseEvent) {
	// TODO implement this
}

// MouseMoved responds to mouse movements within the rectangle pad
func (rp *RectanglePad) MouseMoved(event *desktop.MouseEvent) {
	// TODO implement this
}

// MouseOut responds to mouse movements leaving the rectangle pad
func (rp *RectanglePad) MouseOut() {
	// TODO implement this
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
	rpr.rp.connectionPad.padOwner.GetDiagram().forceRepaint()
}
