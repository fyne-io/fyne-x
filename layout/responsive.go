package layout

import (
	"fmt"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Layout = (*ResponsiveLayout)(nil)

// ResponsiveLayout is the layout that will adapt objects with the responsive rules. See NewResponsiveLayout
// for details.
type ResponsiveLayout struct{}

// Layout will place the size and place the objects following the configured reponsive rules.
//
// Implements: fyne.Layout
func (resp *ResponsiveLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	// yes, it may happen
	if len(objects) == 0 || objects[0] == nil {
		return
	}

	// For each object, place it at the right position (currentPos) and resize it.
	var (
		currentPos             = fyne.Position{X: 0, Y: 0} // current position
		shouldCount    bool    = true                      // should we count the number of object in the same line ?
		lineHeight     float32                             // to know the height of the current line
		appliedPadding int                                 // number of padding applied in the current line
	)

	for currentIndex, currentObject := range objects {
		if currentObject == nil || !currentObject.Visible() {
			continue
		}

		if shouldCount {
			// calculate the number of object that can be placed in the same line
			// starting from the current object
			appliedPadding = resp.computeElementInLine(objects[currentIndex:], containerSize)
			shouldCount = false
		}

		// resize the object width using the number of applied paddings
		// and place it at the right position
		currentObject.Resize(
			currentObject.Size().Add(fyne.NewSize(
				theme.Padding()/float32(appliedPadding-1),
				0,
			)),
		)
		currentObject.Move(currentPos)

		// next element X position is the current X + the current width + padding
		currentPos.X += currentObject.Size().Width + theme.Padding()

		lineHeight = float32(math.Max(
			float64(lineHeight),
			float64(currentObject.Size().Height),
		))
		// Manage end of line, the next position overflows, so go to next line.
		if currentPos.X >= containerSize.Width {
			currentPos.X = 0                             // back to left
			currentPos.Y += lineHeight + theme.Padding() // move to the next line
			lineHeight = 0
			shouldCount = true
		}
	}
}

// MinSize return the minimum size ot the layout.
//
// Implements: fyne.Layout
func (resp *ResponsiveLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}

	var h, w, maxHeight float32
	currentY := objects[0].Position().Y

	for _, o := range objects {
		if o == nil || !o.Visible() {
			continue
		}
		// the min width is the max of one item in the container
		w = theme.Padding() + float32(math.Max(
			float64(o.MinSize().Width),
			float64(w),
		))

		// but the height is the sum of all items in the container of the same line
		if o.Position().Y == currentY {
			maxHeight = float32(math.Max(
				float64(maxHeight),
				float64(o.MinSize().Height),
			))
		} else {
			h += maxHeight + theme.Padding()
			maxHeight = o.MinSize().Height
			currentY = o.Position().Y
		}
	}
	h += maxHeight
	return fyne.NewSize(w, h)
}

// calcSize calculate the size to apply to the object based on the container size.
func (resp *ResponsiveLayout) calcSize(o fyne.CanvasObject, containerSize fyne.Size) fyne.Size {
	ro, ok := o.(*responsiveWidget)
	if !ok {
		// We now allow non responsive objects to be packed inside a ResponsiveLayout. The
		// size of the object will be 100% of the container.
		ro = Responsive(o).(*responsiveWidget)
	}

	objectSize := o.MinSize()

	// adapt object witdh from the configuration
	width := responsiveBreakpoint(containerSize.Width)
	var factor float32
	if width <= ExtraSmall {
		factor = ro.responsiveConfig[ExtraSmall] // extra small
	} else if width <= Small {
		factor = ro.responsiveConfig[Small] // small
	} else if width <= Medium {
		factor = ro.responsiveConfig[Medium] // medium
	} else if width <= Large {
		factor = ro.responsiveConfig[Large] // large
	} else {
		factor = ro.responsiveConfig[ExtraLarge] // extra large
	}

	// resize the object width using the factor
	objectSize.Width = (factor * containerSize.Width) - theme.Padding()
	return objectSize
}

// computeElementInLine resize the objects in the same line and return the number of object contained in the line.
func (resp *ResponsiveLayout) computeElementInLine(objects []fyne.CanvasObject, containerSize fyne.Size) int {
	var lineWidth float32
	count := 1
	for _, o := range objects {
		size := resp.calcSize(o, containerSize)
		o.Resize(size)
		lineWidth += size.Width + theme.Padding()
		if lineWidth > containerSize.Width {
			break
		}
		count++
	}
	return count
}

