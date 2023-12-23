package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	w := app.New().NewWindow("Map Widget")

	m := xwidget.NewMapWithOptions(
		xwidget.WithOsmTiles(),
		xwidget.WithZoomButtons(true),
		xwidget.WithScrollButtons(true),
	)
	m.ZoomIn()
	w.SetContent(m)

	w.SetPadded(false)
	w.Resize(fyne.NewSize(512, 320))
	w.ShowAndRun()
}
