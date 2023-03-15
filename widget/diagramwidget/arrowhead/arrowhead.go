// Package arrowhead implements an arrowhead canvas object.
package arrowhead

import (
	"image/color"
	"math"

	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

const (
	defaultTheta       float32 = 0.5235
	defaultStrokeWidth float32 = 2
	defaultLength      int     = 15
)

// Arrowhead defines a canvas object which renders an arrow pointing in
// a particular direction.
//
//	        Left
//	          \
//	           \  Length
//	      Theta \
//	Base ------- + Tip
//	            /
//	           /
//	          /
//	        Right
type Arrowhead struct {
	// Base is used to define the "base" of the arrow, which thus defines
	// the direction which the arrow faces.
	Base fyne.Position

	// Tip is the point at which the tip of the arrow will be placed.
	Tip fyne.Position

	// StrokeWidth is the width of the arrowhead lines
	StrokeWidth float32

	// StrokeColor is the color of the arrowhead
	StrokeColor color.Color

	// Theta is the angle between the two "tails" that intersect at the
	// tip. This angle is in radians.
	Theta float32

	// Length is the length of the two "tails" that intersect at the tip.
	Length int

	central *canvas.Line
	left    *canvas.Line
	right   *canvas.Line
	visible bool
}

func MakeArrowhead(base, tip fyne.Position) *Arrowhead {
	return &Arrowhead{
		Base:        base,
		Tip:         tip,
		StrokeWidth: defaultStrokeWidth,
		StrokeColor: theme.TextColor(),
		Theta:       defaultTheta,
		Length:      defaultLength,
		central:     canvas.NewLine(theme.TextColor()),
		left:        canvas.NewLine(theme.TextColor()),
		right:       canvas.NewLine(theme.TextColor()),
		visible:     true,
	}
}

func (a *Arrowhead) Refresh() {
	a.central.StrokeWidth = a.StrokeWidth
	a.left.StrokeWidth = a.StrokeWidth
	a.right.StrokeWidth = a.StrokeWidth

	a.central.StrokeColor = a.StrokeColor
	a.left.StrokeColor = a.StrokeColor
	a.right.StrokeColor = a.StrokeColor

	a.central.Position1 = a.Tip
	a.central.Position2 = a.Base

	a.left.Position1 = a.Tip
	a.left.Position2 = a.LeftPoint()

	a.right.Position1 = a.Tip
	a.right.Position2 = a.RightPoint()

	if a.visible {
		a.central.Show()
		a.left.Show()
		a.right.Show()
	} else {
		a.central.Hide()
		a.left.Hide()
		a.right.Hide()
	}

	canvas.Refresh(a.central)
	canvas.Refresh(a.left)
	canvas.Refresh(a.right)

}

func (a *Arrowhead) LeftPoint() fyne.Position {
	// Have to change the sign of Y because the window coordinated Y axis goes down rather than up
	baseVector := r2.Vec2{
		X: float64(a.Base.X - a.Tip.X),
		Y: -float64(a.Base.Y - a.Tip.Y),
	}
	baseAngle := baseVector.Angle()
	leftAngle := r2.AddAngles(baseAngle, -float64(a.Theta))
	// We have to change the sign of Y because the window coordinate Y axis goes down rather than up
	leftPosition := fyne.Position{
		X: float32(float64(a.Length) * math.Cos(leftAngle)),
		Y: -float32(float64(a.Length) * math.Sin(leftAngle)),
	}
	leftPoint := a.Tip.Add(leftPosition)
	return leftPoint
}

func (a *Arrowhead) RightPoint() fyne.Position {
	// Have to change the sign of Y because the window coordinated Y axis goes down rather than up
	baseVector := r2.Vec2{
		X: float64(a.Base.X - a.Tip.X),
		Y: -float64(a.Base.Y - a.Tip.Y),
	}
	baseAngle := baseVector.Angle()
	rightAngle := r2.AddAngles(baseAngle, float64(a.Theta))
	// We have to change the sign of Y because the window coordinate Y axis goes down rather than up
	rightPosition := fyne.Position{
		X: float32(float64(a.Length) * math.Cos(rightAngle)),
		Y: -float32(float64(a.Length) * math.Sin(rightAngle)),
	}
	rightPoint := a.Tip.Add(rightPosition)
	return rightPoint
}

func (a *Arrowhead) Size() fyne.Size {
	lp := a.LeftPoint()
	rp := a.RightPoint()
	points := []r2.Vec2{
		{X: float64(a.Tip.X), Y: float64(a.Tip.Y)},
		{X: float64(a.Base.X), Y: float64(a.Base.Y)},
		{X: float64(lp.X), Y: float64(lp.Y)},
		{X: float64(rp.X), Y: float64(rp.Y)},
	}

	bounding := r2.BoundingBox(points)
	return fyne.Size{
		Width:  float32(bounding.Width()),
		Height: float32(bounding.Height()),
	}
}

func (a *Arrowhead) Resize(s fyne.Size) {
	l := r2.V2(float64(s.Width), float64(s.Height))
	a.Length = int(l.Length())

	tip := r2.V2(float64(a.Tip.X), float64(a.Tip.Y))
	base := r2.V2(float64(a.Base.X), float64(a.Base.Y))
	v := tip.Add(base.Scale(-1))
	v = v.ScaleToLength(l.Length())
	base = tip.Add(v)

	a.Base = fyne.Position{X: float32(base.X), Y: float32(base.Y)}
}

func (a *Arrowhead) Move(p fyne.Position) {
	a.Tip = p

	tip := r2.V2(float64(a.Tip.X), float64(a.Tip.Y))
	base := r2.V2(float64(a.Base.X), float64(a.Base.Y))
	v := tip.Add(base.Scale(-1))
	base = tip.Add(v)

	a.Base = fyne.Position{X: float32(base.X), Y: float32(base.Y)}
}

func (a *Arrowhead) MinSize() fyne.Size {
	return a.Size()
}

func (a *Arrowhead) Visible() bool {
	return a.visible
}

func (a *Arrowhead) Show() {
	a.visible = true
}

func (a *Arrowhead) Hide() {
	a.visible = false
}

func (a *Arrowhead) Position() fyne.Position {
	return a.Tip
}

// temporary hack, don't do this
func (a *Arrowhead) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		a.central,
		a.left,
		a.right,
	}
}
