package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// must be in sync with adwaita_colors_generator.go - getting the colors from the Adwaita document page.
//go:generate go run ./adwaita_colors_generator.go

var _ fyne.Theme = (*Adwaita)(nil)

// AdwaitaTheme returns a new Adwaita theme.
func AdwaitaTheme() fyne.Theme {
	return &Adwaita{}
}

// Adwaita is a theme that follows the Adwaita theme.
// See: https://gnome.pages.gitlab.gnome.org/libadwaita/doc/main/named-colors.html
type Adwaita struct{}

// Color returns the named color for the current theme.
func (a *Adwaita) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch variant {
	case theme.VariantLight:
		if c, ok := adwaitaLightScheme[name]; ok {
			return c
		}
	case theme.VariantDark:
		if c, ok := adwaitaDarkScheme[name]; ok {
			return c
		}
	}
	return theme.DefaultTheme().Color(name, variant)
}

// Font returns the named font for the current theme.
func (a *Adwaita) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns the named resource for the current theme.
func (a *Adwaita) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size of the named resource for the current theme.
func (a *Adwaita) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
