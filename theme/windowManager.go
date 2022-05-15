package theme

import (
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/x/fyne/theme/desktop"
)

// NewWindowManagerTheme returns a new WindowManagerTheme instance for the current desktop session. If the desktop manager is
// not supported or if it is not found, return the default theme
func NewWindowManagerTheme() fyne.Theme {
	wm := os.Getenv("XDG_CURRENT_DESKTOP")
	if wm == "" {
		wm = os.Getenv("DESKTOP_SESSION")
	}
	wm = strings.ToUpper(wm)

	switch wm {
	case "GNOME", "XFCE", "UNITY", "GNOME-SHELL", "GNOME-CLASSIC", "MATE", "GNOME-MATE":
		return desktop.NewGnomeTheme(-1)
	case "KDE", "KDE-PLASMA", "PLASMA":
		return desktop.NewKDETheme()

	}

	log.Println("Window manager not supported:", wm, "using default theme")
	return theme.DefaultTheme()
}
