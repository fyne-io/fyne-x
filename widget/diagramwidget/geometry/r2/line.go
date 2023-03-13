package r2

// Line describes a line in R2
//
//  (1) A.X,A.Y  +
//                \
//                 \
//                  \
//                   \
//                    + A.X+S.X,A.Y+S.Y (2)
//
type Line struct {
	// A defines the basis point of the line
	A Vec2

	// S defines the direction and length of the line
	S Vec2
}

func MakeLine(a, s Vec2) Line {
	return Line{
		A: a,
		S: s,
	}
}

// Return a line which has endpoints a, b
func MakeLineFromEndpoints(a, b Vec2) Line {
	s := b.Add(a.Scale(-1))

	return MakeLine(a, s)
}

func (l Line) Endpoint1() Vec2 {
	return l.A
}

func (l Line) Endpoint2() Vec2 {
	return l.A.Add(l.S)
}

func samesign(a, b float64) bool {
	if (a < 0) && (b < 0) {
		return true
	}

	if (a > 0) && (b > 0) {
		return true
	}

	if a == b {
		return true
	}

	return false
}

func (l Line) Length() float64 {
	return l.S.Length()
}

// This code is transliterated from here:
//
// https://github.com/JulNadeauCA/libagar/blob/master/gui/primitive.co
//
// Which is in turn based on Gem I.2 in Graphics Gems II by James Arvo.
func IntersectLines(l1, l2 Line) (Vec2, bool) {
	x1 := l1.Endpoint1().X
	y1 := l1.Endpoint1().Y
	x2 := l1.Endpoint2().X
	y2 := l1.Endpoint2().Y
	x3 := l2.Endpoint1().X
	y3 := l2.Endpoint1().Y
	x4 := l2.Endpoint2().X
	y4 := l2.Endpoint2().Y

	a1 := y2 - y1
	b1 := x1 - x2
	c1 := x2*y1 - x1*y2

	r3 := a1*x3 + b1*y3 + c1
	r4 := a1*x4 + b1*y4 + c1

	if (r3 != 0) && (r4 != 0) && samesign(r3, r4) {
		return V2(0, 0), false
	}

	a2 := y4 - y3
	b2 := x3 - x4
	c2 := x4*y3 - x3*y4

	r1 := a2*x1 + b2*y1 + c2
	r2 := a2*x2 + b2*y2 + c2

	if (r1 != 0) && (r2 != 0) && samesign(r1, r2) {
		return V2(0, 0), false
	}

	denom := a1*b2 - a2*b1
	if denom == 0 {
		return V2(0, 0), false
	}

	offset := 0.0 - denom/2.0
	if denom < 0 {
		offset = denom / 2.0
	}

	num := b1*c2 - b2*c1
	xi := 0.0
	if num < 0 {
		xi = num - offset
	} else {
		xi = num + offset
	}
	xi /= denom

	num = a2*c1 - a1*c2
	yi := 0.0
	if num < 0 {
		yi = num - offset
	} else {
		yi = num + offset
	}
	yi /= denom

	return V2(xi, yi), true

}
