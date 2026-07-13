package main

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	app := app.NewWithID("gaugetest")
	w := app.NewWindow("Gauge Test")

	gauge := xwidget.NewGauge()
	gauge.Min = 0
	gauge.Max = 280
	gauge.Steps = 28
	gauge.Title = "km/h"
	gauge.SetMinSize(fyne.NewSquareSize(175))

	gauge.TextFormatter = func() string {
		return strconv.FormatFloat(gauge.Value, 'f', 0, 64)
	}

	content := container.NewBorder(
		nil,
		widget.NewButtonWithIcon("Animate", theme.MediaPlayIcon(), func() {
			fyne.NewAnimation(5*time.Second, func(f float32) {
				gauge.SetValue(float64(f * 280))
			}).Start()
		}),
		nil,
		nil,
		gauge,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
