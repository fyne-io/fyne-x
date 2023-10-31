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
		xcontainer.Responsive(
			widget.NewButton("One !", func() {}),
			1, .5, xcontainer.OneThird, // 100% for small, 50% for medium and 33% for larger
		),
		xcontainer.Responsive(
			widget.NewButton("Two !", func() {}),
			1, .5, xcontainer.OneThird, // 100% for small, 50% for medium and 33% for larger
		),
		xcontainer.Responsive(
			widget.NewButton("Three !", func() {}),
			1, 1, xcontainer.OneThird, // 100% for small and medium, 33% for larger
		),
		xcontainer.Responsive(formLayout(), 1, .5), // 100% for small, 50% for others
		xcontainer.Responsive(formLayout(), 1, .5), // 100% for small, 50% for others
		button, // 100% by default
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
			canvas := window.Content()
			if canvas == nil {
				continue
			}
			time.Sleep(time.Millisecond * 100)
			width := canvas.Size().Width
			if width <= float32(layout.ExtraSmall) {
				label.SetText(fmt.Sprintf("Extra small devicce %v <= %v", width, layout.ExtraSmall))
			} else if width <= float32(layout.Small) {
				label.SetText(fmt.Sprintf("Small device %v <= %v", width, layout.Small))
			} else if width <= float32(layout.Medium) {
				label.SetText(fmt.Sprintf("Medium device %v <= %v", width, layout.Medium))
			} else if width <= float32(layout.Large) {
				label.SetText(fmt.Sprintf("Large device %v <= %v", width, layout.Large))
			} else {
				label.SetText(fmt.Sprintf("Extra large device %v >= %v", width, layout.ExtraLarge))
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

	entry1 := widget.NewEntry()
	entry2 := widget.NewEntry()
	entry3 := widget.NewEntry()

	label1 := widget.NewLabel("Give your name")
	label2 := widget.NewLabel("Give your age")
	label3 := widget.NewLabel("Give your email")
	label1.Wrapping = fyne.TextWrapWord
	label1.Truncation = fyne.TextTruncateEllipsis
	label2.Wrapping = fyne.TextWrapWord
	label2.Truncation = fyne.TextTruncateEllipsis
	label3.Wrapping = fyne.TextWrapWord
	label3.Truncation = fyne.TextTruncateEllipsis

	// define the sizes for medium and large devices
	mediumLabelSize := float32(.25) // we can use float32
	mediumEntrySize := float32(.75)
	largeLabelSize := xcontainer.OneThird // or helpers
	largeEntrySize := xcontainer.TwoThird
	return xcontainer.NewResponsive(
		title,
		//                         Small,     Medium,     Large and above
		xcontainer.Responsive(label1, 1, mediumLabelSize, largeLabelSize),
		xcontainer.Responsive(entry1, 1, mediumEntrySize, largeEntrySize),
		xcontainer.Responsive(label2, 1, mediumLabelSize, largeLabelSize),
		xcontainer.Responsive(entry2, 1, mediumEntrySize, largeEntrySize),
		xcontainer.Responsive(label3, 1, mediumLabelSize, largeLabelSize),
		xcontainer.Responsive(entry3, 1, mediumEntrySize, largeEntrySize),
	)
}
