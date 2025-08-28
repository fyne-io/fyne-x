package fynex

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/embedded"
)

type Driver interface {
	Render(image.Image)
	Run()
	ScreenSize() fyne.Size

	Details() (func(image.Image), chan embedded.Event)
}
