package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	xwidget "fyne.io/x/fyne/widget"
)

var (
	counter1 int
	ledbar1  *xwidget.LedBar
	ledbar2  *xwidget.LedBar
)

func run() {
	for {
		counter1 += 1
		ledbar1.Set(counter1 & 0xFF)
		ledbar2.Set(^(counter1 & 0xFF))
		time.Sleep(time.Millisecond * 350)
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Led")

	ledbar1 = xwidget.NewLedBar([]string{"BIT7", "BIT6", "BIT5", "BIT4", "BIT3", "BIT2", "BIT1", "BIT0"})
	ledbar2 = xwidget.NewLedBar([]string{"BIT7", "BIT6", "BIT5", "BIT4", "BIT3", "BIT2", "BIT1", "BIT0"})
	ledbar2.SetOnColor(color.RGBA{0, 255, 0, 255})

	c := container.NewVBox(ledbar1, ledbar2)

	w.SetContent(c)
	go run()
	w.ShowAndRun()
}
