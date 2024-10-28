package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Two State Demo")

	twoState0 := xwidget.NewTwoStateToolbarAction(nil,
		nil, func(on bool) {
			fmt.Println(on)
		})
	sep := widget.NewToolbarSeparator()
	tb := widget.NewToolbar(twoState0, sep)

	toggleButton := widget.NewButton("Toggle State", func() {
		on := twoState0.GetOn()
		twoState0.SetOn(!on)
	})
	offIconButton := widget.NewButton("Set OffIcon", func() {
		twoState0.SetOffIcon(theme.MediaPlayIcon())
	})
	onIconButton := widget.NewButton("Set OnIcon", func() {
		twoState0.SetOnIcon(theme.MediaStopIcon())
	})
	vc := container.NewVBox(toggleButton, offIconButton, onIconButton)
	c := container.NewBorder(tb, vc, nil, nil)
	w.SetContent(c)
	w.ShowAndRun()
}
