package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

var spinnerDisabled bool
var data binding.Float = binding.NewFloat()
var s1 *xwidget.Spinner
var s5 *xwidget.Spinner
var bs *widget.Button

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
		dataValueLabel.Text = fmt.Sprintf("%d", int(val))
		dataValueLabel.Refresh()
	}))

	ls5 := widget.NewLabel("Value set in Spinner 5:")
	s5ValueLabel := widget.NewLabel("")
	floatData := binding.NewFloat()
	floatData.AddListener(binding.NewDataListener(func() {
		val, err := floatData.Get()
		if err != nil {
			return
		}
		s5ValueLabel.Text = strconv.FormatFloat(val, 'f', 3, 64)
		s5ValueLabel.Refresh()
	}))
	c6 := container.NewHBox(ls5, s5ValueLabel)

	c2 := container.NewGridWithColumns(2, ls1, s1ValueLabel, ls2, dataValueLabel)
	l1 := widget.NewLabel("Spinner 1 (0, 100, 1, \"%d %%\"):")
	s1 = xwidget.NewSpinner(0, 100, 1, "%d %%", nil)
	s1.OnChanged = func(val float64) {
		s1ValueLabel.Text = fmt.Sprintf("%d %%", int(s1.GetValue()))
		s1ValueLabel.Refresh()
	}
	// OnChanged has to be called here to display initial value in s1ValueLabel.
	s1.OnChanged(s1.GetValue())
	l2 := widget.NewLabel("Spinner 2 With Data (-2, 16, 3, \"%+d\"):")
	s2 := xwidget.NewSpinnerWithData(-2, 16, 3, "%+d", data)
	c := container.NewGridWithColumns(2, l1, s1)
	c1 := container.NewHBox(l2, s2)
	l3 := widget.NewLabel("Uninitialized Spinner 3:")
	s3 := xwidget.NewSpinnerUninitialized("%d")
	c3 := container.NewHBox(l3, s3)
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
	bs = widget.NewButton("Initialize Spinner 3", func() {
		s3.SetMinMaxStep(1, 10, 1)
		l3.Text = "Initialized Spinner 3 (1, 10, 1, \"%d\"):"
		l3.Refresh()
		s3.Enable()
		bs.Disable()
	})
	bs1 := widget.NewButton("Set Spinner 1 to 5", func() { s1.SetValue(5) })
	bs2 := widget.NewButton("Set Spinner 3 bound value to 12", func() { data.Set(12) })

	l4 := widget.NewLabel("Spinner 4 (-1., 400., 10.3, \"%+.1f\"):")
	s4 := xwidget.NewSpinner(-1., 400., 10.3, "%.1f", nil)
	c4 := container.NewHBox(l4, s4)

	l5 := widget.NewLabel("Spinner 5 (0., 16., 3.215, \"%.2f\")")
	s5 = xwidget.NewSpinnerWithData(0., 16., 3.215, "%.2f", floatData)
	c5 := container.NewHBox(l5, s5)

	v := container.NewVBox(c, c1, c3, b, bs, c2, bs1, bs2, c4, c5, c6)
	w := a.NewWindow("SpinnerDemo")
	w.SetContent(v)
	w.ShowAndRun()
}
