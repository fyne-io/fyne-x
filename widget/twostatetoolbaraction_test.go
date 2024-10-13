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
	assert.Equal(t, theme.MediaPlayIcon().Name(), action.icon0.Name())
	assert.Equal(t, theme.MediaPauseIcon().Name(), action.icon1.Name())
	assert.Equal(t, action.icon0.Name(), action.button.Icon.Name())
}

func TestTwoStateToolbarAction_Activated(t *testing.T) {
	action := NewTwoStateToolbarAction(theme.MediaPlayIcon(),
		theme.MediaPauseIcon(),
		func(_ TwoStateState) {})
	require.Equal(t, action.icon0.Name(), action.button.Icon.Name())
	action.button.Tapped(nil)
	assert.Equal(t, action.icon1.Name(), action.button.Icon.Name())
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

func TestTwoStateToolbarAction_GetSetState(t *testing.T) {
	var ts TwoStateState
	playState := OffState
	test.NewApp()
	action := NewTwoStateToolbarAction(theme.MediaPlayIcon(),
		theme.MediaPauseIcon(),
		func(state TwoStateState) {
			ts = state
		})
	tb := widget.NewToolbar(action)
	w := test.NewWindow(tb)
	defer w.Close()
	assert.Equal(t, playState, action.GetState())
	action.SetState(OnState)
	assert.Equal(t, OnState, action.GetState())
	assert.Equal(t, OnState, ts)
	test.AssertRendersToImage(t, "twostatetoolbaraction/state1.png", w.Canvas())
}

func TestTwoStateToolbarAction_SetState0Icon(t *testing.T) {
	test.NewApp()
	action := NewTwoStateToolbarAction(nil,
		theme.MediaPauseIcon(),
		func(state TwoStateState) {})
	tb := widget.NewToolbar(action)
	w := test.NewWindow(tb)
	defer w.Close()

	action.SetState0Icon(theme.MediaPlayIcon())
	assert.Equal(t, theme.MediaPlayIcon().Name(), action.icon0.Name())
}

func TestTwoStateToolbarAction_SetState1Icon(t *testing.T) {
	test.NewApp()
	action := NewTwoStateToolbarAction(theme.MediaPlayIcon(),
		nil,
		func(state TwoStateState) {})
	tb := widget.NewToolbar(action)
	w := test.NewWindow(tb)
	defer w.Close()

	action.SetState1Icon(theme.MediaPauseIcon())
	assert.Equal(t, theme.MediaPauseIcon().Name(), action.icon1.Name())
}
