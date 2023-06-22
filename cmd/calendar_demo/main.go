package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Calendar")

	i := widget.NewLabel("Please Choose a Date")
	i.Alignment = fyne.TextAlignCenter
	l := widget.NewLabel("")
	l.Alignment = fyne.TextAlignCenter
	d := &date{instruction: i, dateChosen: l}

	// Defines which date you would like the calendar to start
	startingDate := time.Now()
	calendar := xwidget.NewCalendar(startingDate, xwidget.CalendarSingle, d.onChanged)

	selection := widget.NewRadioGroup([]string{"Single", "Multi", "Range"}, func(s string) {
		calendar.ClearSelection()
		switch s {
		case "Single":
			calendar.SelectionMode = xwidget.CalendarSingle
		case "Multi":
			calendar.SelectionMode = xwidget.CalendarMulti
		case "Range":
			calendar.SelectionMode = xwidget.CalendarRange
		}
		calendar.Refresh()
	})
	selection.Horizontal = true
	selection.Required = true
	selection.Selected = "Single"

	scroll := container.NewVScroll(l)
	scroll.SetMinSize(fyne.NewSize(200, 100))

	c := container.NewVBox(
		i,
		scroll,
		calendar,
		selection,
	)

	w.SetContent(c)
	w.ShowAndRun()
}

type date struct {
	instruction *widget.Label
	dateChosen  *widget.Label
}

func (d *date) onChanged(selectedDates []time.Time) {
	// use time object to set text on label with given format
	d.instruction.SetText("Date Selected:")
	var str string
	for _, d := range selectedDates {
		str += fmt.Sprint(d.Format("Mon 02 Jan 2006"), "\n")
	}
	d.dateChosen.SetText(str)
}
