package theme

import (
	"fyne.io/fyne/v2/app"
)

// ExampleFromDesktopEnvironment_simple demonstrates how to use the FromDesktopEnvironment function.
func ExampleFromDesktopEnvironment_simple() {
	app := app.New()
	theme := FromDesktopEnvironment()
	app.Settings().SetTheme(theme)
}
