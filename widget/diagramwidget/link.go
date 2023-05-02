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

// DiagramLink is a directed graphic connection between two DiagramElements that are referred to as the Source
// and Target. The link consists of one or more line segments. By default a single line segment connects the
// Source and Target. The Link connects to ConnectionPads on the DiagramElements.
// There are three key points on a Link: the Source connection point, the Target connection point, and a MidPoint.
// For a single segment, the MidPoint is the middle of the segment. When there are two or more segments, the
// MidPoint is the source end of the next-to-last segment.
// Graphic Decoration widgets may be added at each of these points. Multiple decorations may be added at each point
// Multiple decorations are "stacked" along the line in the order added. These graphic decorations rotate with their
// associated line segments to maintain their orientation with respect to the line segment.
// Textual AnchoredText widgets may be added at each of the key points. These may be moved with respect to their associated
// key points. When the key point is moved, associated anchored texts move by the same amount, maintaining the existing offset
// between the anchored text and the key point. Multiple AnchoredText widgets may be associated with each key point.
// AnchoredText widgets are indexed by string key values provided at the time the AnchoredText is added. These key values
// can be used to retrieve the AnchoredText widget so that the displayed text value (among other things) can be set programatically.
// By default, there is a single ConnectionPad (a PointPad) associated with a Link and located at the MidPoint. Thus a
// Link can connect to another Link using this ConnectionPad.
type DiagramLink struct {
	widget.BaseWidget
	diagramElement
	linkPoints           []*LinkPoint
	linkSegments         []*LinkSegment
	LinkColor            color.Color
	strokeWidth          float32
	sourcePad            ConnectionPad
	midPad               *PointPad
	targetPad            ConnectionPad
	SourceDecorations    []Decoration
	sourceAnchoredText   map[string]*AnchoredText
	TargetDecorations    []Decoration
	targetAnchoredText   map[string]*AnchoredText
	MidpointDecorations  []Decoration
	midpointAnchoredText map[string]*AnchoredText
	showHandles          bool
}

// NewDiagramLink creates a DiagramLink widget connecting the two indicated ConnectionPads. It adds itself to the
// DiagramWidget, indexed by the supplied LinkID. This id must be unique across all of the DiagramElements in the Diagram.
// It can be used to retrieve the DiagramLink from the Diagram. The ID is intended to be used to facilitate mapping the
// DiagramLink to the information it represents in the application. The DiagramLink uses the DiagramWidget's ForegroundColor
// as the default color for the line segments.
func NewDiagramLink(diagram *DiagramWidget, sourcePad, targetPad ConnectionPad, linkID string) *DiagramLink {
	dl := &DiagramLink{
		linkPoints:           []*LinkPoint{},
		linkSegments:         []*LinkSegment{},
		LinkColor:            diagram.GetForegroundColor(),
		strokeWidth:          2,
		sourcePad:            sourcePad,
		targetPad:            targetPad,
		sourceAnchoredText:   make(map[string]*AnchoredText),
		midpointAnchoredText: make(map[string]*AnchoredText),
		targetAnchoredText:   make(map[string]*AnchoredText),
		showHandles:          false,
	}
	dl.diagramElement.initialize(diagram, linkID)
	dl.linkPoints = append(dl.linkPoints, NewLinkPoint(dl))
	dl.linkPoints = append(dl.linkPoints, NewLinkPoint(dl))
	dl.linkSegments = append(dl.linkSegments, NewLinkSegment(dl, dl.linkPoints[0].Position(), dl.linkPoints[1].Position()))
	dl.midPad = NewPointPad(dl)
	dl.midPad.Move(dl.getMidPosition())
	dl.ExtendBaseWidget(dl)

	dl.diagram.addLink(dl)
	dl.diagram.addLinkDependency(dl.sourcePad.GetPadOwner(), dl, dl.sourcePad)
	dl.diagram.addLinkDependency(dl.targetPad.GetPadOwner(), dl, dl.targetPad)
	dl.Refresh()
	return dl
}

// CreateRenderer creates the WidgetRenderer for a DiagramLink
func (dl *DiagramLink) CreateRenderer() fyne.WidgetRenderer {
	dlr := diagramLinkRenderer{
		link: dl,
	}

	(&dlr).Refresh()

	return &dlr
}

