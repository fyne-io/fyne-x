package widget

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

// Test various layouts
func TestSegmentedEntry_Layout(t *testing.T) {
	windowSize := fyne.NewSize(400, 200)

	// Test a basic 10 entry layout
	entry := NewSegmentedEntry(AllRunes, 10)
	window := test.NewWindow(entry)
	window.Resize(windowSize)
	test.AssertImageMatches(t, "segmentedentry/basic.png", window.Canvas().Capture())
	window.Close()

	// Test a delimited layout
	entry = NewSegmentedEntry(AllRunes, 5)
	entry.Delimiter = "—"
	window = test.NewWindow(entry)
	window.Resize(windowSize)
	test.AssertImageMatches(t, "segmentedentry/delimited.png", window.Canvas().Capture())
	window.Close()

	// Test a form layout with delimiters at specific indices (SSN example)
	entry = NewSegmentedEntry(Digits, 9)
	entry.Delimiter = "—"
	entry.DelimitAt = []int{2, 4}
	window = test.NewWindow(widget.NewForm(widget.NewFormItem("SSN", entry)))
	window.Resize(windowSize)
	test.AssertImageMatches(t, "segmentedentry/ssn.png", window.Canvas().Capture())
	window.Close()
}

// Test interaction and value retrieval
func TestSegmentedEntry_Interact(t *testing.T) {
	entry := NewSegmentedEntry(AllRunes, 5)
	window := test.NewWindow(entry)
	defer window.Close()

	window.Canvas().Focus(entry)

	// Test text entry
	for i := 0; i < 5; i++ {
		assert.Equal(t, i, entry.selected)
		test.Type(entry, "a")
		expected := i + 1
		expectedText := strings.Repeat("a", expected)
		if i == len(entry.entries)-1 {
			// The last entry shouldn't focus the next one
			expected = i
		}
		assert.Equal(t, expected, entry.selected)
		assert.Equal(t, expectedText, entry.GetValue())
	}

	// Test navigation
	entry.FocusIndex(0)
	window.Canvas().Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Equal(t, 1, entry.selected)
	window.Canvas().Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	assert.Equal(t, 0, entry.selected)

	// Backspace usage
	entry.FocusIndex(4)
	for i := 4; i >= 0; i-- {
		window.Canvas().Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
		expected := i - 0
		expectedText := strings.Repeat("a", expected)
		if i == 0 {
			expected = i
		}
		// Assert the previous entry got focussed
		assert.Equal(t, expected, entry.selected)
		// Assert we lost a character
		assert.Equal(t, expectedText, entry.GetValue())
	}
}
