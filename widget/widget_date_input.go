package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"time"
)

var _ fyne.Focusable = (*JDateEntry)(nil)
var _ fyne.Tappable = (*JDateEntry)(nil)

// Custom Date Entry using RichText
type JDateEntry struct {
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


func (t *JDateEntry) addTime(secName string, v int) {
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

func (t *JDateEntry) updateDisplay() {

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

// implement : focusable interface
func (t *JDateEntry) FocusGained() {
	t.border.FillColor = theme.FocusColor()
	t.border.StrokeColor = theme.PrimaryColor()
	t.updateDisplay()
	t.Refresh()
}

func (t *JDateEntry) FocusLost() {
	t.border.FillColor = theme.BackgroundColor()
	t.border.StrokeColor = theme.InputBorderColor()
	t.daySeg.Style.TextStyle.Bold = false
	t.monthSeg.Style.TextStyle.Bold = false
	t.yearSeg.Style.TextStyle.Bold = false
	t.Refresh()
}

func (t *JDateEntry) TypedKey(e *fyne.KeyEvent) {
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


// implemnet: Tappable
func (t *JDateEntry) Tapped(p *fyne.PointEvent) {
	t.mainWindow.Canvas().Focus(t)
}

// helper functions
func newTextSegment() *widget.TextSegment {
	t := new(widget.TextSegment)
	t.Style.Inline = true
	return t
}

// Set New Date
func (t *JDateEntry) SetDate(newDate time.Time) {
	t.currentDate = newDate
	t.updateDisplay()
}

// Get Inputed Data
func (t *JDateEntry) GetDate() time.Time {
	return t.currentDate
}

// Get Date As String
func (t *JDateEntry) GetDateString() string {
	return t.currentDate.Format("02-Jan-2006")
}

// widget render
func (t *JDateEntry) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	return widget.NewSimpleRenderer(t.gl)
}

// constructor
func NewJDateInputWidget(app fyne.App, window fyne.Window) *JDateEntry {
	t := new(JDateEntry)
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
