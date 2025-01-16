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

func TestEnableDisabledSpinner(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.SetValue(7)
	s.Disable()
	assert.True(t, s.Disabled())
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
	s.Enable()
	assert.False(t, s.Disabled())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestEnableDisabledSpinner_UpButtonDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.SetValue(10)
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
	s.Disable()
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())

	s.Enable()
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestEnableDisabledSpinner_DownButtonDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.SetValue(4)
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
	s.Disable()
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())

	s.Enable()
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestRunePlus(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = true
	s.TypedRune('+')
	assert.Equal(t, 9, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedRune('+')
	assert.Equal(t, 10, s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestRuneMinus(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = true
	s.SetValue(10)
	s.TypedRune('-')
	assert.Equal(t, 5, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedRune('-')
	assert.Equal(t, 4, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestRunePlus_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = true
	s.Disable()
	s.TypedRune('+')
	assert.Equal(t, 4, s.GetValue())
}

func TestRuneMinus_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = true
	s.SetValue(8)
	s.Disable()
	s.TypedRune('-')
	assert.Equal(t, 8, s.GetValue())
}

func TestRunePlus_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = false
	s.TypedRune('+')
	assert.Equal(t, 4, s.GetValue())
}

func TestRuneMinus_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = true
	s.SetValue(8)
	s.focused = false
	s.TypedRune('-')
	assert.Equal(t, 8, s.GetValue())
}

func TestKeyUp(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = true
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 9, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedKey(&key)
	assert.Equal(t, 10, s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestKeyDown(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = true
	s.SetValue(10)
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 5, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedKey(&key)
	assert.Equal(t, 4, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestKeyUp_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = true
	s.Disable()
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 4, s.GetValue())
}

func TestKeyDown_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = true
	s.SetValue(8)
	s.Disable()
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 8, s.GetValue())
}

func TestKeyUp_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = false
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 4, s.GetValue())
}

func TestKeyDown_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, nil)
	s.focused = true
	s.SetValue(8)
	s.focused = false
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 8, s.GetValue())
}

func TestScrolled(t *testing.T) {
	s := NewSpinner(1, 10, 1, nil)
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 2, s.GetValue())
	s.Scrolled(&e)
	assert.Equal(t, 3, s.GetValue())
	delta = fyne.Delta{DX: 0, DY: -25}
	e = fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 2, s.GetValue())
}

func TestScrolled_Disabled(t *testing.T) {
	s := NewSpinner(1, 10, 1, nil)
	s.Disable()
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 1, s.GetValue())
}
