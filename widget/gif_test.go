package widget

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
)

func TestNewAnimatedGif(t *testing.T) {
	gif, err := NewAnimatedGif(storage.NewFileURI("./testdata/gif/earth.gif"))
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
	f, err := os.Open("./testdata/gif/earth.gif")
	assert.Nil(t, err)

	r, err := io.ReadAll(f)
	assert.Nil(t, err)

	res := fyne.NewStaticResource("earth.gif", r)
	gif, _ := NewAnimatedGifFromResource(res)
	assert.True(t, gif.min.IsZero())

	gif.SetMinSize(fyne.NewSize(10.0, 10.0))
	assert.Equal(t, float32(10), gif.MinSize().Width)
	assert.Equal(t, float32(10), gif.MinSize().Height)
}
