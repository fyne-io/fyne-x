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

	i := widget.NewLabel("OnChanged")
	i.Alignment = fyne.TextAlignCenter
	l := widget.NewLabel("")
	l.Alignment = fyne.TextAlignCenter
	ll := widget.NewLabel("OnSelected: ")
	d := &date{selectedDates: l, dateSelected: ll}

	// Defines which date you would like the calendar to start
	startingDate := time.Now()
	calendar := xwidget.NewCalendarWithMode(startingDate, d.onChanged, xwidget.CalendarSingle)
	calendar.OnSelected = d.onSelected

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
		ll,
		calendar,
		selection,
	)

	w.SetContent(c)
	w.ShowAndRun()
}

type date struct {
	selectedDates *widget.Label
	dateSelected  *widget.Label
}

func (d *date) onChanged(selectedDates []time.Time) {
	// use time object to set text on label with given format
	var str string
	for _, d := range selectedDates {
		str += fmt.Sprint(d.Format("Mon 02 Jan 2006"), "\n")
	}
	d.selectedDates.SetText(str)
}

func (d *date) onSelected(selectedDate time.Time) {
	d.dateSelected.SetText("OnSelected: " + selectedDate.Format("Mon 02 Jan 2006"))
}
