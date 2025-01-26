package main

import (
	"strconv"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

var spinnerDisabled bool
var data binding.Int = binding.NewInt()
var s1 *xwidget.IntSpinner
var s5 *xwidget.Float64Spinner
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
		dataValueLabel.Text = strconv.Itoa(val)
		dataValueLabel.Refresh()
	}))

	ls5 := widget.NewLabel("Value set in Float64Spinner 2:")
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
	l1 := widget.NewLabel("IntSpinner 1 (1, 12, 3):")
	s1 = xwidget.NewIntSpinner(1, 12, 3, nil)
	s1.OnChanged = func(val int) {
		s1ValueLabel.Text = strconv.Itoa(s1.GetValue())
		s1ValueLabel.Refresh()
	}
	// OnChanged has to be called here to display initial value in s1ValueLabel.
	s1.OnChanged(s1.GetValue())
	l2 := widget.NewLabel("IntSpinner 2 (-2, 16, 1):")
	s2 := xwidget.NewIntSpinnerWithData(-2, 16, 1, data)
	c := container.NewGridWithColumns(2, l1, s1)
	c1 := container.NewHBox(l2, s2)
	l3 := widget.NewLabel("Uninitialized IntSpinner 3:")
	s3 := xwidget.NewIntSpinnerUninitialized()
	c3 := container.NewHBox(l3, s3)
	b := widget.NewButton("Disable IntSpinner 1", func() {})
	b.OnTapped = func() {
		spinnerDisabled = !spinnerDisabled
		if spinnerDisabled {
			s1.Disable()
			b.SetText("Enable IntSpinner 1")
		} else {
			s1.Enable()
			b.SetText("Disable IntSpinner 1")
		}
	}
	bs = widget.NewButton("Initialize IntSpinner 3", func() {
		s3.SetMinMaxStep(1, 10, 1)
		l3.Text = "Initialized IntSpinner 3 (1, 10, 1):"
		l3.Refresh()
		s3.Enable()
		bs.Disable()
	})
	bs1 := widget.NewButton("Set IntSpinner 1 to 5", func() { s1.SetValue(5) })
	bs2 := widget.NewButton("Set IntSpinner 3 bound value to 12", func() { data.Set(12) })

	l4 := widget.NewLabel("Float64Spinner 1 (-1., 400., 10.3, 1):")
	s4 := xwidget.NewFloat64Spinner(-1., 400., 10.3, 1, nil)
	c4 := container.NewHBox(l4, s4)

	l5 := widget.NewLabel("Float64Spinner 2 (0., 16., 3.215, 2)")
	s5 = xwidget.NewFloat64SpinnerWithData(0., 16., 3.215, 2, floatData)
	c5 := container.NewHBox(l5, s5)

	v := container.NewVBox(c, c1, c3, b, bs, c2, bs1, bs2, c4, c5, c6)
	w := a.NewWindow("SpinnerDemo")
	w.SetContent(v)
	w.ShowAndRun()
}
