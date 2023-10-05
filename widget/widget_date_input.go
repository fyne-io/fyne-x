package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"time"
)

var _ fyne.Focusable = (*JDateInputWidget)(nil)
var _ fyne.Tappable = (*JDateInputWidget)(nil)

// Custom Date Input Widget using RichText
type JDateInputWidget struct {
	widget.BaseWidget

	// currently selected item data
	currentDate time.Time

	// current selected section
	curSection string

	// select section yes/no
	selectSection int

	// rich-text widget
	richText *widget.RichText

	// textSegments
	daySeg   *widget.TextSegment
	monthSeg *widget.TextSegment
	yearSeg  *widget.TextSegment

	// border
	border *canvas.Rectangle

	// layout
	gl *fyne.Container

	// app
	fyneApp fyne.App

	// main window
	mainWindow fyne.Window
}

func (t *JDateInputWidget) addTime(secName string, v int) {
	if secName == "d" {
		t.currentDate = t.currentDate.AddDate(0, 0, v)
	}

	if secName == "m" {
		t.currentDate = t.currentDate.AddDate(0, v, 0)
	}

	if secName == "y" {
		t.currentDate = t.currentDate.AddDate(v, 0, 0)
	}
	t.updateDisplay()
}

func (t *JDateInputWidget) updateDisplay() {

	t.daySeg.Text = t.currentDate.Format("02-")
	t.monthSeg.Text = t.currentDate.Format("Jan-")
	t.yearSeg.Text = t.currentDate.Format("2006")

	t.daySeg.Style.TextStyle.Bold = false
	t.monthSeg.Style.TextStyle.Bold = false
	t.yearSeg.Style.TextStyle.Bold = false

	if t.selectSection == 1 {
		if t.curSection == "d" {
			t.daySeg.Style.TextStyle.Bold = true
		} else if t.curSection == "m" {
			t.monthSeg.Style.TextStyle.Bold = true
		} else if t.curSection == "y" {
			t.yearSeg.Style.TextStyle.Bold = true
		}
	}
	t.richText.Refresh()
}

// Start: implement focusable interface
func (t *JDateInputWidget) FocusGained() {
	t.border.FillColor = theme.FocusColor()
	t.border.StrokeColor = theme.PrimaryColor()
	t.updateDisplay()
	t.Refresh()
}

func (t *JDateInputWidget) FocusLost() {
	t.border.FillColor = theme.BackgroundColor()
	t.border.StrokeColor = theme.InputBorderColor()
	t.daySeg.Style.TextStyle.Bold = false
	t.monthSeg.Style.TextStyle.Bold = false
	t.yearSeg.Style.TextStyle.Bold = false
	t.Refresh()
}

func (t *JDateInputWidget) TypedRune(k rune) {
	// needed
}

func (t *JDateInputWidget) TypedKey(e *fyne.KeyEvent) {
	if e.Name == fyne.KeyUp || e.Name == fyne.KeyU {
		t.addTime(t.curSection, 1)
		t.updateDisplay()
	} else if e.Name == fyne.KeyDown || e.Name == fyne.KeyD {
		t.addTime(t.curSection, -1)
		t.updateDisplay()
	} else if e.Name == fyne.KeyRight || e.Name == fyne.KeyR {
		if t.curSection == "d" {
			t.curSection = "m"
		} else if t.curSection == "m" {
			t.curSection = "y"
		} else if t.curSection == "y" {
			t.curSection = "d"
		}
		t.updateDisplay()
	} else if e.Name == fyne.KeyLeft || e.Name == fyne.KeyL {
		if t.curSection == "d" {
			t.curSection = "y"
		} else if t.curSection == "m" {
			t.curSection = "d"
		} else if t.curSection == "y" {
			t.curSection = "m"
		}
		t.updateDisplay()
	} else if e.Name == fyne.KeySpace {
		t.currentDate = time.Now()
		t.updateDisplay()
	} else if e.Name == fyne.KeyPageDown {
		t.addTime("d", -10)
		t.updateDisplay()
	} else if e.Name == fyne.KeyPageUp {
		t.addTime("d", 10)
		t.updateDisplay()
	}
}

// End: implement focusable interface

// Start: implemnet Tappable
func (t *JDateInputWidget) Tapped(p *fyne.PointEvent) {
	t.mainWindow.Canvas().Focus(t)
}

// End: implemnet Tappable

// helper functions
func newTextSegment() *widget.TextSegment {
	t := new(widget.TextSegment)
	t.Style.Inline = true
	return t
}

// public interface
func (t *JDateInputWidget) SetDate(newDate time.Time) {
	t.currentDate = newDate
	t.updateDisplay()
}

func (t *JDateInputWidget) GetDate() time.Time {
	return t.currentDate
}

func (t *JDateInputWidget) GetDateString() string {
	return t.currentDate.Format("02-Jan-2006")
}

// widget render
func (t *JDateInputWidget) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	return widget.NewSimpleRenderer(t.gl)
}

// constructor
func NewJDateInputWidget(app fyne.App, window fyne.Window) *JDateInputWidget {
	t := new(JDateInputWidget)
	t.ExtendBaseWidget(t)

	t.fyneApp = app
	t.mainWindow = window

	t.curSection = "d"
	t.currentDate = time.Now()

	t.daySeg = newTextSegment()
	t.monthSeg = newTextSegment()
	t.yearSeg = newTextSegment()

	t.richText = widget.NewRichText(t.daySeg, t.monthSeg, t.yearSeg)

	t.border = canvas.NewRectangle(theme.BackgroundColor())
	t.border.StrokeWidth = float32(0.4)
	t.border.StrokeColor = theme.InputBorderColor()

	t.gl = container.NewHBox(container.NewMax(t.border, t.richText))

	t.selectSection = 0
	t.updateDisplay()
	t.selectSection = 1

	return t
}
