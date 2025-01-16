package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

var spinnerDisabled bool

func main() {
	a := app.New()
	spinner := xwidget.NewSpinner(1, 12, 3, stapped)
	c := container.NewCenter(spinner)
	b := widget.NewButton("(En/Dis)able Spinner", func() {
		spinnerDisabled = !spinnerDisabled
		if spinnerDisabled {
			spinner.Disable()
		} else {
			spinner.Enable()
		}
	})
	e := widget.NewEntry()
	v := container.NewVBox(c, b, e)
	w := a.NewWindow("SpinnerDemo")
	w.Resize(fyne.NewSize(200, 200))
	w.SetContent(v)
	w.ShowAndRun()
}

func stapped(value int) {
	fmt.Printf("sbutton tapped with value: %v\n", value)
}
