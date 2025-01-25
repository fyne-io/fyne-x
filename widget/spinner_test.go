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

func TestNewIntSpinner(t *testing.T) {
	s := NewIntSpinner(1, 5, 2, nil)
	assert.Equal(t, 1, s.min)
	assert.Equal(t, 5, s.max)
	assert.Equal(t, 2, s.step)
	assert.Equal(t, 1, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestNewIntSpinner_BadArgs(t *testing.T) {
	assert.Panics(t, func() { NewIntSpinner(5, 5, 1, nil) }, "Did not panic with min == max")
	assert.Panics(t, func() { NewIntSpinner(5, 4, 1, nil) }, "Did not panic with min < max")
	assert.Panics(t, func() { NewIntSpinner(1, 5, 0, nil) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewIntSpinner(1, 5, -5, nil) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewIntSpinner(1, 5, 5, nil) }, "Did not panic with step > max - min")
}

func TestNewIntSpinnerWithData(t *testing.T) {
	data := binding.NewInt()
	s := NewIntSpinnerWithData(1, 5, 2, data)
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

func TestIntSpinner_Unbind(t *testing.T) {
	data := binding.NewInt()
	s := NewIntSpinnerWithData(1, 5, 2, data)
	waitForBinding()
	s.Unbind()
	s.SetValue(2)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestNewIntSpinnerWithData_BadArgs(t *testing.T) {
	boundValue := binding.NewInt()
	assert.Panics(t, func() { NewIntSpinnerWithData(5, 5, 1, boundValue) }, "Did not panic with min == max")
	assert.Panics(t, func() { NewIntSpinnerWithData(5, 4, 1, boundValue) }, "Did not panic with min < max")
	assert.Panics(t, func() { NewIntSpinnerWithData(1, 5, 0, boundValue) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewIntSpinnerWithData(1, 5, -5, boundValue) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewIntSpinnerWithData(1, 5, 5, boundValue) }, "Did not panic with step > max - min")
}

func TestNewIntSpinnerUninitialized(t *testing.T) {
	s := NewIntSpinnerUninitialized()
	assert.True(t, s.Disabled())
	s.Enable()
	assert.True(t, s.Disabled())
	s.SetMinMaxStep(-1, 2, 1)
	assert.True(t, s.Disabled())
	s.Enable()
	assert.False(t, s.Disabled())
}

func TestIntSpinner_SetValue(t *testing.T) {
	s := NewIntSpinner(1, 5, 2, nil)
	s.SetValue(2)
	assert.Equal(t, 2, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestIntSpinner_SetValue_LessThanMin(t *testing.T) {
	s := NewIntSpinner(4, 22, 5, nil)
	s.SetValue(3)
	assert.Equal(t, 4, s.GetValue())
	assert.True(t, s.downButton.Disabled())
	assert.False(t, s.upButton.Disabled())
}

func TestIntSpinner_SetValue_GreaterThanMax(t *testing.T) {
	s := NewIntSpinner(4, 22, 5, nil)
	s.SetValue(23)
	assert.Equal(t, 22, s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestIntSpinner_SetValue_DisabledSpinner(t *testing.T) {
	s := NewIntSpinner(4, 22, 5, nil)
	s.Disable()
	s.SetValue(10)
	assert.Equal(t, 4, s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestIntSpinner_SetMinMaxStep(t *testing.T) {
	s := NewIntSpinner(1, 6, 2, nil)
	s.SetMinMaxStep(0, 10, 1)
	assert.Equal(t, 0, s.min)
	assert.Equal(t, 10, s.max)
	assert.Equal(t, 1, s.step)
}

func TestIntSpinner_SetMinMaxStep_BadArgs(t *testing.T) {
	s := NewIntSpinner(1, 10, 1, nil)
	assert.Panics(t, func() { s.SetMinMaxStep(11, 10, 2) })
	assert.Panics(t, func() { s.SetMinMaxStep(1, 10, 10) })
	assert.Panics(t, func() { s.SetMinMaxStep(1, 10, -1) })
}

func TestIntSpinner_SetMinMaxStep_OutsideRange(t *testing.T) {
	s := NewIntSpinner(-2, 20, 1, nil)
	s.SetValue(19)
	s.SetMinMaxStep(-1, 10, 1)
	assert.Equal(t, 10, s.GetValue())
	s.SetValue(-1)
	s.SetMinMaxStep(1, 10, 1)
	assert.Equal(t, 1, s.GetValue())
}

func TestIntSpinner_SetMinMaxStep_DataAboveRange(t *testing.T) {
	data := binding.NewInt()
	s := NewIntSpinnerWithData(-2, 20, 1, data)
	data.Set(19)
	waitForBinding()
	assert.Equal(t, 19, s.GetValue())
	s.SetMinMaxStep(-1, 10, 1)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 10, val)
}

func TestIntSpinner_SetMinMaxStep_DataBelowRange(t *testing.T) {
	data := binding.NewInt()
	s := NewIntSpinnerWithData(-2, 20, 1, data)
	data.Set(-2)
	waitForBinding()
	assert.Equal(t, -2, s.GetValue())
	s.SetMinMaxStep(-1, 10, 1)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, -1, val)
}

func TestIntSpinner_UpButtonTapped(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
	s.upButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 9, s.GetValue())
	s.upButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 10, s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestIntSpinner_DownButtonTapped(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
	s.SetValue(10)
	s.downButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 5, s.GetValue())
	s.downButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 4, s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestIntSpinner_EnableDisabledSpinner(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
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

func TestIntSpinner_EnableDisabledSpinner_UpButtonDisabled(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
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

func TestIntSpinner_EnableDisabledSpinner_DownButtonDisabled(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
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

func TestIntSpinner_RunePlus(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
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

func TestIntSpinner_RuneMinus(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
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

func TestIntSpinner_RunePlus_SpinnerDisabled(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
	s.focused = true
	s.Disable()
	s.TypedRune('+')
	assert.Equal(t, 4, s.GetValue())
}

func TestIntSpinner_RuneMinus_SpinnerDisabled(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
	s.focused = true
	s.SetValue(8)
	s.Disable()
	s.TypedRune('-')
	assert.Equal(t, 8, s.GetValue())
}

func TestIntSpinner_RunePlus_SpinnerNotFocused(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
	s.focused = false
	s.TypedRune('+')
	assert.Equal(t, 4, s.GetValue())
}

func TestIntSpinner_RuneMinus_SpinnerNotFocused(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
	s.focused = true
	s.SetValue(8)
	s.focused = false
	s.TypedRune('-')
	assert.Equal(t, 8, s.GetValue())
}

func TestIntSpinner_KeyUp(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
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

func TestIntSpinner_KeyDown(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
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

func TestIntSpinner_KeyUp_SpinnerDisabled(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
	s.focused = true
	s.Disable()
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 4, s.GetValue())
}

func TestIntSpinner_KeyDown_SpinnerDisabled(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
	s.focused = true
	s.SetValue(8)
	s.Disable()
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 8, s.GetValue())
}

func TestIntSpinner_KeyUp_SpinnerNotFocused(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
	s.focused = false
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 4, s.GetValue())
}

func TestIntSpinner_KeyDown_SpinnerNotFocused(t *testing.T) {
	s := NewIntSpinner(4, 10, 5, nil)
	s.focused = true
	s.SetValue(8)
	s.focused = false
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 8, s.GetValue())
}

func TestIntSpinner_Scrolled(t *testing.T) {
	s := NewIntSpinner(1, 10, 1, nil)
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

func TestIntSpinner_Scrolled_Disabled(t *testing.T) {
	s := NewIntSpinner(1, 10, 1, nil)
	s.focused = true
	s.Disable()
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 1, s.GetValue())
}

func TestIntSpinner_Scrolled_NotFocused(t *testing.T) {
	s := NewIntSpinner(1, 10, 1, nil)
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 1, s.GetValue())
}

func TestIntSpinner_OnChanged(t *testing.T) {
	var v int
	s := NewIntSpinner(1, 10, 1, func(newVal int) {
		v = newVal
	})
	s.SetValue(3)
	assert.Equal(t, 3, v)
}

func TestIntSpinner_OnChanged_Disabled(t *testing.T) {
	var v int
	s := NewIntSpinner(1, 10, 1, func(newVal int) {
		v = newVal
	})
	s.Disable()
	s.SetValue(3)
	assert.Equal(t, 1, v)
}

func TestIntSpinner_Binding_OutsideRange(t *testing.T) {
	val := binding.NewInt()
	s := NewIntSpinnerWithData(1, 5, 2, val)
	waitForBinding()
	err := val.Set(7)
	assert.NoError(t, err)
	waitForBinding()

	assert.Equal(t, 5, s.GetValue())

	v, err := val.Get()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
}

func TestNewFloat64Spinner(t *testing.T) {
	s := NewFloat64Spinner(1., 5., 1.5, 2, nil)
	assert.Equal(t, 1., s.min)
	assert.Equal(t, 5., s.max)
	assert.Equal(t, 1.5, s.step)
	assert.Equal(t, uint(2), s.precision)
	assert.Equal(t, 1., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestNewFloat64Spinner_BadArgs(t *testing.T) {
	assert.Panics(t, func() { NewFloat64Spinner(5., 5., 1., 1, nil) }, "Did not panic with min == max")
	assert.Panics(t, func() { NewFloat64Spinner(5., 4., 1., 1, nil) }, "Did not panic with min < max")
	assert.Panics(t, func() { NewFloat64Spinner(1., 5., 0., 2, nil) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewFloat64Spinner(1., 5., -5., 1, nil) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewFloat64Spinner(1., 5., 5., 1, nil) }, "Did not panic with step > max - min")
}

func TestNewFloat64SpinnerWithData(t *testing.T) {
	data := binding.NewFloat()
	s := NewFloat64SpinnerWithData(1., 5., 2., 1, data)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1., val)

	s.SetValue(1.52)
	waitForBinding()
	val, err = data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1.52, val)

	err = data.Set(3.1)
	assert.NoError(t, err)
	waitForBinding()
	assert.Equal(t, 3.1, s.GetValue())
}

func TestFloat64Spinner_Unbind(t *testing.T) {
	data := binding.NewFloat()
	s := NewFloat64SpinnerWithData(1., 5., 2., 1, data)
	waitForBinding()
	s.Unbind()
	s.SetValue(2.)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1., val)
}

func TestNewFloat64SpinnerWithData_BadArgs(t *testing.T) {
	boundValue := binding.NewFloat()
	assert.Panics(t, func() { NewFloat64SpinnerWithData(5., 5., 1., 1, boundValue) }, "Did not panic with min == max")
	assert.Panics(t, func() { NewFloat64SpinnerWithData(5., 4., 1., 1, boundValue) }, "Did not panic with min < max")
	assert.Panics(t, func() { NewFloat64SpinnerWithData(1., 5., 0., 1, boundValue) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewFloat64SpinnerWithData(1., 5., -5., 1, boundValue) }, "Did not panic with step = 0")
	assert.Panics(t, func() { NewFloat64SpinnerWithData(1., 5., 5., 1, boundValue) }, "Did not panic with step > max - min")
}

func TestNewFloat64SpinnerUninitialized(t *testing.T) {
	s := NewFloat64SpinnerUninitialized()
	assert.True(t, s.Disabled())
	s.Enable()
	assert.True(t, s.Disabled())
	s.SetMinMaxStep(-1., 2., 1.1)
	assert.True(t, s.Disabled())
	s.Enable()
	assert.False(t, s.Disabled())
}

func TestFloat64Spinner_SetValue(t *testing.T) {
	s := NewFloat64Spinner(1, 5, 2, 1, nil)
	s.SetValue(2)
	assert.Equal(t, 2., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestFloat64Spinner_SetValue_LessThanMin(t *testing.T) {
	s := NewFloat64Spinner(4, 22, 5, 0, nil)
	s.SetValue(3)
	assert.Equal(t, 4., s.GetValue())
	assert.True(t, s.downButton.Disabled())
	assert.False(t, s.upButton.Disabled())
}

func TestFloat64Spinner_SetValue_GreaterThanMax(t *testing.T) {
	s := NewFloat64Spinner(4, 22, 5, 1, nil)
	s.SetValue(23.)
	assert.Equal(t, 22., s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestFloat64Spinner_SetValue_DisabledSpinner(t *testing.T) {
	s := NewFloat64Spinner(4, 22, 5, 1, nil)
	s.Disable()
	s.SetValue(10.)
	assert.Equal(t, 4., s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestFloat64Spinner_SetMinMaxStep(t *testing.T) {
	s := NewFloat64Spinner(1., 6., 2., 1, nil)
	s.SetMinMaxStep(0., 10., 1.)
	assert.Equal(t, 0., s.min)
	assert.Equal(t, 10., s.max)
	assert.Equal(t, 1., s.step)
}

func TestFloat64Spinner_SetMinMaxStep_BadArgs(t *testing.T) {
	s := NewFloat64Spinner(1, 10, 1, 1, nil)
	assert.Panics(t, func() { s.SetMinMaxStep(11, 10, 2) })
	assert.Panics(t, func() { s.SetMinMaxStep(1, 10, 10) })
	assert.Panics(t, func() { s.SetMinMaxStep(1, 10, -1) })
}

func TestFloat64Spinner_SetMinMaxStep_OutsideRange(t *testing.T) {
	s := NewFloat64Spinner(-2, 20, 1, 1, nil)
	s.SetValue(19.)
	s.SetMinMaxStep(-1., 10., 1.2)
	assert.Equal(t, 10., s.GetValue())
	s.SetValue(-1.)
	s.SetMinMaxStep(1., 10., 1.)
	assert.Equal(t, 1., s.GetValue())
}

func TestFloat64Spinner_SetMinMaxStep_DataAboveRange(t *testing.T) {
	data := binding.NewFloat()
	s := NewFloat64SpinnerWithData(-2, 20, 1, 1, data)
	data.Set(19.)
	waitForBinding()
	assert.Equal(t, 19., s.GetValue())
	s.SetMinMaxStep(-1., 10., 1.)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 10., val)
}

func TestFloat64Spinner_SetMinMaxStep_DataBelowRange(t *testing.T) {
	data := binding.NewFloat()
	s := NewFloat64SpinnerWithData(-2, 20, 1, 2, data)
	data.Set(-2.)
	waitForBinding()
	assert.Equal(t, -2., s.GetValue())
	s.SetMinMaxStep(-1., 10., 1.)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, -1., val)
}

func TestFloat64Spinner_UpButtonTapped(t *testing.T) {
	s := NewFloat64Spinner(4., 10., 5., 1, nil)
	s.upButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 9., s.GetValue())
	s.upButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 10., s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestFloat64Spinner_DownButtonTapped(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.SetValue(10.)
	s.downButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 5., s.GetValue())
	s.downButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 4., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestFloa4Spinner_EnableDisabledSpinner(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.SetValue(7.)
	s.Disable()
	assert.True(t, s.Disabled())
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
	s.Enable()
	assert.False(t, s.Disabled())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestFloat64Spinner_EnableDisabledSpinner_UpButtonDisabled(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.SetValue(10.)
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
	s.Disable()
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())

	s.Enable()
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestFloat64Spinner_EnableDisabledSpinner_DownButtonDisabled(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.SetValue(4.)
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
	s.Disable()
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())

	s.Enable()
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestFloat64Spinner_RunePlus(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = true
	s.TypedRune('+')
	assert.Equal(t, 9., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedRune('+')
	assert.Equal(t, 10., s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestFloat64Spinner_RuneMinus(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = true
	s.SetValue(10.)
	s.TypedRune('-')
	assert.Equal(t, 5., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedRune('-')
	assert.Equal(t, 4., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestFloat64Spinner_RunePlus_SpinnerDisabled(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = true
	s.Disable()
	s.TypedRune('+')
	assert.Equal(t, 4., s.GetValue())
}

func TestFloat64Spinner_RuneMinus_SpinnerDisabled(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = true
	s.SetValue(8.)
	s.Disable()
	s.TypedRune('-')
	assert.Equal(t, 8., s.GetValue())
}

func TestFloat64Spinner_RunePlus_SpinnerNotFocused(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = false
	s.TypedRune('+')
	assert.Equal(t, 4., s.GetValue())
}

func TestFloat64Spinner_RuneMinus_SpinnerNotFocused(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = true
	s.SetValue(8)
	s.focused = false
	s.TypedRune('-')
	assert.Equal(t, 8., s.GetValue())
}

func TestFloat64Spinner_KeyUp(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = true
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 9., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedKey(&key)
	assert.Equal(t, 10., s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestFloat64Spinner_KeyDown(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = true
	s.SetValue(10)
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 5., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedKey(&key)
	assert.Equal(t, 4., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestFloat64Spinner_KeyUp_SpinnerDisabled(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = true
	s.Disable()
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 4., s.GetValue())
}

func TestFloat64Spinner_KeyDown_SpinnerDisabled(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = true
	s.SetValue(8.)
	s.Disable()
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 8., s.GetValue())
}

func TestFloat64Spinner_KeyUp_SpinnerNotFocused(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = false
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 4., s.GetValue())
}

func TestFloat64Spinner_KeyDown_SpinnerNotFocused(t *testing.T) {
	s := NewFloat64Spinner(4, 10, 5, 1, nil)
	s.focused = true
	s.SetValue(8.)
	s.focused = false
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 8., s.GetValue())
}

func TestFloat64Spinner_Scrolled(t *testing.T) {
	s := NewFloat64Spinner(1, 10, 1, 1, nil)
	s.focused = true
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 2., s.GetValue())
	s.Scrolled(&e)
	assert.Equal(t, 3., s.GetValue())
	delta = fyne.Delta{DX: 0, DY: -25}
	e = fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 2., s.GetValue())
}

func TestFloat64Spinner_Scrolled_Disabled(t *testing.T) {
	s := NewFloat64Spinner(1, 10, 1, 1, nil)
	s.focused = true
	s.Disable()
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 1., s.GetValue())
}

func TestFloat64Spinner_Scrolled_NotFocused(t *testing.T) {
	s := NewFloat64Spinner(1, 10, 1, 1, nil)
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 1., s.GetValue())
}

func TestFloat64Spinner_OnChanged(t *testing.T) {
	var v float64
	s := NewFloat64Spinner(1, 10, 1, 1, func(newVal float64) {
		v = newVal
	})
	s.SetValue(3.)
	assert.Equal(t, 3., v)
}

func TestFloat64Spinner_OnChanged_Disabled(t *testing.T) {
	var v float64
	s := NewFloat64Spinner(1, 10, 1, 1, func(newVal float64) {
		v = newVal
	})
	s.Disable()
	s.SetValue(3.)
	assert.Equal(t, 1., v)
}

func TestFloat64Spinner_Binding_OutsideRange(t *testing.T) {
	val := binding.NewFloat()
	s := NewFloat64SpinnerWithData(1, 5, 2, 1, val)
	waitForBinding()
	err := val.Set(7)
	assert.NoError(t, err)
	waitForBinding()

	assert.Equal(t, 5., s.GetValue())

	v, err := val.Get()
	assert.NoError(t, err)
	assert.Equal(t, 5., v)
}
