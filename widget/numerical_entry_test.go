package widget

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestNumericalnEntry_Int(t *testing.T) {
	entry := NewNumericalEntry(false)

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := "123456789"
	test.Type(entry, number)
	assert.Equal(t, number, entry.Text)
}

func TestNumericalnEntry_Float(t *testing.T) {
	entry := NewNumericalEntry(true)

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := "123.456789"
	test.Type(entry, number)
	assert.Equal(t, number, entry.Text)
}
