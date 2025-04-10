package widget

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestNumericalnEntry_Int(t *testing.T) {
	entry := NewNumericalEntry()

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := "123456789"
	test.Type(entry, number)
	assert.Equal(t, number, entry.Text)
}

func TestNumericalnEntry_Float(t *testing.T) {
	entry := NewNumericalEntry()
	entry.AllowFloat = true

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := "123.456789"
	test.Type(entry, number)
	assert.Equal(t, number, entry.Text)
}

func TestNumericalEntry_NegInt(t *testing.T) {
	entry := NewNumericalEntry()

	test.Type(entry, "-2")
	assert.Equal(t, "-2", entry.Text)

	entry.Text = ""
	test.Type(entry, "24-")
	assert.Equal(t, "24", entry.Text)
}

func TestNumericalEntry_NegFloat(t *testing.T) {
	entry := NewNumericalEntry()
	entry.AllowFloat = true

	test.Type(entry, "-2.4")
	assert.Equal(t, "-2.4", entry.Text)

	entry.Text = ""
	test.Type(entry, "24.-5")
	assert.Equal(t, "24.5", entry.Text)
}
