package widget

import (
	"math"
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
	daysPerWeek      = 7
	maxWeeksPerMonth = 6
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

	onSelected func(time.Time)
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

	for d := start; d.Month() == start.Month(); d = d.AddDate(0, 0, 1) {

		dayNum := d.Day()
		s := strconv.Itoa(dayNum)
		b := widget.NewButton(s, func() {

			selectedDate := c.dateForButton(dayNum)

			c.onSelected(selectedDate)
		})
		b.Importance = widget.LowImportance

		buttons = append(buttons, b)
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
func NewCalendar(cT time.Time, onSelected func(time.Time)) *Calendar {
	c := &Calendar{
		currentTime: cT,
		onSelected:  onSelected,
	}

	c.ExtendBaseWidget(c)

	return c
}
