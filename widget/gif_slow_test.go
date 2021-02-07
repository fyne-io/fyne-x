// +build slowtests

package widget

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/storage"
)

func TestNewAnimatedGif_Once(t *testing.T) {
	gif, err := NewAnimatedGif(storage.NewFileURI("./testdata/gif/minions-once.gif"))
	assert.Nil(t, err)

	gif.Start()
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, 1, gif.remaining)
	time.Sleep(time.Second * 2)
	assert.Equal(t, 0, gif.remaining)
}

func TestNewAnimatedGif_RunTwice(t *testing.T) {
	gif, err := NewAnimatedGif(storage.NewFileURI("./testdata/gif/minions-once.gif"))
	assert.Nil(t, err)

	gif.Start()
	time.Sleep(time.Millisecond * 10)
	assert.True(t, gif.running)
	time.Sleep(time.Second * 2)
	assert.False(t, gif.running)

	gif.Start()
	gif.Start()
	time.Sleep(time.Millisecond * 10)
	assert.True(t, gif.running)
	time.Sleep(time.Second * 2)
	assert.False(t, gif.running)
}
