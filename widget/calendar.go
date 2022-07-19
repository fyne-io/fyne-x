package widget

import (
	"fmt"
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
var (
	_        fyne.Layout = (*calendarLayout)(nil)
	padding  float32     = 0
	cellSize float64     = 32
)

type calendarLayout struct {
	cols int
}

func newCalendarLayout(s float64) fyne.Layout {
	cellSize = s
	return &calendarLayout{cols: 7}
}

func (g *calendarLayout) countRows(objects []fyne.CanvasObject) int {
	count := 0
	for _, child := range objects {
		if child.Visible() {
			count++
		}
	}

	return int(math.Ceil(float64(count) / float64(g.cols)))
}

// Get the leading (top or left) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getLeading(size float64, offset int) float32 {
	ret := (size + float64(padding)) * float64(offset)

	return float32(math.Round(ret))
}

// Get the trailing (bottom or right) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getTrailing(size float64, offset int) float32 {
	return getLeading(size, offset+1) - padding
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

		x1 := getLeading(cellSize, col)
		y1 := getLeading(cellSize, row)
		x2 := getTrailing(cellSize, col)
		y2 := getTrailing(cellSize, row)

		child.Move(fyne.NewPos(x1, y1))
		child.Resize(fyne.NewSize(x2-x1, y2-y1))

		if (i+1)%g.cols == 0 {
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
	return fyne.NewSize(float32(cellSize+float64(padding))*7, float32(cellSize+float64(padding))*float32(rows))
}

type calendar struct {
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

	callback func(time.Time)
}

func (c *calendar) daysOfMonth() []fyne.CanvasObject {
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
		s := fmt.Sprint(dayNum)
		var b fyne.CanvasObject = widget.NewButton(s, func() {

			selectedDate := c.dateForButton(dayNum)

			c.callback(selectedDate)

			c.hideOverlay()
		})

		buttons = append(buttons, b)
	}

	return buttons
}

func (c *calendar) dateForButton(dayNum int) time.Time {
	oldName, off := c.calendarTime.Zone()
	return time.Date(c.year, time.Month(c.month), dayNum, c.calendarTime.Hour(), c.calendarTime.Minute(), 0, 0, time.FixedZone(oldName, off)).In(c.calendarTime.Location())
}

func (c *calendar) hideOverlay() {
	overlayList := c.canvas.Overlays().List()
	overlayList[0].Hide()
}

func (c *calendar) monthYear() string {
	return time.Month(c.month).String() + " " + strconv.Itoa(c.year)
}

func (c *calendar) calendarObjects() []fyne.CanvasObject {
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
func (c *calendar) CreateRenderer() fyne.WidgetRenderer {
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

	c.dates = container.New(newCalendarLayout(32), c.calendarObjects()...)

	dateContainer := container.NewVBox(nav, c.dates)

	return widget.NewSimpleRenderer(dateContainer)
}

// NewCalendar creates a calendar instance
func NewCalendar(cT time.Time, callback func(time.Time)) *calendar {
	c := &calendar{day: cT.Day(), month: int(cT.Month()), year: cT.Year(), calendarTime: cT, callback: callback}
	c.ExtendBaseWidget(c)

	return c
}
