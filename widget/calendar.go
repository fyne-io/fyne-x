package widget

import (
	"fmt"
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

const daysPerWeek int = 7

type calendarLayout struct {
	cellSize float32
}

func newCalendarLayout() fyne.Layout {
	return &calendarLayout{}
}

// Get the leading (top or left) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func (g *calendarLayout) getLeading(offset int) float32 {
	ret := (g.cellSize) * float32(offset)

	return float32(math.Round(float64(ret)))
}

// Get the trailing (bottom or right) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func (g *calendarLayout) getTrailing(offset int) float32 {
	return g.getLeading(offset + 1)
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *calendarLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	g.cellSize = size.Width / float32(daysPerWeek)
	row, col := 0, 0
	i := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		x1 := g.getLeading(col)
		y1 := g.getLeading(row)
		x2 := g.getTrailing(col)
		y2 := g.getTrailing(row)

		child.Move(fyne.NewPos(x1, y1))
		child.Resize(fyne.NewSize(x2-x1, y2-y1))

		if (i+1)%daysPerWeek == 0 {
			row++
			col = 0
		} else {
			col++
		}
		i++
	}
}

//MinSize sets the minimum size for the calendar
func (g *calendarLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(250, 250)
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
	fmt.Println(start)
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
	return c.currentTime.Month().String() + " " + strconv.Itoa(c.currentTime.Year())
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

	dateContainer := container.NewVBox(nav, c.dates)

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
