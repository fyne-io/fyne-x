package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/dialog"
)

func main() {
	app := app.New()
	w := app.NewWindow("Hello")

	button := widget.NewButton("Click Me", func() {
		dialog.NewFileSelector(func(file fyne.URIReadCloser, err error) {
			log.Println("File selected", file.URI(), err)
		}, w).Open()

	})

	w.SetContent(button)

	w.ShowAndRun()

}
