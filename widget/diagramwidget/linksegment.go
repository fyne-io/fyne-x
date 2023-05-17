package diagramwidget

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type LinkSegment struct {
	widget.BaseWidget
	link *BaseDiagramLink
	p1   fyne.Position
	p2   fyne.Position
}

func NewLinkSegment(link *BaseDiagramLink, p1 fyne.Position, p2 fyne.Position) *LinkSegment {
	ls := &LinkSegment{
		link: link,
		p1:   p1,
		p2:   p2,
	}
	ls.BaseWidget.ExtendBaseWidget(ls)
	return ls
}

func (ls *LinkSegment) CreateRenderer() fyne.WidgetRenderer {
	lsr := &linkSegmentRenderer{
		ls:   ls,
		line: canvas.NewLine(ls.link.diagram.GetForegroundColor()),
	}
	return lsr
}

func (ls *LinkSegment) SetPoints(p1 fyne.Position, p2 fyne.Position) {
	ls.p1 = p1
	ls.p2 = p2
	ls.Refresh()
}

// linkSegmentRenderer
type linkSegmentRenderer struct {
	ls   *LinkSegment
	line *canvas.Line
}

func (lsr *linkSegmentRenderer) Destroy() {

}

func (lsr *linkSegmentRenderer) Layout(size fyne.Size) {
}

func (lsr *linkSegmentRenderer) MinSize() fyne.Size {
	return fyne.NewSize(float32(math.Abs(float64(lsr.ls.p1.X-lsr.ls.p2.X))), float32(math.Abs(float64(lsr.ls.p1.Y-lsr.ls.p2.Y))))
}

func (lsr *linkSegmentRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{
		lsr.line,
	}
	return obj
}

func (lsr *linkSegmentRenderer) Refresh() {
	lsr.line.Position1 = lsr.ls.p1
	lsr.line.Position2 = lsr.ls.p2
	lsr.line.StrokeColor = lsr.ls.link.LinkColor
	lsr.line.StrokeWidth = lsr.ls.link.strokeWidth
	lsr.ls.link.GetDiagram().ForceRepaint()
}
