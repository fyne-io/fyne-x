//go:build !linux
// +build !linux

package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// NewKdeTheme returns the KDE theme. If the current OS is not Linux, it returns the default theme.
func NewKDETheme() fyne.Theme {
	return theme.DefaultTheme()
}
