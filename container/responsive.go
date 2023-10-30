package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/x/fyne/layout"
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
// The number of sizes can be up to 4, for small, medium, large and extra large.
// If more than 4 sizes are provided, the extra sizes are ignored.
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
