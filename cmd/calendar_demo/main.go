package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	xWidget "fyne.io/x/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Calendar")

	i := widget.NewLabel("Please Choose a Date")
	i.Alignment = fyne.TextAlignCenter
	l := widget.NewLabel("")
	l.Alignment = fyne.TextAlignCenter
	d := &date{instruction: i, dateChosen: l}

	calendar := xWidget.NewCalendar(time.Now(), d.onSelected, 32)

	c := container.NewVBox(i, l, calendar)
	w.SetContent(c)
	w.ShowAndRun()
}

type date struct {
	instruction *widget.Label
	dateChosen  *widget.Label
}

func (d *date) onSelected(t time.Time) {
	// use time object to set text on label with given format
	d.instruction.SetText("Date Selected:")
	d.dateChosen.SetText(t.Format("Mon 02 Jan 2006"))
}
