package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	xwidget "fyne.io/x/fyne/widget"

	"image/color"
	"strconv"
)

func main() {
	app := app.New()

	h1 := xwidget.NewHexWidget()
	h2 := xwidget.NewHexWidget()
	h3 := xwidget.NewHexWidget()
	h4 := xwidget.NewHexWidget()
	h5 := xwidget.NewHexWidget()
	h6 := xwidget.NewHexWidget()
	h7 := xwidget.NewHexWidget()
	h8 := xwidget.NewHexWidget()

	hexes := []*xwidget.HexWidget{h1, h2, h3, h4, h5, h6, h7, h8}

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

	slantSlider := widget.NewSlider(0, 30)
	slantSlider.SetValue(10)
	slantSlider.OnChanged = func(v float64) {
		for _, w := range hexes {
			w.SetSlant(float32(v))
		}
	}

	w := app.NewWindow("Hex Widget Demo")

	colorOnButton := widget.NewButton("change active color", func() {
		cd := dialog.NewColorPicker(
			"choose a new active color",
			"choose a new active color",
			func(c color.Color) {
				for _, w := range hexes {
					r, g, b, a := c.RGBA()
					w.SetOnColor(color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
				}
			},
			w)

		cd.Advanced = true
		cd.Show()
	})

	colorOffButton := widget.NewButton("change inactive color", func() {
		cd := dialog.NewColorPicker(
			"choose a new inactive color",
			"choose a new inactive color",
			func(c color.Color) {
				for _, w := range hexes {
					r, g, b, a := c.RGBA()
					w.SetOffColor(color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
				}
			},
			w)

		cd.Advanced = true
		cd.Show()
	})

	size := h1.MinSize()

	widthSlider := widget.NewSlider(10, 200)
	widthSlider.SetValue(float64(size.Width))
	widthSlider.OnChanged = func(v float64) {
		size.Width = float32(v)
		for _, w := range hexes {
			w.SetSize(size)
		}
	}

	heightSlider := widget.NewSlider(10, 200)
	heightSlider.SetValue(float64(size.Height))
	heightSlider.OnChanged = func(v float64) {
		size.Height = float32(v)
		for _, w := range hexes {
			w.SetSize(size)
		}
	}

	w.SetContent(
		container.NewVBox(
			container.NewHBox(h8, h7, h6, h5, h4, h3, h2, h1),
			container.NewAdaptiveGrid(2, e, b),
			container.NewAdaptiveGrid(2, widget.NewLabel("Slide to change hex slant:"), slantSlider),
			container.NewAdaptiveGrid(2, colorOnButton, colorOffButton),
			container.NewAdaptiveGrid(2, widget.NewLabel("Hex width"), widthSlider),
			container.NewAdaptiveGrid(2, widget.NewLabel("Hex height"), heightSlider),
		),
	)
	w.ShowAndRun()
}