// AddSourceAnchoredText creates a new AnchoredText widget and adds it to the DiagramLink at the Source
// position. It uses the supplied key to index the widget so that it can be retrieved later.
// Multiple AnchoredText widgets can be added.
func (dl *DiagramLink) AddSourceAnchoredText(key string, displayedText string) {
	at := NewAnchoredText(displayedText)
	at.link = dl
	dl.sourceAnchoredText[key] = at
	at.SetReferencePosition(dl.getSourcePosition())
	at.Move(dl.getSourcePosition())
	dl.Refresh()
}

// AddSourceDecoration adds the supplied Decoration widget at the Source position. Multiple
// calls to this function will stack the decorations along the line segment at the Source position.
func (dl *DiagramLink) AddSourceDecoration(decoration Decoration) {
	decoration.setLink(dl)
	dl.SourceDecorations = append(dl.SourceDecorations, decoration)
	dl.Refresh()
}

// AddMidpointAnchoredText creates a new AnchoredText widget and adds it to the DiagramLink at the Midpoint
// position. It uses the supplied key to index the widget so that it can be retrieved later.
// Multiple AnchoredText widgets can be added.
func (dl *DiagramLink) AddMidpointAnchoredText(key string, displayedText string) {
	at := NewAnchoredText(displayedText)
	at.link = dl
	dl.midpointAnchoredText[key] = at
	at.SetReferencePosition(dl.getMidPosition())
	at.Move(dl.getMidPosition())
	dl.Refresh()
}

// AddMidpointDecoration adds the supplied Decoration widget at the Midpoint position. Multiple
// calls to this function will stack the decorations along the line segment at the Midpoint position.
func (dl *DiagramLink) AddMidpointDecoration(decoration Decoration) {
	decoration.setLink(dl)
	dl.MidpointDecorations = append(dl.MidpointDecorations, decoration)
	dl.Refresh()
}

// AddTargetAnchoredText creates a new AnchoredText widget and adds it to the DiagramLink at the Target
// position. It uses the supplied key to index the widget so that it can be retrieved later.
// Multiple AnchoredText widgets can be added.
func (dl *DiagramLink) AddTargetAnchoredText(key string, displayedText string) {
	at := NewAnchoredText(displayedText)
	at.link = dl
	dl.targetAnchoredText[key] = at
	at.SetReferencePosition(dl.getTargetPosition())
	at.Move(dl.getTargetPosition())
	dl.Refresh()
}

// AddTargetDecoration adds the supplied Decoration widget at the Target position. Multiple
// calls to this function will stack the decorations along the line segment at the Target position.
func (dl *DiagramLink) AddTargetDecoration(decoration Decoration) {
	decoration.setLink(dl)
	dl.TargetDecorations = append(dl.TargetDecorations, decoration)
	dl.Refresh()
}

// GetDefaultConnectionPad returns the midPad of the Link
func (dl *DiagramLink) GetDefaultConnectionPad() ConnectionPad {
	return dl.GetMidPad()
}

// GetMidPad returns the PointPad at the midpoint so that it can be used as either the Source or Target
// pad for another Link.
func (dl *DiagramLink) GetMidPad() ConnectionPad {
	return dl.midPad
}

func (dl *DiagramLink) getMidPosition() fyne.Position {
	// TODO update when additional points are introduced
	sourcePoint := dl.linkPoints[0].Position()
	targetPoint := dl.linkPoints[len(dl.linkPoints)-1].Position()
	midPoint := fyne.NewPos((sourcePoint.X+targetPoint.X)/2, (sourcePoint.Y+targetPoint.Y)/2)
	return midPoint
}

func (dl *DiagramLink) getSourcePosition() fyne.Position {
	return dl.linkPoints[0].Position()
}

func (dl *DiagramLink) getTargetPosition() fyne.Position {
	return dl.linkPoints[len(dl.linkPoints)-1].Position()
}

func (dl *DiagramLink) handleDragged(handle *Handle, event *fyne.DragEvent) {
	// TODO implement this
}

// HideHandles prevents the handles from being displayed
func (dl *DiagramLink) HideHandles() {
	dl.showHandles = false
	dl.Refresh()
}

