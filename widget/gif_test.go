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
	time.Sleep(time.Millisecond * 100)
	assert.Less(t, gif.remaining, -1)
}

func TestNewAnimatedGif_Once(t *testing.T) {
	gif, err := NewAnimatedGif(storage.NewFileURI("./testdata/gif/minions-once.gif"))
	assert.Nil(t, err)

	gif.Start()
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, 1, gif.remaining)
	time.Sleep(time.Second * 2)
	assert.Equal(t, 0, gif.remaining)
}
