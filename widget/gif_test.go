package widget

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
)

func TestNewAnimatedGif(t *testing.T) {
	gif, err := NewAnimatedGif(storage.NewFileURI("./testdata/gif/minions.gif"))
	assert.Nil(t, err)

	w := test.NewWindow(gif)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(128, 128))

	test.AssertImageMatches(t, "gif/initial.png", w.Canvas().Capture())

	gif.Start()
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, -1, gif.remaining)
}

func TestAnimatedGif_MinSize(t *testing.T) {
	gif, _ := NewAnimatedGif(storage.NewFileURI("./testdata/gif/minions.gif"))
	assert.True(t, gif.min.IsZero())

	gif.SetMinSize(fyne.NewSize(10.0, 10.0))
	assert.Equal(t, float32(10), gif.MinSize().Width)
	assert.Equal(t, float32(10), gif.MinSize().Height)
}
