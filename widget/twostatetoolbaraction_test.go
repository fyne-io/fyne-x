package widget

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTwoStateToolbarAction(t *testing.T) {
	action := NewTwoStateToolbarAction(theme.MediaPlayIcon(),
		theme.MediaPauseIcon(),
		func(_ TwoStateState) {})
	assert.Equal(t, theme.MediaPlayIcon().Name(), action.Icon0.Name())
	assert.Equal(t, theme.MediaPauseIcon().Name(), action.Icon1.Name())
	assert.Equal(t, action.Icon0.Name(), action.button.Icon.Name())
}

func TestTwoStateToolbarAction_Activated(t *testing.T) {
	action := NewTwoStateToolbarAction(theme.MediaPlayIcon(),
		theme.MediaPauseIcon(),
		func(_ TwoStateState) {})
	require.Equal(t, action.Icon0.Name(), action.button.Icon.Name())
	action.button.Tapped(nil)
	assert.Equal(t, action.Icon1.Name(), action.button.Icon.Name())
}

func TestTwoStateToolbarAction_Tapped(t *testing.T) {
	test.NewApp()
	action := NewTwoStateToolbarAction(theme.MediaPlayIcon(),
		theme.MediaPauseIcon(),
		func(_ TwoStateState) {})
	tb := widget.NewToolbar(action)
	w := test.NewWindow(tb)
	defer w.Close()
	test.AssertRendersToImage(t, "twostatetoolbaraction/state0.png", w.Canvas())
	action.button.Tapped(nil)
	test.AssertRendersToImage(t, "twostatetoolbaraction/state1.png", w.Canvas())
}
