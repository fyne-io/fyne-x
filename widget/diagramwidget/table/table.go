// package table implements a simple table widget.
//
// At present Fyne does not have it's own native table widget. This one intends
// to bridge the gap until a better one is available.
//
// Current limitations:
//   - Not very pretty
//   - No re-sizing
//   - No editing cell contents
//   - Very inefficient (re-generates all table cells each time the widget is
//     refreshed)
//   - No sorting
package table

import (
	"fmt"
	"image/color"

	"github.com/rocketlaunchr/dataframe-go"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// TODO: should calculate row height based on font size.
const textHeight int = 20

const maxWidth int = 300

type tableRenderer struct {
	table   *TableWidget
	objects []fyne.CanvasObject
}

func (t *tableRenderer) MinSize() fyne.Size {
	width := theme.Padding() * float32(2+len(t.table.df.Series))

	t.table.updateColumnWidthsIfNeeded()

	for _, v := range t.table.columnWidths {
		width += float32(v)
	}

	return fyne.Size{Width: width, Height: float32(1+t.table.df.NRows()) * (float32(textHeight) + theme.Padding())}

}

func (t *tableRenderer) Layout(size fyne.Size) {
	// TODO: in the future, it would be better for the table to report a
	// minimum size smaller than the column widths would suggest, then
	// shrink or grow the columns automatically when Layout is called.
}

func (t *tableRenderer) ApplyTheme() {
}

func (t *tableRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (t *tableRenderer) Refresh() {
	// TODO: this is not efficient since it destroys and re-creates all the
	// table cells.

	t.table.updateColumnWidthsIfNeeded()

	t.objects = make([]fyne.CanvasObject, 0)

	xpos := theme.Padding()
	for col, v := range t.table.df.Names() {
		cell := canvas.NewText(fmt.Sprintf("%v", v), theme.ForegroundColor())
		cell.Move(fyne.Position{X: xpos, Y: theme.Padding()})
		// TODO confirm the default values for the last two parameters (new)
		cell.TextStyle = fyne.TextStyle{Bold: true, Italic: false, Monospace: false, Symbol: false, TabWidth: 10}
		t.objects = append(t.objects, cell)
		xpos += theme.Padding() + float32(t.table.columnWidths[col])
	}

	iterator := t.table.df.ValuesIterator(dataframe.ValuesOptions{InitialRow: 0, Step: 1, DontReadLock: true}) // Don't apply read lock because we are write locking from outside.
	t.table.df.Lock()
	for {
		row, vals, _ := iterator()
		if row == nil {
			break
		}
		fmt.Println(*row, vals)

		xpos := theme.Padding()
		ypos := theme.Padding() + (theme.Padding()+float32(textHeight))*float32(*row+1)

		for col := 0; col < len(t.table.df.Series); col++ {
			v := vals[col]

			cell := canvas.NewText(fmt.Sprintf("%v", v), theme.ForegroundColor())
			cell.Move(fyne.Position{X: xpos, Y: ypos})

			t.objects = append(t.objects, cell)

			xpos += theme.Padding() + float32(t.table.columnWidths[col])
		}
	}
	t.table.df.Unlock()

}

func (t *tableRenderer) Destroy() {
}

func (t *tableRenderer) Objects() []fyne.CanvasObject {
	return t.objects
}

type TableWidget struct {
	widget.BaseWidget
	df           *dataframe.DataFrame
	columnWidths []int
}

// updateColumnWidthsIfNeeded will guarantee that t.columnWidths is non-nil,
// and contains the same number of column widths as there are columns in t.df.
func (table *TableWidget) updateColumnWidthsIfNeeded() {
	if (table.columnWidths == nil) || (len(table.columnWidths) != len(table.df.Series)) {
		table.CalculateColumnWidths(maxWidth)
	}
}

// CalculateColumnWidths will replace t.columnWidths with appropriate widths
// that accommodate the full string ivied content of the largest element
// in a column. maxWidth is the widest that any single column can be. Specify
// a maxWidth of 0 for an unlimited maximum width.
func (table *TableWidget) CalculateColumnWidths(maxWidth int) {
	table.columnWidths = make([]int, len(table.df.Series))

	for col, v := range table.df.Names() {
		strwidth := fyne.MeasureText(v, theme.TextSize(), fyne.TextStyle{Bold: true, Italic: false, Monospace: false, Symbol: false, TabWidth: 10}).Width
		if strwidth > float32(maxWidth) {
			strwidth = float32(maxWidth)
		}
		if table.columnWidths[col] < int(strwidth) {
			table.columnWidths[col] = int(strwidth)
		}
	}

	iterator := table.df.ValuesIterator(dataframe.ValuesOptions{InitialRow: 0, Step: 1, DontReadLock: true}) // Don't apply read lock because we are write locking from outside.
	table.df.Lock()
	for {
		row, vals, _ := iterator()
		if row == nil {
			break
		}

		for col := 0; col < len(table.df.Series); col++ {
			s := fmt.Sprintf("%v", vals[col])
			strwidth := fyne.MeasureText(s, theme.TextSize(), fyne.TextStyle{Bold: false, Italic: false, Monospace: false, Symbol: false, TabWidth: 10}).Width
			if strwidth > float32(maxWidth) {
				strwidth = float32(maxWidth)
			}
			if table.columnWidths[col] < int(strwidth) {
				table.columnWidths[col] = int(strwidth)
			}
		}
	}
	table.df.Unlock()

}

func (table *TableWidget) Tapped(ev *fyne.PointEvent) {
}

func (table *TableWidget) TappedSecondary(ev *fyne.PointEvent) {
}

func (table *TableWidget) CreateRenderer() fyne.WidgetRenderer {
	r := tableRenderer{
		table: table,
	}

	r.Refresh()

	return &r
}

func NewTableWidget(df *dataframe.DataFrame) *TableWidget {
	table := &TableWidget{df: df}
	table.CalculateColumnWidths(maxWidth)
	table.ExtendBaseWidget(table)
	return table
}

func (table *TableWidget) ReplaceDataFrame(newdf *dataframe.DataFrame) {
	table.df = newdf
	table.CalculateColumnWidths(maxWidth)
	table.Refresh()
}
