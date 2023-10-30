package layout

import (
	"fmt"
	"log"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
// Example:
//    layout := NewResponsiveLayout(
//         Responsive(label, 1, .5, .25, .5), // small, medium, large, xlarge ratio
//    }
// Note that a responsive layout can handle others layouts, responsive or not.

// responsiveBreakpoint is a integer representing a breakpoint size as defined in Bootstrap.
//
// See: https://getbootstrap.com/docs/4.0/layout/overview/#responsive-breakpoints
type responsiveBreakpoint = float32

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

	// this will be updatad for each element to know where to place
	// the next object.
	pos := fyne.NewPos(theme.Padding(), 0)

	// to calculate the next pos.Y when a new line is needed
	var maxHeight float32

	// objects in a line
	line := []fyne.CanvasObject{}

	// For each object, place it at the right position (pos) and resize it.
	for _, o := range objects {
		if o == nil || !o.Visible() {
			continue
		}

		// get tht configuration
		ro, ok := o.(*responsiveWidget)
		if !ok {
			// We now allow non responsive objects to be packed inside a ResponsiveLayout. The
			// size of the object will be 100% of the container.
			ro = Responsive(o).(*responsiveWidget)
		}
		conf := ro.responsiveConfig

		line = append(line, o) // add the container to the line
		size := o.MinSize()    // get some informations

		// adapt object witdh from the configuration
		if containerSize.Width <= SMALL {
			size.Width = conf[SMALL] * containerSize.Width
		} else if containerSize.Width <= MEDIUM {
			size.Width = conf[MEDIUM] * containerSize.Width
		} else if containerSize.Width <= LARGE {
			size.Width = conf[LARGE] * containerSize.Width
		} else {
			size.Width = conf[XLARGE] * containerSize.Width
		}

		// place and resize the element
		size = size.Subtract(fyne.NewSize(theme.Padding(), 0))
		o.Resize(size)
		o.Move(pos)

		// next element X position
		pos = pos.Add(fyne.NewPos(size.Width+theme.Padding(), 0))

		maxHeight = resp.maxFloat32(maxHeight, size.Height)

		// Manage end of line, the next position overflows, so go to next line.
		if pos.X >= containerSize.Width {
			line = []fyne.CanvasObject{}
			pos.X = theme.Padding()              // back to left
			pos.Y += maxHeight + theme.Padding() // move to the next line
			maxHeight = 0
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

// math.Max only works with float64, so let's make our own
func (resp *ResponsiveLayout) maxFloat32(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
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
//	                                    // => 1, 1, 1
//	)
func NewResponsiveLayout(o ...fyne.CanvasObject) *fyne.Container {
	r := &ResponsiveLayout{}

	objects := []fyne.CanvasObject{}
	for _, unknowObject := range o {
		if _, ok := unknowObject.(*responsiveWidget); !ok {
			unknowObject = Responsive(unknowObject)
		}
		objects = append(objects, unknowObject)
	}

	return container.New(r, objects...)
}

type responsiveWidget struct {
	widget.BaseWidget

	render           fyne.CanvasObject
	responsiveConfig responsiveConfig
}

var _ fyne.Widget = (*responsiveWidget)(nil)

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

func (ro *responsiveWidget) CreateRenderer() fyne.WidgetRenderer {
	if ro.render == nil {
		return nil
	}
	return widget.NewSimpleRenderer(ro.render)
}
