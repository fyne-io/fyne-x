package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

var spinnerDisabled bool
var data binding.Int = binding.NewInt()
var s1 *xwidget.Spinner

func main() {
	a := app.New()

	ls1 := widget.NewLabel("Value set in Spinner 1:")
	s1ValueLabel := widget.NewLabel("")
	ls2 := widget.NewLabel("Data value bound to Spinner 2:")
	dataValueLabel := widget.NewLabel("")
	data.AddListener(binding.NewDataListener(func() {
		val, err := data.Get()
		if err != nil {
			return
		}
		dataValueLabel.Text = strconv.Itoa(val)
		dataValueLabel.Refresh()
	}))

	c2 := container.NewGridWithColumns(2, ls1, s1ValueLabel, ls2, dataValueLabel)
	l1 := widget.NewLabel("Spinner 1:")
	s1 = xwidget.NewSpinner(1, 12, 3, nil)
	s1.OnChanged = func(val int) {
		s1ValueLabel.Text = strconv.Itoa(s1.GetValue())
		s1ValueLabel.Refresh()
	}
	// OnChanged has to be called here to display initial value in s1ValueLabel.
	s1.OnChanged(s1.GetValue())
	l2 := widget.NewLabel("Spinner 2:")
	s2 := xwidget.NewSpinnerWithData(-2, 16, 1, data)
	c := container.NewGridWithColumns(2, l1, s1, l2, s2)
	b := widget.NewButton("Disable Spinner 1", func() {})
	b.OnTapped = func() {
		spinnerDisabled = !spinnerDisabled
		if spinnerDisabled {
			s1.Disable()
			b.SetText("Enable Spinner 1")
		} else {
			s1.Enable()
			b.SetText("Disable Spinner 1")
		}
	}
	bs1 := widget.NewButton("Set Spinner 1 to 5", func() { s1.SetValue(5) })
	bs2 := widget.NewButton("Set bound value to 12", func() { data.Set(12) })

	v := container.NewVBox(c, b, c2, bs1, bs2)
	w := a.NewWindow("SpinnerDemo")
	w.Resize(fyne.NewSize(200, 200))
	w.SetContent(v)
	w.ShowAndRun()
}
