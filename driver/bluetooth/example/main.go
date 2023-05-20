package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")
	buildMainMenu(w)
	w.ShowAndRun()
}

func buildMainMenu(w fyne.Window) {
	status := widget.NewLabel("status!")
	b1 := widget.NewButton("server!", nil)
	b2 := widget.NewButton("client!", nil)
	con := container.NewVBox(status, b1, b2)
	b1.OnTapped = func() {
		status.SetText("start server")
		status.Refresh()
		time.Sleep(time.Second * 5)
		con.Hide()
		runServerReturnError(w)
	}
	b2.OnTapped = func() {
		status.SetText("start client")
		status.Refresh()
		time.Sleep(time.Second * 5)
		con.Hide()
		runClientReturnError(w)
	}
	w.SetContent(con)
}
