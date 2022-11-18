package wrapper

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestMouseable(t *testing.T) {
	app := test.NewApp()
	win := app.NewWindow("Test")

	label := widget.NewLabel("Label 1, not wrapped")
	label2 := widget.NewLabel("Label, mouseable")

	pos := fyne.NewPos(0, 0)
	in := false
	out := false
	mouseable := SetHoverable(label, func(e *desktop.MouseEvent) {
		in = true
	}, func(e *desktop.MouseEvent) {
		pos = e.Position
	}, func() {
		out = true
	})

	win.SetContent(container.NewHBox(label2, mouseable))

	moveTo := fyne.NewPos(15, 15)
	expectedPos := mouseable.Position().Add(moveTo)
	test.MoveMouse(win.Canvas(), mouseable.Position().Add(fyne.NewPos(5, 5))) // to place the mouse
	if !in {
		t.Error("MouseIn was not called")
	}
	test.MoveMouse(win.Canvas(), expectedPos) // to move the mouse
	if pos != moveTo.Subtract(fyne.NewPos(theme.Padding(), theme.Padding())) {
		t.Error("MouseMoved was not called", pos, moveTo)
	}
	if out {
		t.Error("MouseOut was called")
	}

	test.MoveMouse(win.Canvas(), fyne.NewPos(0, 0)) // go out now
	if !out {
		t.Error("MouseOut was not called")
	}

}
