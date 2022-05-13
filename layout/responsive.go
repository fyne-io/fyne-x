package layout

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
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

// ResponsiveBreakpoint is a integer representing a breakpoint size as defined in Bootstrap.
//
// See: https://getbootstrap.com/docs/4.0/layout/overview/#responsive-breakpoints
type ResponsiveBreakpoint uint16

const (
	SMALL  ResponsiveBreakpoint = 576
	MEDIUM ResponsiveBreakpoint = 768
	LARGE  ResponsiveBreakpoint = 992
	XLARGE ResponsiveBreakpoint = 1200
	// Aliases
	SM ResponsiveBreakpoint = SMALL
	MD ResponsiveBreakpoint = MEDIUM
	LG ResponsiveBreakpoint = LARGE
	XL ResponsiveBreakpoint = XLARGE
)

// responsiveConfiguration is the configuration for a responsive object. It's
// a simple map from the breakpoint to the size ratio from it's container.
// Breakpoint is a uint16 that should be set from const SMALL, MEDIUM, LARGE and XLARGE.
type responsiveConfiguration map[ResponsiveBreakpoint]float32

// NewResponsiveConf return a new responsive configuration.
// The optional ratios must
// be 0 < ratio <= 1 and  passed in this order:
//      Responsive(object, smallRatio, mediumRatio, largeRatio, xlargeRatio)
// They are set to previous value if a value is not passed, or 1.0 if there is no previous value.
func NewResponsiveConf(ratios ...float32) responsiveConfiguration {
	responsive := responsiveConfiguration{}

	if len(ratios) > 4 {
		log.Println("Responsive: you declared more than 4 ratios, only the first 4 will be used")
	}

	// basic check
	for _, i := range ratios {
		if i <= 0 || i > 1 {
			message := "Responsive: size must be > 0 and <= 1, got: %d"
			panic(errors.New(fmt.Sprintf(message, i)))
		}
	}

	// Set default values
	for index, bp := range []ResponsiveBreakpoint{SMALL, MEDIUM, LARGE, XLARGE} {
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

// temporary register for the responsiveConfigurations in this variable. This
// map is filled by Responsive() function. It's not thread safe, but it should not be a problem.
// SetResponsiveConfig and GetResponsiveConfig() managed the configuration pulling from this variable.
var responsiveConfigurations = make(map[fyne.CanvasObject]responsiveConfiguration)

// ResponsiveLayout is the layout that will adapt objects with the responsive rules. See NewResponsiveLayout
// for details.
type ResponsiveLayout struct {
	configuration map[fyne.CanvasObject]responsiveConfiguration
	mutex         *sync.Mutex
}

// GetResponsiveConfig returns the configuration for the object.
// If the object was not registered, the function returns an error.
func (resp *ResponsiveLayout) GetResponsiveConfig(object fyne.CanvasObject) (responsiveConfiguration, error) {

	// case of the user didn't use the NewResponsiveLayout method, we need to
	// drop global configuration
	if _, ok := responsiveConfigurations[object]; ok {
		resp.SetResponsiveConfig(object, responsiveConfigurations[object])
	}

	if conf, ok := resp.configuration[object]; ok {
		return conf, nil
	}
	return nil, errors.New("object not registered")
}

// Layout will place the size and place the objects following the configured reponsive rules.
//
// Implements: fyne.Layout
func (resp *ResponsiveLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	// yes, it may happen
	if len(objects) == 0 || objects[0] == nil {
		return
	}

	// Responsive is based on the window size, so we need to get it
	window := fyne.CurrentApp().Driver().CanvasForObject(objects[0])
	if window == nil {
		return
	}

	windowSize := window.Size()

	// this will be updatad for each element to know where to place
	// the next object.
	pos := fyne.NewPos(0, 0)

	// to calculate the next pos.Y when a new line is needed
	maxHeight := float32(0)

	// For each object, place it at the right position (pos) and resize it.
	for _, o := range objects {

		if o == nil || !o.Visible() {
			continue
		}

		// get tht configuration
		conf, err := resp.GetResponsiveConfig(o)
		if err != nil {
			// should never happen!
			log.Fatal(err)
		}

		// get some informations
		size := o.MinSize()

		// Time to adapt the size and position of the object based on the window size.
		// The object size is based on containerSize, but the decision is taken from the window size.
		ww := ResponsiveBreakpoint(windowSize.Width)
		if ww <= SMALL {
			size.Width = conf[SMALL] * containerSize.Width
		} else if ww <= MEDIUM {
			size.Width = conf[MEDIUM] * containerSize.Width
		} else if ww <= LARGE {
			size.Width = conf[LARGE] * containerSize.Width
		} else {
			size.Width = conf[XLARGE] * containerSize.Width
		}

		// place and resize the element
		o.Move(pos)
		o.Resize(size)

		// to know where to place the next line
		maxHeight = resp.maxFloat32(maxHeight, size.Height)

		// move the next element to the right
		pos.X += size.Width + theme.Padding()

		// manage end of line, the next position overflows, so go to next line.
		if pos.X >= containerSize.Width {
			pos.X = 0
			pos.Y += maxHeight
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
		if o.Position().Y != currentY {
			currentY = o.Position().Y
			// new line, so we can add the maxHeight to h
			h += maxHeight

			// drop the line
			maxHeight = 0
		}
		maxHeight = resp.maxFloat32(maxHeight, o.MinSize().Height)
	}
	h += maxHeight + theme.Padding()

	return fyne.NewSize(w, h)
}

// SetResponsiveConfig sets the configuration for the object.
// It creates a new configuration if the object was not registered.
func (resp *ResponsiveLayout) SetResponsiveConfig(object fyne.CanvasObject, conf responsiveConfiguration) {
	resp.mutex.Lock()
	defer resp.mutex.Unlock()

	if _, ok := responsiveConfigurations[object]; ok {
		delete(responsiveConfigurations, object)
	}

	resp.configuration[object] = conf
}

// math.Max only works with float64, so let's make our own
func (resp *ResponsiveLayout) maxFloat32(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
}

// NewResponsiveLayout return a responsive layout that will adapt objects with the responsive rules. To
// configure the rule, each object could be encapsulated by a "Responsive" object.
//
// Example:
//      container := NewResponsiveLayout(
//          Responsive(label, 1, .5, .25),  // 100% for small, 50% for medium, 25% for large
//          Responsive(button, 1, .5, .25), // ...
//          label2,                         // this will be placed and resized with default behaviors
//                                          // => 1, 1, 1
//      )
func NewResponsiveLayout(o ...fyne.CanvasObject) *fyne.Container {

	r := &ResponsiveLayout{
		configuration: make(map[fyne.CanvasObject]responsiveConfiguration),
		mutex:         &sync.Mutex{},
	}

	for _, ob := range o {
		if conf, ok := responsiveConfigurations[ob]; !ok {
			r.SetResponsiveConfig(ob, NewResponsiveConf(1, 1, 1, 1))
		} else {
			r.SetResponsiveConfig(ob, conf)
		}
	}

	return fyne.NewContainerWithLayout(r, o...)
}

// Responsive register the object with a responsive configuration.
// The optional ratios must
// be 0 < ratio <= 1 and  passed in this order:
//      Responsive(object, smallRatio, mediumRatio, largeRatio, xlargeRatio)
// They are set to previous value if a value is not passed, or 1.0 if there is no previous value.
// The returned object is not modified.
func Responsive(object fyne.CanvasObject, breakpointRatio ...float32) fyne.CanvasObject {

	// Create configuration for this object
	responsiveConfigurations[object] = NewResponsiveConf(breakpointRatio...)

	// return the original object
	return object
}
