package layout

import (
	"math"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func ensureGlobalConfigIsEmpty(t *testing.T) {
	if len(responsiveConfigurations) != 0 {
		t.Errorf("Expected responsiveConfigurations to be empty, but got %d", len(responsiveConfigurations))
	}
}

// Check if a simple widget is responsive to fill 100% of the layout.
func TestResponsive_SimpleLayout(t *testing.T) {
	defer ensureGlobalConfigIsEmpty(t)
	padding := theme.Padding()
	w, h := float32(SMALL), float32(300)

	// build
	label := widget.NewLabel("Hello World")
	layout := NewResponsiveLayout(label)

	win := test.NewWindow(layout)
	defer win.Close()
	win.Resize(fyne.NewSize(w, h))

	size := layout.Size()

	assert.Equal(t, w-padding*2, size.Width)

}

// Test is a basic responsive layout is correctly configured. This test 2 widgets
// with 100% for small size and 50% for medium size or taller.
func TestResponsive_Responsive(t *testing.T) {
	defer ensureGlobalConfigIsEmpty(t)
	padding := theme.Padding()

	// build
	label1 := Responsive(widget.NewLabel("Hello World"), 1, .5)
	label2 := Responsive(widget.NewLabel("Hello World"), 1, .5)

	win := test.NewWindow(
		NewResponsiveLayout(label1, label2),
	)
	win.SetPadded(true)
	defer win.Close()

	// First, we are at w < SMALL so the labels should be sized to 100% of the layout
	w, h := float32(SMALL), float32(300)
	win.Resize(fyne.NewSize(w, h))
	size1 := label1.Size()
	size2 := label2.Size()
	assert.Equal(t, w-padding*2, size1.Width)
	assert.Equal(t, w-padding*2, size2.Width)

	// Then resize to w > SMALL so the labels should be sized to 50% of the layout
	w = float32(MEDIUM)
	win.Resize(fyne.NewSize(w, h))
	size1 = label1.Size()
	size2 = label2.Size()
	assert.Equal(t, w/2-padding, size1.Width)
	assert.Equal(t, w/2-padding, size2.Width)
}

// Check if a widget that overflows the container goes to the next line.
func TestResponsive_GoToNextLine(t *testing.T) {
	defer ensureGlobalConfigIsEmpty(t)
	w, h := float32(200), float32(300)

	// build
	label1 := Responsive(widget.NewLabel("Hello World"), .5)
	label2 := Responsive(widget.NewLabel("Hello World"), .5)
	label3 := Responsive(widget.NewLabel("Hello World"), .5)

	layout := NewResponsiveLayout(label1, label2, label3)
	win := test.NewWindow(layout)
	defer win.Close()
	w = float32(MEDIUM)
	win.Resize(fyne.NewSize(w, h))

	// label 1 and 2 are on the same line
	assert.Equal(t, label1.Position().Y, label2.Position().Y)

	// but not the label 3
	assert.NotEqual(t, label1.Position().Y, label3.Position().Y)

	// just to be sure...
	// the label3 should be at label1.Position().Y + label1.Size().Height
	assert.Equal(t, label1.Position().Y+label1.Size().Height, label3.Position().Y)
}

// Check if sizes are correctly computed for responsive widgets when the window size
// changes. There are some corner cases with padding. So, it needs to be improved...
func TestResponsive_SwitchAllSizes(t *testing.T) {
	defer ensureGlobalConfigIsEmpty(t)

	// build
	n := 4
	labels := make([]fyne.CanvasObject, n)
	for i := 0; i < n; i++ {
		labels[i] = widget.NewLabel("Hello World")
		Responsive(labels[i], 1, 1/float32(2), 1/float32(3), 1/float32(4))
	}
	layout := NewResponsiveLayout(labels...)
	win := test.NewWindow(layout)
	defer win.Close()

	h := float32(1200)
	p := theme.Padding()
	// First, we are at w < SMALL so the labels should be sized to 100% of the layout
	w := float32(SMALL)
	win.Resize(fyne.NewSize(w, h))
	win.Content().Refresh()
	for i := 0; i < n; i++ {
		size := labels[i].Size()
		assert.Equal(t, w-p*2, size.Width)
	}

	// Then resize to w > SMALL so the labels should be sized to 50% of the layout
	w = float32(MEDIUM)
	win.Resize(fyne.NewSize(w, h))
	for i := 0; i < n; i++ {
		size := labels[i].Size()
		assert.Equal(t, w/2-p, size.Width)
	}

	// Then resize to w > MEDIUM so the labels should be sized to 33% of the layout
	w = float32(LARGE)
	win.Resize(fyne.NewSize(w, h))
	for i := 0; i < n; i++ {
		size := labels[i].Size()
		// weird but it "1px" error in calculation
		round := float32(math.Floor(float64(w/3))) - p/2
		assert.Equal(t, round, size.Width)
	}

	// Then resize to w > LARGE so the labels should be sized to 25% of the layout
	w = float32(XLARGE * 2)
	win.Resize(fyne.NewSize(w, h))
	for i := 0; i < n; i++ {
		size := labels[i].Size()
		assert.Equal(t, w/4-p/2, size.Width)
	}
}

// Test if a widget is responsive to fill 100% of the layout
// when we don't provides rsponsive ratios.
func TestResponsive_NoArgs(t *testing.T) {
	defer ensureGlobalConfigIsEmpty(t)
	label := widget.NewLabel("Hello World")
	resp := NewResponsiveLayout(Responsive(label))
	conf := resp.Layout.(*ResponsiveLayout).configuration
	for _, s := range []ResponsiveBreakpoint{SMALL, MEDIUM, LARGE, XLARGE} {
		assert.Equal(t, float32(1), conf[label][s])
	}
}

// Check that we manage some errors.
func TestResponsive_Errors(t *testing.T) {
	defer ensureGlobalConfigIsEmpty(t)

	// assert that a Responsive with size > 1 panic
	assert.Panics(t, func() {
		Responsive(widget.NewLabel("Hello World"), 1.2)
	})
	assert.Panics(t, func() {
		Responsive(widget.NewLabel("Hello World"), 1, 1.2)
	})

	// assert that GetResponsiveConfig return nil, err != nil
	// for a non-responsive widget
	resp := &ResponsiveLayout{}
	label := widget.NewLabel("Hello World")
	o, err := resp.GetResponsiveConfig(label)
	assert.NotNil(t, err)
	assert.Nil(t, o)

}
