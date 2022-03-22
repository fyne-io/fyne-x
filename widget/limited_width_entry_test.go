package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

// Check typed chars are limited
func TestLimitedWidthEntry_TypedRune(t *testing.T) {
	entry := NewLimitedWidthEntry()

	win := test.NewWindow(entry)
	win.Resize(fyne.NewSize(200, 100))
	defer win.Close()

	for i := 1; i < 20; i++ {
		// reset the widget for `i` chars
		entry.CharsWide = i
		entry.SetText("")

		// Get the width before we add runes
		width := entry.Size().Width

		for j := 0; j < 30; j++ {
			entry.TypedRune(rune('W' + j))
			assert.LessOrEqual(t, len(entry.Text), i)
			assert.EqualValues(t, width, entry.Size().Width)
		}
	}
}
