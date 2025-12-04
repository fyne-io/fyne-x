package widget

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

	data     binding.Float
	listener binding.DataListener
	mPr      *message.Printer
}

// NewNumericalEntry returns an extended entry that only allows numerical input.
func NewNumericalEntry() *NumericalEntry {
	entry := &NumericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

// NewNumericalEntryWithData creates a numerical entry that is bound to a
// data source.
func NewNumericalEntryWithData(data binding.Float) *NumericalEntry {
	e := NewNumericalEntry()
	e.Bind(data)

	return e
}

// Append appends text to the entry, filtering out non-numerical characters
// based on the current locale and allowed input types (negative, float).
func (e *NumericalEntry) Append(text string) {
	if e.Validator == nil {
		e.Validator = e.validateText
	}
	s := e.getValidText(e.Text, text)
	e.Entry.Append(s)
}

// Bind connects the specified data source to this Spinner widget.
// The current value will be displayed and any changes in the data will cause the widget
// to update.
func (e *NumericalEntry) Bind(data binding.Float) {
	e.data = data
	e.listener = binding.NewDataListener(e.updateFromData)
	data.AddListener(e.listener)
	e.OnChanged = e.writeData
}

func (e *NumericalEntry) FocusLost() {
	if e.Validator == nil {
		e.Validator = e.validateText
	}
	e.Entry.FocusLost()
}

// Value parses the text content of the entry as a float64.
// It returns the parsed float and an error if parsing fails.
func (e *NumericalEntry) Value() (float64, error) {
	t, err := e.makeParsable(e.Text)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(t, 64)
}

// SetValue sets the entry's text to the string representation of the given float64 value,
// formatted according to the entry's locale.
func (e *NumericalEntry) SetValue(value float64) {
	// if radixSep == 0, then setup has not been called, so call it now.
	if e.radixSep == 0 {
		e.setup()
	}
	var numStr string
	if e.AllowFloat {
		numStr = e.mPr.Sprintf("%f", value)
	} else {
		numStr = e.mPr.Sprintf("%d", int(value))
	}
	e.SetText(numStr)
}

// SetText manually sets the text of the Entry.
// The text will be filtered to allow only numerical input
// according to the current locale.
func (e *NumericalEntry) SetText(text string) {
	if e.Validator == nil {
		e.Validator = e.validateText
	}
	s := e.getValidText("", text)
	e.Entry.SetText(s)
	e.Refresh()
}

// TypedRune is called when this item receives a char event.
//
// Implements: fyne.Focusable
func (e *NumericalEntry) TypedRune(r rune) {
	if e.Validator == nil {
		e.Validator = e.validateText
	}
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
	if e.Validator == nil {
		e.Validator = e.validateText
	}
	runes := []rune(e.Text)
	_, ok := shortcut.(*fyne.ShortcutPaste)
	if ok && len(runes) > 0 && e.CursorColumn == 0 && runes[0] == e.minus {
		return
	}

	e.Entry.TypedShortcut(shortcut)
	// now reprocess the NumericalEntry's Text to change characters to locale-specific values
	// and delete those that are not valid.
	t := e.Text
	e.SetText(t)
}

// Keyboard sets up the right keyboard to use on mobile.
//
// Implements: mobile.Keyboardable
func (e *NumericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

// Unbind disconnects the entry from the bound data.
// This will remove the data updates and allow the entry to be used independently.
func (e *NumericalEntry) Unbind() {
	if e.data == nil {
		return
	}
	e.data.RemoveListener(e.listener)
	e.data = nil
	e.listener = nil
	e.OnChanged = nil
}

// validateText checks if the entered text is a valid numerical input
// according to the system locale.
func (e *NumericalEntry) validateText(text string) error {
	if len(text) == 0 {
		return nil
	}
	if e.radixSep == 0 {
		e.setup()
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

// getRuneForLocale checks if a rune is valid for the entry,
// and returns the correct rune for the locale.
func (e *NumericalEntry) getRuneForLocale(r rune) (rune, bool) {
	// if radixSep == 0, then setup has not been called, so call it now.
	if e.radixSep == 0 {
		e.setup()
	}
	if unicode.IsDigit(r) {
		return r, true
	}

	switch r {
	case '-', 0x2212:
		if e.AllowNegative {
			return e.minus, true
		} else {
			return 0, false
		}
	case '.', ',':
		if r == e.radixSep || r == e.thouSep {
			return r, true
		}
	case ' ', 0xa0:
		if e.thouSep == ' ' || e.thouSep == 0xa0 {
			return e.thouSep, true
		}
	case '\'', 0x2019:
		if e.thouSep == '\'' || e.thouSep == 0x2019 {
			return e.thouSep, true
		}
	default:
		if r == e.radixSep || r == e.thouSep {
			return r, true
		}
	}
	return 0, false
}

// getValidText filters the input text to allow only valid characters
// based on the current locale and the entry's configuration.
func (e *NumericalEntry) getValidText(curText, text string) string {
	var s strings.Builder
	for i, r := range text {
		rn, ok := e.getRuneForLocale(r)
		if !ok {
			continue
		}
		if rn == e.minus {
			if !e.AllowNegative || !(len(curText) == 0) || i != 0 {
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
	return s.String()
}

// makeParsable prepares the text for parsing by removing thousand separators,
// replacing the minus sign with a standard minus, and replacing the radix
// separator with a period. It also validates the text before processing.
func (e *NumericalEntry) makeParsable(text string) (string, error) {
	err := e.validateText(text)
	e.SetValidationError(err)
	if err != nil {
		return "", err
	}
	t := text
	t = strings.ReplaceAll(text, string(e.thouSep), "")
	t = strings.Replace(t, string(e.minus), "-", -1)
	t = strings.Replace(t, string(e.radixSep), ".", -1)
	return t, nil
}

// setup determines the system locale and calls setupRunes to determine the
// locale's minus sign, radix, and thousands separator.
func (e *NumericalEntry) setup() {
	locale := lang.SystemLocale().String()
	e.setupRunes(locale)
}

// setupRunes determines the minus sign, radix, and thousand separator for
// the specified locale.
func (e *NumericalEntry) setupRunes(locale string) {
	lang, err := language.Parse(locale)
	if err != nil {
		fyne.LogError("Language parse error: ", err)
		lang = language.English // Fallback to English
		locale = "en-US"
	}
	e.mPr = message.NewPrinter(lang)
	e.minus = '-'
	e.thouSep = ','
	e.radixSep = '.'
	p := message.NewPrinter(lang)
	numStr := p.Sprintf("%f", -12345.5678901)
	runes := []rune(numStr)

	// Define helper function to find separator
	findSeparator := func(runes []rune) rune {
		for _, r := range runes {
			if !unicode.IsDigit(r) && r != e.minus {
				return r
			}
		}
		return 0 // Return 0 if no separator is found
	}

	// first rune is the "minus" sign
	if len(runes) > 0 {
		e.minus = runes[0]
	}

	// look for thousands separator
	if len(runes) > 1 {
		e.thouSep = findSeparator(runes[1:min(6, len(runes))])
	}

	// look for radix separator
	if len(runes) > 6 {
		e.radixSep = findSeparator(runes[6:])
		if e.radixSep == 0 && e.thouSep != 0 {
			// If no radix separator is found, and thousand separator is present,
			// use the last non-digit character as radix separator
			e.radixSep = findSeparator(runes[len(runes)-7:])
		}
	}
}

// updateFromData updates the entry's text with the value from the data source.
// It checks if the current value is different before updating to prevent unnecessary refreshes.
func (e *NumericalEntry) updateFromData() {
	if e.listener == nil {
		return
	}
	if e.data == nil {
		return
	}
	val, err := e.data.Get()
	if err != nil {
		fyne.LogError("Error getting current data value", err)
		return
	}

	currentVal, _ := e.Value()
	if currentVal == val {
		return
	}

	e.SetValue(val)
}

// writeData writes the entry's parsed float value to the specified data binding.
// It validates the text, parses it as a float, and updates the binding if the value has changed.
func (e *NumericalEntry) writeData(_ string) {
	if e.data == nil {
		return
	}
	if err := e.Validate(); err != nil {
		return
	}
	val, err := e.Value()
	if err != nil {
		return
	}
	dVal, err := e.data.Get()
	if err != nil {
		fyne.LogError("Error getting float value: ", err)
		return
	}

	if dVal == val {
		return
	}

	if err := e.data.Set(val); err != nil {
		fyne.LogError("Error setting float value: ", err)
	}
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
