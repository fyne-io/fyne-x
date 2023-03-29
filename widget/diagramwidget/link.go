package diagramwidget

import (
	"image/color"
	"math"

	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Validate Hoverable Implementation
var _ desktop.Hoverable = (*DiagramLink)(nil)
var _ DiagramElement = (*DiagramLink)(nil)

type DiagramLink struct {
	widget.BaseWidget
	diagramElement
	linkPoints           []*LinkPoint
	linkSegments         []*LinkSegment
	LinkColor            color.Color
	Width                float32
	midPad               *PointPad
	sourcePad            ConnectionPad
	targetPad            ConnectionPad
	SourceDecorations    []Decoration
	sourceAnchoredText   map[string]*AnchoredText
	TargetDecorations    []Decoration
	targetAnchoredText   map[string]*AnchoredText
	MidpointDecorations  []Decoration
	midpointAnchoredText map[string]*AnchoredText
}

func NewDiagramLink(diagram *DiagramWidget, sourcePad, targetPad ConnectionPad, linkID string) *DiagramLink {
	dl := &DiagramLink{
		linkPoints:           []*LinkPoint{},
		linkSegments:         []*LinkSegment{},
		LinkColor:            diagram.GetForegroundColor(),
		Width:                2,
		sourcePad:            sourcePad,
		targetPad:            targetPad,
		sourceAnchoredText:   make(map[string]*AnchoredText),
		midpointAnchoredText: make(map[string]*AnchoredText),
		targetAnchoredText:   make(map[string]*AnchoredText),
	}
	dl.diagramElement.initialize(diagram, linkID)
	dl.linkPoints = append(dl.linkPoints, NewLinkPoint(dl))
	dl.linkPoints = append(dl.linkPoints, NewLinkPoint(dl))
	dl.linkSegments = append(dl.linkSegments, NewLinkSegment(dl, dl.linkPoints[0].Position(), dl.linkPoints[1].Position()))
	dl.midPad = NewPointPad(dl)
	dl.midPad.Move(dl.GetMidPosition())
	dl.ExtendBaseWidget(dl)

	dl.diagram.AddLink(dl)
	dl.diagram.addLinkDependency(dl.sourcePad.GetPadOwner(), dl, dl.linkPoints[0])
	dl.diagram.addLinkDependency(dl.targetPad.GetPadOwner(), dl, dl.linkPoints[1])
	dl.Refresh()
	return dl
}

func (dl *DiagramLink) CreateRenderer() fyne.WidgetRenderer {
	dlr := diagramLinkRenderer{
		link: dl,
	}

	(&dlr).Refresh()

	return &dlr
}

func (dl *DiagramLink) AddSourceAnchoredText(key string, displayedText string) {
	at := NewAnchoredText(displayedText)
	dl.sourceAnchoredText[key] = at
	at.SetReferencePosition(dl.GetSourcePosition())
	at.Move(dl.GetSourcePosition())
	dl.Refresh()
}

func (dl *DiagramLink) AddSourceDecoration(decoration Decoration) {
	decoration.setLink(dl)
	dl.SourceDecorations = append(dl.SourceDecorations, decoration)
	dl.Refresh()
}

func (dl *DiagramLink) AddMidpointAnchoredText(key string, displayedText string) {
	at := NewAnchoredText(displayedText)
	dl.midpointAnchoredText[key] = at
	at.SetReferencePosition(dl.GetMidPosition())
	at.Move(dl.GetMidPosition())
	dl.Refresh()
}

func (dl *DiagramLink) AddMidpointDecoration(decoration Decoration) {
	decoration.setLink(dl)
	dl.MidpointDecorations = append(dl.MidpointDecorations, decoration)
	dl.Refresh()
}

func (dl *DiagramLink) AddTargetAnchoredText(key string, displayedText string) {
	at := NewAnchoredText(displayedText)
	dl.targetAnchoredText[key] = at
	at.SetReferencePosition(dl.GetTargetPosition())
	at.Move(dl.GetTargetPosition())
	dl.Refresh()
}

func (dl *DiagramLink) AddTargetDecoration(decoration Decoration) {
	decoration.setLink(dl)
	dl.TargetDecorations = append(dl.TargetDecorations, decoration)
	dl.Refresh()
}

func (dl *DiagramLink) GetSourcePosition() fyne.Position {
	return dl.linkPoints[0].Position()
}

func (dl *DiagramLink) GetMidPad() ConnectionPad {
	return dl.midPad
}

func (dl *DiagramLink) GetMidPosition() fyne.Position {
	// TODO update when additional points are introduced
	sourcePoint := dl.linkPoints[0].Position()
	targetPoint := dl.linkPoints[len(dl.linkPoints)-1].Position()
	midPoint := fyne.NewPos((sourcePoint.X+targetPoint.X)/2, (sourcePoint.Y+targetPoint.Y)/2)
	return midPoint
}

func (dl *DiagramLink) GetTargetPosition() fyne.Position {
	return dl.linkPoints[len(dl.linkPoints)-1].Position()
}

func (dl *DiagramLink) handleDragged(handle *Handle, event *fyne.DragEvent) {
	// TODO implement this
}

func (dl *DiagramLink) HideHandles() {

}

func (dl *DiagramLink) MouseIn(event *desktop.MouseEvent) {
}

func (dl *DiagramLink) MouseMoved(event *desktop.MouseEvent) {

}

func (dl *DiagramLink) MouseOut() {
}

func (dl *DiagramLink) ShowHandles() {

}

// diagramLinkRenderer
type diagramLinkRenderer struct {
	link *DiagramLink
}

func (dlr *diagramLinkRenderer) Destroy() {
}

func (dlr *diagramLinkRenderer) MinSize() fyne.Size {
	var xMin, xMax, yMin, yMax float32
	for i, point := range dlr.link.linkPoints {
		if i == 0 {
			xMin = point.Position().X
			xMax = point.Position().X
			yMin = point.Position().Y
			yMax = point.Position().Y
		} else {
			xMin = float32(math.Min(float64(xMin), float64(point.Position().X)))
			xMax = float32(math.Max(float64(xMax), float64(point.Position().X)))
			yMin = float32(math.Min(float64(yMin), float64(point.Position().Y)))
			yMax = float32(math.Max(float64(yMax), float64(point.Position().Y)))
		}
	}
	return fyne.Size{Width: float32(math.Abs(float64(xMax - xMin))), Height: float32(math.Abs(float64(yMax - yMin)))}
}

func (dlr *diagramLinkRenderer) Layout(size fyne.Size) {
}

func (dlr *diagramLinkRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{}
	for i := 0; i < len(dlr.link.linkSegments); i++ {
		obj = append(obj, dlr.link.linkSegments[i])
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
	obj = append(obj, dlr.link.midPad)
	return obj
}

func (dlr *diagramLinkRenderer) Refresh() {
	padBasedSourceReferencePoint := dlr.link.sourcePad.GetCenter()
	padBasedTargetReferencePoint := dlr.link.targetPad.GetCenter()
	padBasedSourcePosition := dlr.link.sourcePad.GetConnectionPoint(padBasedTargetReferencePoint)
	padBasedTargetPosition := dlr.link.targetPad.GetConnectionPoint(padBasedSourceReferencePoint)
	// The Position of the link is the upper left hand corner of a bounding box surrounding the source and target positions
	linkPosition := fyne.NewPos(float32(math.Min(float64(padBasedSourcePosition.X), float64(padBasedTargetPosition.X))),
		float32(math.Min(float64(padBasedSourcePosition.Y), float64(padBasedTargetPosition.Y))))
	dlr.link.Move(linkPosition)

	dlr.link.linkPoints[0].Move(padBasedSourcePosition.Subtract(linkPosition))
	dlr.link.linkPoints[len(dlr.link.linkPoints)-1].Move(padBasedTargetPosition.Subtract(linkPosition))
	// Position segments only after all points have been positioned
	for i := 0; i < len(dlr.link.linkPoints)-1; i++ {
		linkSegment := dlr.link.linkSegments[i]
		linkSegment.SetPoints(dlr.link.linkPoints[i].Position(), dlr.link.linkPoints[i+1].Position())
		// linkSegment.Refresh()
	}

	// Have to change the sign of Y since the window inverts the Y axis
	lineVector := r2.Vec2{X: float64(dlr.link.linkPoints[1].Position().X - dlr.link.linkPoints[0].Position().X),
		Y: -float64(dlr.link.linkPoints[1].Position().Y - dlr.link.linkPoints[0].Position().Y)}
	sourceAngle := lineVector.Angle()
	sourceOffset := 0.0
	for _, decoration := range dlr.link.SourceDecorations {
		decorationReferencePoint := fyne.Position{
			X: float32(float64(dlr.link.linkPoints[0].Position().X) + math.Cos(sourceAngle)*sourceOffset),
			Y: float32(float64(dlr.link.linkPoints[0].Position().Y) - math.Sin(sourceAngle)*sourceOffset),
		}
		decoration.Move(decorationReferencePoint)
		decoration.SetReferenceAngle(sourceAngle)
		sourceOffset = sourceOffset + float64(decoration.GetReferenceLength())
	}

	// TODO Update target angle for multiple segments
	targetAngle := r2.AddAngles(sourceAngle, math.Pi)

	midOffset := 0.0
	for _, decoration := range dlr.link.MidpointDecorations {
		decorationReferencePoint := fyne.Position{
			X: float32(float64(dlr.link.GetMidPosition().X) + math.Cos(targetAngle)*midOffset),
			Y: float32(float64(dlr.link.GetMidPosition().Y) - math.Sin(targetAngle)*midOffset),
		}
		decoration.Move(decorationReferencePoint)
		decoration.SetReferenceAngle(targetAngle)
		midOffset = midOffset + float64(decoration.GetReferenceLength())
	}
	dlr.link.midPad.Move(dlr.link.GetMidPosition())

	targetOffset := 0.0
	for _, decoration := range dlr.link.TargetDecorations {
		decorationReferencePoint := fyne.Position{
			X: float32(float64(dlr.link.linkPoints[len(dlr.link.linkPoints)-1].Position().X) + math.Cos(targetAngle)*targetOffset),
			Y: float32(float64(dlr.link.linkPoints[len(dlr.link.linkPoints)-1].Position().Y) - math.Sin(targetAngle)*targetOffset),
		}
		decoration.Move(decorationReferencePoint)
		decoration.SetReferenceAngle(targetAngle)
		targetOffset = targetOffset + float64(decoration.GetReferenceLength())
	}
	for _, anchoredText := range dlr.link.sourceAnchoredText {
		anchoredText.SetReferencePosition(dlr.link.GetSourcePosition())
	}
	for _, anchoredText := range dlr.link.midpointAnchoredText {
		anchoredText.SetReferencePosition(dlr.link.GetMidPosition())
	}
	for _, anchoredText := range dlr.link.targetAnchoredText {
		anchoredText.SetReferencePosition(dlr.link.GetTargetPosition())
	}

	// Now we take care of property changes.
	for _, linkSegment := range dlr.link.linkSegments {
		linkSegment.Refresh()
	}
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
	dlr.link.diagram.refreshDependentLinks(dlr.link)
	ForceRepaint()
}
