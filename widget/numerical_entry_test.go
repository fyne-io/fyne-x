package widget

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestNumericalEntry_IntHyphenStopComma(t *testing.T) {
	entry := NewNumericalEntry()
	entry.minus = '-'
	entry.radixSep = '.'
	entry.thouSep = ','

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := "1" + string(rune(0x2019)) + " 23,4" + string(rune(0xa0)) + "56,789"
	test.Type(entry, number)
	assert.Equal(t, "123,456,789", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "3")
	assert.Equal(t, "3123,456,789", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "-4")
	assert.Equal(t, "43123,456,789", entry.Text)
}

func TestNumericalEntry_IntHyphenCommaStop(t *testing.T) {
	entry := NewNumericalEntry()
	entry.minus = '-'
	entry.radixSep = ','
	entry.thouSep = '.'

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := "123.456.789"
	test.Type(entry, number)
	assert.Equal(t, number, entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "3")
	assert.Equal(t, "3123.456.789", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "-4")
	assert.Equal(t, "43123.456.789", entry.Text)
}

func TestNumericalnEntry_FloatHyphenStopComma(t *testing.T) {
	entry := NewNumericalEntry()
	entry.minus = '-'
	entry.radixSep = '.'
	entry.thouSep = ','

	entry.AllowFloat = true

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := "123.456,789"
	test.Type(entry, number)
	assert.Equal(t, number, entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "3")
	assert.Equal(t, "3123.456,789", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "-4")
	assert.Equal(t, "43123.456,789", entry.Text)
}

func TestNumericalnEntry_FloatHyphenCommaStop(t *testing.T) {
	entry := NewNumericalEntry()
	entry.minus = '-'
	entry.radixSep = ','
	entry.thouSep = '.'

	entry.AllowFloat = true

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := "123.456,789"
	test.Type(entry, number)
	assert.Equal(t, number, entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "3")
	assert.Equal(t, "3123.456,789", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "-4")
	assert.Equal(t, "43123.456,789", entry.Text)
}

func TestNumericalEntry_IntMathMinusStopComma(t *testing.T) {
	entry := NewNumericalEntry()
	entry.minus = 0x2212
	entry.radixSep = '.'
	entry.thouSep = ','

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := "123,456,789"
	test.Type(entry, number)

	entry.CursorColumn = 0
	test.Type(entry, "3")
	assert.Equal(t, "3123,456,789", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, string(entry.minus)+"4")
	assert.Equal(t, "43123,456,789", entry.Text)
}

func TestNumericalEntry_IntMathMinusCommaStop(t *testing.T) {
	entry := NewNumericalEntry()
	entry.minus = 0x2212
	entry.radixSep = ','
	entry.thouSep = '.'

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := "123.456.789"
	test.Type(entry, number)
	assert.Equal(t, number, entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "3")
	assert.Equal(t, "3123.456.789", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, string(entry.minus)+"-4")
	assert.Equal(t, "43123.456.789", entry.Text)
}

func TestNumericalnEntry_FloatMathMinusStopComma(t *testing.T) {
	entry := NewNumericalEntry()
	entry.minus = 0x2212
	entry.radixSep = '.'
	entry.thouSep = ','

	entry.AllowFloat = true

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := string(rune(0x2212)) + "1" + string(rune(0xa0)) + "23.456,789"
	test.Type(entry, number)
	assert.Equal(t, "123.456,789", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "3")
	assert.Equal(t, "3123.456,789", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, string(entry.minus)+"4")
	assert.Equal(t, "43123.456,789", entry.Text)
}

func TestNumericalnEntry_FloatMathMinusCommaStop(t *testing.T) {
	entry := NewNumericalEntry()
	entry.minus = 0x2212
	entry.radixSep = ','
	entry.thouSep = '.'

	entry.AllowFloat = true

	test.Type(entry, "Not a number")
	assert.Empty(t, entry.Text)

	number := string(rune(0x2212)) + "1" + string(rune(0xa0)) + "23.456,789"
	test.Type(entry, number)
	assert.Equal(t, "123.456,789", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "3")
	assert.Equal(t, "3123.456,789", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, string(entry.minus)+"4")
	assert.Equal(t, "43123.456,789", entry.Text)
}

