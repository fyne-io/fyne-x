//go:build !linux
// +build !linux

package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// NewGnomeTheme returns the GNOME theme. If the current OS is not Linux, it returns the default theme.
func NewGnomeTheme() fyne.Theme {
	return theme.DefaultTheme()
}
