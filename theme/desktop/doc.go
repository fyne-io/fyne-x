// Package desktop provides theme for Linux (for now) desktop environment like Gnome, KDE, Plasma...
//
// To be fully used, the system need to have gsettings and gjs for all GTK/Gnome based desktop.
// KDE/Plasma theme only works when the user has already initialize a session to create ~/.config/kdeglobals
//
// The package will try to use fontconfig ("fc-match" and "fontconfig" commands).
// This is not required but recommended to be able to generate TTF font if the user as configure a non TTF font.
//
// The package also tries to use "inkscape" or "convert" command (from ImageMagick) to generate SVG icons when they
// cannot be parsed by Fyne (it happens when oksvg package fails to load icons).
//
// Some recent desktop environment now use Adwaita as default theme. If this theme is applied, the desktop package
// loads the default Fyne theme colors. It only try to change the scaling factor, font and icons of the applications.
//
// The easiest way to use this package is to call the FromDesktopEnvironment function from "theme" package.
//
// Example:
//
//  app := app.New()
//  theme := FromDesktopEnvironment()
//  app.Settings().SetTheme(theme)
//
// This loads the theme from the current detected desktop environment. For Windows and MacOS, and mobile devices
// it will return the default Fyne theme.
package desktop
