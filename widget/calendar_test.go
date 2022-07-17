package widget

import (
	"strconv"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestNewCalendar(t *testing.T) {
	now := time.Now()
	c := NewCalendar(now, func(time.Time) {})
	assert.Equal(t, now.Day(), c.day)
	assert.Equal(t, int(now.Month()), c.month)
	assert.Equal(t, now.Year(), c.year)

	_ = test.WidgetRenderer(c) // and render
	assert.Equal(t, now.Format("January 2006"), c.monthLabel.Segments[0].(*widget.TextSegment).Text)
}

func TestNewCalendar_ButtonDate(t *testing.T) {
	date := time.Now()
	c := NewCalendar(date, func(time.Time) {})
	_ = test.WidgetRenderer(c) // and render

	endNextMonth := date.AddDate(0, 1, 0).AddDate(0, 0, -(date.Day() - 1))
	last := endNextMonth.AddDate(0, 0, -1)

	firstDate := firstDateButton(c.dates)
	assert.Equal(t, "1", firstDate.Text)
	lastDate := c.dates.Objects[len(c.dates.Objects)-1].(*widget.Button)
	assert.Equal(t, strconv.Itoa(last.Day()), lastDate.Text)
}

func TestNewCalendar_ButtonTap(t *testing.T) {
	tapped := time.Now()
	c := NewCalendar(time.Now(), func(t time.Time) {
		tapped = t
	})
	win := test.NewCanvas()
	pos := fyne.NewPos(10, 10)
	c.ShowAtPos(win, pos)
	_ = test.WidgetRenderer(c) // and render

	firstDate := firstDateButton(c.dates)
	test.Tap(firstDate)
	assert.Equal(t, 1, tapped.Day())
	assert.Equal(t, 0, len(win.Overlays().List()))
}

func TestNewCalendar_Next(t *testing.T) {
	date := time.Now()
	c := NewCalendar(date, func(time.Time) {})
	_ = test.WidgetRenderer(c) // and render

	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Segments[0].(*widget.TextSegment).Text)

	test.Tap(c.monthNext)
	date = date.AddDate(0, 1, 0)
	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Segments[0].(*widget.TextSegment).Text)
}

func TestNewCalendar_Previous(t *testing.T) {
	date := time.Now()
	c := NewCalendar(date, func(time.Time) {})
	_ = test.WidgetRenderer(c) // and render

	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Segments[0].(*widget.TextSegment).Text)

	test.Tap(c.monthPrevious)
	date = date.AddDate(0, -1, 0)
	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Segments[0].(*widget.TextSegment).Text)
}

func TestNewCalendar_ShowHide(t *testing.T) {
	win := test.NewCanvas()
	c := NewCalendar(time.Now(), func(time.Time) {})
	pos := fyne.NewPos(10, 10)

	c.ShowAtPos(win, pos)
	assert.Equal(t, 1, len(win.Overlays().List()))

	c.hideOverlay()
	assert.Equal(t, 0, len(win.Overlays().List()))
}

func firstDateButton(c *fyne.Container) *widget.Button {
	for _, b := range c.Objects {
		if nonBlank, ok := b.(*widget.Button); ok {
			return nonBlank
		}
	}

	return nil
}
