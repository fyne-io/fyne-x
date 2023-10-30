package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	xcontainer "fyne.io/x/fyne/container"
	"fyne.io/x/fyne/layout"
)

func main() {
	app := app.New()

	window := app.NewWindow("Responsive")
	window.Resize(fyne.Size{Width: 320, Height: 480})

	// just a button
	button := widget.NewButton("Click me", func() {
		dialog.NewInformation("Hello", "Hello World", window).Show()
	})

	resp := xcontainer.NewResponsive(
		presentation(),       // 100% by default
		winSizeLabel(window), // 100% by default
		layout.Responsive(
			widget.NewButton("One !", func() {}),
			1, .5, layout.OneThird, // 100% for small, 50% for medium and 33% for larger
		),
		layout.Responsive(
			widget.NewButton("Two !", func() {}),
			1, .5, layout.OneThird, // 100% for small, 50% for medium and 33% for larger
		),
		layout.Responsive(
			widget.NewButton("Three !", func() {}),
			1, 1, layout.OneThird, // 100% for small and medium, 33% for larger
		),
		layout.Responsive(formLayout(), 1, .5), // 100% for small, 50% for others
		layout.Responsive(formLayout(), 1, .5), // 100% for small, 50% for others
		button,                                 // 100% by default
	)

	window.SetContent(
		container.NewVScroll(resp),
	)

	window.ShowAndRun()
}

// winSizeLabel returns a label with the current window size inside
func winSizeLabel(window fyne.Window) fyne.CanvasObject {
	label := widget.NewLabel("")
	label.Wrapping = fyne.TextWrapWord
	label.Alignment = fyne.TextAlignCenter
	go func() {
		// Continuously update the label with the window size
		for {
			canvas := window.Canvas()
			if canvas == nil {
				continue
			}
			time.Sleep(time.Millisecond * 100)
			if canvas.Size().Width <= float32(layout.SMALL) {
				label.SetText(fmt.Sprintf("Extra small devicce %v <= %v", canvas.Size().Width, layout.SMALL))
			} else if canvas.Size().Width <= float32(layout.MEDIUM) {
				label.SetText(fmt.Sprintf("Small device %v <= %v", canvas.Size().Width, layout.MEDIUM))
			} else if canvas.Size().Width <= float32(layout.LARGE) {
				label.SetText(fmt.Sprintf("Medium device %v <= %v", canvas.Size().Width, layout.LARGE))
			} else if canvas.Size().Width <= float32(layout.XLARGE) {
				label.SetText(fmt.Sprintf("Large device %v <= %v", canvas.Size().Width, layout.XLARGE))
			} else {
				label.SetText(fmt.Sprintf("Extra large device %v > %v", canvas.Size().Width, layout.LARGE))
			}
		}
	}()

	return label
}

// presentation returns a container with a title text in bold / italic
func presentation() fyne.CanvasObject {
	label := widget.NewLabel("Example of responsive layout")
	label.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	label.Alignment = fyne.TextAlignCenter
	return label
}

// formLayout returns responsive layout where label and entries width ratios are set.
// Each label will:
// - be 100% width for small device
// - be 25% for medium device
// - be 33% for larger device
// And to make entry to be adapted
// - be 100% width for small device
// - be 75% for medium device (100 - 25% from the label)
// - be 67% for larger device (100 - 33% from the label)
func formLayout() fyne.CanvasObject {
	title := widget.NewLabel(
		"This container should be 100% width of small device and 50% for larger.\n" +
			"The labels are sized to 100% width for small devices, 25% for medium and 33% for larger")
	title.Alignment = fyne.TextAlignCenter
	title.Wrapping = fyne.TextWrapWord

	label := widget.NewLabel("Give your name")
	label.Wrapping = fyne.TextWrapWord
	entry := widget.NewEntry()
	label2 := widget.NewLabel("Give your age")
	label2.Wrapping = fyne.TextWrapWord
	entry2 := widget.NewEntry()
	label3 := widget.NewLabel("Give your email")
	label3.Wrapping = fyne.TextWrapWord
	entry3 := widget.NewEntry()

	labelw := float32(.25)
	entryw := float32(.75)
	labelx := layout.OneThird
	entryx := layout.TwoThird
	return xcontainer.NewResponsive(
		title,
		layout.Responsive(label, 1, 1, labelw, labelx),
		layout.Responsive(entry, 1, 1, entryw, entryx),
		layout.Responsive(label2, 1, 1, labelw, labelx),
		layout.Responsive(entry2, 1, 1, entryw, entryx),
		layout.Responsive(label3, 1, 1, labelw, labelx),
		layout.Responsive(entry3, 1, 1, entryw, entryx),
	)
}
