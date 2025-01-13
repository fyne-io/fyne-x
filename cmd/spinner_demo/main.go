package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	a := app.New()
	spinner := xwidget.NewSpinner()
	//	sButton := xwidget.NewSpinnerButton(nil, stapped)
	c := container.NewCenter(spinner)
	b := widget.NewButton("Button", nil)
	v := container.NewVBox(c, b)
	w := a.NewWindow("SpinnerDemo")
	w.Resize(fyne.NewSize(200, 200))
	w.SetContent(v)
	w.ShowAndRun()
}

func stapped() {
	fmt.Println("sbutton tapped")
}
