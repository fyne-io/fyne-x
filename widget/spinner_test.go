package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewSpinner(t *testing.T) {
	s := NewSpinner(1, 5, 2, nil)
	assert.Equal(t, 1, s.min)
	assert.Equal(t, 5, s.max)
	assert.Equal(t, 2, s.step)
	assert.Equal(t, 1, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestSetValue(t *testing.T) {
	s := NewSpinner(1, 5, 2, nil)
	s.SetValue(2)
	assert.Equal(t, 2, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestSetValue_LessThanMin(t *testing.T) {
	s := NewSpinner(4, 22, 5, nil)
	s.SetValue(3)
	assert.Equal(t, 4, s.GetValue())
	assert.True(t, s.downButton.Disabled())
	assert.False(t, s.upButton.Disabled())
}

func TestSetValue_GreaterThanMax(t *testing.T) {
	s := NewSpinner(4, 22, 5, nil)
	s.SetValue(23)
	assert.Equal(t, 22, s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestSetValue_DisabledSpinner(t *testing.T) {
	s := NewSpinner(4, 22, 5, nil)
	s.Disable()
	s.SetValue(10)
	assert.Equal(t, 4, s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestUpButtonTapped(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.upButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 9, s.GetValue())
	s.upButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 10, s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestDownButtonTapped(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.SetValue(10)
	s.downButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 5, s.GetValue())
	s.downButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 4, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}
