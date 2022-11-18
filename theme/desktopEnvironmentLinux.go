//go:build linux
// +build linux

package theme

import (
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/x/fyne/theme/desktop"
)

// FromDesktopEnvironment returns a new WindowManagerTheme instance for the current desktop session.
// If the desktop manager is not supported or if it is not found, return the default theme
func FromDesktopEnvironment() fyne.Theme {
	wm := os.Getenv("XDG_CURRENT_DESKTOP")
	if wm == "" {
		wm = os.Getenv("DESKTOP_SESSION")
	}
	wm = strings.ToLower(wm)

	switch wm {
	case "gnome", "xfce", "unity", "gnome-shell", "gnome-classic", "mate", "gnome-mate":
		return desktop.NewGnomeTheme(-1, desktop.GnomeFlagAutoReload)
	case "kde", "kde-plasma", "plasma":
		return desktop.NewKDETheme()

	}

	log.Println("Window manager not supported:", wm, "using default theme")
	return theme.DefaultTheme()
}
