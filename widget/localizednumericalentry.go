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
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/cloudfoundry-attic/jibber_jabber"
)

// LocalizedNumericalEntry is an extended entry that only allows numerical input.
// Only integers are allowed by default. Support for floats can be enabled by setting AllowFloat.
type LocalizedNumericalEntry struct {
	widget.Entry
	AllowFloat bool
	// AllowNegative determines if negative numbers can be entered.
	AllowNegative bool
	minus         rune
	radixSep      rune
	thouSep       rune

	binder basicBinder
	locale string
	mPr    *message.Printer
}

// NewLocalizedNumericalEntry returns an extended entry that only allows numerical input.
func NewLocalizedNumericalEntry() *LocalizedNumericalEntry {
	entry := &LocalizedNumericalEntry{}
	var err error
	entry.locale, err = jibber_jabber.DetectIETF()
	if err != nil {
		fyne.LogError("DetectIETF error: %s\n", err)
		return nil
	} else {
		lang, err := language.Parse(entry.locale)
		if err != nil {
			fyne.LogError("Language parse error: ", err)
			return nil
		}
		entry.mPr = message.NewPrinter(lang)

		entry.getLocaleRunes(entry.locale)
	}
	entry.ExtendBaseWidget(entry)
	entry.Validator = entry.ValidateText
	return entry
}

// NewLocalizedNumericalEntryWithData creates a numerical entry that is bound to a
// data source and can allow or disallow float and negative numbers.
func NewLocalizedNumericalEntryWithData(allowFloat bool, allowNegative bool, data binding.Float) *LocalizedNumericalEntry {
	e := NewLocalizedNumericalEntry()
	e.AllowFloat = allowFloat
	e.AllowNegative = allowNegative
	e.Bind(data)
	e.OnChanged = func(string) {
		e.binder.CallWithData(e.writeData)
	}

	return e
}

// Append appends text to the entry, filtering out non-numerical characters
// based on the current locale and allowed input types (negative, float).
func (e *LocalizedNumericalEntry) Append(text string) {
	s := e.getValidText(e.Text, text)
	e.Entry.Append(s)
}

// Bind connects the specified data source to this Spinner widget.
// The current value will be displayed and any changes in the data will cause the widget
// to update.
func (e *LocalizedNumericalEntry) Bind(data binding.Float) {
	e.binder.SetCallback(e.updateFromData)
	e.binder.Bind(data)
}

// ParseFloat parses the text content of the entry as a float64.
// It returns the parsed float and an error if parsing fails.
func (e *LocalizedNumericalEntry) ParseFloat() (float64, error) {
	t, err := e.makeParsable(e.Text)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(t, 64)
}

// SetValue sets the entry's text to the string representation of the given float64 value,
// formatted according to the entry's locale.
func (e *LocalizedNumericalEntry) SetValue(value float64) {
	if e.mPr == nil {
		return
	}
	var numStr string
	if e.AllowFloat {
		numStr = e.mPr.Sprintf("%f", value)
	} else {
		numStr = e.mPr.Sprintf("&d", int(value))
	}
	e.SetText(numStr)
}

// SetText manually sets the text of the Entry.
// The text will be filtered to allow only numerical input
// according to the current locale.
func (e *LocalizedNumericalEntry) SetText(text string) {
	s := e.getValidText("", text)
	e.Entry.SetText(s)
}

// TypedRune is called when this item receives a char event.
//
// Implements: fyne.Focusable
func (e *LocalizedNumericalEntry) TypedRune(r rune) {
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
func (e *LocalizedNumericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	runes := []rune(e.Text)
	_, ok := shortcut.(*fyne.ShortcutPaste)
	if ok && len(runes) > 0 && e.CursorColumn == 0 && runes[0] == e.minus {
		return
	}

	e.Entry.TypedShortcut(shortcut)
	// now reprocess the LocalizedNumericalEntry's Text to change characters to locale-specific values
	// and delete those that are not valid.
	t := e.Text
	e.SetText(t)
}

// Keyboard sets up the right keyboard to use on mobile.
//
// Implements: mobile.Keyboardable
func (e *LocalizedNumericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

// Unbind disconnects the entry from the bound data.
// This will remove the data updates and allow the entry to be used independently.
func (e *LocalizedNumericalEntry) Unbind() {
	e.binder.Unbind()
}

// ValidateText checks if the entered text is a valid numerical input
// according to the system locale.
func (e *LocalizedNumericalEntry) ValidateText(text string) error {
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

// getRuneForLocale checks if a rune is valid for the entry,
// and returns the correct rune for the locale.
func (e *LocalizedNumericalEntry) getRuneForLocale(r rune) (rune, bool) {
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
	case '\'': // single quote
		fallthrough
	case 0x2019: // right single quote mark
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
func (e *LocalizedNumericalEntry) getValidText(curText, text string) string {
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
func (e *LocalizedNumericalEntry) makeParsable(text string) (string, error) {
	err := e.ValidateText(text)
	if err != nil {
		return "", err
	}
	t := text
	t = strings.ReplaceAll(text, string(e.thouSep), "")
	t = strings.Replace(t, string(e.minus), "-", -1)
	t = strings.Replace(t, string(e.radixSep), ".", -1)
	return t, nil
}

// minusRadixThou determines the minus sign, radix, and thousand separator
// characters for a given locale by formatting a number and extracting
// the relevant characters. It returns the minus sign, radix, and thousand
// separator runes.
func (e *LocalizedNumericalEntry) getLocaleRunes(locale string) {
	e.minus = '-'
	e.thouSep = ','
	e.radixSep = '.'
	lang, err := language.Parse(locale)
	if err != nil {
		fyne.LogError("Language parse error: ", err)
		return
	}
	p := message.NewPrinter(lang)
	numStr := p.Sprintf("%f", -12345.5678901)
	runes := []rune(numStr)
	// first rune is the "minus" sign
	e.minus = runes[0]
	// look for thousands separator
	for _, r := range runes[1:5] {
		if !unicode.IsDigit(r) {
			e.thouSep = r
			break
		}
	}
	// look for radix separator
	for _, r := range runes[5:] {
		if !unicode.IsDigit(r) {
			e.radixSep = r
			break
		}
	}
}

// updateFromData updates the entry's text with the value from the data source.
func (e *LocalizedNumericalEntry) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	fltSource, ok := data.(binding.Float)
	if !ok {
		return
	}
	val, err := fltSource.Get()
	if err != nil {
		fyne.LogError("Error getting current data value", err)
		return
	}
	e.SetValue(val)
}

// writeData writes the entry's parsed float value to the specified data binding.
func (e *LocalizedNumericalEntry) writeData(data binding.DataItem) {
	if data == nil {
		return
	}
	if err := e.Validator(e.Text); err != nil {
		return
	}
	val, err := e.ParseFloat()
	if err != nil {
		return
	}
	flt, ok := data.(binding.Float)
	if !ok {
		return
	}

	dVal, err := flt.Get()
	if err != nil || dVal == val {
		return
	}

	flt.Set(val)
}