// NewResponsiveLayout return a responsive layout that will adapt objects with the responsive rules. To
// configure the rule, each object could be encapsulated by a "Responsive" object.
//
// Example:
//
//	container := NewResponsiveLayout(
//	    Responsive(label, 1, .5, .25),  // 100% for small, 50% for medium, 25% for large
//	    Responsive(button, 1, .5, .25), // ...
//	    label2,                         // this will be placed and resized with default behaviors
//	                                    // => 1, 1, 1, 1
//	)
func NewResponsiveLayout() fyne.Layout {
	return &ResponsiveLayout{}
}

var _ fyne.Widget = (*responsiveWidget)(nil)

type responsiveWidget struct {
	widget.BaseWidget

	render           fyne.CanvasObject
	responsiveConfig responsiveConfig
}

func (ro *responsiveWidget) CreateRenderer() fyne.WidgetRenderer {
	if ro.render == nil {
		return nil
	}
	return widget.NewSimpleRenderer(ro.render)
}

// Responsive register the object with a responsive configuration.
// The optional ratios must
// be 0 < ratio <= 1 and  passed in this order:
//
//	Responsive(object, smallRatio, mediumRatio, largeRatio, xlargeRatio)
//
// They are set to previous value if a value is not passed, or 1.0 if there is no previous value.
// The returned object is not modified.
func Responsive(object fyne.CanvasObject, breakpointRatio ...float32) fyne.CanvasObject {
	ro := &responsiveWidget{render: object, responsiveConfig: newResponsiveConf(breakpointRatio...)}
	ro.ExtendBaseWidget(ro)
	return ro
}

// responsiveBreakpoint is a integer representing a breakpoint size as defined in Bootstrap.
//
// See: https://getbootstrap.com/docs/4.0/layout/overview/#responsive-breakpoints
type responsiveBreakpoint float32

const (
	// ExtraSmall is below the smallest breakpoint (mobile vertical).
	ExtraSmall responsiveBreakpoint = 576

	// Small is the smallest breakpoint (mobile vertical).
	Small responsiveBreakpoint = 768

	// Medium is the medium breakpoint (mobile horizontal, tablet vertical).
	Medium responsiveBreakpoint = 992

	// Large is the largest breakpoint (tablet horizontal, small desktop).
	Large responsiveBreakpoint = 1200

	// ExtraLarge is the largest breakpoint (large desktop).
	ExtraLarge responsiveBreakpoint = Large + 1

	// XS is an alias for ExtraSmall
	XS responsiveBreakpoint = ExtraSmall

	// SM is an alias for Small
	SM responsiveBreakpoint = Small

	// MD is an alias for Medium
	MD responsiveBreakpoint = Medium

	// LG is an alias for Large
	LG responsiveBreakpoint = Large

	// XL is an alias for ExtraLarge
	XL responsiveBreakpoint = ExtraLarge
)

// ResponsiveConfiguration is the configuration for a responsive object. It's
// a simple map from the breakpoint to the size ratio from it's container.
// Breakpoint is a uint16 that should be set from const SMALL, MEDIUM, LARGE and XLARGE.
type responsiveConfig map[responsiveBreakpoint]float32

// newResponsiveConf return a new responsive configuration.
// The optional ratios must
// be 0 < ratio <= 1 and  passed in this order:
//
//	Responsive(object, smallRatio, mediumRatio, largeRatio, xlargeRatio)
//
// They are set to previous value if a value is not passed, or 1.0 if there is no previous value.
func newResponsiveConf(ratios ...float32) responsiveConfig {
	if len(ratios) > 5 {
		err := fmt.Errorf("responsive: you declared more than 5 ratios, only the first 5 will be used")
		fyne.LogError("Too many responsive ratios", err)
	}

	responsive := responsiveConfig{}

	// basic check
	for _, i := range ratios {
		if i <= 0 || i > 1 {
			message := "Responsive: size must be > 0 and <= 1, got: %f"
			panic(fmt.Errorf(message, i))
		}
	}

	// Set default values
	for index, bp := range []responsiveBreakpoint{ExtraSmall, Small, Medium, Large, ExtraLarge} {
		if len(ratios) <= index {
			if index == 0 {
				ratios = append(ratios, 1)
			} else {
				ratios = append(ratios, ratios[index-1])
			}
		}
		responsive[bp] = ratios[index]
	}
	return responsive
}
