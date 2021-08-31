package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

// Test a SingleRuneEntry allowing all runes.
func TestSingleRuneEntry_Basic(t *testing.T) {
	entry := NewSingleRuneEntry(AllRunes)
	win := test.NewWindow(entry)
	defer win.Close()

	for _, r := range AllRunes {
		test.Type(entry, string(r))
		// Assert this rune is only one in the entry
		assert.Equal(t, string(r), entry.Text)
	}
}

// Test a SingleRuneEntry only allowing specific runes.
func TestSingleRuneEntry_Filtered(t *testing.T) {
	entry := NewSingleRuneEntry(Lowercase)
	win := test.NewWindow(entry)
	defer win.Close()

	for _, r := range Lowercase {
		test.Type(entry, string(r))
		// Assert this rune is only one in the entry
		assert.Equal(t, string(r), entry.Text)
	}

	// Set the entry to a character that should not change
	test.Type(entry, "a")

	// Entry should ignore all other runes

	for _, r := range Uppercase {
		test.Type(entry, string(r))
		assert.Equal(t, "a", entry.Text)
	}
	for _, r := range Digits {
		test.Type(entry, string(r))
		assert.Equal(t, "a", entry.Text)
	}
	for _, r := range Special {
		test.Type(entry, string(r))
		assert.Equal(t, "a", entry.Text)
	}
}

// Test OnNavigate events
func TestSingleRuneEntry_OnNavigate(t *testing.T) {
	entry := NewSingleRuneEntry(AllRunes)
	win := test.NewWindow(entry)
	defer win.Close()

	win.Canvas().Focus(entry)

	var expected *fyne.KeyEvent = nil

	called := false

	entry.OnNavigate = func(ev *fyne.KeyEvent) {
		called = true
		if expected == nil {
			assert.Equal(t, expected, ev)
			return
		}
		assert.Equal(t, expected.Name, ev.Name)
	}

	test.Type(entry, "a")
	assert.Equal(t, true, called)

	called = false
	expected = &fyne.KeyEvent{Name: fyne.KeyBackspace}
	win.Canvas().Focused().TypedKey(expected)

	// Make sure the text cleared
	assert.Equal(t, "", entry.Text)
	// Should not have navigated because there was text to clear
	assert.Equal(t, false, called)

	win.Canvas().Focused().TypedKey(expected)
	// No text in field should have called to navigate
	assert.Equal(t, true, called)

	// Direction events

	called = false
	expected = &fyne.KeyEvent{Name: fyne.KeyLeft}
	win.Canvas().Focused().TypedKey(expected)
	assert.Equal(t, true, called)

	called = false
	expected = &fyne.KeyEvent{Name: fyne.KeyRight}
	win.Canvas().Focused().TypedKey(expected)
	assert.Equal(t, true, called)
}
