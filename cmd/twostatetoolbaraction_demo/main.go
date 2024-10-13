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
		nil, func(state xwidget.TwoStateState) {
			fmt.Println(state)
		})
	sep := widget.NewToolbarSeparator()
	tb := widget.NewToolbar(twoState0, sep)

	toggleButton := widget.NewButton("Toggle State", func() {
		state := twoState0.GetState()
		twoState0.SetState(!state)
	})
	offIconButton := widget.NewButton("Set OffIcon", func() {
		twoState0.SetOffStateIcon(theme.MediaPlayIcon())
	})
	onIconButton := widget.NewButton("Set OnIcon", func() {
		twoState0.SetState1Icon(theme.MediaPauseIcon())
	})
	vc := container.NewVBox(toggleButton, offIconButton, onIconButton)
	c := container.NewBorder(tb, vc, nil, nil)
	w.SetContent(c)
	w.ShowAndRun()
}
