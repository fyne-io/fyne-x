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
	w.SetContent(CreateDateInputWidget())
	w.ShowAndRun()
}


func CreateDateInputWidget() *fyne.Container {

	dateEntry := xwidget.NewMyDateEntry()
	label := widget.NewLabel("Current Selection")
	button := widget.NewButton("Show", func() { on_show_click(dateEntry, label) } )

	helpText := `
	1. Press UP/Down Arrow to change the part of date.
	2. Space to set current date.
    3. Delete to clear
    4. Press enter key to update the date
	4. You can enter part of date (like Only day, day-month)
    5. *Assuming first part is Day*
    `
	helpTextWidget := widget.NewTextGridFromString(helpText)

	return container.NewVBox(dateEntry, button, label, widget.NewSeparator(), helpTextWidget)

}

func on_show_click(d *xwidget.MyDateEntry, l *widget.Label) {
	var msg string
	if d.ToDate().IsZero() == true {
		msg = "No Input"
	} else {
		msg = fmt.Sprintf("Your input is : %s", d.ToDate().Format(time.DateOnly))
	}
	l.SetText(msg)
}
