package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/layout"
)

func main() {
	a := app.New()
	w := a.NewWindow("Portions")

	long := widget.NewButton("A very long", nil)
	shorter := widget.NewButton("Short", nil)
	long2 := widget.NewButton("I am also long", nil)
	btn := widget.NewButton("123", nil)
	w.SetContent(container.New(layout.NewHPortion([]float64{30, 20, 30, 10}), long, shorter, long2, btn))
	w.ShowAndRun()
}
