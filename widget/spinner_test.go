package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpinner(t *testing.T) {
	s := NewSpinner(1, 5, 2, nil)
	assert.Equal(t, 1, s.min)
	assert.Equal(t, 5, s.max)
	assert.Equal(t, 2, s.step)
	assert.Equal(t, 1, s.GetValue())
}

func TestSetValue(t *testing.T) {
	s := NewSpinner(1, 5, 2, nil)
	s.SetValue(2)
	assert.Equal(t, 2, s.GetValue())
}

func TestSetValue_LessThanMin(t *testing.T) {
	s := NewSpinner(4, 22, 5, nil)
	s.SetValue(3)
	assert.Equal(t, 4, s.GetValue())
}

func TestSetValue_GreaterThanMax(t *testing.T) {
	s := NewSpinner(4, 22, 5, nil)
	s.SetValue(23)
	assert.Equal(t, 22, s.GetValue())

}