func TestNumericalEntry_NegIntMinus(t *testing.T) {
	entry := NewNumericalEntry()
	entry.minus = '-'
	entry.radixSep = ','
	entry.thouSep = '.'
	entry.AllowNegative = true

	test.Type(entry, "-2")
	assert.Equal(t, "-2", entry.Text)

	entry.Text = ""
	test.Type(entry, "24-")
	assert.Equal(t, "24", entry.Text)
	entry.CursorColumn = 0
	test.Type(entry, "-")
	assert.Equal(t, "-24", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "4")
	assert.Equal(t, "-24", entry.Text)
}

func TestNumericalEntry_NegFloatMinus(t *testing.T) {
	entry := NewNumericalEntry()
	entry.minus = '-'
	entry.radixSep = ','
	entry.thouSep = '.'
	entry.AllowNegative = true
	entry.AllowFloat = true

	test.Type(entry, "-2.4")
	assert.Equal(t, "-2.4", entry.Text)

	entry.Text = ""
	test.Type(entry, "-24.-5")
	assert.Equal(t, "-24.5", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "-")
	assert.Equal(t, "-24.5", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "4")
	assert.Equal(t, "-24.5", entry.Text)

}

func TestNumericalEntry_NegIntMathMinus(t *testing.T) {
	entry := NewNumericalEntry()
	entry.minus = 0x2212
	entry.radixSep = ','
	entry.thouSep = '.'
	entry.AllowNegative = true

	test.Type(entry, string(entry.minus)+"2")
	fmt.Println(entry.Text)
	assert.Equal(t, string(entry.minus)+"2", entry.Text)

	entry.Text = ""
	test.Type(entry, "2-4"+string(rune(0x2212)))
	assert.Equal(t, "24", entry.Text)
	entry.CursorColumn = 0
	test.Type(entry, string(rune(0x2212)))
	assert.Equal(t, string(rune(0x2212))+"24", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "4")
	assert.Equal(t, string(rune(0x2212))+"24", entry.Text)
}

func TestNumericalEntry_NegFloat(t *testing.T) {
	entry := NewNumericalEntry()
	entry.AllowNegative = true
	entry.AllowFloat = true

	test.Type(entry, "-2.4")
	assert.Equal(t, "-2.4", entry.Text)

	entry.Text = ""
	test.Type(entry, "-24.-5")
	assert.Equal(t, "-24.5", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "-")
	assert.Equal(t, "-24.5", entry.Text)

	entry.CursorColumn = 0
	test.Type(entry, "4")
	assert.Equal(t, "-24.5", entry.Text)

}

