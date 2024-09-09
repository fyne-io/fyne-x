package widget

import (
	"testing"

	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestNewTwoStateToolbarAction(t *testing.T) {
	action := NewTwoStateToolbarAction(theme.MediaPlayIcon(),
		theme.MediaPauseIcon(),
		func(_ TwoStateState) {})
	assert.Equal(t, theme.MediaPlayIcon().Name(), action.Icon0.Name())
	assert.Equal(t, theme.MediaPauseIcon().Name(), action.Icon1.Name())
	assert.Equal(t, action.Icon0.Name(), action.button.Icon.Name())
}
