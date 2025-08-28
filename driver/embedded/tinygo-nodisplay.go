//go:build (tinygo || noos) && nodisplay

package fynex

import (
	"image"

	"fyne.io/fyne/v2"
)

func nextKey() uint16 {
	return noKey
}

func (t *tinygo) Render(_ image.Image) {
}

func (t *tinygo) ScreenSize() fyne.Size {
	return fyne.NewSize(320, 240)
}
