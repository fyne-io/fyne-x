//go:build (tinygo || noos) && !nodisplay

package fynex

import (
	"image"
	"image/draw"

	"fyne.io/fyne/v2"
	"github.com/sago35/tinydisplay"
)

var display *tinydisplay.Client

func nextKey() uint16 {
	return display.GetPressedKey()
}

func (t *tinygo) Render(img image.Image) {
	// tinydisplay does not like our *image.NRGBA, so copy it
	out := image.NewRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.ZP, draw.Src)

	display.SetImage(out)
}

func (t *tinygo) ScreenSize() fyne.Size {
	display, _ = tinydisplay.NewClient("127.0.0.1", 9812, 0, 0)
	w, h := display.Size()

	return fyne.NewSize(float32(w), float32(h))
}
