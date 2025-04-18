package widget

import (
	"fmt"
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
	minus = runes[0]
	for _, r := range runes[1:5] {
		if !unicode.IsDigit(r) && r != '-' && r != 0x2212 {
			/*			if r != ' ' && // space
						r != 0xa0 && // non-breaking space
						r == 0xe2a8af {
						continue
					}*/
			thou = r
			break
		}
	}
	for _, r := range runes[5:] {
		if !unicode.IsDigit(r) {
			radix = r
			break
		}
		fmt.Printf("locale: %s\n", locale)
		fmt.Printf("minus = '%v'  0x%x\n", string(minus), minus)
		fmt.Printf("radix = '%v'  0x%x\n", string(radix), radix)
		fmt.Printf("thou = '%v'  0x%x\n", string(thou), thou)
	}
	fmt.Println(numStr)
	return minus, radix, thou
}
