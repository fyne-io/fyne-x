package widget

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestNewCalendar(t *testing.T) {
	now := time.Now()
	c := NewCalendar(now, nil)
	assert.Equal(t, now.Day(), c.currentTime.Day())
	assert.Equal(t, int(now.Month()), int(c.currentTime.Month()))
	assert.Equal(t, now.Year(), c.currentTime.Year())

	_ = test.WidgetRenderer(c) // and render
	assert.Equal(t, now.Format("January 2006"), c.monthLabel.Text)
}

func TestNewCalendar_ButtonDate(t *testing.T) {
	date := time.Now()
	c := NewCalendar(date, nil)
	_ = test.WidgetRenderer(c) // and render

	endNextMonth := date.AddDate(0, 1, 0).AddDate(0, 0, -(date.Day() - 1))
	last := endNextMonth.AddDate(0, 0, -1)

	firstDate := firstDateButton(c.dates)
	assert.Equal(t, "1", firstDate.Text)
	lastDate := c.dates.Objects[len(c.dates.Objects)-1].(*widget.Button)
	assert.Equal(t, strconv.Itoa(last.Day()), lastDate.Text)
}

func TestNewCalendar_Next(t *testing.T) {
	date := time.Now()
	c := NewCalendar(date, nil)
	_ = test.WidgetRenderer(c) // and render

	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Text)

	test.Tap(c.monthNext)
	date = date.AddDate(0, 1, 0)
	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Text)
}

func TestNewCalendar_Previous(t *testing.T) {
	date := time.Now()
	c := NewCalendar(date, nil)
	_ = test.WidgetRenderer(c) // and render

	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Text)

	test.Tap(c.monthPrevious)
	date = date.AddDate(0, -1, 0)
	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Text)
}

func TestNewCalendar_Resize(t *testing.T) {
	date := time.Now()
	c := NewCalendar(date, func(time.Time) {})
	r := test.WidgetRenderer(c) // and render
	layout := c.dates.Layout.(*calendarLayout)

	baseSize := c.MinSize()
	r.Layout(baseSize)
	min := layout.cellSize

	r.Layout(baseSize.AddWidthHeight(100, 0))
	assert.Greater(t, layout.cellSize.Width, min.Width)
	assert.Equal(t, layout.cellSize.Height, min.Height)

	r.Layout(baseSize.AddWidthHeight(0, 100))
	assert.Equal(t, layout.cellSize.Width, min.Width)
	assert.Greater(t, layout.cellSize.Height, min.Height)

	r.Layout(baseSize.AddWidthHeight(100, 100))
	assert.Greater(t, layout.cellSize.Width, min.Width)
	assert.Greater(t, layout.cellSize.Height, min.Height)
}

func TestNewCalendar_Single(t *testing.T) {
	date := time.Date(2023, time.June, 22, 13, 48, 45, 0, time.UTC)
	c := NewCalendar(date, nil)
	_ = test.WidgetRenderer(c) // and render

	btn := getDateButton(c.dates, 14)
	test.Tap(btn)

	assert.Equal(t, widget.HighImportance, btn.Importance)
	if assert.Len(t, c.SelectedDates, 1) {
		tYear, tMonth, tDay := c.SelectedDates[0].Date()
		assert.Equal(t, 2023, tYear)
		assert.Equal(t, time.June, tMonth)
		assert.Equal(t, 14, tDay)
	}

	test.Tap(btn)

	assert.Equal(t, widget.LowImportance, btn.Importance)
	assert.Empty(t, c.SelectedDates)
}

