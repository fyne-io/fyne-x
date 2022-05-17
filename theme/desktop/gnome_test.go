package desktop

import "fyne.io/fyne/v2/app"

func ExampleNewGnomeTheme() {
	app := app.New()
	app.Settings().SetTheme(NewGnomeTheme(0))
}

// Force GTK version to 3
func ExampleNewGnomeTheme_forceGtkVersion() {
	app := app.New()
	app.Settings().SetTheme(NewGnomeTheme(3))
}

// Will reload theme when it changes in Gnome (or other GTK environment)
// connecting to DBus signal.
func ExampleNewGnomeTheme_autoReload() {
	app := app.New()
	app.Settings().SetTheme(NewGnomeTheme(0, GnomeFlagAutoReload))
}
