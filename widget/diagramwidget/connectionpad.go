package diagramwidget

import (
	"image/color"

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
	rpr.rp.Resize(rpr.rp.padOwner.Size())
	rpr.rect.Resize(rpr.rp.padOwner.Size())
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
	ForceRefresh()
}
