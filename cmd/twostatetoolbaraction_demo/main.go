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

	twoState0 := xwidget.NewTwoStateToolbarAction(theme.MediaPlayIcon(),
		theme.MediaPauseIcon(), func(state xwidget.TwoStateState) {
			fmt.Println(state)
		})
	sep := widget.NewToolbarSeparator()
	ta := widget.NewToolbarAction(theme.MediaPhotoIcon(), nil)
	tb := widget.NewToolbar(twoState0, sep, ta)

	toggleButton := widget.NewButton("Toggle State", func() {
		state := twoState0.GetState()
		twoState0.SetState(!state)
	})
	c := container.NewBorder(tb, toggleButton, nil, nil)
	w.SetContent(c)
	w.ShowAndRun()
}
