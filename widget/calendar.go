package widget

import (
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*calendarLayout)(nil)

const (
	daysPerWeek           = 7
	maxWeeksPerMonth      = 6
	calendarSelectColor   = widget.HighImportance
	calendarUnselectColor = widget.LowImportance

	CalendarSingle = iota
	CalendarMulti  = iota
	CalendarRange  = iota
)

type calendarLayout struct {
	cellSize fyne.Size
}

func newCalendarLayout() fyne.Layout {
	return &calendarLayout{}
}

// Get the leading edge position of a grid cell.
// The row and col specify where the cell is in the calendar.
func (g *calendarLayout) getLeading(row, col int) fyne.Position {
	x := (g.cellSize.Width) * float32(col)
	y := (g.cellSize.Height) * float32(row)

	return fyne.NewPos(float32(math.Round(float64(x))), float32(math.Round(float64(y))))
}

// Get the trailing edge position of a grid cell.
// The row and col specify where the cell is in the calendar.
func (g *calendarLayout) getTrailing(row, col int) fyne.Position {
	return g.getLeading(row+1, col+1)
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *calendarLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	weeks := 1
	day := 0
	for i, child := range objects {
		if !child.Visible() {
			continue
		}

		if day%daysPerWeek == 0 && i >= daysPerWeek {
			weeks++
		}
		day++
	}

	g.cellSize = fyne.NewSize(size.Width/float32(daysPerWeek),
		size.Height/float32(weeks))
	row, col := 0, 0
	i := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		lead := g.getLeading(row, col)
		trail := g.getTrailing(row, col)
		child.Move(lead)
		child.Resize(fyne.NewSize(trail.X, trail.Y).Subtract(lead))

		if (i+1)%daysPerWeek == 0 {
			row++
			col = 0
		} else {
			col++
		}
		i++
	}
}

// MinSize sets the minimum size for the calendar
func (g *calendarLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	pad := theme.Padding()
	largestMin := widget.NewLabel("22").MinSize()
	return fyne.NewSize(largestMin.Width*daysPerWeek+pad*(daysPerWeek-1),
		largestMin.Height*maxWeeksPerMonth+pad*(maxWeeksPerMonth-1))
}

// Calendar creates a new date time picker which returns a time object
type Calendar struct {
	widget.BaseWidget
	currentTime time.Time

	monthPrevious *widget.Button
	monthNext     *widget.Button
	monthLabel    *widget.Label

	dates *fyne.Container

	SelectionMode int
	SelectedDates []time.Time
	OnChanged     func([]time.Time)
}

