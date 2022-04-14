// Package charts provides a simple way to create charts.
//
// Each chart can change the stroke color, line width, and fille color. They are all redreshed when you change data (via `SetData()` method).
//
// It's possible to draw over the chart. Use `GetDrawable()` method to get the canvas, and use `fyne.Canvas` API to add lines, circles, square, text...
// See the example in fyne-x repository to see some examples.
//
// The charts types are scaled from minimal value to maximal value.
package charts
