package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/x/fyne/layout"
)

// NewResponsive returns a container with a responsive layout. The objects
// can be copmmon containers or responsive objects using the Responsive()
// function. Note that the content size is computed from the container size and not
// from the window size.
//
// Example:
//
//	NewResponsive(
//		widget.NewLabel("Hello World"), // 100% by default
//		Responsive(widget.NewLabel("Hello World"), 1, .5), // 100% for small, 50% for others
//		Responsive(widget.NewLabel("Hello World"), 1, .5), // 100% for small, 50% for others
//	)
func NewResponsive(objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewResponsiveLayout(), objects...)
}

// Responsive returns a responsive object configured with breakpoint sizes.
// The number of ratios can be up to 5, for extra small, small, medium, large and extra large and above.
// If no size is provided, the object will be 100% of the layout for the whole possible size breakpoints.
// If more than 5 ratios are provided, the extra ratios are ignored.
// If less than 5 raiios are provided, the last ratio is used for the missing sizes.
//
// This sizes are used for the following breakpoints:
//
//   - extra small: 0px to 576px
//   - small: 567px to 767px
//   - medium: 768px to 992 px
//   - large: 992px to 1200
//   - extra large: 1281px and above
//
// Example:
//
//	Responsive(widget.NewLabel("Hello World"), 1, .5)
//
// It's commonly use in a responsive container, like this:
//
//	NewResponsive(
//		Responsive(widget.NewLabel("Hello World"), 1, .5),
//		Responsive(widget.NewLabel("Hello World"), 1, .5, .25),
//	)
//
// Or with the Add() method:
//
//	ctn := NewResponsive()
//	ctn.Add(Responsive(widget.NewLabel("Hello World"), 1, .5))
//	ctn.Add(Responsive(widget.NewLabel("Hello World"), 1, .5))
func Responsive(object fyne.CanvasObject, ratios ...Ratio) fyne.CanvasObject {
	return layout.Responsive(object, ratios...)
}
