package diagramwidget

import (
	"image"
	"image/color"
	"math"

	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/math/fixed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

var _ Decoration = (*Polygon)(nil)

// Polygon defines a canvas object which renders a polygon defined with respect to
// a reference point (the Position of the canvas object). The polygon is intended
// to be a decoration on a line. The nominal definition of the polygon assumes that
// the nominal line lays along the X axis extending to the right of the reference point.
// The polygon is defined as a series of points, with (0,0) corresponting to the
// the polygon's Position. When rendered, the polygon is rotated by the BaseAngle.
// The BaseAngle is supplied by the Link for which the Polygon is a decoration.
// By default, the polygon is a closed structure that can optionally be filled by
// providing p fill color (the default is no fill). SetClosed(false) can be used
// to make the structure open, which makes it a polyline. Fill color is ignored in this case.
// By default, stroke color and fill color are color.Black
type Polygon struct {
	widget.BaseWidget
	link *BaseDiagramLink
	// baseAngle is used to define the rotation of the polygon from the nominal position
	// Base fyne.Position
	baseAngle float64
	// StrokeWidth is the width of the polygon perimeter line
	StrokeWidth float32
	// StrokeColor is the color of the polygon perimeter line
	StrokeColor color.Color
	// FillColor is the color if the polygon interior. A nil value indicates
	// no fill
	FillColor      color.Color
	visible        bool
	definingPoints []fyne.Position
	closed         bool
	solid          bool
}

// NewPolygon creates a Polygon as defined by the supplied points
func NewPolygon(definingPoints []fyne.Position) *Polygon {
	p := &Polygon{
		StrokeWidth:    defaultStrokeWidth,
		StrokeColor:    color.Black,
		FillColor:      color.Black,
		visible:        true,
		definingPoints: definingPoints,
		closed:         true,
		solid:          false,
	}
	p.ExtendBaseWidget(p)
	return p
}

// CreateRenderer creates an instance of the renderer for the Polygon
func (p *Polygon) CreateRenderer() fyne.WidgetRenderer {
	pr := polygonRenderer{
		polygon: p,
		image:   *canvas.NewImageFromImage(nil),
	}
	return &pr
}

// GetReferenceLength returns the length of the decoration along the reference axis
func (p *Polygon) GetReferenceLength() float32 {
	xMin := float64(0)
	xMax := float64(0)
	for _, point := range p.definingPoints {
		xMin = math.Min(xMin, float64(point.X))
		xMax = math.Max(xMax, float64(point.X))
	}
	return float32(math.Abs(xMax-xMin)) + p.StrokeWidth
}

// getRenderingData returns the defining points rotated to the correct orientation and
// translated so that all points have positive coordinates (as required by the image rendering).
// It also returns an offset vector indicating the translation required for the rendered image
// so that the reference point on the defining points (the 0,0 coordinate) will align with the
// position of the polygon. Thus the rendered polygon will appear on the link with the correct
// orientation and position.
func (p *Polygon) getRenderingData() ([]fyne.Position, fyne.Position) {
	rotatedPoints := p.getRotatedPoints()
	var xMin, yMin float64
	for i, point := range rotatedPoints {
		if i == 0 {
			xMin = float64(point.X)
			yMin = float64(point.Y)
		} else {
			xMin = math.Min(xMin, float64(point.X))
			yMin = math.Min(yMin, float64(point.Y))
		}
	}
	// For positioning, treat the reference coordinate of 0,0 as one of the rotated points
	// even if it is not part of the polygon definition
	xMin = math.Min(xMin, 0)
	yMin = math.Min(yMin, 0)
	// the inverse of xMin and yMin now represent the point translation required to
	// normalize the coordinates for rendering into a pixmap, ignoring the stroke width.
	// We have to further shift coordinates by the stroke width so that the points render
	// at the correct positions in a pixmap. We have to add the stroke width as well
	normalizationVectorWithStrokeOfset := r2.MakeVec2(float64(-xMin)+float64(p.StrokeWidth/2), float64(-yMin)+float64(p.StrokeWidth/2))
	renderingPoints := []fyne.Position{}
	for _, point := range rotatedPoints {
		pointVector := r2.MakeVec2(float64(point.X), float64(point.Y))
		normalizedPointVector := pointVector.Add(normalizationVectorWithStrokeOfset)
		normalizedPoint := fyne.Position{X: float32(normalizedPointVector.X), Y: float32(normalizedPointVector.Y)}
		renderingPoints = append(renderingPoints, normalizedPoint)
	}
	// The offset vector is the inverse of the normalization vector without the stroke width
	offsetVector := fyne.Position{X: float32(xMin), Y: float32(yMin)}
	return renderingPoints, offsetVector
}

// getRotatedPoints returns the points after the nominal points
// have been rotated by the reference angle
func (p *Polygon) getRotatedPoints() []fyne.Position {
	rotatedPoints := []fyne.Position{}
	var rotX float32
	var rotY float32
	for _, point := range p.definingPoints {
		v2Point := r2.V2(float64(point.X), float64(point.Y))
		len := v2Point.Length()
		if len == 0 {
			rotX = 0
			rotY = 0
		} else {
			ang := v2Point.Angle()
			rotAng := r2.AddAngles(ang, p.baseAngle)
			rotX = float32(math.Cos(rotAng) * len)
			rotY = -float32(math.Sin(rotAng) * len)
		}
		rotatedPoint := fyne.Position{
			X: rotX,
			Y: rotY,
		}
		rotatedPoints = append(rotatedPoints, rotatedPoint)
	}
	return rotatedPoints
}

// MinSize returns the minimum size based on nominal polygon points, base angle, and stroke width
func (p *Polygon) MinSize() fyne.Size {
	// // The origin is always one of the points regardless of whether the polygon uses that point
	// points := []r2.Vec2{{X: 0.0, Y: 0.0}}
	points := []r2.Vec2{}
	for _, point := range p.getRotatedPoints() {
		points = append(points, r2.Vec2{X: float64(point.X), Y: float64(point.Y)})
	}
	bounding := r2.BoundingBox(points)
	return fyne.Size{
		Width:  float32(bounding.Width()) + p.StrokeWidth,
		Height: float32(bounding.Height()) + p.StrokeWidth,
	}
}

// setBaseAngle sets the angle (in radians) of the reference axis
func (p *Polygon) setBaseAngle(angle float64) {
	p.baseAngle = angle
}

// SetClosed determines whether the polygon is open or closed. A value of false
// makes the polygon a polyline
func (p *Polygon) SetClosed(closed bool) {
	p.closed = closed
	p.Refresh()
}

// SetFillColor sets the color that will be used for the polygon interior
// A nil value indicates there is no fill
func (p *Polygon) SetFillColor(fillColor color.Color) {
	p.FillColor = fillColor
}

// setLink sets the Link with which the polygon is associated
func (p *Polygon) setLink(link *BaseDiagramLink) {
	p.link = link
}

// SetStrokeColor sets the color that will be used for the polygon perimeter
func (p *Polygon) SetStrokeColor(strokeColor color.Color) {
	p.StrokeColor = strokeColor
}

// SetStrokeWidth sets the width of the linke that will be used for the
// polygon perimeter
func (p *Polygon) SetStrokeWidth(strokeWidth float32) {
	p.StrokeWidth = strokeWidth
}

// SetSolid determines whether the foreground color should be used to fill the polygon.
// It has no effect if the polygon is open instead of closed
func (p *Polygon) SetSolid(value bool) {
	p.solid = value
}

// polygonRenderer is a renderer for the Polygon
type polygonRenderer struct {
	polygon *Polygon
	image   canvas.Image
}

func (pr *polygonRenderer) Destroy() {
}

func (pr *polygonRenderer) MinSize() fyne.Size {
	return pr.polygon.Size()
}

// Layout is a noop for the Polygon
func (pr *polygonRenderer) Layout(size fyne.Size) {
}

func (pr *polygonRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{
		&pr.image,
	}
	return obj
}

func (pr *polygonRenderer) Refresh() {
	renderingPoints, offsetVector := pr.polygon.getRenderingData()
	// For efficiency, only get the size once
	polygonSize := pr.polygon.MinSize()
	stroke := pr.polygon.StrokeWidth
	width := int(polygonSize.Width + pr.polygon.StrokeWidth)
	height := int(polygonSize.Height + pr.polygon.StrokeWidth)
	pr.polygon.Resize(polygonSize)
	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(polygonSize.Width), int(polygonSize.Height), raw, raw.Bounds())

	if pr.polygon.closed && pr.polygon.FillColor != nil {
		filler := rasterx.NewFiller(width, height, scanner)
		if pr.polygon.solid {
			filler.SetColor(pr.polygon.StrokeColor)
		} else {
			filler.SetColor(pr.polygon.FillColor)
		}
		for i, point := range renderingPoints {
			if i == 0 {
				filler.Start(rasterx.ToFixedP(float64(point.X), float64(point.Y)))
			} else {
				filler.Line(rasterx.ToFixedP(float64(point.X), float64(point.Y)))
			}
		}
		filler.Stop(true)
		filler.Draw()
	}

	if pr.polygon.StrokeColor != nil && pr.polygon.StrokeWidth > 0 {
		dasher := rasterx.NewDasher(width, height, scanner)
		dasher.SetColor(pr.polygon.StrokeColor)
		dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
		for i, point := range renderingPoints {
			if i == 0 {
				dasher.Start(rasterx.ToFixedP(float64(point.X), float64(point.Y)))
			} else {
				dasher.Line(rasterx.ToFixedP(float64(point.X), float64(point.Y)))
			}
		}
		dasher.Stop(true)
		dasher.Draw()
	}

	pr.image.Image = raw
	pr.image.Resize(polygonSize)
	pr.image.Move(offsetVector)
	pr.image.Refresh()
}
