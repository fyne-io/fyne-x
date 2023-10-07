package main

import (
	"fmt"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/app"
	xwidget "fyne.io/x/fyne/widget"
)


func main() {
	a := app.New()
	w := a.NewWindow("Date Input Widget")
	w.SetContent(CreateDateInputWidget(a, w))
	w.ShowAndRun()
}


func CreateDateInputWidget(app fyne.App, win fyne.Window) *fyne.Container {

	dateEntry := xwidget.NewJDateInputWidget(app, win)
	label := widget.NewLabel("Current Selection")
	button := widget.NewButton("Show", func() { on_show_click(dateEntry, label) } )

	helpText := `
	1. Press UP/Down/U/D/PgDw/PgUp to change part of a date.
	2. Left/Right/R/L to move cursor to the part of a date.
	2. Space to set current date.
    `
	helpTextWidget := widget.NewTextGridFromString(helpText)

	return container.NewVBox(dateEntry, button, label, widget.NewSeparator(), helpTextWidget)

}

func on_show_click(d *xwidget.JDateEntry, l *widget.Label) {
	var msg string
	if d.GetDate().IsZero() == true {
		msg = "No Input"
	} else {
		msg = fmt.Sprintf("Your input is : %s", d.GetDate().Format(time.DateOnly))
	}
	l.SetText(msg)
}
