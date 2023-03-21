package diagramwidget

import (
	"image/color"
	"log"
	"math"

	"fyne.io/x/fyne/widget/diagramwidget/decoration"
	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var hoverable desktop.Hoverable

type DiagramLink struct {
	widget.BaseWidget
	Diagram              *DiagramWidget
	LinkColor            color.Color
	Width                float32
	Origin               *DiagramNode
	sourcePoint          r2.Vec2
	Target               *DiagramNode
	targetPoint          r2.Vec2
	midPoint             r2.Vec2
	SourceDecorations    []decoration.Decoration
	sourceAnchoredText   map[string]*AnchoredText
	TargetDecorations    []decoration.Decoration
	targetAnchoredText   map[string]*AnchoredText
	MidpointDecorations  []decoration.Decoration
	midpointAnchoredText map[string]*AnchoredText
}

func NewDiagramLink(g *DiagramWidget, v, u *DiagramNode) *DiagramLink {
	dl := &DiagramLink{
		Diagram:              g,
		LinkColor:            theme.TextColor(),
		Width:                2,
		Origin:               v,
		Target:               u,
		sourceAnchoredText:   make(map[string]*AnchoredText),
		midpointAnchoredText: make(map[string]*AnchoredText),
		targetAnchoredText:   make(map[string]*AnchoredText),
	}

	dl.ExtendBaseWidget(dl)

	hoverable = dl
	return dl
}

func (dl *DiagramLink) AddSourceAnchoredText(key string, displayedText string) {
	at := NewAnchoredText(displayedText)
	dl.sourceAnchoredText[key] = at
	at.Move(fyne.Position{X: float32(dl.sourcePoint.X), Y: float32(dl.sourcePoint.Y)})
}

func (dl *DiagramLink) AddMidpointAnchoredText(key string, displayedText string) {
	at := NewAnchoredText(displayedText)
	dl.midpointAnchoredText[key] = at
	at.Move(fyne.Position{X: float32(dl.midPoint.X), Y: float32(dl.midPoint.Y)})
}

func (dl *DiagramLink) AddTargetAnchoredText(key string, displayedText string) {
	at := NewAnchoredText(displayedText)
	dl.targetAnchoredText[key] = at
	at.Move(fyne.Position{X: float32(dl.targetPoint.X), Y: float32(dl.targetPoint.Y)})
}

func (dl *DiagramLink) CreateRenderer() fyne.WidgetRenderer {
	dlr := diagramLinkRenderer{
		link: dl,
		line: canvas.NewLine(dl.LinkColor),
	}

	(&dlr).Refresh()

	return &dlr
}

func (dl *DiagramLink) R2Line() r2.Line {
	return r2.MakeLineFromEndpoints(dl.Origin.R2Center(), dl.Target.R2Center())
}

func (dl *DiagramLink) MouseIn(event *desktop.MouseEvent) {
	log.Printf("MouseIn DiagramLink Text %p", dl)
}

func (dl *DiagramLink) MouseMoved(event *desktop.MouseEvent) {

}

func (dl *DiagramLink) MouseOut() {
	log.Printf("MouseOut DiagramLink Text %p", dl)
}

type diagramLinkRenderer struct {
	link *DiagramLink
	line *canvas.Line
}

func (dlr *diagramLinkRenderer) MinSize() fyne.Size {
	xdelta := dlr.link.Origin.Position().X - dlr.link.Target.Position().X
	if xdelta < 0 {
		xdelta *= -1
	}

	ydelta := dlr.link.Origin.Position().Y - dlr.link.Target.Position().Y
	if ydelta < 0 {
		ydelta *= -1
	}
	return fyne.Size{Width: xdelta, Height: ydelta}
}

func (dlr *diagramLinkRenderer) Layout(size fyne.Size) {
	l := dlr.link.R2Line()
	sourceBox := dlr.link.Origin.R2Box()
	targetBox := dlr.link.Target.R2Box()
	dlr.link.sourcePoint, _ = sourceBox.Intersect(l)
	dlr.link.targetPoint, _ = targetBox.Intersect(l)
	dlr.line.Position1 = fyne.Position{
		X: float32(dlr.link.sourcePoint.X),
		Y: float32(dlr.link.sourcePoint.Y),
	}
	dlr.line.Position2 = fyne.Position{
		X: float32(dlr.link.targetPoint.X),
		Y: float32(dlr.link.targetPoint.Y),
	}
	dlr.line.StrokeColor = dlr.link.LinkColor
	dlr.line.StrokeWidth = dlr.link.Width
	canvas.Refresh(dlr.line)
	// Have to change the sign of Y since the window inverts the Y axis
	lineVector := r2.Vec2{X: float64(dlr.line.Position2.X - dlr.line.Position1.X), Y: -float64(dlr.line.Position2.Y - dlr.line.Position1.Y)}
	sourceAngle := lineVector.Angle()
	targetAngle := r2.AddAngles(sourceAngle, math.Pi)
	sourceOffset := 0.0
	for _, decoration := range dlr.link.SourceDecorations {
		decorationReferencePoint := fyne.Position{
			X: float32(float64(dlr.line.Position1.X) + math.Cos(sourceAngle)*sourceOffset),
			Y: float32(float64(dlr.line.Position1.Y) - math.Sin(sourceAngle)*sourceOffset),
		}
		decoration.Move(decorationReferencePoint)
		decoration.SetReferenceAngle(sourceAngle)
		sourceOffset = sourceOffset + float64(decoration.GetReferenceLength())
	}
	dlr.link.midPoint = r2.Vec2{
		X: float64((dlr.link.sourcePoint.X + dlr.link.targetPoint.X) / 2),
		Y: float64((dlr.link.sourcePoint.Y + dlr.link.targetPoint.Y) / 2),
	}
	midOffset := 0.0
	for _, decoration := range dlr.link.MidpointDecorations {
		decorationReferencePoint := fyne.Position{
			X: float32(float64(dlr.link.midPoint.X) + math.Cos(targetAngle)*midOffset),
			Y: float32(float64(dlr.link.midPoint.Y) - math.Sin(targetAngle)*midOffset),
		}
		decoration.Move(decorationReferencePoint)
		decoration.SetReferenceAngle(targetAngle)
		midOffset = midOffset + float64(decoration.GetReferenceLength())
	}
	targetOffset := 0.0
	for _, decoration := range dlr.link.TargetDecorations {
		decorationReferencePoint := fyne.Position{
			X: float32(float64(dlr.line.Position2.X) + math.Cos(targetAngle)*targetOffset),
			Y: float32(float64(dlr.line.Position2.Y) - math.Sin(targetAngle)*targetOffset),
		}
		decoration.Move(decorationReferencePoint)
		decoration.SetReferenceAngle(targetAngle)
		targetOffset = targetOffset + float64(decoration.GetReferenceLength())
	}
	for _, anchoredText := range dlr.link.sourceAnchoredText {
		anchoredText.SetReferencePosition(fyne.Position{X: float32(dlr.link.sourcePoint.X), Y: float32(dlr.link.sourcePoint.Y)})
	}
	for _, anchoredText := range dlr.link.midpointAnchoredText {
		anchoredText.SetReferencePosition(fyne.Position{X: float32(dlr.link.midPoint.X), Y: float32(dlr.link.midPoint.Y)})
	}
	for _, anchoredText := range dlr.link.targetAnchoredText {
		anchoredText.SetReferencePosition(fyne.Position{X: float32(dlr.link.targetPoint.X), Y: float32(dlr.link.targetPoint.Y)})
	}
}

func (dlr *diagramLinkRenderer) ApplyTheme(size fyne.Size) {
}

func (dlr *diagramLinkRenderer) Refresh() {
	for _, decoration := range dlr.link.SourceDecorations {
		decoration.SetStrokeColor(dlr.link.LinkColor)
		decoration.SetStrokeWidth(dlr.link.Width)
		decoration.Refresh()
	}
	for _, decoration := range dlr.link.MidpointDecorations {
		decoration.SetStrokeColor(dlr.link.LinkColor)
		decoration.SetStrokeWidth(dlr.link.Width)
		decoration.Refresh()
	}
	for _, decoration := range dlr.link.TargetDecorations {
		decoration.SetStrokeColor(dlr.link.LinkColor)
		decoration.SetStrokeWidth(dlr.link.Width)
		decoration.Refresh()
	}
}

func (dlr *diagramLinkRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (dlr *diagramLinkRenderer) Destroy() {
}

func (dlr *diagramLinkRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{
		dlr.line,
	}
	for _, sourceDecoration := range dlr.link.SourceDecorations {
		if sourceDecoration != nil {
			obj = append(obj, sourceDecoration)
		}
	}
	for _, sourceAnchoredText := range dlr.link.sourceAnchoredText {
		obj = append(obj, sourceAnchoredText)
	}
	for _, midpointDecoration := range dlr.link.MidpointDecorations {
		if midpointDecoration != nil {
			obj = append(obj, midpointDecoration)
		}
	}
	for _, midpointAnchoredText := range dlr.link.midpointAnchoredText {
		obj = append(obj, midpointAnchoredText)
	}
	for _, targetDecoration := range dlr.link.TargetDecorations {
		if targetDecoration != nil {
			obj = append(obj, targetDecoration)
		}
	}
	for _, targetAnchoredText := range dlr.link.targetAnchoredText {
		obj = append(obj, targetAnchoredText)
	}

	return obj
}
