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

func main() {
	a := app.New()
	w := a.NewWindow("Pagination")
	w.Resize(fyne.NewSize(500, 330))

	pager := widget2.NewPagination(10)
	pager.SetTotalRows(50)
	pager.Resize(fyne.NewSize(500, 30))
	pager.OnChange = func(p, s int) {
		fmt.Printf("page: %d, pageSize: %d\n", p, s)
	}
	btn := widget.NewButton("reset total", func() {
		pager.SetTotalRows(getTotalRows())
	})
	c := container.NewVBox(pager, btn)
	w.SetContent(c)
	w.ShowAndRun()
}

func getTotalRows() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int(100 * r.Float64())
}
