package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"fyne.io/x/fyne/widget/hexwidget"

	"strconv"
)

func main() {
	app := app.New()

	h1 := hexwidget.NewHexWidget()
	h2 := hexwidget.NewHexWidget()
	h3 := hexwidget.NewHexWidget()
	h4 := hexwidget.NewHexWidget()
	h5 := hexwidget.NewHexWidget()
	h6 := hexwidget.NewHexWidget()
	h7 := hexwidget.NewHexWidget()
	h8 := hexwidget.NewHexWidget()

	e := widget.NewEntry()
	e.PlaceHolder = "enter a 32-bit number"
	e.Validator = func(s string) error {
		_, err := strconv.Atoi(s)
		return err
	}

	b := widget.NewButton("update", func() {
		i, _ := strconv.Atoi(e.Text)
		u := uint(i)
		h1.Set((u & 0x0000000f) >> 0)
		h2.Set((u & 0x000000f0) >> 4)
		h3.Set((u & 0x00000f00) >> 8)
		h4.Set((u & 0x0000f000) >> 12)
		h5.Set((u & 0x000f0000) >> 16)
		h6.Set((u & 0x00f00000) >> 20)
		h7.Set((u & 0x0f000000) >> 24)
		h8.Set((u & 0xf0000000) >> 28)
	},
	)

	w := app.NewWindow("Hello")
	w.SetContent(
		container.NewVBox(
			container.NewHBox(h8, h7, h6, h5, h4, h3, h2, h1),
			container.NewAdaptiveGrid(2, e, b),
		),
	)
	w.ShowAndRun()
}
