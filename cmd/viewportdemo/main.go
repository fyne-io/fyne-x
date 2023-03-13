package main

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/x/fyne/widget/diagramwidget/viewport"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/pkg/profile"
)

func main() {

	defer profile.Start(profile.MemProfile).Stop()

	app := app.New()
	w := app.NewWindow("Viewport Demo")

	w.SetMaster()

	stepSize := 1.0

	stepSizeEntry := widget.NewEntry()
	stepSizeEntry.OnChanged = func(text string) {
		var err error
		stepSize, err = strconv.ParseFloat(text, 64)
		if err != nil {
			stepSizeEntry.SetText(fmt.Sprintf("%f", stepSize))
		}
	}
	stepSizeEntry.SetText("1.0")

	vp := viewport.NewViewportWidget(800, 600)

	vp.Objects = append(vp.Objects, &viewport.ViewportLine{
		X1:          20,
		Y1:          20,
		X2:          200,
		Y2:          400,
		StrokeColor: color.RGBA{255, 255, 255, 255},
		StrokeWidth: 1,
	})

	vp.Objects = append(vp.Objects, &viewport.ViewportLine{
		X1:          40,
		Y1:          20,
		X2:          220,
		Y2:          400,
		StrokeColor: color.RGBA{255, 128, 128, 255},
		StrokeWidth: 3,
	})

	vp.Objects = append(vp.Objects, &viewport.ViewportLine{
		X1:          10,
		Y1:          200,
		X2:          300,
		Y2:          10,
		StrokeColor: color.RGBA{64, 255, 64, 255},
		StrokeWidth: 0.5,
	})

	w.SetContent(container.NewHSplit(
		vp,
		container.NewVBox(
			stepSizeEntry,
			widget.NewButton("Pan Left", func() {
				fmt.Printf("vp.XOffset %v", vp.XOffset)
				vp.XOffset += vp.Zoom * stepSize
				vp.Refresh()
				fmt.Printf(" -> %v\n", vp.XOffset)
			}),
			widget.NewButton("Pan Right", func() {
				fmt.Printf("vp.XOffset %v", vp.XOffset)
				vp.XOffset -= vp.Zoom * stepSize
				vp.Refresh()
				fmt.Printf(" -> %v\n", vp.XOffset)
			}),
			widget.NewButton("Pan Up", func() {
				fmt.Printf("vp.YOffset %v", vp.YOffset)
				vp.YOffset += vp.Zoom * stepSize
				vp.Refresh()
				fmt.Printf(" -> %v\n", vp.YOffset)
			}),
			widget.NewButton("Pan Down", func() {
				fmt.Printf("vp.YOffset %v", vp.YOffset)
				vp.YOffset -= vp.Zoom * stepSize
				vp.Refresh()
				fmt.Printf(" -> %v\n", vp.YOffset)
			}),
			widget.NewButton("Zoom In", func() {
				fmt.Printf("vp.Zoom %v", vp.Zoom)
				vp.Zoom *= 1.15
				vp.Refresh()
				fmt.Printf(" -> %v\n", vp.Zoom)
			}),
			widget.NewButton("Zoom Out", func() {
				fmt.Printf("vp.Zoom %v", vp.Zoom)
				vp.Zoom *= 0.85
				vp.Refresh()
				fmt.Printf(" -> %v\n", vp.Zoom)
			}),
		),
	))

	w.ShowAndRun()

}
