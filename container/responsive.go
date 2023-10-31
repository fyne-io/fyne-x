package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/x/fyne/layout"
)

type FractionHelper = float32

const (
	// Full is the full size of the container.
	Full FractionHelper = 1.0
	// Half is half the size of the container.
	Half FractionHelper = 0.5
	// OneThird is one third the size of the container.
	OneThird FractionHelper = 1.0 / 3.0
	// TwoThird is two third the size of the container.
	TwoThird FractionHelper = 2.0 / 3.0
	// OneQuarter is one quarter the size of the container.
	OneQuarter FractionHelper = 0.25
	// OneFifth is five twelfths the size of the container.
	OneFifth FractionHelper = 0.2
	// OneSixth is one sixth the size of the container.
	OneSixth FractionHelper = 1.0 / 6.0
)

// NewResponsive returns a container with a responsive layout. The objects
// can be copmmon containers or responsive objects using the Responsive()
// function.
//
// Example:
//
//	NewResponsive(
//		widget.NewLabel("Hello World"), // 100% by default
//		Responsive(widget.NewLabel("Hello World"), 1, .5), // 100% for small, 50% for others
//		Responsive(widget.NewLabel("Hello World"), 1, .5), // 100% for small, 50% for others
//	)
func NewResponsive(objects ...fyne.CanvasObject) *fyne.Container {
	container := container.New(layout.NewResponsiveLayout())
	container.Objects = objects
	return container
}

// Responsive returns a responsive object configured with breakpoint sizes.
// If no size is provided, the object will be 100% of the layout.
// The number of sizes can be up to 5, for extra small, small, medium, large and extra large and above.
// If more than 5 sizes are provided, the extra sizes are ignored.
//
// The sizes are used for the following breakpoints:
//
//   - extra small: 0px to 479px
//   - small: 480px to 767px
//   - medium: 768px to 1023px
//   - large: 1024px to 1279px
//   - extra large: 1280px and above
//
// Example:
//
//	Responsive(widget.NewLabel("Hello World"), 1, .5)
//
// It's commonly use in a responsive container, like this:
//
//	NewResponsive(
//		Responsive(widget.NewLabel("Hello World"), 1, .5),
//		Responsive(widget.NewLabel("Hello World"), 1, .5),
//	)
//
// Or with the Add() method:
//
//	ctn := NewResponsive()
//	ctn.Add(Responsive(widget.NewLabel("Hello World"), 1, .5))
//	ctn.Add(Responsive(widget.NewLabel("Hello World"), 1, .5))
//
// This function is a shortcut for layout.Responsive()
func Responsive(object fyne.CanvasObject, sizes ...float32) fyne.CanvasObject {
	return layout.Responsive(object, sizes...)
}
