package widget

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

// NumericalEntry is an extended entry that only allows numerical input.
// Only integers are allowed by default. Support for floats can be enabled by setting AllowFloat.
type NumericalEntry struct {
	widget.Entry
	AllowFloat bool
	// AllowNegative determines if negative numbers can be entered.
	AllowNegative bool
}

// NewNumericalEntry returns an extended entry that only allows numerical input.
func NewNumericalEntry() *NumericalEntry {
	entry := &NumericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

// TypedRune is called when this item receives a char event.
//
// Implements: fyne.Focusable
func (e *NumericalEntry) TypedRune(r rune) {
	if e.Entry.CursorColumn == 0 && e.Entry.CursorRow == 0 {
		if e.AllowNegative {
			if len(e.Text) > 0 && e.Text[0] == '-' {
				return
			} else if r == '-' {
				e.Entry.TypedRune(r)
				return
			}
		}
	}

	if r >= '0' && r <= '9' {
		e.Entry.TypedRune(r)
		return
	}

	if e.AllowFloat && (r == '.' || r == ',') {
		e.Entry.TypedRune(r)
		return
	}
}

// TypedShortcut handles the registered shortcuts.
//
// Implements: fyne.Shortcutable
func (e *NumericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	if e.isNumber(paste.Clipboard.Content()) {
		e.Entry.TypedShortcut(shortcut)
	}
}

// Keyboard sets up the right keyboard to use on mobile.
//
// Implements: mobile.Keyboardable
func (e *NumericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

func (e *NumericalEntry) isNumber(content string) bool {
	if e.AllowFloat {
		_, err := strconv.ParseFloat(content, 64)
		return err == nil
	}

	_, err := strconv.Atoi(content)
	return err == nil
}
