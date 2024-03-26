package main

import (
	"net/url"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/dialog"
)

func main() {
	a := app.New()
	w := a.NewWindow("Fyne-x demo")

	docURL, _ := url.Parse("https://docs.fyne.io")
	links := []*widget.Hyperlink{
		widget.NewHyperlink("Docs", docURL),
	}
	w.SetContent(container.NewGridWithColumns(1,
		widget.NewButton("About", func() {
			dialog.ShowAbout("Some **cool** stuff", links, a, w)
		}),
		widget.NewButton("About window", func() {
			dialog.ShowAboutWindow("Some **cool** stuff", links, a)
		}),
	))

	w.ShowAndRun()
}
