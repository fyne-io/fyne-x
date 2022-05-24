package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xtheme "fyne.io/x/fyne/theme"
	"fyne.io/x/fyne/theme/desktop"
)

func main() {
	app := app.New()
	app.Settings().SetTheme(xtheme.FromDesktopEnvironment())
	win := app.NewWindow("Desktop integration demo")
	win.Resize(fyne.NewSize(550, 390))
	win.CenterOnScreen()

	// Gnome/GTK theme invertion
	invertButton := widget.NewButton("Invert Gnome theme", func() {
		if t, ok := app.Settings().Theme().(*desktop.GnomeTheme); ok {
			t.Invert()
			win.Content().Refresh()
		}
	})

	// the invertButton can only work on Gnome / GTK theme.
	if _, ok := app.Settings().Theme().(*desktop.GnomeTheme); !ok {
		invertButton.Disable()
		invertButton.SetText("Invert only works on Gnome/GTK")
	}

	var switched bool
	switchThemeButton := widget.NewButton("Switch theme", func() {
		if switched {
			app.Settings().SetTheme(xtheme.FromDesktopEnvironment())
		} else {
			app.Settings().SetTheme(theme.DefaultTheme())
		}
		switched = !switched
		win.Content().Refresh()
	})

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Example of text entry...")
	win.SetContent(container.NewBorder(
		nil,
		container.NewHBox(
			widget.NewButtonWithIcon("Home icon button", theme.HomeIcon(), nil),
			widget.NewButtonWithIcon("Info icon button", theme.InfoIcon(), nil),
			widget.NewButtonWithIcon("Example file dialog", theme.FolderIcon(), func() {
				dialog.ShowFileSave(func(fyne.URIWriteCloser, error) {}, win)
			}),
			invertButton,
		),
		nil,
		nil,
		container.NewVBox(
			createExplanationLabel(app),
			entry,
			widget.NewLabel("Try to switch theme"),
			switchThemeButton,
		),
	))

	win.ShowAndRun()
}

func createExplanationLabel(app fyne.App) fyne.CanvasObject {

	var current string

	switch app.Settings().Theme().(type) {
	case *desktop.GnomeTheme:
		current = "Gnome / GTK"
	case *desktop.KDETheme:
		current = "KDE / Plasma"
	default:
		current = "This window manager is not supported for now"
	}

	text := "Current Desktop: " + current + "\n"
	text += `

This window should be styled to look like a desktop application. It works with GTK/Gnome based desktops and KDE/Plasma at this time
For the others desktops, the application will look like a normal window with default theme.

You may try to change icon theme or GTK/KDE theme in your desktop settings, as font, font scaling...

Note that you need to have fontforge package to make Fyne able to convert non-ttf fonts to ttf.
`
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}