// MouseIn responds to the mouse entering the bounding rectangle of the Link
func (dl *DiagramLink) MouseIn(event *desktop.MouseEvent) {
	// TODO implement this
}

// MouseMoved responds to the mouse moving while within the bounding rectangle of the Link
func (dl *DiagramLink) MouseMoved(event *desktop.MouseEvent) {
	// TODO implement this
}

// MouseOut responds to the mouse leaving the bounding rectangle of the Link
func (dl *DiagramLink) MouseOut() {
}

// ShowHandles causes the handles of the Link to be displayed
func (dl *DiagramLink) ShowHandles() {
	dl.showHandles = true
	dl.Refresh()
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
	padBasedSourcePosition := dlr.link.sourcePad.getConnectionPoint(padBasedTargetReferencePoint)
	padBasedTargetPosition := dlr.link.targetPad.getConnectionPoint(padBasedSourceReferencePoint)
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
		decoration.setBaseAngle(sourceAngle)
		sourceOffset = sourceOffset + float64(decoration.GetReferenceLength())
	}

	// TODO Update target angle for multiple segments
	targetAngle := r2.AddAngles(sourceAngle, math.Pi)

	midOffset := 0.0
	for _, decoration := range dlr.link.MidpointDecorations {
		decorationReferencePoint := fyne.Position{
			X: float32(float64(dlr.link.getMidPosition().X) + math.Cos(targetAngle)*midOffset),
			Y: float32(float64(dlr.link.getMidPosition().Y) - math.Sin(targetAngle)*midOffset),
		}
		decoration.Move(decorationReferencePoint)
		decoration.setBaseAngle(sourceAngle)
		midOffset = midOffset + float64(decoration.GetReferenceLength())
	}
	dlr.link.midPad.Move(dlr.link.getMidPosition())

	targetOffset := 0.0
	for _, decoration := range dlr.link.TargetDecorations {
		decorationReferencePoint := fyne.Position{
			X: float32(float64(dlr.link.linkPoints[len(dlr.link.linkPoints)-1].Position().X) + math.Cos(targetAngle)*targetOffset),
			Y: float32(float64(dlr.link.linkPoints[len(dlr.link.linkPoints)-1].Position().Y) - math.Sin(targetAngle)*targetOffset),
		}
		decoration.Move(decorationReferencePoint)
		decoration.setBaseAngle(targetAngle)
		targetOffset = targetOffset + float64(decoration.GetReferenceLength())
	}
	for _, anchoredText := range dlr.link.sourceAnchoredText {
		anchoredText.SetReferencePosition(dlr.link.getSourcePosition())
	}
	for _, anchoredText := range dlr.link.midpointAnchoredText {
		anchoredText.SetReferencePosition(dlr.link.getMidPosition())
	}
	for _, anchoredText := range dlr.link.targetAnchoredText {
		anchoredText.SetReferencePosition(dlr.link.getTargetPosition())
	}

	// Now we take care of property changes.
	for _, linkSegment := range dlr.link.linkSegments {
		linkSegment.Refresh()
	}
	for _, decoration := range dlr.link.SourceDecorations {
		decoration.SetStrokeColor(dlr.link.LinkColor)
		decoration.SetStrokeWidth(dlr.link.strokeWidth)
		decoration.SetFillColor(dlr.link.diagram.GetBackgroundColor())
		decoration.Refresh()
	}
	for _, decoration := range dlr.link.MidpointDecorations {
		decoration.SetStrokeColor(dlr.link.LinkColor)
		decoration.SetStrokeWidth(dlr.link.strokeWidth)
		decoration.SetFillColor(dlr.link.diagram.GetBackgroundColor())
		decoration.Refresh()
	}
	for _, decoration := range dlr.link.TargetDecorations {
		decoration.SetStrokeColor(dlr.link.LinkColor)
		decoration.SetStrokeWidth(dlr.link.strokeWidth)
		decoration.SetFillColor(dlr.link.diagram.GetBackgroundColor())
		decoration.Refresh()
	}
	dlr.link.diagram.refreshDependentLinks(dlr.link)
	dlr.link.GetDiagram().forceRepaint()
}
