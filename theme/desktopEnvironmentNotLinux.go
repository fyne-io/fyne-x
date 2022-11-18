//go:build !linux
// +build !linux

package theme

import (
	"fyne.io/fyne/v2"
	fynetheme "fyne.io/fyne/v2/theme"
)

func FromDesktopEnvironment() fyne.Theme {
	return fynetheme.DefaultTheme()
}
