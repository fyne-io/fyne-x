// Package validation provides validation for data inside widgets
package validation // import "fyne.io/x/fyne/data/validation"

import (
	gpv "github.com/wagslane/go-password-validator"

	"fyne.io/fyne/v2"
)

// NewPassword returns a new validator for validating passwords.
// Validate returns nil if the password entropy is greater than or equal
// to the minimum entropy. If not, an error is returned that explains
// how the password can be strengthened. Advice on entropy value:
// https://github.com/wagslane/go-password-validator/tree/main#what-entropy-value-should-i-use
func NewPassword(minEntropy float64) fyne.StringValidator {
	return func(text string) error {
		return gpv.Validate(text, minEntropy)
	}
}
