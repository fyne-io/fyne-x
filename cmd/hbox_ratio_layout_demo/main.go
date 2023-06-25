package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/app"	
	"fyne.io/x/fyne/layout"
	
)


func main() {
	a := app.New()
	w := a.NewWindow("HBoxLayoutRatio")
	w.SetContent(CreateHBoxRatioLayoutWidget())
	w.ShowAndRun()
}



func CreateHBoxRatioLayoutWidget() *fyne.Container {
	lbl1 := widget.NewLabel("Name:")
	entry1 := widget.NewEntry()
	lbl2 := widget.NewLabel("Number:")
	entry2 := widget.NewEntry()
	btn  := widget.NewButton("okay", nil)
	content := container.New(layout.NewHBoxRatioLayout([]float32{10, 20, 10, 20, 5, 15}), lbl1, entry1, lbl2, entry2, widget.NewSeparator(), btn)
	return content
}