func (c *Calendar) daysOfMonth() []fyne.CanvasObject {
	start := time.Date(c.currentTime.Year(), c.currentTime.Month(), 1, 0, 0, 0, 0, c.currentTime.Location())
	buttons := []fyne.CanvasObject{}

	//account for Go time pkg starting on sunday at index 0
	dayIndex := int(start.Weekday())
	if dayIndex == 0 {
		dayIndex += daysPerWeek
	}

	//add spacers if week doesn't start on Monday
	for i := 0; i < dayIndex-1; i++ {
		buttons = append(buttons, layout.NewSpacer())
	}

	currentMonth := start.Month()
	currentYear := start.Year()
	// Stores all buttons of the current month
	dateButtons := []*widget.Button{nil}

	for d := start; d.Month() == start.Month(); d = d.AddDate(0, 0, 1) {

		dayNum := d.Day()
		s := strconv.Itoa(dayNum)
		var b *widget.Button
		b = widget.NewButton(s, func() {

			selectedDate := c.dateForButton(dayNum)

			if c.SelectionMode == CalendarSingle {
				// Unselect all currently selected buttons
				isSelected := false
				for _, t := range c.SelectedDates {
					tYear, tMonth, tDay := t.Date()

					if tYear == currentYear && tMonth == currentMonth {
						dateButtons[tDay].Importance = calendarUnselectColor
						dateButtons[tDay].Refresh()
						if tDay == dayNum {
							isSelected = true
						}
					}
				}

				// Toggle the selection of a button
				if !isSelected {
					c.SelectedDates = []time.Time{selectedDate}
					b.Importance = calendarSelectColor
					b.Refresh()
				} else {
					c.SelectedDates = c.SelectedDates[:0]
				}
			} else if c.SelectionMode == CalendarMulti {
				// Only add it to the slice if it does not contain it
				index := math.MaxInt
				for i, t := range c.SelectedDates {
					if dateEquals(t, selectedDate) {
						index = i
						break
					}
				}

				if index == math.MaxInt {
					c.SelectedDates = append(c.SelectedDates, selectedDate)
					b.Importance = calendarSelectColor
				} else {
					c.SelectedDates = append(c.SelectedDates[:index], c.SelectedDates[index+1:]...)
					b.Importance = calendarUnselectColor
				}

				// Sort the dates
				sort.Sort(&dateSort{c.SelectedDates})

				b.Refresh()
			} else if c.SelectionMode == CalendarRange {
				if len(c.SelectedDates) == 0 {
					sDay := selectedDate.Day()
					dateButtons[sDay].Importance = calendarSelectColor
					dateButtons[sDay].Refresh()

					c.SelectedDates = append(c.SelectedDates, selectedDate)
				} else if len(c.SelectedDates) == 1 {
					if dateEquals(c.SelectedDates[0], selectedDate) {
						sDay := selectedDate.Day()
						dateButtons[sDay].Importance = calendarUnselectColor
						dateButtons[sDay].Refresh()

						c.SelectedDates = c.SelectedDates[:0]
					} else {
						var beginDate, endDate time.Time

						if c.SelectedDates[0].Before(selectedDate) {
							beginDate = c.SelectedDates[0]
							endDate = selectedDate
						} else {
							beginDate = selectedDate
							endDate = c.SelectedDates[0]
						}

						// Add all other dates in between the beginning and the end
						var allDates []time.Time
						for t := beginDate; !dateEquals(t, endDate); t = t.AddDate(0, 0, 1) {
							tYear, tMonth, tDay := t.Date()
							if tYear == currentYear && tMonth == currentMonth {
								dateButtons[tDay].Importance = calendarSelectColor
								dateButtons[tDay].Refresh()
							}

							allDates = append(allDates, t)
						}
						allDates = append(allDates, endDate)

						eYear, eMonth, eDay := endDate.Date()
						if eYear == currentYear && eMonth == currentMonth {
							dateButtons[eDay].Importance = calendarSelectColor
							dateButtons[eDay].Refresh()
						}

						c.SelectedDates = allDates
					}
				} else {
					if dateEquals(c.SelectedDates[0], selectedDate) {
						// Clear the entire selection except for the last
						for _, t := range c.SelectedDates[:len(c.SelectedDates)-1] {
							tYear, tMonth, tDay := t.Date()
							if tYear == currentYear && tMonth == currentMonth {
								dateButtons[tDay].Importance = calendarUnselectColor
								dateButtons[tDay].Refresh()
							}
						}

						c.SelectedDates = []time.Time{c.SelectedDates[len(c.SelectedDates)-1]}
					} else if dateEquals(c.SelectedDates[len(c.SelectedDates)-1], selectedDate) {
						// Clear the entire selection except for the first
						for _, t := range c.SelectedDates[1:] {
							tYear, tMonth, tDay := t.Date()
							if tYear == currentYear && tMonth == currentMonth {
								dateButtons[tDay].Importance = calendarUnselectColor
								dateButtons[tDay].Refresh()
							}
						}

						c.SelectedDates = []time.Time{c.SelectedDates[0]}
					} else if c.SelectedDates[0].Before(selectedDate) && selectedDate.Before(c.SelectedDates[len(c.SelectedDates)-1]) {
						// If a date between the start and end is clicked

						index := math.MaxInt
						for i, t := range c.SelectedDates {
							if dateEquals(t, selectedDate) {
								index = i
								break
							}
						}

						// Should not happen
						if index == math.MaxInt {
							fyne.LogError("Calendar", errors.New("index of selected date could not be found in the selected dates of the calendar"))
							return
						}

						// Wether the selected date is closer to the start or end of the range
						if selectedDate.Sub(c.SelectedDates[0]) < c.SelectedDates[len(c.SelectedDates)-1].Sub(selectedDate) {
							// Change the start to the current one
							for _, t := range c.SelectedDates[:index] {
								tYear, tMonth, tDay := t.Date()
								if tYear == currentYear && tMonth == currentMonth {
									dateButtons[tDay].Importance = calendarUnselectColor
									dateButtons[tDay].Refresh()
								}
							}

							c.SelectedDates = c.SelectedDates[index:]
						} else {
							// Change the end to the current one
							for _, t := range c.SelectedDates[index+1:] {
								tYear, tMonth, tDay := t.Date()
								if tYear == currentYear && tMonth == currentMonth {
									dateButtons[tDay].Importance = calendarUnselectColor
									dateButtons[tDay].Refresh()
								}
							}

							c.SelectedDates = c.SelectedDates[:index+1]
						}
					} else {
						// if a date outside start and end is clicked

						if selectedDate.Before(c.SelectedDates[0]) {
							// Change the start to the current one
							for t := c.SelectedDates[0].AddDate(0, 0, -1); !dateEquals(t, selectedDate); t = t.AddDate(0, 0, -1) {
								tYear, tMonth, tDay := t.Date()
								if tYear == currentYear && tMonth == currentMonth {
									dateButtons[tDay].Importance = calendarSelectColor
									dateButtons[tDay].Refresh()
								}

								c.SelectedDates = append([]time.Time{t}, c.SelectedDates...)
							}

							sDay := selectedDate.Day()
							dateButtons[sDay].Importance = calendarSelectColor
							dateButtons[sDay].Refresh()

							c.SelectedDates = append([]time.Time{selectedDate}, c.SelectedDates...)
						} else {
							// Change the end to the current one
							for t := c.SelectedDates[len(c.SelectedDates)-1].AddDate(0, 0, 1); !dateEquals(t, selectedDate); t = t.AddDate(0, 0, 1) {
								tYear, tMonth, tDay := t.Date()
								if tYear == currentYear && tMonth == currentMonth {
									dateButtons[tDay].Importance = calendarSelectColor
									dateButtons[tDay].Refresh()
								}

								c.SelectedDates = append(c.SelectedDates, t)
							}

							sDay := selectedDate.Day()
							dateButtons[sDay].Importance = calendarSelectColor
							dateButtons[sDay].Refresh()

							c.SelectedDates = append(c.SelectedDates, selectedDate)
						}
					}
				}
			}

			if c.OnChanged != nil {
				c.OnChanged(c.SelectedDates)
			}
		})

		// Give the selected dates a different importance
		isSelected := false
		for _, t := range c.SelectedDates {
			if dateEquals(t, d) {
				isSelected = true
				break
			}
		}

		if isSelected {
			b.Importance = calendarSelectColor
		} else {
			b.Importance = calendarUnselectColor
		}

		buttons = append(buttons, b)
		dateButtons = append(dateButtons, b)
	}

	return buttons
}

