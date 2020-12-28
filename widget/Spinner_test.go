package widget_test

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	xwidget "fyne.io/x/fyne/widget"
)

// This Example demonstrates a few uses for the Spinner widget.
func Example() {
	var (
		a = app.NewWithID("x.fyne.demo.spinner")
		w = a.NewWindow("Spinner Demo")

		basic     = xwidget.NewSpinner(0, 1)
		integer   = xwidget.NewIntSpinner(0, 1)
		negative  = xwidget.NewSpinner(-5, 0.637)
		positive  = xwidget.NewSpinner(5, 0.637)
		precision = xwidget.NewSpinner(0, 0.28671058)
	)

	negative.SetMax(0)
	positive.SetMin(0)

	// shows the value as accurately as possible (default is to round to SpinnerDefaultPrecision)
	precision.SetPrecision(-1)

	w.SetContent(container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Basic", basic),
			widget.NewFormItem("Integer", integer),
			widget.NewFormItem("Max", negative),
			widget.NewFormItem("Min", positive),
			widget.NewFormItem("Precision", precision),
		),
		fyne.NewContainerWithLayout(
			layout.NewGridLayout(2),
			widget.NewButton("Disable", func() {
				basic.Disable()
				integer.Disable()
				negative.Disable()
				positive.Disable()
				precision.Disable()
			}),
			widget.NewButton("Enable", func() {
				basic.Enable()
				integer.Enable()
				negative.Enable()
				positive.Enable()
				precision.Enable()
			}),
		),
	))

	w.CenterOnScreen()
	w.ShowAndRun()
}
