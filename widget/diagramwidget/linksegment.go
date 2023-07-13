package diagramwidget

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

type LinkSegment struct {
	widget.BaseWidget
	link *BaseDiagramLink
	// p1 and p2 are coordinates in the link's coordinate space
	p1                fyne.Position
	p2                fyne.Position
	mouseDownPosition fyne.Position
}

func NewLinkSegment(link *BaseDiagramLink, p1 fyne.Position, p2 fyne.Position) *LinkSegment {
	ls := &LinkSegment{
		link: link,
		p1:   p1,
		p2:   p2,
	}
	ls.BaseWidget.ExtendBaseWidget(ls)
	ls.Resize(ls.MinSize())
	return ls
}

func (ls *LinkSegment) CreateRenderer() fyne.WidgetRenderer {
	lsr := &linkSegmentRenderer{
		ls:   ls,
		line: canvas.NewLine(ls.link.GetForegroundColor()),
	}
	return lsr
}

// MouseDown behavior depends upon the mouse event. If it is the primary button, it records the locateion of the
// MouseDown in preparation for a MouseUp at the same location, which will trigger Tapped() behavior. Otherwise, if
// it is the seconday button and a callback is present, it will invoke the callback
func (ls *LinkSegment) MouseDown(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonPrimary {
		ls.mouseDownPosition = event.Position
	} else if event.Button == desktop.MouseButtonSecondary && ls.link.diagram.LinkSegmentMouseDownSecondaryCallback != nil {
		ls.link.diagram.LinkSegmentMouseDownSecondaryCallback(ls.link.typedLink, event)
	}
}

// MouseUp behavior depends on the mouse event. If it is the primary button and it is at the same location as the MouseDown,
// the Tapped() behavior is invoked. Otherwise, if there is a callback present, the callback is invoked.
func (ls *LinkSegment) MouseUp(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonPrimary && ls.mouseDownPosition == event.Position {
		clickPoint := geom.Coord{float64(event.Position.X), float64(event.Position.Y)}
		p1 := geom.Coord{float64(ls.p1.X), float64(ls.p1.Y)}
		p2 := geom.Coord{float64(ls.p2.X), float64(ls.p2.Y)}
		if xy.DistanceFromPointToLine(clickPoint, p1, p2) <= float64(ls.link.properties.StrokeWidth/2)+1 {
			ls.link.diagram.DiagramElementTapped(ls.link)
		}
	} else if ls.link.diagram.LinkSegmentMouseUpCallback != nil {
		ls.link.diagram.LinkSegmentMouseUpCallback(ls.link.typedLink, event)
	}
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
	minX := math.Min(float64(lsr.ls.p1.X), float64(lsr.ls.p2.X))
	minY := math.Min(float64(lsr.ls.p1.Y), float64(lsr.ls.p2.Y))
	widgetPosition := fyne.NewPos(float32(minX), float32(minY))
	lsr.ls.Move(widgetPosition)
	lsr.ls.Resize(lsr.MinSize())
	lsr.line.Position1 = lsr.ls.p1.AddXY(-widgetPosition.X, -widgetPosition.Y)
	lsr.line.Position2 = lsr.ls.p2.AddXY(-widgetPosition.X, -widgetPosition.Y)
	lsr.line.StrokeColor = lsr.ls.link.properties.ForegroundColor
	lsr.line.StrokeWidth = lsr.ls.link.properties.StrokeWidth
}
