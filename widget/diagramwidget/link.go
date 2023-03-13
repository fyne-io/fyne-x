package diagramwidget

import (
	"image/color"

	"fyne.io/x/fyne/widget/diagramwidget/arrowhead"
	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type diagramLinkRenderer struct {
	edge  *DiagramLink
	line  *canvas.Line
	arrow *arrowhead.Arrowhead
}

type DiagramLink struct {
	widget.BaseWidget

	Diagram *DiagramWidget

	LinkColor color.Color

	Width float32

	Origin *DiagramNode
	Target *DiagramNode

	Directed bool
}

func (r *diagramLinkRenderer) MinSize() fyne.Size {
	xdelta := r.edge.Origin.Position().X - r.edge.Target.Position().X
	if xdelta < 0 {
		xdelta *= -1
	}

	ydelta := r.edge.Origin.Position().Y - r.edge.Target.Position().Y
	if ydelta < 0 {
		ydelta *= -1
	}

	return fyne.Size{Width: xdelta, Height: ydelta}
}

func (r *diagramLinkRenderer) Layout(size fyne.Size) {
}

func (r *diagramLinkRenderer) ApplyTheme(size fyne.Size) {
}

func (r *diagramLinkRenderer) Refresh() {
	l := r.edge.R2Line()
	b1 := r.edge.Origin.R2Box()
	b2 := r.edge.Target.R2Box()

	p1, _ := b1.Intersect(l)
	p2, _ := b2.Intersect(l)

	r.line.Position1 = fyne.Position{
		X: float32(p1.X),
		Y: float32(p1.Y),
	}

	r.line.Position2 = fyne.Position{
		X: float32(p2.X),
		Y: float32(p2.Y),
	}

	r.line.StrokeColor = r.edge.LinkColor
	r.line.StrokeWidth = r.edge.Width

	if r.edge.Directed {
		r.arrow.Show()
		r.arrow.Tip = r.line.Position1
		r.arrow.Base = r.line.Position2
		r.arrow.StrokeColor = r.edge.LinkColor
		r.arrow.StrokeWidth = r.edge.Width
	} else {
		r.arrow.Hide()
	}

	canvas.Refresh(r.line)
	canvas.Refresh(r.arrow)
	r.arrow.Refresh()
}

func (r *diagramLinkRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *diagramLinkRenderer) Destroy() {
}

func (r *diagramLinkRenderer) Objects() []fyne.CanvasObject {
	// obj := []fyne.CanvasObject{
	//         r.line,
	//         // r.arrow,
	// }

	// XXX: temporary hack because otherwise I can't get my canvas object
	// to show up???
	obj := r.arrow.Objects()
	obj = append(obj, r.line)
	return obj
}

func (e *DiagramLink) CreateRenderer() fyne.WidgetRenderer {
	r := diagramLinkRenderer{
		edge:  e,
		line:  canvas.NewLine(e.LinkColor),
		arrow: arrowhead.MakeArrowhead(fyne.Position{X: 0, Y: 0}, fyne.Position{X: 0, Y: 0}),
	}

	(&r).Refresh()

	return &r
}

func (e *DiagramLink) R2Line() r2.Line {
	return r2.MakeLineFromEndpoints(e.Origin.R2Center(), e.Target.R2Center())
}

func NewDiagramEdge(g *DiagramWidget, v, u *DiagramNode) *DiagramLink {
	e := &DiagramLink{
		Diagram:   g,
		LinkColor: theme.TextColor(),
		Width:     2,
		Origin:    v,
		Target:    u,
		Directed:  false,
	}

	e.ExtendBaseWidget(e)

	return e
}
