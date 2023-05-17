package diagramwidget

import (
	"image/color"

	"fyne.io/fyne/v2"
)

const (
	defaultStrokeWidth float32 = 1
)

// Decoration is a widget intended to be used as a decoration on a Link widget
// The graphical representation of the widget is defined along a reference axis with
// one point on that axis designated as the reference point (generally the origin).
// Depending on the Link widget's use of the decoration, the reference point will either
// be aligned with one of the endpoints of the link or with some intermediate point on the
// link. The Link will move the Decoration's reference point as the link itself is modified.
// The Link will also determine the slope of the Link's line at the reference point and
// direct the Decoration to rotate about the reference point to achieve the correct alignment
// of the decoration with respect to the Link's line.
// The Link may have more than one decoration stacked along the line at the reference point.
// To accomplish this, it needs to know the length of the decoration along the reference axis
// so that it can adjust the position of the next decoration appropriately.
type Decoration interface {
	fyne.Widget
	setLink(link *BaseDiagramLink)
	// setBaseAngle sets the angle of the reference axis
	setBaseAngle(angle float64) // Angle in radians
	SetFillColor(color color.Color)
	// SetSolid determines whether the stroke color is used to fill the decoration
	// It has no impact if the decoration is open
	SetSolid(bool)
	// SetStrokeColor sets the color to be used for lines in the decoration
	SetStrokeColor(color color.Color)
	// SetStrokeWidth sets the width of the lines to be used in the decoration
	SetStrokeWidth(width float32)
	// GetReferenceLength returns the length of the decoration along the reference axis
	GetReferenceLength() float32
}
