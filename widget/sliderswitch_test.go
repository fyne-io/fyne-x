package widget

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

// Test tap events
func TestSliderSwitch_Tapped(t *testing.T) {
	sw := NewSliderSwitch()
	win := test.NewWindow(sw)
	defer win.Close()

	test.Tap(sw)
	assert.Equal(t, float64(2), sw.Value)

	test.Tap(sw)
	assert.Equal(t, float64(0), sw.Value)
}

// Test OnToggle events
func TestSliderSwitch_OnToggle(t *testing.T) {
	sw := NewSliderSwitch()

	called := false
	expected := true

	sw.OnToggle = func(val bool) {
		called = true
		assert.Equal(t, expected, val)
	}

	win := test.NewWindow(sw)
	defer win.Close()

	test.Tap(sw)
	assert.Equal(t, true, called)
	assert.Equal(t, float64(2), sw.Value)

	called = false
	expected = false

	test.Tap(sw)
	assert.Equal(t, true, called)
	assert.Equal(t, float64(0), sw.Value)
}

// Test OnTransition events
func TestSliderSwitch_OnTransition(t *testing.T) {
	sw := NewSliderSwitch()

	called := false
	expected := true

	sw.OnTransition = func(val bool) {
		called = true
		assert.Equal(t, float64(1), sw.Value)
		assert.Equal(t, expected, val)
	}

	win := test.NewWindow(sw)
	defer win.Close()

	test.Tap(sw)
	assert.Equal(t, true, called)
	assert.Equal(t, float64(2), sw.Value)

	called = false
	expected = false

	test.Tap(sw)
	assert.Equal(t, true, called)
	assert.Equal(t, float64(0), sw.Value)
}
