package widget

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// LimitedWidthEntry is an extended entry that sets a minimum width, and
// limits the entry/pasting to a given number of characters `CharsWide`
type LimitedWidthEntry struct {
	widget.Entry
	CharsWide int
}

// LimitedWidthEntry returns an extended entry that has a minimum width and is limited
// to a CharsWide number of characters
func NewLimitedWidthEntry() *LimitedWidthEntry {
	entry := &LimitedWidthEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

// MinSize returns the size that this widget should not shrink below.
//
// Implements: fyne.Widget
func (e *LimitedWidthEntry) MinSize() fyne.Size {
	min := e.Entry.MinSize()
	if e.CharsWide > 0 {
		width := fyne.MeasureText(strings.Repeat("W", e.CharsWide), theme.TextSize(), e.TextStyle).Width
		if e.Validator != nil {
			width += theme.IconInlineSize() + theme.Padding()
		}
		if min.Width < width {
			min.Width = width
		}
	}
	return min
}

// TypedRune is called when this item receives a char event.
//
// Implements: fyne.Focusable
func (e *LimitedWidthEntry) TypedRune(r rune) {
	if e.CharsWide < 1 || len(e.Text) < e.CharsWide {
		e.Entry.TypedRune(r)
		return
	}
}

// TypedShortcut handles the registered shortcuts.
//
// Implements: fyne.Shortcutable
func (e *LimitedWidthEntry) TypedShortcut(shortcut fyne.Shortcut) {
	e.Entry.TypedShortcut(shortcut)
	if e.CharsWide > 0 {
		// Limit the text length after the paste as the paste might
		// have inserted or appended to existing content
		if len(e.Text) > e.CharsWide {
			e.SetText(e.Text[:e.CharsWide])
		}
	}
}