func TestNumericalEntry_getRuneForLocale(t *testing.T) {
	entry := NewNumericalEntry()
	entry.AllowNegative = true
	entry.AllowFloat = true
	entry.minus = '-'
	entry.radixSep = '.'
	entry.thouSep = ','

	testCases := []struct {
		name         string
		minus        rune
		radixSep     rune
		thouSep      rune
		input        rune
		expectedRune rune
		expectedBool bool
	}{
		{
			name:         "digit",
			minus:        '-',
			radixSep:     '.',
			thouSep:      ',',
			input:        '1',
			expectedRune: '1',
			expectedBool: true,
		},
		{
			name:         "minus - minus",
			minus:        '-',
			radixSep:     '.',
			thouSep:      ',',
			input:        '-',
			expectedRune: '-',
			expectedBool: true,
		},
		{
			name:         "minus - alt minus",
			minus:        0x2212,
			radixSep:     '.',
			thouSep:      ',',
			input:        '-',
			expectedRune: 0x2212,
			expectedBool: true,
		},
		{
			name:         "alternative minus - minus",
			minus:        '-',
			radixSep:     '.',
			thouSep:      ',',
			input:        0x2212,
			expectedRune: '-',
			expectedBool: true,
		},
		{
			name:         "alternative minus - alt minus",
			minus:        0x2212,
			radixSep:     '.',
			thouSep:      ',',
			input:        0x2212,
			expectedRune: 0x2212,
			expectedBool: true,
		},
		{
			name:         "radix stop - stop",
			minus:        '-',
			radixSep:     '.',
			thouSep:      ',',
			input:        '.',
			expectedRune: '.',
			expectedBool: true,
		},
		{
			name:         "radix comma - comma",
			minus:        '-',
			radixSep:     ',',
			thouSep:      '.',
			input:        '.',
			expectedRune: '.',
			expectedBool: true,
		},
		{
			name:         "thou stop - stop",
			minus:        '-',
			radixSep:     ',',
			thouSep:      '.',
			input:        '.',
			expectedRune: '.',
			expectedBool: true,
		},
		{
			name:         "thou stop - space",
			minus:        '-',
			radixSep:     ',',
			thouSep:      '.',
			input:        ' ',
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou stop - alt space",
			minus:        '-',
			radixSep:     ',',
			thouSep:      '.',
			input:        0xa0,
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou stop - apostrophe",
			minus:        '-',
			radixSep:     ',',
			thouSep:      '.',
			input:        0x2019,
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou comma - comma",
			minus:        '-',
			radixSep:     '.',
			thouSep:      ',',
			input:        ',',
			expectedRune: ',',
			expectedBool: true,
		},
		{
			name:         "thou comma - space",
			minus:        '-',
			radixSep:     '.',
			thouSep:      ',',
			input:        ' ',
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou comma - alt space",
			minus:        '-',
			radixSep:     '.',
			thouSep:      ',',
			input:        0xa0,
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou comma - apostrophe",
			minus:        '-',
			radixSep:     '.',
			thouSep:      ',',
			input:        0x2019,
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou space - space",
			minus:        '-',
			radixSep:     '.',
			thouSep:      ' ',
			input:        ' ',
			expectedRune: ' ',
			expectedBool: true,
		},
		{
			name:         "thou space - stop",
			minus:        '-',
			radixSep:     ',',
			thouSep:      ' ',
			input:        '.',
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou space - comma",
			minus:        '-',
			radixSep:     '.',
			thouSep:      ' ',
			input:        ',',
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou space - apostrophe",
			minus:        '-',
			radixSep:     ',',
			thouSep:      ' ',
			input:        0x2019,
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou space - non-breaking space",
			minus:        '-',
			radixSep:     '.',
			thouSep:      ' ',
			input:        0xa0,
			expectedRune: ' ',
			expectedBool: true,
		},
		{
			name:         "thou non-breaking space - non-breaking space",
			minus:        '-',
			radixSep:     '.',
			thouSep:      0xa0,
			input:        0xa0,
			expectedRune: 0xa0,
			expectedBool: true,
		},
		{
			name:         "thou non-breaking space - space",
			minus:        '-',
			radixSep:     '.',
			thouSep:      0xa0,
			input:        ' ',
			expectedRune: 0xa0,
			expectedBool: true,
		},
		{
			name:         "thou non-breaking space - stop",
			minus:        '-',
			radixSep:     ',',
			thouSep:      0xa0,
			input:        '.',
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou non-breaking space - comma",
			minus:        '-',
			radixSep:     '.',
			thouSep:      0xa0,
			input:        ',',
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou non-breaking space - apostrophe",
			minus:        '-',
			radixSep:     ',',
			thouSep:      0xa0,
			input:        0x2019,
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou apostrophe - apostrophe",
			minus:        '-',
			radixSep:     '.',
			thouSep:      0x2019,
			input:        0x2019,
			expectedRune: 0x2019,
			expectedBool: true,
		},
		{
			name:         "thou apostrophe - stop",
			minus:        '-',
			radixSep:     ',',
			thouSep:      0x2019,
			input:        '.',
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou apostrophe - comma",
			minus:        '-',
			radixSep:     '.',
			thouSep:      0x2019,
			input:        ',',
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou apostrophe - space",
			minus:        '-',
			radixSep:     ',',
			thouSep:      0x2019,
			input:        ' ',
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "thou apostrophe - non-breaking space",
			minus:        '-',
			radixSep:     ',',
			thouSep:      0x2019,
			input:        0xa0,
			expectedRune: 0,
			expectedBool: false,
		},
		{
			name:         "invalid rune",
			minus:        '-',
			radixSep:     '.',
			thouSep:      0x2019,
			input:        'a',
			expectedRune: 0,
			expectedBool: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry.minus = tc.minus
			entry.radixSep = tc.radixSep
			entry.thouSep = tc.thouSep
			actualRune, actualBool := entry.getRuneForLocale(tc.input)
			if actualRune != tc.expectedRune {
				t.Errorf("Expected rune %q, but got %q", tc.expectedRune, actualRune)
			}
			if actualBool != tc.expectedBool {
				t.Errorf("Expected bool %v, but got %v", tc.expectedBool, actualBool)
			}
		})
	}
}

