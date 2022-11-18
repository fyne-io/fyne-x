package theme

import (
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/theme/desktop"
)

func setup() (tmp, home string) {
	// create a false home directory
	var err error
	tmp, err = os.MkdirTemp("", "fyne-test-")
	if err != nil {
		panic(err)
	}
	home = os.Getenv("HOME")
	os.Setenv("HOME", tmp)

	// creat a false KDE configuration
	os.MkdirAll(tmp+"/.config", 0755)
	os.WriteFile(tmp+"/.config/kdeglobals", []byte("[General]\nwidgetStyle=GTK"), 0644)

	return
}

func teardown(tmp, home string) {
	os.Unsetenv("XDG_CURRENT_DESKTOP")
	os.RemoveAll(tmp)
	os.Setenv("HOME", home)
}

// ExampleFromDesktopEnvironment_simple demonstrates how to use the FromDesktopEnvironment function.
func ExampleFromDesktopEnvironment_simple() {
	app := app.New()
	theme := FromDesktopEnvironment()
	app.Settings().SetTheme(theme)
}

// Test to load from desktop environment.
func TestLoadFromEnvironment(t *testing.T) {
	tmp, home := setup()
	defer teardown(tmp, home)

	// Set XDG_CURRENT_DESKTOP to "GNOME"
	envs := []string{"GNOME", "KDE", "FAKE"}
	for _, env := range envs {
		// chante desktop environment
		os.Setenv("XDG_CURRENT_DESKTOP", env)
		app := test.NewApp()
		app.Settings().SetTheme(FromDesktopEnvironment())
		win := app.NewWindow("Test")
		defer win.Close()
		win.Resize(fyne.NewSize(200, 200))
		win.SetContent(widget.NewLabel("Hello"))

		// check if the theme is loaded
		current := app.Settings().Theme()
		// Check if the type of the theme is correct
		if current == nil {
			t.Error("Theme is nil")
		}
		switch env {
		case "GNOME":
			switch v := current.(type) {
			case *desktop.GnomeTheme:
				// OK
			default:
				t.Error("Theme is not GnomeTheme")
				t.Logf("Theme is %T\n", v)
			}
		case "KDE":
			switch v := current.(type) {
			case *desktop.KDETheme:
				// OK
			default:
				t.Error("Theme is not KDETheme")
				t.Logf("Theme is %T\n", v)
			}
		case "FAKE":
			if current != theme.DefaultTheme() {
				t.Error("Theme is not DefaultTheme")
			}
		}

	}
}