func (c *Calendar) dateForButton(dayNum int) time.Time {
	oldName, off := c.currentTime.Zone()
	return time.Date(c.currentTime.Year(), c.currentTime.Month(), dayNum, c.currentTime.Hour(), c.currentTime.Minute(), 0, 0, time.FixedZone(oldName, off)).In(c.currentTime.Location())
}

func (c *Calendar) monthYear() string {
	return c.currentTime.Format("January 2006")
}

func (c *Calendar) calendarObjects() []fyne.CanvasObject {
	columnHeadings := []fyne.CanvasObject{}
	for i := 0; i < daysPerWeek; i++ {
		j := i + 1
		if j == daysPerWeek {
			j = 0
		}

		t := widget.NewLabel(strings.ToUpper(time.Weekday(j).String()[:3]))
		t.Alignment = fyne.TextAlignCenter
		columnHeadings = append(columnHeadings, t)
	}
	columnHeadings = append(columnHeadings, c.daysOfMonth()...)

	return columnHeadings
}

// CreateRenderer returns a new WidgetRenderer for this widget.
// This should not be called by regular code, it is used internally to render a widget.
func (c *Calendar) CreateRenderer() fyne.WidgetRenderer {
	c.monthPrevious = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		c.currentTime = c.currentTime.AddDate(0, -1, 0)
		// Dates are 'normalised', forcing date to start from the start of the month ensures move from March to February
		c.currentTime = time.Date(c.currentTime.Year(), c.currentTime.Month(), 1, 0, 0, 0, 0, c.currentTime.Location())
		c.monthLabel.SetText(c.monthYear())
		c.dates.Objects = c.calendarObjects()
	})
	c.monthPrevious.Importance = widget.LowImportance

	c.monthNext = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		c.currentTime = c.currentTime.AddDate(0, 1, 0)
		c.monthLabel.SetText(c.monthYear())
		c.dates.Objects = c.calendarObjects()
	})
	c.monthNext.Importance = widget.LowImportance

	c.monthLabel = widget.NewLabel(c.monthYear())

	nav := container.New(layout.NewBorderLayout(nil, nil, c.monthPrevious, c.monthNext),
		c.monthPrevious, c.monthNext, container.NewCenter(c.monthLabel))

	c.dates = container.New(newCalendarLayout(), c.calendarObjects()...)

	dateContainer := container.NewBorder(nav, nil, nil, nil, c.dates)

	return widget.NewSimpleRenderer(dateContainer)
}