func TestNumericalEntry_SetText(t *testing.T) {
	e := NewNumericalEntry()
	e.AllowNegative = true

	// Test with valid numerical input
	e.SetText("123")
	if e.Text != "123" {
		t.Errorf("Expected '123', got '%s'", e.Text)
	}

	// Test with invalid characters - should be filtered out
	e.SetText("1a2b3c")
	if e.Text != "123" {
		t.Errorf("Expected '123', got '%s'", e.Text)
	}

	// Test with negative sign
	e.SetText("-456")
	if e.Text != "-456" {
		t.Errorf("Expected '-456', got '%s'", e.Text)
	}

	// Test with negative sign when AllowNegative is false
	e.AllowNegative = false
	e.SetText("-789")
	if e.Text != "789" {
		t.Errorf("Expected '789', got '%s'", e.Text)
	}

	// Test with empty string
	e.SetText("")
	if e.Text != "" {
		t.Errorf("Expected '', got '%s'", e.Text)
	}

	// Test with minus sign in the middle of the number when AllowNegative is true
	e.AllowNegative = true
	e.SetText("1-23")
	if e.Text != "123" {
		t.Errorf("Expected '123', got '%s'", e.Text)
	}

	// Test with leading and trailing spaces
	e.SetText("  123  ")
	if e.Text != "123" {
		t.Errorf("Expected '123', got '%s'", e.Text)
	}

	// Test with only spaces
	e.SetText("   ")
	if e.Text != "" {
		t.Errorf("Expected '', got '%s'", e.Text)
	}
}

func TestNumericalEntry_SetText_Locale(t *testing.T) {
	e := NewNumericalEntry()
	e.AllowNegative = true
	e.minus = '−' // Different minus sign

	// Test with different minus sign
	e.SetText("−123")
	if e.Text != "−123" {
		t.Errorf("Expected '−123', got '%s'", e.Text)
	}

	// Test with regular minus sign when custom minus is set
	e.SetText("-456")
	if e.Text != "−456" {
		t.Errorf("Expected '−456', got '%s'", e.Text)
	}
}

func TestNumericalEntry_SetText_Filtering(t *testing.T) {
	e := NewNumericalEntry()

	e.SetText("abc123def456")
	if e.Text != "123456" {
		t.Errorf("Expected '123456', got '%s'", e.Text)
	}

	e.SetText("123.45") // Assuming only integers are allowed
	if e.Text != "12345" {
		t.Errorf("Expected '12345', got '%s'", e.Text)
	}
}

func TestNumericalEntry_SetText_Callbacks(t *testing.T) {
	e := NewNumericalEntry()
	var callbackCalled bool
	e.OnChanged = func(string) {
		callbackCalled = true
	}

	e.SetText("123")
	if !callbackCalled {
		t.Error("Expected OnChanged callback to be called")
	}

	callbackCalled = false
	e.SetText("abc") // Should still call callback even if text is filtered
	if !callbackCalled {
		t.Error("Expected OnChanged callback to be called even with invalid input")
	}
}

func TestNumericalEntry_Append(t *testing.T) {
	e := NewNumericalEntry()
	e.AllowFloat = true
	e.AllowNegative = true
	e.minus = '-'
	e.radixSep = '.'
	e.thouSep = ','

	e.SetText("")
	e.Append("123")
	if e.Text != "123" {
		t.Errorf("expected '123', got '%s'", e.Text)
	}

	e.Append(".45")
	if e.Text != "123.45" {
		t.Errorf("expected '123.45', got '%s'", e.Text)
	}

	e = NewNumericalEntry()
	e.AllowFloat = false
	e.AllowNegative = false
	e.minus = '-'
	e.radixSep = '.'
	e.thouSep = ','

	e.SetText("")
	e.Append("123.45")
	if e.Text != "12345" {
		t.Errorf("expected '12345', got '%s'", e.Text)
	}

	e = NewNumericalEntry()
	e.AllowFloat = false
	e.AllowNegative = true
	e.minus = '-'
	e.radixSep = '.'
	e.thouSep = ','

	e.SetText("")
	e.Append("-123")
	if e.Text != "-123" {
		t.Errorf("expected '-123', got '%s'", e.Text)
	}

	e.Append(".45")
	if e.Text != "-12345" {
		t.Errorf("expected '-12345', got '%s'", e.Text)
	}

	e = NewNumericalEntry()
	e.AllowFloat = true
	e.AllowNegative = true
	e.minus = '-'
	e.radixSep = '.'
	e.thouSep = ','
	e.Text = ""

	e.Append("-")
	if e.Text != "-" {
		t.Errorf("expected '-', got '%s'", e.Text)
	}

	e.Append("1")
	if e.Text != "-1" {
		t.Errorf("expected '-1', got '%s'", e.Text)
	}
}
