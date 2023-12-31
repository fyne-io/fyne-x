package r2

// Box defines a box in R2
//
//	            A
//	            |
//	            |
//	            |
//	            v
//	(1) A.X,A.Y +------+ A.X+S.X,A.Y (2)
//	            |\     |
//	            | \    |
//	            |  \   |
//	            |   \S |
//	            |    \ |
//	            |     \|
//
// (3)A.X,A.Y+S.Y +------+ A.X+S.X,A.Y+S.Y (4)
type Box struct {

	// A defines the top-left corner of the box
	A Vec2

	// S defines the size of the box
	S Vec2
}

// MakeBox creates an r2 Box
func MakeBox(a, s Vec2) Box {
	return Box{
		A: a,
		S: s,
	}
}

// Area returns the area of the Box
func (b Box) Area() float64 {
	return b.S.X * b.S.Y
}

// FindPerimeterPointNearestContainedPoint returns the perimiter point closest to the contained point.
// If the point is not actually within the Box, it returns a (0,0) vector
func (b Box) FindPerimeterPointNearestContainedPoint(containedPoint Vec2) Vec2 {
	if !b.Contains(containedPoint) {
		return MakeVec2(0, 0)
	}
	top := b.GetCorner1().Y
	left := b.GetCorner1().X
	bottom := b.GetCorner4().Y
	right := b.GetCorner4().X
	topDistance := containedPoint.Y - top
	leftDistance := containedPoint.X - left
	bottomDistance := bottom - containedPoint.Y
	rightDistance := right - containedPoint.X
	if bottomDistance > topDistance {
		// top is closer
		if rightDistance > leftDistance {
			// left is closer
			if leftDistance > topDistance {
				// top is the closest
				return MakeVec2(containedPoint.X, top)
			}
			// left is the closest
			return MakeVec2(left, containedPoint.Y)
		}
		// right is closer
		if rightDistance > topDistance {
			// top is the closest
			return MakeVec2(containedPoint.X, top)
		}
		// right is the closest
		return MakeVec2(right, containedPoint.Y)
	}
	// bottom is closer
	if rightDistance > leftDistance {
		// left is closer
		if leftDistance > bottomDistance {
			// bottom is the closest
			return MakeVec2(containedPoint.X, bottom)
		}
		// left is the closest
		return MakeVec2(left, containedPoint.Y)
	}
	// right is closer
	if rightDistance > bottomDistance {
		// bottom is the closest
		return MakeVec2(containedPoint.X, bottom)
	}
	// right is the closest
	return MakeVec2(right, containedPoint.Y)
}

// GetCorner1 returns the top left corner of the box
func (b Box) GetCorner1() Vec2 {
	return b.A
}

// GetCorner2 returns the top right corner of the box
func (b Box) GetCorner2() Vec2 {
	return b.A.Add(V2(b.S.X, 0))
}

// GetCorner3 returns the bottom left corner of the box.
func (b Box) GetCorner3() Vec2 {
	return b.A.Add(V2(0, b.S.Y))
}

// GetCorner4 returns the bottom right corner of the box.
func (b Box) GetCorner4() Vec2 {
	return b.A.Add(V2(b.S.X, b.S.Y))
}

// Intersect returns the intersection of the box and the line, and a Boolean indicating
// if the box and vector intersect. If they do not collide, the zero vector is
// returned.
func (b Box) Intersect(l Line) (Vec2, bool) {
	// This is transliterated in part from:
	//
	// https://github.com/JulNadeauCA/libagar/blob/master/gui/primitive.c

	faces := []Line{
		b.Top(),
		b.Left(),
		b.Right(),
		b.Bottom(),
	}

	dists := []float64{-1, -1, -1, -1}
	intersects := []bool{false, false, false, false}
	intersectPoints := make([]Vec2, 4)

	shortestDist := float64(-1.0)
	best := -1

	for i := range faces {
		in, ok := IntersectLines(faces[i], l)
		if !ok {
			continue
		}
		dists[i] = in.Length()
		intersects[i] = ok
		intersectPoints[i] = in

		if (dists[i] < shortestDist) || (shortestDist == float64(-1)) {
			shortestDist = dists[i]
			best = i
		}
	}

	if shortestDist < 0 {
		return V2(0, 0), false
	}

	return intersectPoints[best], true
}

// Top returns the top face of the box.
func (b Box) Top() Line {
	return MakeLineFromEndpoints(b.GetCorner1(), b.GetCorner2())
}

// Left returns the left face of the box.
func (b Box) Left() Line {
	return MakeLineFromEndpoints(b.GetCorner1(), b.GetCorner3())
}

// Right returns the right face of the box.
func (b Box) Right() Line {
	return MakeLineFromEndpoints(b.GetCorner2(), b.GetCorner4())
}

// Bottom returns the bottom face of the box.
func (b Box) Bottom() Line {
	return MakeLineFromEndpoints(b.GetCorner3(), b.GetCorner4())
}

// Center returns the center of the Box as an r2 vector
func (b Box) Center() Vec2 {
	return b.A.Add(b.S.Scale(0.5))
}

// Contains returns true if the point v is within the box b.
func (b Box) Contains(v Vec2) bool {
	if (v.X < b.GetCorner1().X) || (v.X > b.GetCorner2().X) {
		return false
	}

	if (v.Y < b.GetCorner1().Y) || (v.Y > b.GetCorner3().Y) {
		return false
	}

	return true
}

// BoundingBox creates a minimum axis-aligned bounding box for the given list
// of points.
func BoundingBox(points []Vec2) Box {
	if len(points) < 2 {
		return MakeBox(V2(0, 0), V2(0, 0))
	}
	var xMin, xMax, yMin, yMax float64
	for i, p := range points {
		if i == 0 {
			xMin = p.X
			xMax = p.X
			yMin = p.Y
			yMax = p.Y
		} else {
			if p.X < xMin {
				xMin = p.X
			}
			if p.Y < yMin {
				yMin = p.Y
			}
			if p.X > xMax {
				xMax = p.X
			}
			if p.Y > yMax {
				yMax = p.Y
			}
		}
	}
	// MakeBox expects the first point to be the upper left, second the bottom right
	return MakeBox(V2(xMin, yMax), V2(xMax-xMin, yMin-yMax))
}

// Width returns the width of the Box
func (b Box) Width() float64 {
	return b.Top().Length()
}

// Height returns the height of the box
func (b Box) Height() float64 {
	return b.Left().Length()
}