// NewCalendar creates a calendar instance
func NewCalendar(cT time.Time, selectionMode int, onChanged func([]time.Time)) *Calendar {
	c := &Calendar{
		currentTime:   cT,
		SelectionMode: selectionMode,
		OnChanged:     onChanged,
	}

	c.ExtendBaseWidget(c)

	return c
}

func (c *Calendar) ClearSelection() {
	c.SelectedDates = c.SelectedDates[:0]

	for _, d := range c.dates.Objects {
		if b, ok := d.(*widget.Button); ok {
			b.Importance = calendarUnselectColor
		}
	}

	if c.OnChanged != nil {
		c.OnChanged(c.SelectedDates)
	}
}

func (c *Calendar) Refresh() {
	cYear, cMonth, _ := c.currentTime.Date()

	// Color the date buttons according to the selected dates
	for _, d := range c.dates.Objects {
		if b, ok := d.(*widget.Button); ok {
			day, err := strconv.Atoi(b.Text)
			if err != nil {
				continue
			}

			isSelected := false
			for _, t := range c.SelectedDates {
				tYear, tMonth, tDay := t.Date()
				if tDay == day && tMonth == cMonth && tYear == cYear {
					isSelected = true
					break
				}
			}

			if isSelected {
				b.Importance = calendarSelectColor
			} else {
				b.Importance = calendarUnselectColor
			}
		}
	}

	c.BaseWidget.Refresh()
}

type dateSort struct {
	dates []time.Time
}

func (d *dateSort) Len() int {
	return len(d.dates)
}

func (d *dateSort) Less(i, j int) bool {
	return d.dates[i].Before(d.dates[j])
}

func (d *dateSort) Swap(i, j int) {
	d.dates[i], d.dates[j] = d.dates[j], d.dates[i]
}

func dateEquals(t1, t2 time.Time) bool {
	d1, m1, y1 := t1.Date()
	d2, m2, y2 := t2.Date()
	return d1 == d2 && m1 == m2 && y1 == y2
}
