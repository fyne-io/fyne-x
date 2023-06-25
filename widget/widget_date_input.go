package widget

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
	"time"
)

const ConstDefaultDateFormat = "02-Jan-2006"

var DayPos = []int{0, 2}
var MonthPos = []int{3, 6}
var YearPos = []int{7, 11}

type MyDateEntry struct {
	widget.Entry
	currentDateValue time.Time
}

func NewMyDateEntry() *MyDateEntry {
	entry := new(MyDateEntry)
	entry.ExtendBaseWidget(entry)
	entry.SetPlaceHolder("DD-MMM-YYYY")
	return entry
}

// Get current cursor position to Section (Day Month year)
func (e *MyDateEntry) cursorPosToSection() string {
	if e.CursorColumn >= DayPos[0] && e.CursorColumn <= DayPos[1] {
		return "d"
	}
	if e.CursorColumn >= MonthPos[0] && e.CursorColumn <= MonthPos[1] {
		return "m"
	}
	if e.CursorColumn >= YearPos[0] && e.CursorColumn <= YearPos[1] {
		return "y"
	}
	return ""
}

// Get Day Month year to cursor postion
func (e *MyDateEntry) sectionToCursorPos(lsec string) int {
	if lsec == "d" {
		return DayPos[0]
	}

	if lsec == "m" {
		return MonthPos[0]
	}

	if lsec == "y" {
		return YearPos[0]
	}
	return -1
}

// add year, month, day to current date
func (e *MyDateEntry) addTime(v int, cur_section string) {
	if e.currentDateValue.IsZero() == true {
		return
	}

	if cur_section == "d" {
		e.currentDateValue = e.currentDateValue.AddDate(0, 0, v)
	}

	if cur_section == "m" {
		e.currentDateValue = e.currentDateValue.AddDate(0, v, 0)
	}

	if cur_section == "y" {
		e.currentDateValue = e.currentDateValue.AddDate(v, 0, 0)
	}
	e.updateDisplay()
}

// set current date on space key
func (e *MyDateEntry) setCurrentDate() {
	e.currentDateValue = time.Now()
	e.updateDisplay()
}

// update current display
func (e *MyDateEntry) updateDisplay() {
	e.SetText(e.currentDateValue.Format(ConstDefaultDateFormat))
}

// handle key events
func (e *MyDateEntry) TypedKey(key *fyne.KeyEvent) {

	if key.Name == fyne.KeyDelete {
		e.SetText("")
		return
	}

	if key.Name == fyne.KeyUp {
		e.addTime(1, e.cursorPosToSection())
		return
	}

	if key.Name == fyne.KeyDown {
		e.addTime(-1, e.cursorPosToSection())
		return
	}

	if key.Name == fyne.KeySpace {
		e.setCurrentDate()
		return
	}

	if key.Name == fyne.KeyEnter {
		e.parse_and_update_date()
		e.addTime(0, e.cursorPosToSection())
		return
	}

	if key.Name == fyne.KeyReturn {
		e.parse_and_update_date()
		e.addTime(0, e.cursorPosToSection())
		return
	}

	e.Entry.TypedKey(key)
}

// this where we are converting current text to date
func (e *MyDateEntry) FocusLost() {
	e.parse_and_update_date()
	e.Entry.FocusLost()
}

// Date string to time.Time conversion
// we assume 1st part is always Day
// input = 1 -> 1-CurMonth-CurYear
// input = 1.5, 1/5, 1-5 -> 1-5-CurYear
func (e *MyDateEntry) parse_and_update_date() {
	var date_str = e.Text

	e.TextStyle.Bold = false

	if len(date_str) == 0 {
		e.SetText("")
		return
	}

	var y, m int

	y = time.Now().Year()
	m = int(time.Now().Month())

	date_str = strings.Replace(date_str, ".", "-", -1)
	date_str = strings.Replace(date_str, "/", "-", -1)
	dt := strings.Split(date_str, "-")

	if len(dt) == 1 {
		dt[0] = strings.TrimSpace(dt[0])
		date_str = fmt.Sprintf("%s-%d-%d", dt[0], m, y)
	}

	if len(dt) == 2 {
		dt[0] = strings.TrimSpace(dt[0])
		dt[1] = strings.TrimSpace(dt[1])

		if len(dt[1]) == 0 {
			dt[1] = strconv.Itoa(m)
		}

		date_str = fmt.Sprintf("%s-%s-%d", dt[0], dt[1], y)
	}

	date_str = strings.TrimSpace(date_str)

	var allowed_formats = []string{"02-01-2006", "2-1-2006", "2006-01-02", "2006-1-2", "2-Jan-2006"}
	for _, v := range allowed_formats {
		e.currentDateValue, _ = time.Parse(v, date_str)
		if e.currentDateValue.IsZero() == false {
			break
		}
	}

	if e.currentDateValue.IsZero() == true {
		e.SetText("")
	} else {
		e.SetText(e.currentDateValue.Format(ConstDefaultDateFormat))
		e.TextStyle.Bold = true
	}
}

// return current date
func (e *MyDateEntry) ToDate() time.Time {
	return e.currentDateValue
}
