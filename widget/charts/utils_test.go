package charts

import (
	"image"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

type rasterizer interface {
	rasterize(int, int) image.Image
	Size() fyne.Size
}

func makeRasterize(win fyne.Window, graph rasterizer) image.Image {
	win.Resize(fyne.NewSize(500, 300))
	img := graph.rasterize(int(graph.Size().Width), int(graph.Size().Height))
	return img
}

func assertSize(t *testing.T, img image.Image, graph rasterizer) {
	assert.Greater(t, img.Bounds().Size().X, 0)
	assert.Greater(t, img.Bounds().Size().Y, 0)
	assert.Equal(t, img.Bounds().Size().X, int(graph.Size().Width))
	assert.Equal(t, img.Bounds().Size().Y, int(graph.Size().Height))
}
