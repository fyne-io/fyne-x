package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	spinner := xwidget.NewSpinner()
	c := container.NewCenter(spinner)

	a := app.New()
	w := a.NewWindow("SpinnerDemo")
	w.SetContent(c)
	w.Resize(fyne.NewSize(200, 200))
	w.ShowAndRun()

}
