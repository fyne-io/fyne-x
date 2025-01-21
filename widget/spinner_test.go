package widget

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/stretchr/testify/assert"
)

func waitForBinding() {
	time.Sleep(time.Millisecond * 100) // data resolves on background thread
}

func TestNewSpinner(t *testing.T) {
	s := NewSpinner(1, 5, 2, nil)
	assert.Equal(t, 1, s.min)
	assert.Equal(t, 5, s.max)
	assert.Equal(t, 2, s.step)
	assert.Equal(t, 1, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestNewSpinner_BadArgs(t *testing.T) {
	assert.Panics(t, func() { NewSpinner(5, 5, 1, nil) }, "Did not panic with min == max")
	assert.Panics(t, func() { NewSpinner(5, 4, 1, nil) }, "Did not panic with min < max")
	assert.Panics(t, func() { NewSpinner(1, 5, 0, nil) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewSpinner(1, 5, -5, nil) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewSpinner(1, 5, 5, nil) }, "Did not panic with step > max - min")
}

func TestNewSpinnerWithData(t *testing.T) {
	data := binding.NewInt()
	s := NewSpinnerWithData(1, 5, 2, data)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	s.SetValue(2)
	waitForBinding()
	val, err = data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 2, val)

	err = data.Set(3)
	assert.NoError(t, err)
	waitForBinding()
	assert.Equal(t, 3, s.GetValue())
}

func TestUnbind(t *testing.T) {
	data := binding.NewInt()
	s := NewSpinnerWithData(1, 5, 2, data)
	waitForBinding()
	s.Unbind()
	s.SetValue(2)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestNewSpinnerWithData_BadArgs(t *testing.T) {
	boundValue := binding.NewInt()
	assert.Panics(t, func() { NewSpinnerWithData(5, 5, 1, boundValue) }, "Did not panic with min == max")
	assert.Panics(t, func() { NewSpinnerWithData(5, 4, 1, boundValue) }, "Did not panic with min < max")
	assert.Panics(t, func() { NewSpinnerWithData(1, 5, 0, boundValue) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewSpinnerWithData(1, 5, -5, boundValue) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewSpinnerWithData(1, 5, 5, boundValue) }, "Did not panic with step > max - min")
}

func TestNewUninitializedSpinner(t *testing.T) {
	s := NewSpinnerUninitialized()
	assert.True(t, s.Disabled())
	s.Enable()
	assert.True(t, s.Disabled())
	s.SetMinMaxStep(-1, 2, 1)
	assert.True(t, s.Disabled())
	s.Enable()
	assert.False(t, s.Disabled())
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

func TestSetMinMaxStep(t *testing.T) {
	s := NewSpinner(1, 6, 2, nil)
	s.SetMinMaxStep(0, 10, 1)
	assert.Equal(t, 0, s.min)
	assert.Equal(t, 10, s.max)
	assert.Equal(t, 1, s.step)
}

func TestSetMinMaxStep_BadArgs(t *testing.T) {
	s := NewSpinner(1, 10, 1, nil)
	assert.Panics(t, func() { s.SetMinMaxStep(11, 10, 2) })
	assert.Panics(t, func() { s.SetMinMaxStep(1, 10, 10) })
	assert.Panics(t, func() { s.SetMinMaxStep(1, 10, -1) })
}

func TestSetMinMaxStep_OutsideRange(t *testing.T) {
	s := NewSpinner(-2, 20, 1, nil)
	s.SetValue(19)
	s.SetMinMaxStep(-1, 10, 1)
	assert.Equal(t, 10, s.GetValue())
	s.SetValue(-1)
	s.SetMinMaxStep(1, 10, 1)
	assert.Equal(t, 1, s.GetValue())
}

func TestSetMinMaxStep_DataAboveRange(t *testing.T) {
	data := binding.NewInt()
	s := NewSpinnerWithData(-2, 20, 1, data)
	data.Set(19)
	waitForBinding()
	assert.Equal(t, 19, s.GetValue())
	s.SetMinMaxStep(-1, 10, 1)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 10, val)
}

func TestSetMinMaxStep_DataBelowRange(t *testing.T) {
	data := binding.NewInt()
	s := NewSpinnerWithData(-2, 20, 1, data)
	data.Set(-2)
	waitForBinding()
	assert.Equal(t, -2, s.GetValue())
	s.SetMinMaxStep(-1, 10, 1)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, -1, val)
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
	s.focused = true
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
	s.focused = true
	s.Disable()
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 1, s.GetValue())
}

func TestScrolled_NotFocused(t *testing.T) {
	s := NewSpinner(1, 10, 1, nil)
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 1, s.GetValue())
}

func TestOnChanged(t *testing.T) {
	var v int
	s := NewSpinner(1, 10, 1, func(newVal int) {
		v = newVal
	})
	s.SetValue(3)
	assert.Equal(t, 3, v)
}

func TestOnChanged_Disabled(t *testing.T) {
	var v int
	s := NewSpinner(1, 10, 1, func(newVal int) {
		v = newVal
	})
	s.Disable()
	s.SetValue(3)
	assert.Equal(t, 1, v)
}

func TestBinding_OutsideRange(t *testing.T) {
	val := binding.NewInt()
	s := NewSpinnerWithData(1, 5, 2, val)
	waitForBinding()
	err := val.Set(7)
	assert.NoError(t, err)
	waitForBinding()

	assert.Equal(t, 5, s.GetValue())

	v, err := val.Get()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
}
