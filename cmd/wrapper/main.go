package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/wrapper"
)

func main() {
	app := app.New()
	win := app.NewWindow("Wrapper example")
	win.Resize(fyne.NewSize(400, 400))

	label1 := widget.NewLabel("Label 1, not wrapped")
	label2 := widget.NewLabel("Label 2, click me")
	label3 := widget.NewLabel("Label 3, move mouse over me")
	positionLabel := widget.NewLabel("Informations will be displayed here")

	wrapped := wrapper.MakeTappable(label2, func(e *fyne.PointEvent) {
		dialog.ShowInformation("Tapped", "Label 1 was tapped", win)
	})

	mousable := wrapper.MakeHoverable(label3,
		func(e *desktop.MouseEvent) {
			positionLabel.SetText("Mouse in")
		}, func(e *desktop.MouseEvent) {
			posx := strconv.FormatFloat(float64(e.Position.X), 'f', 2, 32)
			posy := strconv.FormatFloat(float64(e.Position.Y), 'f', 2, 32)
			positionLabel.SetText("Mouse moved at:" + posx + ", " + posy)
		}, func() {
			positionLabel.SetText("Mouse out")
		},
	)

	mainContainer := container.NewGridWithColumns(2, label1, wrapped, mousable, positionLabel)

	win.SetContent(mainContainer)
	win.ShowAndRun()
}
