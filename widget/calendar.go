package widget

import (
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*calendarLayout)(nil)

type calendarLayout struct {
	cellSize float32
	padding  float32
}

func newCalendarLayout(c float32, p float32) fyne.Layout {
	return &calendarLayout{cellSize: c, padding: p}
}

func (g *calendarLayout) countRows(objects []fyne.CanvasObject) int {
	count := 0
	for _, child := range objects {
		if child.Visible() {
			count++
		}
	}

	return int(math.Ceil(float64(count) / float64(7)))
}

// Get the leading (top or left) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func (g *calendarLayout) getLeading(offset int) float32 {
	ret := (g.cellSize + g.padding) * float32(offset)

	return float32(math.Round(float64(ret)))
}

// Get the trailing (bottom or right) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func (g *calendarLayout) getTrailing(offset int) float32 {
	return g.getLeading(offset+1) - g.padding
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *calendarLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
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

		if (i+1)%7 == 0 {
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
	rows := g.countRows(objects)
	return fyne.NewSize((float32(g.cellSize)+g.padding)*7, (float32(g.cellSize)+g.padding)*float32(rows))
}

// Calendar creates a new date time picker which returns a time object
type Calendar struct {
	widget.BaseWidget
	canvas       fyne.Canvas
	calendarTime time.Time

	monthPrevious *widget.Button
	monthNext     *widget.Button
	monthLabel    *widget.RichText

	day   int
	month int
	year  int

	dates *fyne.Container

	onSelected func(time.Time)

	cellSize float32
	padding  float32
}

func (c *Calendar) daysOfMonth() []fyne.CanvasObject {
	start, _ := time.Parse("2006-1-2", strconv.Itoa(c.year)+"-"+strconv.Itoa(c.month)+"-"+strconv.Itoa(1))

	buttons := []fyne.CanvasObject{}

	//account for Go time pkg starting on sunday at index 0
	dayIndex := int(start.Weekday())
	if dayIndex == 0 {
		dayIndex += 7
	}

	//add spacers if week doesn't start on Monday
	for i := 0; i < dayIndex-1; i++ {
		buttons = append(buttons, layout.NewSpacer())
	}

	for d := start; d.Month() == start.Month(); d = d.AddDate(0, 0, 1) {

		dayNum := d.Day()
		s := strconv.Itoa(dayNum)
		var b fyne.CanvasObject = widget.NewButton(s, func() {

			selectedDate := c.dateForButton(dayNum)

			c.onSelected(selectedDate)

			c.hideOverlay()
		})

		buttons = append(buttons, b)
	}

	return buttons
}

func (c *Calendar) dateForButton(dayNum int) time.Time {
	oldName, off := c.calendarTime.Zone()
	return time.Date(c.year, time.Month(c.month), dayNum, c.calendarTime.Hour(), c.calendarTime.Minute(), 0, 0, time.FixedZone(oldName, off)).In(c.calendarTime.Location())
}

func (c *Calendar) hideOverlay() {
	overlayList := c.canvas.Overlays().List()
	overlayList[0].Hide()
}

func (c *Calendar) monthYear() string {
	return time.Month(c.month).String() + " " + strconv.Itoa(c.year)
}

func (c *Calendar) calendarObjects() []fyne.CanvasObject {
	textSize := float32(8)
	columnHeadings := []fyne.CanvasObject{}
	for i := 0; i < 7; i++ {
		j := i + 1
		if j == 7 {
			j = 0
		}

		var canvasObject fyne.CanvasObject = canvas.NewText(strings.ToUpper(time.Weekday(j).String()[:3]), color.NRGBA{0xFF, 0xFF, 0xFF, 0xBF})
		canvasObject.(*canvas.Text).TextSize = textSize
		canvasObject.(*canvas.Text).Alignment = fyne.TextAlignCenter
		columnHeadings = append(columnHeadings, canvasObject)
	}
	columnHeadings = append(columnHeadings, c.daysOfMonth()...)

	return columnHeadings
}

// CreateRenderer returns a new WidgetRenderer for this widget.
// This should not be called by regular code, it is used internally to render a widget.
func (c *Calendar) CreateRenderer() fyne.WidgetRenderer {
	c.monthPrevious = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		c.month--
		if c.month < 1 {
			c.month = 12
			c.year--
		}
		c.monthLabel.ParseMarkdown(c.monthYear())

		c.dates.Objects = c.calendarObjects()
	})
	c.monthNext = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		c.month++
		if c.month > 12 {
			c.month = 1
			c.year++
		}
		c.monthLabel.ParseMarkdown(c.monthYear())

		c.dates.Objects = c.calendarObjects()
	})

	c.monthLabel = widget.NewRichTextFromMarkdown(c.monthYear())

	nav := container.New(layout.NewBorderLayout(nil, nil, c.monthPrevious, c.monthNext),
		c.monthPrevious, c.monthNext, container.NewCenter(c.monthLabel))

	c.dates = container.New(newCalendarLayout(c.cellSize, c.padding), c.calendarObjects()...)

	dateContainer := container.NewVBox(nav, c.dates)

	return widget.NewSimpleRenderer(dateContainer)
}

// NewCalendar creates a calendar instance
func NewCalendar(cT time.Time, onSelected func(time.Time), cellSize float32, padding float32) *Calendar {
	c := &Calendar{day: cT.Day(),
		month:        int(cT.Month()),
		year:         cT.Year(),
		calendarTime: cT,
		onSelected:   onSelected,
		cellSize:     cellSize,
		padding:      padding}

	c.ExtendBaseWidget(c)

	return c
}
