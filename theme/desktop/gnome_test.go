//go:build linux
// +build linux

package desktop

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

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

// Check if  the GnomeTheme can be loaded.
func TestGnomeTheme(t *testing.T) {
	app := test.NewApp()
	app.Settings().SetTheme(NewGnomeTheme(0))
	win := app.NewWindow("Test")
	defer win.Close()
	win.Resize(fyne.NewSize(200, 200))
	win.SetContent(widget.NewLabel("Hello"))
}
