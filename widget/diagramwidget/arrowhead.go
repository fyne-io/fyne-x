package diagramwidget

import (
	"image/color"
	"math"

	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	defaultTheta  float64 = 0.5235 // 30 degrees in radians
	defaultLength int     = 15
)

// Arrowhead defines a canvas object which renders an arrow. The arrowhead is defined with reference
// to X axis, with the tip of the arrow being at the origin. When rendered, the arrowhead is rotated
// to match the angle of the link's line segment with which it is oriented, indicated by the baseAngle.
type Arrowhead struct {
	widget.BaseWidget
	link *BaseDiagramLink
	// baseAngle is used to define direction in which the arrowhead points
	// Base fyne.Position
	baseAngle float64
	// Position() is the point at which the tip of the arrow will be placed.
	// StrokeWidth is the width of the arrowhead lines
	StrokeWidth float32
	// StrokeColor is the color of the arrowhead
	StrokeColor color.Color
	// Theta is the angle between each of the tails and the nominal reference axis.
	// This angle is in radians.
	Theta float64
	// Length is the length of the two "tails" that intersect at the tip.
	Length int
	// central *canvas.Line
	// left    *canvas.Line
	// right   *canvas.Line
	visible bool
}

// NewArrowhead creates an arrowhead with defaults
func NewArrowhead() *Arrowhead {
	a := &Arrowhead{
		baseAngle:   0.0,
		StrokeWidth: defaultStrokeWidth,
		StrokeColor: theme.ForegroundColor(),
		Theta:       defaultTheta,
		Length:      defaultLength,
		visible:     true,
	}
	a.ExtendBaseWidget(a)
	return a
}

// CreateRenderer creates a renderer for the Arrowhead
func (a *Arrowhead) CreateRenderer() fyne.WidgetRenderer {
	ar := arrowheadRenderer{
		arrowhead: a,
		left:      canvas.NewLine(a.link.LinkColor),
		right:     canvas.NewLine(a.link.LinkColor),
	}
	return &ar
}

// GetReferenceLength returns the length of the decoration along the reference axis
func (a *Arrowhead) GetReferenceLength() float32 {
	return float32(math.Abs(math.Cos(float64(a.Theta)) * float64(a.Length)))
}

// LeftPoint returns the position of the end of the left half of the arrowhead
func (a *Arrowhead) LeftPoint() fyne.Position {
	leftAngle := r2.AddAngles(a.baseAngle, -a.Theta)
	// We have to change the sign of Y because the window coordinate Y axis goes down rather than up
	leftPosition := fyne.Position{
		X: float32(float64(a.Length) * math.Cos(leftAngle)),
		Y: -float32(float64(a.Length) * math.Sin(leftAngle)),
	}
	return leftPosition
}

// MinSize returns the minimum size which is the actual size of the arrowhead
func (a *Arrowhead) MinSize() fyne.Size {
	return a.Size()
}

// Resize scales the arrowhead
func (a *Arrowhead) Resize(s fyne.Size) {
	// We get the current size and scale the length based on the difference between sizes
	currentSize := a.Size()
	currentLengthVector := r2.V2(float64(currentSize.Width), float64(currentSize.Height))
	currentLength := currentLengthVector.Length()
	newLengthVector := r2.V2(float64(s.Width), float64(s.Height))
	newLength := newLengthVector.Length()
	a.Length = int(float64(a.Length) * newLength / currentLength)
}

// RightPoint returns the position of the end of the right half of the arrowhead
func (a *Arrowhead) RightPoint() fyne.Position {
	rightAngle := r2.AddAngles(a.baseAngle, a.Theta)
	// We have to change the sign of Y because the window coordinate Y axis goes down rather than up
	rightPosition := fyne.Position{
		X: float32(float64(a.Length) * math.Cos(rightAngle)),
		Y: -float32(float64(a.Length) * math.Sin(rightAngle)),
	}
	return rightPosition
}

// setBaseAngle sets the angle (in radians) of the reference axis
func (a *Arrowhead) setBaseAngle(angle float64) {
	a.baseAngle = angle
}

// setLink sets the DiagramLink on which this arrowhead appears
func (a *Arrowhead) setLink(link *BaseDiagramLink) {
	a.link = link
}

// SetFillColor is a noop for the arrowhead
func (a *Arrowhead) SetFillColor(fillColor color.Color) {

}

// SetSolid is a noop because the arrowhead is an open structure
func (a *Arrowhead) SetSolid(bool) {
}

// SetStrokeColor sets the color used to draw the arrowhead
func (a *Arrowhead) SetStrokeColor(strokeColor color.Color) {
	a.StrokeColor = strokeColor
}

// SetStrokeWidth sets the width of the lines used to render the arrowhead
func (a *Arrowhead) SetStrokeWidth(strokeWidth float32) {
	a.StrokeWidth = strokeWidth
}

// Size returns the size of the arrowhead
func (a *Arrowhead) Size() fyne.Size {
	lp := a.LeftPoint()
	rp := a.RightPoint()
	points := []r2.Vec2{
		{X: float64(a.Position().X), Y: float64(a.Position().Y)},
		{X: float64(lp.X), Y: float64(lp.Y)},
		{X: float64(rp.X), Y: float64(rp.Y)},
	}

	bounding := r2.BoundingBox(points)
	return fyne.Size{
		Width:  float32(bounding.Width()),
		Height: float32(bounding.Height()),
	}
}

type arrowheadRenderer struct {
	arrowhead *Arrowhead
	left      *canvas.Line
	right     *canvas.Line
}

func (ar *arrowheadRenderer) Destroy() {
}

func (ar *arrowheadRenderer) MinSize() fyne.Size {
	return ar.arrowhead.Size()
}

func (ar *arrowheadRenderer) Layout(size fyne.Size) {
	ar.left.Position1 = fyne.Position{X: 0, Y: 0}
	ar.left.Position2 = ar.arrowhead.LeftPoint()
	ar.right.Position1 = fyne.Position{X: 0, Y: 0}
	ar.right.Position2 = ar.arrowhead.RightPoint()
}

func (ar *arrowheadRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{
		ar.left,
		ar.right,
	}
	return obj
}

func (ar *arrowheadRenderer) Refresh() {
	ar.left.StrokeWidth = ar.arrowhead.StrokeWidth
	ar.right.StrokeWidth = ar.arrowhead.StrokeWidth
	ar.left.StrokeColor = ar.arrowhead.StrokeColor
	ar.right.StrokeColor = ar.arrowhead.StrokeColor
	if ar.arrowhead.visible {
		ar.left.Show()
		ar.right.Show()
	} else {
		ar.left.Hide()
		ar.right.Hide()
	}
}
