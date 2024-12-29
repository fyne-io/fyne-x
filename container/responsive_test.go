package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

// This is a simple test to check if the responsive layout is correctly configured.
// The others tests are in the layout package.

// Check if a simple widget is responsive to fill 100% of the layout.
func TestResponsive_Responsive(t *testing.T) {
	responsive := NewResponsive(
		widget.NewLabel("Hello World"),                    // 100% by default
		Responsive(widget.NewLabel("Hello World"), 1, .5), // 100% for small, 50% for others
		Responsive(widget.NewLabel("Hello World"), 1, .5), // 100% for small, 50% for others
	)
	win := test.NewWindow(responsive)
	win.Resize(fyne.NewSize(320, 480))
	size1 := responsive.Objects[1].Size()
	size2 := responsive.Objects[2].Size()
	assert.Equal(t, size1.Width, size2.Width)
}
