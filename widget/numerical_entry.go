package widget

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"

	"github.com/cloudfoundry-attic/jibber_jabber"
)

// NumericalEntry is an extended entry that only allows numerical input.
// Only integers are allowed by default. Support for floats can be enabled by setting AllowFloat.
type NumericalEntry struct {
	widget.Entry
	AllowFloat bool
	// AllowNegative determines if negative numbers can be entered.
	AllowNegative bool
	minus         rune
	radixSep      rune
	thouSep       rune
}

// NewNumericalEntry returns an extended entry that only allows numerical input.
func NewNumericalEntry() *NumericalEntry {
	entry := &NumericalEntry{}
	userLocale, err := jibber_jabber.DetectIETF()
	if err != nil {
		fyne.LogError("DetectIETF error: %s\n", err)
	} else {
		entry.minus, entry.radixSep, entry.thouSep = minusRadixThou(userLocale)
	}
	entry.ExtendBaseWidget(entry)
	entry.Validator = entry.ValidateText
	return entry
}

// Append appends text to the entry, filtering out non-numerical characters
// based on the current locale and allowed input types (negative, float).
func (e *NumericalEntry) Append(text string) {
	var s strings.Builder
	for _, r := range text {
		rn, ok := e.getRuneForLocale(r)
		if !ok {
			continue
		}
		if rn == e.minus {
			if !e.AllowNegative || !(len(e.Text) == 0) {
				continue
			}
		}
		if rn == e.radixSep {
			if !e.AllowFloat {
				continue

			}

		}
		s.WriteRune(rn)
	}
	e.Entry.Append(s.String())
}

// SetText manually sets the text of the Entry.
// The text will be filtered to allow only numerical input
// according to the current locale.
func (e *NumericalEntry) SetText(text string) {
	var s strings.Builder
	for _, r := range text {
		rn, ok := e.getRuneForLocale(r)
		if !ok {
			continue
		}
		if rn == e.minus {
			if !e.AllowNegative || !(len(s.String()) == 0) {
				continue
			}
		}
		if rn == e.radixSep {
			if !e.AllowFloat {
				continue
			}
		}
		s.WriteRune(rn)
	}
	e.Entry.SetText(s.String())
}

// TypedRune is called when this item receives a char event.
//
// Implements: fyne.Focusable
func (e *NumericalEntry) TypedRune(r rune) {
	rn, ok := e.getRuneForLocale(r)
	if !ok {
		return
	}
	if e.Entry.CursorColumn == 0 && e.Entry.CursorRow == 0 {
		if e.AllowNegative {
			if len(e.Text) > 0 && []rune(e.Text)[0] == e.minus {
				return
			} else if rn == e.minus {
				e.Entry.TypedRune(rn)
				return
			}
		}
	}

	if unicode.IsDigit(rn) {
		e.Entry.TypedRune(rn)
		return
	}

	if e.AllowFloat && rn == e.radixSep {
		e.Entry.TypedRune(rn)
		return
	}

	if rn == e.thouSep {
		e.Entry.TypedRune(rn)
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

// getRuneForLocale checks if a rune is valid for the entry,
// and returns the correct rune for the locale.
func (e *NumericalEntry) getRuneForLocale(r rune) (rune, bool) {
	if unicode.IsDigit(r) {
		return r, true
	}

	switch r {
	case '-': // hyphen - minus
		fallthrough
	case 0x2212: //mathematical minus
		if e.AllowNegative {
			return e.minus, true
		} else {
			return 0, false
		}
	case '.': // full stop
		fallthrough
	case ',': // comma
		if r == e.radixSep || r == e.thouSep {
			return r, true
		}
	case ' ': // space
		fallthrough
	case 0xa0: // non-breaking space
		if e.thouSep == ' ' || e.thouSep == 0xa0 {
			return e.thouSep, true
		}
	case 0x2019: // right single quote mark
		if r == e.thouSep {
			return r, true
		}
	}
	return 0, false
}

// ValidateText checks if the entered text is a valid numerical input
// according to the system locale.
func (e *NumericalEntry) ValidateText(text string) error {
	if len(text) == 0 {
		return nil
	}
	runes := []rune(text)
	if !e.AllowNegative && runes[0] == e.minus {
		return errors.New("negative numbers are not allowed")
	}
	if !e.AllowFloat && strings.Contains(text, string(e.radixSep)) {
		return errors.New("floating point numbers are not allowed")
	}
	radixCount := 0
	for i, r := range runes {
		if unicode.IsDigit(r) {
			continue
		}
		if r != e.minus && r != e.radixSep && r != e.thouSep {
			return fmt.Errorf("invalid character: %q", r)
		}
		if r == e.minus && i != 0 {
			return errors.New("minus must be the first character")
		}
		if r == e.radixSep {
			radixCount++
			if radixCount > 1 {
				return errors.New("only one radix separator is allowed")
			}
			if i > 0 && runes[i-1] == e.thouSep {
				return errors.New("thousand separator cannot immediately precede radix separator")
			}
		}
		if r == e.thouSep {
			if i == 0 {
				return errors.New("thousand separator cannot be the first character")
			} else if i == 1 && runes[0] == e.minus {
				return errors.New("thousand separator cannot be the first character after minus")
			} else if runes[i-1] == e.radixSep {
				return errors.New("thousand separator cannot be immediately after radix separator")
			} else if runes[i-1] == e.thouSep {
				return errors.New("thousand separator cannot be immediately after thousand separator")
			}

		}
	}
	return nil
}
