package wrapper

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestTappable(t *testing.T) {
	app := test.NewApp()
	win := app.NewWindow("Test")
	win.Resize(fyne.NewSize(400, 400))
	defer win.Close()

	label1 := widget.NewLabel("Label 1, not wrapped")
	label2 := widget.NewLabel("Label 2, click me")

	tapped := false
	wrapped := MakeTappable(label2, func(e *fyne.PointEvent) {
		tapped = true
	})

	mainContainer := container.NewGridWithColumns(2, label1, wrapped)
	win.SetContent(mainContainer)

	// Click on wrapped label
	test.Tap(wrapped.(fyne.Tappable))
	if !tapped {
		t.Error("Tapped event not fired")
	}
}
