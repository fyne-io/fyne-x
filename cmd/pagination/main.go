package main

import (
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	widget2 "fyne.io/x/fyne/widget"
)

var (
	data            [][]string
	table           *widget.Table
	defaultPageSize = 10
	defaultPage     = 1
)

func main() {
	a := app.New()
	w := a.NewWindow("Pagination")
	w.Resize(fyne.NewSize(500, 400))

	pager := widget2.NewPagination(defaultPageSize)
	pager.SetTotalRows(getTotalRows())
	pager.OnChange = func(p, s int) {
		fmt.Printf("page: %d, pageSize: %d\n", p, s)
		updateData(p, s)
	}
	btn := widget.NewButton("reset total", func() {
		pager.SetTotalRows(getTotalRows())
	})
	updateData(defaultPage, defaultPageSize)
	createTable()
	table.Resize(fyne.NewSize(500, 300))

	c := container.NewVBox(btn, table, pager)
	w.SetContent(c)
	w.ShowAndRun()
}

func getTotalRows() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int(100 * r.Float64())
}

func createTable() {
	tb := widget.NewTable(
		func() (int, int) {
			return len(data), 1
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*widget.Label)
			l.SetText(data[i.Row][i.Col])
		})
	table = tb
	table.Refresh()
}

// updateData will do db query or http request to get a range of data
func updateData(page, pageSize int) {
	var result [][]string
	for i := 0; i < pageSize; i++ {
		result = append(result, []string{
			fmt.Sprintf("page%d", page),
		})
	}
	data = result
}
