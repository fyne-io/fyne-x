//go:build !linux
// +build !linux

package theme

import fynetheme "fyne.io/x/fyne/v2/theme"

func FromDesktopEnvironment() fyne.Theme {
	return fynetheme.DefaultTheme()
}
