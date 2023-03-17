package diagramwidget

import (
	"image/color"
	"math"

	"fyne.io/x/fyne/widget/diagramwidget/decoration"
	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DiagramLink struct {
	widget.BaseWidget
	Diagram             *DiagramWidget
	LinkColor           color.Color
	Width               float32
	Origin              *DiagramNode
	Target              *DiagramNode
	SourceDecorations   []decoration.Decoration
	TargetDecorations   []decoration.Decoration
	MidpointDecorations []decoration.Decoration
	// Directed           bool
}

func NewDiagramLink(g *DiagramWidget, v, u *DiagramNode) *DiagramLink {
	e := &DiagramLink{
		Diagram:   g,
		LinkColor: theme.TextColor(),
		Width:     2,
		Origin:    v,
		Target:    u,
		// Directed:  false,
	}

	e.ExtendBaseWidget(e)

	return e
}

func (e *DiagramLink) CreateRenderer() fyne.WidgetRenderer {
	r := diagramLinkRenderer{
		edge: e,
		line: canvas.NewLine(e.LinkColor),
		// arrow: arrowhead.MakeArrowhead(fyne.Position{X: 0, Y: 0}, fyne.Position{X: 0, Y: 0}),
	}

	(&r).Refresh()

	return &r
}

func (e *DiagramLink) R2Line() r2.Line {
	return r2.MakeLineFromEndpoints(e.Origin.R2Center(), e.Target.R2Center())
}

type diagramLinkRenderer struct {
	edge *DiagramLink
	line *canvas.Line
	// arrow *arrowhead.Arrowhead
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
	l := r.edge.R2Line()
	sourceBox := r.edge.Origin.R2Box()
	targetBox := r.edge.Target.R2Box()
	sourcePoint, _ := sourceBox.Intersect(l)
	targetPoint, _ := targetBox.Intersect(l)
	r.line.Position1 = fyne.Position{
		X: float32(sourcePoint.X),
		Y: float32(sourcePoint.Y),
	}
	r.line.Position2 = fyne.Position{
		X: float32(targetPoint.X),
		Y: float32(targetPoint.Y),
	}
	r.line.StrokeColor = r.edge.LinkColor
	r.line.StrokeWidth = r.edge.Width
	canvas.Refresh(r.line)
	// Have to change the sign of Y since the window inverts the Y axis
	lineVector := r2.Vec2{X: float64(r.line.Position2.X - r.line.Position1.X), Y: -float64(r.line.Position2.Y - r.line.Position1.Y)}
	sourceAngle := lineVector.Angle()
	targetAngle := r2.AddAngles(sourceAngle, math.Pi)
	sourceOffset := 0.0
	for _, decoration := range r.edge.SourceDecorations {
		decorationReferencePoint := fyne.Position{
			X: float32(float64(r.line.Position1.X) + math.Cos(sourceAngle)*sourceOffset),
			Y: float32(float64(r.line.Position1.Y) - math.Sin(sourceAngle)*sourceOffset),
		}
		decoration.SetReferencePoint(decorationReferencePoint)
		decoration.SetReferenceAngle(sourceAngle)
		sourceOffset = sourceOffset + float64(decoration.GetReferenceLength())
	}
	midPosition := fyne.Position{
		X: float32((sourcePoint.X + targetPoint.X) / 2),
		Y: float32((sourcePoint.Y + targetPoint.Y) / 2),
	}
	midOffset := 0.0
	for _, decoration := range r.edge.MidpointDecorations {
		decorationReferencePoint := fyne.Position{
			X: float32(float64(midPosition.X) + math.Cos(targetAngle)*midOffset),
			Y: float32(float64(midPosition.Y) - math.Sin(targetAngle)*midOffset),
		}
		decoration.SetReferencePoint(decorationReferencePoint)
		decoration.SetReferenceAngle(targetAngle)
		midOffset = midOffset + float64(decoration.GetReferenceLength())
	}
	targetOffset := 0.0
	for _, decoration := range r.edge.TargetDecorations {
		decorationReferencePoint := fyne.Position{
			X: float32(float64(r.line.Position2.X) + math.Cos(targetAngle)*targetOffset),
			Y: float32(float64(r.line.Position2.Y) - math.Sin(targetAngle)*targetOffset),
		}
		decoration.SetReferencePoint(decorationReferencePoint)
		decoration.SetReferenceAngle(targetAngle)
		targetOffset = targetOffset + float64(decoration.GetReferenceLength())
	}
}

func (r *diagramLinkRenderer) ApplyTheme(size fyne.Size) {
}

func (r *diagramLinkRenderer) Refresh() {
	for _, decoration := range r.edge.SourceDecorations {
		decoration.SetStrokeColor(r.edge.LinkColor)
		decoration.SetStrokeWidth(r.edge.Width)
		decoration.Refresh()
	}
	for _, decoration := range r.edge.MidpointDecorations {
		decoration.SetStrokeColor(r.edge.LinkColor)
		decoration.SetStrokeWidth(r.edge.Width)
		decoration.Refresh()
	}
	for _, decoration := range r.edge.TargetDecorations {
		decoration.SetStrokeColor(r.edge.LinkColor)
		decoration.SetStrokeWidth(r.edge.Width)
		decoration.Refresh()
	}
}

func (r *diagramLinkRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *diagramLinkRenderer) Destroy() {
}

func (r *diagramLinkRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{
		r.line,
	}
	return obj
}
