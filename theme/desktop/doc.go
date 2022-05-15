// desktop package provides theme for desktop manager like Gnome, KDE, Plasma...
// To be fully used, the system need to have gsettings and gjs for all GTK/Gnome based desktop.
// KDE/Plasma theme only works when the user has already initialize a session to create ~/.config/kdeglobals
//
// For all desktop, we also need fontconfig package installed to have "fc-match" and "fontconfif" commands.
// This is not required but recommended to be able to generate TTF font if the user as configure a non TTF font.
package desktop
