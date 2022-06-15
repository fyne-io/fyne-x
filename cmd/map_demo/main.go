package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Map Widget")

	m := xwidget.NewMap()
	m.EnableOsmDisclaimer = true
	m.EnableMoveButtons = false
	m.EnableZoomButtons = false
	w.SetContent(m)

	w.SetPadded(false)
	w.Resize(fyne.NewSize(256, 256))
	w.ShowAndRun()
}
