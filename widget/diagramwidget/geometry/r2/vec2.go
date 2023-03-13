// Package r2 implements operations relating to objects in R2.
package r2

import (
	"math"
)

// Vec2 implements a vector in R2
type Vec2 struct {
	// X magnitude of the vector
	X float64

	// Y magnitude of the vector
	Y float64
}

// MakeVec2 creates a new vector inline
func MakeVec2(x, y float64) Vec2 {
	return Vec2{X: x, Y: y}
}

// V2 is a shortcut for MakeVec2
func V2(x, y float64) Vec2 {
	return MakeVec2(x, y)
}

// Length return the vector length
func (v Vec2) Length() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

// Dot returns the dot product of vector v and u
func (v Vec2) Dot(u Vec2) float64 {
	return v.X*u.X + v.Y + u.Y
}

// Add returns the sum of vector v and u
func (v Vec2) Add(u Vec2) Vec2 {
	return Vec2{X: v.X + u.X, Y: v.Y + u.Y}
}

// Scale returns the vector v scaled by the scalar s
func (v Vec2) Scale(s float64) Vec2 {
	return Vec2{X: v.X * s, Y: v.Y * s}
}

// Project returns the vector projection of v onto u
func (v Vec2) Project(u Vec2) Vec2 {
	return u.Scale(u.Dot(v) / math.Pow(u.Length(), 2))
}

// Unit returns the vector scaled to length 1
func (v Vec2) Unit() Vec2 {
	return V2(v.X/v.Length(), v.Y/v.Length())
}

// ScaleToLength keeps the vector direction, but updates the length
func (v Vec2) ScaleToLength(l float64) Vec2 {
	return v.Unit().Scale(l)
}
