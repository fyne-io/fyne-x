package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Lowercase represents all lowercase characters
var Lowercase = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

// Uppercase represents all uppercase characters
var Uppercase = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

// Digits represents all digits.
var Digits = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

// Special represents all special characters
var Special = []rune{'`', '~', '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '=', '+', '[', ']', '\\', ';', ':', '\'', '"', ',', '<', '.', '>', '?', '/'}

// Letters includes all lower and upper case letters.
var Letters = append(Lowercase, Uppercase...)

// LettersAndDigits includes all letters and digits.
var LettersAndDigits = append(Letters, Digits...)

// AllRunes contains all possibule runes.
var AllRunes = append(LettersAndDigits, Special...)

// SingleRuneEntry is an entry that only allows a single rune
// from an optional list of allowed runes. It exposes methods
// to respond to navigation events while the entry is in focus.
type SingleRuneEntry struct {
	widget.Entry

	// OnNavigate is called under the following circumstances:
	//
	//  KeyLeft/KeyRight
	//    - Called with the corresponding key
	//
	//  Backspace
	//    - Called with the key if the length of this entry is
	//      already zero.
	//
	//  TextAdded
	//    - Called with key set to a nil value
	//      TODO: Probably have a separate OnChanged or fall up to
	//      the widget.Entry implementation. Or declare new typed
	//      constants.
	OnNavigate func(key *fyne.KeyEvent)

	allowedRunes []rune
}

// NewSingleRuneEntry creates a new SingleRuneEntry configured
// to allow the given runes.
func NewSingleRuneEntry(allowedRunes []rune) *SingleRuneEntry {
	entry := widget.NewEntry()
	entry.TextStyle.Monospace = true
	se := &SingleRuneEntry{
		Entry: widget.Entry{
			TextStyle:    fyne.TextStyle{Monospace: true},
			CursorRow:    1,
			CursorColumn: 1,
		},
		allowedRunes: allowedRunes,
	}
	se.ExtendBaseWidget(se)
	return se
}

// TypedKey overrides the default Entry TypedKey behavior
// and responds to backspace and direction events.
func (s *SingleRuneEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyBackspace:
		switch len(s.Text) {
		case 1:
			// Remove the current content
			s.Entry.SetText("")
		case 0:
			if s.OnNavigate != nil {
				s.OnNavigate(key)
			}
		}
	case fyne.KeyLeft, fyne.KeyRight:
		if s.OnNavigate != nil {
			s.OnNavigate(key)
		}
	default:
		// drop everything else
	}
}

// TypedRune overrides the default Entry TypedRune.
// It ensures only one rune in the entry and filters
// out those that are not in the allowed runes.
func (s *SingleRuneEntry) TypedRune(r rune) {
	// Drop runes that are not allowed
	if !runeSliceContains(s.allowedRunes, r) {
		return
	}
	// If there is any text we are going to replace it
	// with the next rune
	if len(s.Text) != 0 {
		s.SetText("")
	}

	// Pass the rune to the entry
	s.Entry.TypedRune(r)

	// Call to navigate
	if s.OnNavigate != nil {
		s.OnNavigate(nil)
	}
}

func runeSliceContains(rr []rune, r rune) bool {
	for _, x := range rr {
		if x == r {
			return true
		}
	}
	return false
}