func TestNewCalendar_Multi(t *testing.T) {
	date := time.Date(2023, time.June, 22, 13, 48, 45, 0, time.UTC)
	c := NewCalendar(date, nil)
	c.SelectionMode = CalendarMulti
	_ = test.WidgetRenderer(c) // and render

	test.Tap(getDateButton(c.dates, 4))
	test.Tap(getDateButton(c.dates, 10))
	test.Tap(getDateButton(c.dates, 3))
	test.Tap(getDateButton(c.dates, 9))

	if assert.Len(t, c.SelectedDates, 4) {
		y, m, d := c.SelectedDates[0].Date()
		assert.Equal(t, 2023, y)
		assert.Equal(t, time.June, m)
		assert.Equal(t, 3, d)

		y, m, d = c.SelectedDates[1].Date()
		assert.Equal(t, 2023, y)
		assert.Equal(t, time.June, m)
		assert.Equal(t, 4, d)

		y, m, d = c.SelectedDates[2].Date()
		assert.Equal(t, 2023, y)
		assert.Equal(t, time.June, m)
		assert.Equal(t, 9, d)

		y, m, d = c.SelectedDates[3].Date()
		assert.Equal(t, 2023, y)
		assert.Equal(t, time.June, m)
		assert.Equal(t, 10, d)
	}

	test.Tap(getDateButton(c.dates, 3))

	assert.Len(t, c.SelectedDates, 3)

	assert.Equal(t, widget.LowImportance, getDateButton(c.dates, 3).Importance)
	assert.Equal(t, widget.HighImportance, getDateButton(c.dates, 4).Importance)
	assert.Equal(t, widget.HighImportance, getDateButton(c.dates, 9).Importance)
	assert.Equal(t, widget.HighImportance, getDateButton(c.dates, 10).Importance)
}

func TestNewCalendar_Range(t *testing.T) {
	date := time.Date(2023, time.June, 22, 13, 48, 45, 0, time.UTC)
	c := NewCalendar(date, nil)
	c.SelectionMode = CalendarRange
	_ = test.WidgetRenderer(c) // and render

	test.Tap(getDateButton(c.dates, 5))
	test.Tap(getDateButton(c.dates, 20))

	for i := 5; i <= 20; i++ {
		assert.Equal(t, widget.HighImportance, getDateButton(c.dates, i).Importance)
	}

	if assert.Len(t, c.SelectedDates, 16) {
		for i, s := range c.SelectedDates {
			y, m, d := s.Date()
			assert.Equal(t, 2023, y)
			assert.Equal(t, time.June, m)
			assert.Equal(t, i+5, d)
		}
	}

	test.Tap(getDateButton(c.dates, 7))

	assert.Len(t, c.SelectedDates, 14)
	assert.Equal(t, 7, c.SelectedDates[0].Day())

	test.Tap(getDateButton(c.dates, 19))

	assert.Len(t, c.SelectedDates, 13)
	assert.Equal(t, 19, c.SelectedDates[len(c.SelectedDates)-1].Day())

	test.Tap(getDateButton(c.dates, 7))
	assert.Len(t, c.SelectedDates, 1)
	assert.Equal(t, 19, c.SelectedDates[0].Day())

	test.Tap(getDateButton(c.dates, 7))
	test.Tap(getDateButton(c.dates, 4))

	assert.Len(t, c.SelectedDates, 16)
	assert.Equal(t, 4, c.SelectedDates[0].Day())

	test.Tap(getDateButton(c.dates, 21))
	assert.Len(t, c.SelectedDates, 18)
	assert.Equal(t, 21, c.SelectedDates[len(c.SelectedDates)-1].Day())

	test.Tap(getDateButton(c.dates, 21))
	assert.Len(t, c.SelectedDates, 1)
	assert.Equal(t, 4, c.SelectedDates[0].Day())
}

func getDateButton(c *fyne.Container, day int) *widget.Button {
	var counter int
	for _, b := range c.Objects {
		if nonBlank, ok := b.(*widget.Button); ok {
			if counter+1 == day {
				return nonBlank
			}
			counter++
		}
	}

	return nil
}

func firstDateButton(c *fyne.Container) *widget.Button {
	for _, b := range c.Objects {
		if nonBlank, ok := b.(*widget.Button); ok {
			return nonBlank
		}
	}

	return nil
}
