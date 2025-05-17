package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("NumericalEntry Demo")

	nE := xwidget.NewNumericalEntry()
	nE.AllowFloat = true
	nE.AllowNegative = true
	a.Clipboard().SetContent("92.f65")

	errL := widget.NewLabel("Error: ")
	eText := widget.NewLabel("")
	errC := container.NewHBox(errL, eText)

	valL := widget.NewLabel("Value as float:")
	// Should call nE.SetText
	nE.Text = "1!23.45"
	nE.OnChanged = func(s string) {
		err := nE.Validate()
		if err == nil {
			eText.Text = ""
		} else {
			eText.Text = err.Error()
		}
		eText.Refresh()
		f, _ := nE.Float()
		valL.SetText(fmt.Sprintf("Value as float: %f", f))
		valL.Refresh()
	}

	sep := widget.NewSeparator()
	setB := widget.NewButton("Set Entry to -2.35,4", func() {
		nE.SetText("-2.35,4")
	})
	appB := widget.NewButton("Append -3,6.5 to Entry", func() {
		nE.Append("-3,6.5")
	})
	pasteL := widget.NewLabel("Clipboard contains \"92.f65\".\n" +
		"Paste this into the Entry widget at\ndifferent locations to see the effect.")

	vc := container.NewVBox(nE, errC, valL, sep, setB, appB, pasteL)
	w.SetContent(vc)
	w.ShowAndRun()
}
