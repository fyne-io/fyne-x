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
		theme.MediaPauseIcon(), func(state xwidget.TwoStateState) {
			fmt.Println(state)
		})
	sep := widget.NewToolbarSeparator()
	tb := widget.NewToolbar(twoState0, sep)

	toggleButton := widget.NewButton("Toggle State", func() {
		state := twoState0.GetState()
		twoState0.SetState(!state)
	})
	icon0Button := widget.NewButton("Set Icon0", func() {
		twoState0.SetState0Icon(theme.MediaPlayIcon())
	})
	vc := container.NewVBox(toggleButton, icon0Button)
	c := container.NewBorder(tb, vc, nil, nil)
	w.SetContent(c)
	w.ShowAndRun()
}
