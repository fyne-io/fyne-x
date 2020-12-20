package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestNewAnimatedGif(t *testing.T) {
	gif, err := NewAnimatedGif(storage.NewFileURI("./testdata/minions.gif"))
	assert.Nil(t, err)

	w := test.NewWindow(gif)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(128, 128))

	test.AssertImageMatches(t, "gif/initial.png", w.Canvas().Capture())
}
