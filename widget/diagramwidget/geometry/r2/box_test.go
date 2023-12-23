package r2

import (
	"testing"
)

func TestContains(t *testing.T) {
	var top, bottom, left, right float64
	top = 100
	left = 100
	right = 200
	bottom = 200
	upperLeft := MakeVec2(left, top)
	sVector := MakeVec2(right-left, bottom-top)
	box := MakeBox(upperLeft, sVector)
	// test point not being contained
	cp := MakeVec2(50, 50)
	if box.Contains(cp) {
		// this should have returned false
		t.Errorf("Contains returned true for a point not in the box")
	}
	cp.X = 150
	cp.Y = 150
	if !box.Contains(cp) {
		// this should have returned true
		t.Errorf("Contains returned false for a point in the box")
	}
}

func TestFindPerimeterPointNearestContainedPoint(t *testing.T) {
	var top, bottom, left, right float64
	top = 100
	left = 100
	right = 200
	bottom = 200
	upperLeft := MakeVec2(left, top)
	sVector := MakeVec2(right-left, bottom-top)
	box := MakeBox(upperLeft, sVector)
	// test point not being contained
	cp := MakeVec2(50, 50)
	result := box.FindPerimeterPointNearestContainedPoint(cp)
	if result.X != 0 || result.Y != 0 {
		t.Errorf("Non-contained point did not return 0,0, got %f, %f", result.X, result.Y)
	}
	// test top
	cp.X = 140
	cp.Y = 125
	result = box.FindPerimeterPointNearestContainedPoint(cp)
	if result.X != 140 || result.Y != top {
		t.Errorf("Point on top not returned, expected 140, 100, got %f, %f", result.X, result.Y)
	}
	// test left
	cp.X = 125
	cp.Y = 140
	result = box.FindPerimeterPointNearestContainedPoint(cp)
	if result.X != left || result.Y != 140 {
		t.Errorf("Point on left not returned, expected 100, 140, got %f, %f", result.X, result.Y)
	}
	// test bottom
	cp.X = 140
	cp.Y = 175
	result = box.FindPerimeterPointNearestContainedPoint(cp)
	if result.X != 140 || result.Y != bottom {
		t.Errorf("Point on bottom not returned, expected 140, 200, got %f, %f", result.X, result.Y)
	}
	// test right
	cp.X = 175
	cp.Y = 140
	result = box.FindPerimeterPointNearestContainedPoint(cp)
	if result.X != right || result.Y != 140 {
		t.Errorf("Point on right not returned, expected 200, 140, got %f, %f", result.X, result.Y)
	}
}
