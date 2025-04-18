package widget

import (
	"unicode"

	"fyne.io/fyne/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// minusRadixThou determines the minus sign, radix, and thousand separator
// characters for a given locale by formatting a number and extracting
// the relevant characters. It returns the minus sign, radix, and thousand
// separator runes.
func minusRadixThou(locale string) (rune, rune, rune) {
	minus := '-'
	thou := ','
	radix := '.'
	lang, err := language.Parse(locale)
	if err != nil {
		fyne.LogError("Parse error: %s\n", err)
		return minus, radix, thou
	}
	p := message.NewPrinter(lang)
	numStr := p.Sprintf("%f", -12345.5678901)
	runes := []rune(numStr)
	// first rune is the "minus" sign
	minus = runes[0]
	// look for thousands separator
	for _, r := range runes[1:5] {
		if !unicode.IsDigit(r) {
			thou = r
			break
		}
	}
	// look for radix separator
	for _, r := range runes[5:] {
		if !unicode.IsDigit(r) {
			radix = r
			break
		}
	}
	return minus, radix, thou
}
