package diagramwidget

import (
	"image/color"
	"log"
	"math"

	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

var _ fyne.Tappable = (*BaseDiagramLink)(nil)
var _ desktop.Hoverable = (*BaseDiagramLink)(nil)
var _ DiagramElement = (*BaseDiagramLink)(nil)

type LinkEnd int

var LinkEnds [2]LinkEnd = [2]LinkEnd{SOURCE, TARGET}

const (
	SOURCE LinkEnd = iota
	TARGET
)

func (le LinkEnd) ToString() string {
	switch le {
	case SOURCE:
		return "source"
	case TARGET:
		return "target"
	}
	return ""
}

type DiagramLink interface {
	DiagramElement
	getBaseDiagramLink() *BaseDiagramLink
	getLinkPoints() []*LinkPoint
	GetSourcePad() ConnectionPad
	GetTargetPad() ConnectionPad
	isConnectionAllowed(*LinkPoint, ConnectionPad) bool
	SetSourcePad(ConnectionPad)
	SetTargetPad(ConnectionPad)
}

// BaseDiagramLink is a directed graphic connection between two DiagramElements that are referred to as the Source
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
type BaseDiagramLink struct {
	diagramElement
	linkPoints           []*LinkPoint
	linkSegments         []*LinkSegment
	LinkColor            color.Color
	strokeWidth          float32
	sourcePad            ConnectionPad
	targetPad            ConnectionPad
	SourceDecorations    []Decoration
	sourceAnchoredText   map[string]*AnchoredText
	TargetDecorations    []Decoration
	targetAnchoredText   map[string]*AnchoredText
	MidpointDecorations  []Decoration
	midpointAnchoredText map[string]*AnchoredText
	// We keep the typed link so that when extensions are created the callbacks are called with the correct type
	typedLink DiagramLink
}

// NewDiagramLink creates a DiagramLink widget connecting the two indicated ConnectionPads. It adds itself to the
// DiagramWidget, indexed by the supplied LinkID. This id must be unique across all of the DiagramElements in the Diagram.
// It can be used to retrieve the DiagramLink from the Diagram. The ID is intended to be used to facilitate mapping the
// DiagramLink to the information it represents in the application. The DiagramLink uses the DiagramWidget's ForegroundColor
// as the default color for the line segments.
func NewDiagramLink(diagram *DiagramWidget, linkID string) *BaseDiagramLink {
	bdl := &BaseDiagramLink{}
	InitializeBaseDiagramLink(bdl, diagram, linkID)
	return bdl
}

// InitializeBaseDiagramLink initializes the BaseDiagramLink. It must be called by any extensions to BaseDiagramLink
func InitializeBaseDiagramLink(diagramLink DiagramLink, diagram *DiagramWidget, linkID string) {
	bdl := diagramLink.getBaseDiagramLink()
	bdl.linkPoints = []*LinkPoint{}
	bdl.linkSegments = []*LinkSegment{}
	bdl.LinkColor = diagram.GetForegroundColor()
	bdl.strokeWidth = 2
	bdl.sourceAnchoredText = make(map[string]*AnchoredText)
	bdl.midpointAnchoredText = make(map[string]*AnchoredText)
	bdl.targetAnchoredText = make(map[string]*AnchoredText)
	bdl.diagramElement.initialize(diagram, linkID)
	bdl.linkPoints = append(bdl.linkPoints, NewLinkPoint(bdl))
	bdl.linkPoints = append(bdl.linkPoints, NewLinkPoint(bdl))
	bdl.linkSegments = append(bdl.linkSegments, NewLinkSegment(bdl, bdl.linkPoints[0].Position(), bdl.linkPoints[1].Position()))
	bdl.pads["default"] = NewPointPad(bdl)
	bdl.pads["default"].Move(bdl.getMidPosition())
	bdl.pads["default"].Hide()
	bdl.ExtendBaseWidget(diagramLink)
	bdl.typedLink = diagramLink
	for _, linkEnd := range LinkEnds {
		newHandle := NewHandle(bdl)
		bdl.handles[linkEnd.ToString()] = newHandle
		newHandle.Hide()
	}
	bdl.diagram.addLink(diagramLink)
	diagramLink.Refresh()
}

// CreateRenderer creates the WidgetRenderer for a DiagramLink
func (bdl *BaseDiagramLink) CreateRenderer() fyne.WidgetRenderer {
	dlr := diagramLinkRenderer{
		link: bdl,
	}

	(&dlr).Refresh()

	return &dlr
}

// AddSourceAnchoredText creates a new AnchoredText widget and adds it to the DiagramLink at the Source
// position. It uses the supplied key to index the widget so that it can be retrieved later.
// Multiple AnchoredText widgets can be added.
func (bdl *BaseDiagramLink) AddSourceAnchoredText(key string, displayedText string) *AnchoredText {
	at := NewAnchoredText(displayedText)
	at.link = bdl
	bdl.sourceAnchoredText[key] = at
	at.SetReferencePosition(bdl.getSourcePosition())
	at.Move(bdl.getSourcePosition())
	bdl.Refresh()
	return at
}

// AddSourceDecoration adds the supplied Decoration widget at the Source position. Multiple
// calls to this function will stack the decorations along the line segment at the Source position.
func (bdl *BaseDiagramLink) AddSourceDecoration(decoration Decoration) {
	decoration.setLink(bdl)
	bdl.SourceDecorations = append(bdl.SourceDecorations, decoration)
	bdl.Refresh()
}

// AddMidpointAnchoredText creates a new AnchoredText widget and adds it to the DiagramLink at the Midpoint
// position. It uses the supplied key to index the widget so that it can be retrieved later.
// Multiple AnchoredText widgets can be added.
func (bdl *BaseDiagramLink) AddMidpointAnchoredText(key string, displayedText string) *AnchoredText {
	at := NewAnchoredText(displayedText)
	at.link = bdl
	bdl.midpointAnchoredText[key] = at
	at.SetReferencePosition(bdl.getMidPosition())
	at.Move(bdl.getMidPosition())
	bdl.Refresh()
	return at
}

// AddMidpointDecoration adds the supplied Decoration widget at the Midpoint position. Multiple
// calls to this function will stack the decorations along the line segment at the Midpoint position.
func (bdl *BaseDiagramLink) AddMidpointDecoration(decoration Decoration) {
	decoration.setLink(bdl)
	bdl.MidpointDecorations = append(bdl.MidpointDecorations, decoration)
	bdl.Refresh()
}

// AddTargetAnchoredText creates a new AnchoredText widget and adds it to the DiagramLink at the Target
// position. It uses the supplied key to index the widget so that it can be retrieved later.
// Multiple AnchoredText widgets can be added.
func (bdl *BaseDiagramLink) AddTargetAnchoredText(key string, displayedText string) *AnchoredText {
	at := NewAnchoredText(displayedText)
	at.link = bdl
	bdl.targetAnchoredText[key] = at
	at.SetReferencePosition(bdl.getTargetPosition())
	at.Move(bdl.getTargetPosition())
	bdl.Refresh()
	return at
}

// AddTargetDecoration adds the supplied Decoration widget at the Target position. Multiple
// calls to this function will stack the decorations along the line segment at the Target position.
func (bdl *BaseDiagramLink) AddTargetDecoration(decoration Decoration) {
	decoration.setLink(bdl)
	bdl.TargetDecorations = append(bdl.TargetDecorations, decoration)
	bdl.Refresh()
}

// getBaseDiagramLink returns a pointer to the BaseDiagramLink
func (bdl *BaseDiagramLink) getBaseDiagramLink() *BaseDiagramLink {
	return bdl
}

// GetDefaultConnectionPad returns the midPad of the Link
func (bdl *BaseDiagramLink) GetDefaultConnectionPad() ConnectionPad {
	return bdl.pads["default"]
}

// getHandleKey returns the key for the given handle
func (bdl *BaseDiagramLink) getHandleKey(handle *Handle) string {
	for key, h := range bdl.handles {
		if h == handle {
			return key
		}
	}
	return ""
}

func (bdl *BaseDiagramLink) getLinkPoints() []*LinkPoint {
	return bdl.linkPoints
}

// GetMidPad returns the PointPad at the midpoint so that it can be used as either the Source or Target
// pad for another Link.
func (bdl *BaseDiagramLink) GetMidPad() ConnectionPad {
	return bdl.pads["default"]
}

func (bdl *BaseDiagramLink) getMidPosition() fyne.Position {
	// TODO update when additional points are introduced
	sourcePoint := bdl.linkPoints[0].Position()
	targetPoint := bdl.linkPoints[len(bdl.linkPoints)-1].Position()
	midPoint := fyne.NewPos((sourcePoint.X+targetPoint.X)/2, (sourcePoint.Y+targetPoint.Y)/2)
	return midPoint
}

func (bdl *BaseDiagramLink) GetMidpointAnchoredText(key string) *AnchoredText {
	return bdl.midpointAnchoredText[key]
}

func (bdl *BaseDiagramLink) GetSourceAnchoredText(key string) *AnchoredText {
	return bdl.sourceAnchoredText[key]
}

func (bdl *BaseDiagramLink) GetTargetAnchoredText(key string) *AnchoredText {
	return bdl.targetAnchoredText[key]
}

func (bdl *BaseDiagramLink) GetSourcePad() ConnectionPad {
	return bdl.sourcePad
}

func (bdl *BaseDiagramLink) getSourcePosition() fyne.Position {
	return bdl.linkPoints[0].Position()
}

func (bdl *BaseDiagramLink) GetTargetPad() ConnectionPad {
	return bdl.targetPad
}

func (bdl *BaseDiagramLink) getTargetPosition() fyne.Position {
	return bdl.linkPoints[len(bdl.linkPoints)-1].Position()
}

func (bdl *BaseDiagramLink) handleDragged(handle *Handle, event *fyne.DragEvent) {
	handleKey := bdl.getHandleKey(handle)
	var linkPoint *LinkPoint
	var pad ConnectionPad
	switch handleKey {
	case SOURCE.ToString():
		linkPoint = bdl.linkPoints[0]
		pad = bdl.sourcePad
		bdl.sourcePad = nil
	case TARGET.ToString():
		linkPoint = bdl.linkPoints[len(bdl.linkPoints)-1]
		pad = bdl.targetPad
		bdl.targetPad = nil
	}
	if linkPoint == nil {
		return
	}
	connTrans := bdl.diagram.connectionTransaction
	if connTrans == nil {
		connTrans = NewConnectionTransaction(linkPoint, bdl, pad, linkPoint.Position())
		bdl.diagram.connectionTransaction = connTrans
		// TODO remove this after fyne Issue #3906 has been resolved
		bdl.diagram.showAllPads()

	} else if connTrans.linkPoint != linkPoint {
		// The existing transaction is for a different linkPoint
		return
	}
	currentPosition := linkPoint.Position()
	newPosition := fyne.NewPos(currentPosition.X+event.Dragged.DX, currentPosition.Y+event.Dragged.DY)
	linkPoint.Move(newPosition)
	bdl.Refresh()
}

func (bdl *BaseDiagramLink) handleDragEnd(handle *Handle) {
	connTrans := bdl.diagram.connectionTransaction
	handleKey := bdl.getHandleKey(handle)
	if connTrans != nil {
		if connTrans.pendingPad != nil {
			// We have a new pad for connection
			if connTrans.initialPad != nil {
				bdl.diagram.removeLinkDependency(connTrans.initialPad.GetPadOwner(), bdl, connTrans.initialPad)
			}
			bdl.diagram.addLinkDependency(connTrans.pendingPad.GetPadOwner(), bdl, connTrans.pendingPad)
			bdl.pads[handleKey] = connTrans.pendingPad
			switch handleKey {
			case SOURCE.ToString():
				bdl.sourcePad = connTrans.pendingPad
			case TARGET.ToString():
				bdl.targetPad = connTrans.pendingPad
			}
			if bdl.diagram.LinkConnectionChangedCallback != nil {
				bdl.diagram.LinkConnectionChangedCallback(bdl.typedLink, handleKey, connTrans.initialPad, connTrans.pendingPad)
			}
		} else {
			// We revert to the original pad.
			bdl.pads[handleKey] = connTrans.initialPad
			switch handleKey {
			case SOURCE.ToString():
				bdl.sourcePad = connTrans.initialPad
			case TARGET.ToString():
				bdl.targetPad = connTrans.initialPad
			}
		}
		bdl.diagram.connectionTransaction = nil
		bdl.diagram.hideAllPads()
		bdl.diagram.SelectDiagramElement(bdl)
		bdl.Refresh()
	}
}

func (bdl *BaseDiagramLink) isConnectionAllowed(linkPoint *LinkPoint, pad ConnectionPad) bool {
	pointIndex := -1
	for i, lp := range bdl.linkPoints {
		if lp == linkPoint {
			pointIndex = i
		}
	}
	if pointIndex == -1 {
		// the point doesn't belong to this link
		return false
	}
	if pointIndex != 0 && pointIndex != len(bdl.linkPoints)-1 {
		// the point is not the source or target point
		return false
	}
	if bdl.diagram.IsConnectionAllowedCallback != nil {
		var linkEnd LinkEnd
		if pointIndex == 0 {
			linkEnd = SOURCE
		} else if pointIndex == len(bdl.linkPoints)-1 {
			linkEnd = TARGET
		}
		return bdl.diagram.IsConnectionAllowedCallback(bdl, linkEnd, pad)
	}
	// By default, we accept any connection
	return true
}

// IsLink returns true because this is a link
func (bdl *BaseDiagramLink) IsLink() bool {
	return true
}

// IsNode returns false because this is a link
func (bdl *BaseDiagramLink) IsNode() bool {
	return false
}

// MouseIn responds to the mouse entering the bounding rectangle of the Link
func (bdl *BaseDiagramLink) MouseIn(event *desktop.MouseEvent) {
	log.Print("Entered Link")
	// TODO implement this
}

// MouseMoved responds to the mouse moving while within the bounding rectangle of the Link
func (bdl *BaseDiagramLink) MouseMoved(event *desktop.MouseEvent) {
	// TODO implement this
}

// MouseOut responds to the mouse leaving the bounding rectangle of the Link
func (bdl *BaseDiagramLink) MouseOut() {
	log.Printf("Left Link")
}

// SetSourcePad sets the source pad and adds the link dependency to the diagram
func (bdl *BaseDiagramLink) SetSourcePad(pad ConnectionPad) {
	oldPad := bdl.sourcePad
	if oldPad != pad {
		if oldPad != nil {
			bdl.diagram.removeLinkDependency(oldPad.GetPadOwner(), bdl, oldPad)
		}
		bdl.sourcePad = pad
		bdl.diagram.addLinkDependency(bdl.sourcePad.GetPadOwner(), bdl, bdl.sourcePad)
		if bdl.diagram.LinkConnectionChangedCallback != nil {
			bdl.diagram.LinkConnectionChangedCallback(bdl.typedLink, SOURCE.ToString(), oldPad, pad)
		}
		bdl.Refresh()
	}
}

// SetTargetPad sets the target pad and adds the link dependency to the diagram
func (bdl *BaseDiagramLink) SetTargetPad(pad ConnectionPad) {
	oldPad := bdl.targetPad
	if oldPad != pad {
		if oldPad != nil {
			bdl.diagram.removeLinkDependency(oldPad.GetPadOwner(), bdl, oldPad)
		}
		bdl.targetPad = pad
		bdl.diagram.addLinkDependency(bdl.targetPad.GetPadOwner(), bdl, bdl.targetPad)
		if bdl.diagram.LinkConnectionChangedCallback != nil {
			bdl.diagram.LinkConnectionChangedCallback(bdl.typedLink, TARGET.ToString(), oldPad, pad)
		}
		bdl.Refresh()
	}
}

// Tapped handles tap events
func (bdl *BaseDiagramLink) Tapped(event *fyne.PointEvent) {
	log.Print("Link tapped")
}

// diagramLinkRenderer
type diagramLinkRenderer struct {
	link *BaseDiagramLink
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
	for _, pad := range dlr.link.pads {
		obj = append(obj, pad)
	}
	for _, handle := range dlr.link.handles {
		obj = append(obj, handle)
	}
	return obj
}

func (dlr *diagramLinkRenderer) Refresh() {
	// The pads to which the link is connected can be nil during a connection transaction, in which case we leave the end points
	// at their present location. Note that the initial computations are done in diagram coordinates,
	// then link coordinates
	var sourceDiagramCoordinateReferencePoint fyne.Position
	var targetDiagramCoordinateReferencePoint fyne.Position
	var sourceDiagramCoordinatePosition fyne.Position
	var targetDiagramCoordinatePosition fyne.Position
	currentSourceDiagramCoordinatePosition := dlr.link.getSourcePosition().Add(dlr.link.Position())
	currentTargetDiagramCoordinatePosition := dlr.link.getTargetPosition().Add(dlr.link.Position())
	if dlr.link.sourcePad != nil {
		sourceDiagramCoordinateReferencePoint = dlr.link.sourcePad.GetCenterInDiagramCoordinates()
	} else {
		// we have to translate the source position back to diagram coordinates
		sourceDiagramCoordinateReferencePoint = currentSourceDiagramCoordinatePosition
	}
	if dlr.link.targetPad != nil {
		targetDiagramCoordinateReferencePoint = dlr.link.targetPad.GetCenterInDiagramCoordinates()
	} else {
		// we have to translate the target position back to diagram coordinates
		targetDiagramCoordinateReferencePoint = currentTargetDiagramCoordinatePosition
	}
	if dlr.link.sourcePad != nil {
		sourceDiagramCoordinatePosition = dlr.link.sourcePad.getConnectionPointInDiagramCoordinates(targetDiagramCoordinateReferencePoint)
	} else {
		sourceDiagramCoordinatePosition = currentSourceDiagramCoordinatePosition
	}
	if dlr.link.targetPad != nil {
		targetDiagramCoordinatePosition = dlr.link.targetPad.getConnectionPointInDiagramCoordinates(sourceDiagramCoordinateReferencePoint)
	} else {
		targetDiagramCoordinatePosition = currentTargetDiagramCoordinatePosition
	}
	// The Position of the link is the upper left hand corner of a bounding box surrounding the source and target positions
	linkPosition := fyne.NewPos(float32(math.Min(float64(sourceDiagramCoordinatePosition.X), float64(targetDiagramCoordinatePosition.X))),
		float32(math.Min(float64(sourceDiagramCoordinatePosition.Y), float64(targetDiagramCoordinatePosition.Y))))
	dlr.link.Move(linkPosition)

	// Now we put the source and target positions back into link coordinates by subtracting the linkPosition
	dlr.link.linkPoints[0].Move(sourceDiagramCoordinatePosition.Subtract(linkPosition))
	dlr.link.linkPoints[len(dlr.link.linkPoints)-1].Move(targetDiagramCoordinatePosition.Subtract(linkPosition))
	// TODO adjust the position of the other points based on the change in link position
	// Now resize the link - note that MinSize is derived from the point positions
	dlr.link.Resize(dlr.MinSize())

	// Position segments only after all points have been positioned
	for i := 0; i < len(dlr.link.linkPoints)-1; i++ {
		linkSegment := dlr.link.linkSegments[i]
		linkSegment.SetPoints(dlr.link.linkPoints[i].Position(), dlr.link.linkPoints[i+1].Position())
		// linkSegment.Refresh()
	}

	// Have to change the sign of Y since the window inverts the Y axis
	lineVector := r2.Vec2{X: float64(dlr.link.linkPoints[1].Position().X - dlr.link.linkPoints[0].Position().X),
		Y: -float64(dlr.link.linkPoints[1].Position().Y - dlr.link.linkPoints[0].Position().Y)}
	var sourceAngle float64
	if lineVector.Length() == 0 {
		sourceAngle = 0
	} else {
		sourceAngle = lineVector.Angle()
	}
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
	defaultPadPosition := dlr.link.getMidPosition().AddXY(-pointPadSize/2, -pointPadSize/2)
	dlr.link.pads["default"].Move(defaultPadPosition)
	dlr.link.pads["default"].Resize(fyne.NewSize(pointPadSize, pointPadSize))
	dlr.link.pads["default"].Refresh()

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

	// calculate the handle positions
	for key, handle := range dlr.link.handles {
		switch key {
		case SOURCE.ToString():
			handle.Move(dlr.link.linkPoints[0].Position())
		case TARGET.ToString():
			handle.Move(dlr.link.linkPoints[len(dlr.link.linkPoints)-1].Position())
		}
		handle.Resize(fyne.NewSize(handle.handleSize, handle.handleSize))
		handle.Refresh()
	}

	dlr.link.diagram.refreshDependentLinks(dlr.link)
	dlr.link.GetDiagram().ForceRepaint()
}

type connectionTransaction struct {
	linkPoint       *LinkPoint
	link            DiagramLink
	initialPad      ConnectionPad
	initialPosition fyne.Position
	pendingPad      ConnectionPad
}

func NewConnectionTransaction(linkPoint *LinkPoint, link DiagramLink, initialPad ConnectionPad, initialPosition fyne.Position) *connectionTransaction {
	ct := &connectionTransaction{
		linkPoint:       linkPoint,
		link:            link,
		initialPad:      initialPad,
		initialPosition: initialPosition,
	}
	return ct
}
