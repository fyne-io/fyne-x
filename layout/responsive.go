package layout

import (
	"fmt"
	"log"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// responsive layout provides a fyne.Layout that is responsive to the window size.
// All fyne.CanvasObject are resized and positionned following the rules you decide.
//
// It is strongly inspired by Bootstrap's grid system. But instead of using 12 columns,
// we use width ratio.
// By default, a standard fyne.CanvasObject will always be width to 1 * containerSize and place in vertical.
// If you want to change the behavior, you can use Responsive() function that registers the layout configuration.
//
// To avoid using fractions, some constants are defined like Half, OneThird, OneQuarter, OneFifth and OneSixth.
//
// Responsive() function takes a fyne.CanvasObject and a list of ratios. The number of ratios could be 1, 2, 3 or 4.
//
// Example:
//    layout := NewResponsiveLayout(
//         Responsive(label, 1, .5, .25, .5), // small, medium, large, xlarge ratio
//    }
// Note that a responsive layout can handle others layouts, responsive or not.

// responsiveBreakpoint is a integer representing a breakpoint size as defined in Bootstrap.
//
// See: https://getbootstrap.com/docs/4.0/layout/overview/#responsive-breakpoints
type responsiveBreakpoint float32

const (
	// Full is the full size of the container.
	Full float32 = 1.0
	// Half is half the size of the container.
	Half float32 = 0.5
	// OneThird is one third the size of the container.
	OneThird float32 = 1.0 / 3.0
	// TwoThird is two third the size of the container.
	TwoThird float32 = 2.0 / 3.0
	// OneQuarter is one quarter the size of the container.
	OneQuarter float32 = 0.25
	// OneFifth is five twelfths the size of the container.
	OneFifth float32 = 0.2
	// OneSixth is one sixth the size of the container.
	OneSixth float32 = 1.0 / 6.0
)

const (
	// SMALL is the smallest breakpoint (mobile vertical).
	SMALL responsiveBreakpoint = 576

	// MEDIUM is the medium breakpoint (mobile horizontal, tablet vertical).
	MEDIUM responsiveBreakpoint = 768

	// LARGE is the largest breakpoint (tablet horizontal, small desktop).
	LARGE responsiveBreakpoint = 992

	// XLARGE is the largest breakpoint (large desktop).
	XLARGE responsiveBreakpoint = 1200

	// SM is an alias for SMALL
	SM responsiveBreakpoint = SMALL

	// MD is an alias for MEDIUM
	MD responsiveBreakpoint = MEDIUM

	// LG is an alias for LARGE
	LG responsiveBreakpoint = LARGE

	// XL is an alias for XLARGE
	XL responsiveBreakpoint = XLARGE
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
	responsive := responsiveConfig{}

	if len(ratios) > 4 {
		log.Println("Responsive: you declared more than 4 ratios, only the first 4 will be used")
	}

	// basic check
	for _, i := range ratios {
		if i <= 0 || i > 1 {
			message := "Responsive: size must be > 0 and <= 1, got: %f"
			panic(fmt.Errorf(message, i))
		}
	}

	// Set default values
	for index, bp := range []responsiveBreakpoint{SMALL, MEDIUM, LARGE, XLARGE} {
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
		currentPos = currentPos.Add(fyne.NewPos(
			currentObject.Size().Width+theme.Padding(), 0,
		))

		lineHeight = resp.maxFloat32(lineHeight, currentObject.Size().Height)
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
		w = resp.maxFloat32(o.MinSize().Width, w) + theme.Padding()
		if o.Position().Y != currentY {
			currentY = o.Position().Y
			// new line, so we can add the maxHeight to h
			h += maxHeight + theme.Padding()

			// drop the line
			maxHeight = 0
		}
		maxHeight = resp.maxFloat32(maxHeight, o.MinSize().Height)
	}
	h += maxHeight + theme.Padding()
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
	if width <= SMALL {
		factor = ro.responsiveConfig[SMALL]
	} else if width <= MEDIUM {
		factor = ro.responsiveConfig[MEDIUM]
	} else if width <= LARGE {
		factor = ro.responsiveConfig[LARGE]
	} else {
		factor = ro.responsiveConfig[XLARGE]
	}

	objectSize.Width = factor * containerSize.Width
	// set the size (without padding adaptation)
	objectSize = objectSize.Subtract(fyne.NewSize(theme.Padding(), 0))
	return objectSize
}

// math.Max only works with float64, so let's make our own
func (resp *ResponsiveLayout) maxFloat32(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
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
	r := &ResponsiveLayout{}

	//objects := []fyne.CanvasObject{}
	//for _, unknowObject := range o {
	//	if _, ok := unknowObject.(*responsiveWidget); !ok {
	//		unknowObject = Responsive(unknowObject)
	//	}
	//	objects = append(objects, unknowObject)
	//}

	return r
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
	if len(breakpointRatio) > 4 {
		fyne.LogError(
			"Too many arguments in Responsive()",
			fmt.Errorf("the function can take at most 4 arguments, you provided %d", len(breakpointRatio)),
		)
	}
	ro := &responsiveWidget{render: object, responsiveConfig: newResponsiveConf(breakpointRatio...)}
	ro.ExtendBaseWidget(ro)
	return ro
}
