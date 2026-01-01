package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Aswanidev-vs/fyne-multitab/multitab"
)

func main() {
	a := app.New()
	w := a.NewWindow("Enhanced MultiTab Example")
	w.Resize(fyne.NewSize(900, 600))

	tabs := multitab.New()

	// Configure tabs
	tabs.SetAllowClose(true)
	tabs.SetKeyboardNavigation(true)

	// Tab change callback
	tabs.OnTabChange(func(index int) {
		fmt.Printf("Switched to tab %d: %s\n", index, tabs.GetTabTitle(index))
	})

	// Tab removed callback - close window if no tabs left
	tabs.OnTabRemoved(func() {
		if tabs.TabCount() == 0 {
			w.Close()
		}
	})

	// Home tab with icon and tooltip
	homeContent := container.NewVBox(
		widget.NewLabel("Welcome to Enhanced MultiTab"),
		widget.NewLabel("Features:"),
		widget.NewLabel("• Tab icons and tooltips"),
		widget.NewLabel("• Keyboard navigation (Tab, Escape)"),
		widget.NewLabel("• Close buttons"),
		widget.NewLabel("• Event callbacks"),
		widget.NewButton("Add New Tab", func() {
			count := tabs.TabCount() + 1
			tabs.AddTab(fmt.Sprintf("Tab %d", count), widget.NewLabel(fmt.Sprintf("Content for tab %d", count)))
		}),
	)

	tabs.AddTabWithIconAndTooltip("Home", homeContent, theme.HomeIcon(), "Welcome page")

	// Settings tab
	settingsContent := container.NewVBox(
		widget.NewLabel("Application Settings"),
		widget.NewCheck("Enable notifications", func(b bool) {}),
		widget.NewCheck("Dark mode", func(b bool) {}),
		widget.NewSelect([]string{"English", "Spanish", "French"}, func(s string) {}),
		widget.NewButton("Save Settings", func() {}),
		widget.NewSeparator(),
		widget.NewLabel("Tab Configuration:"),
		widget.NewCheck("Allow closing tabs", func(b bool) {
			tabs.SetAllowClose(b)
		}),
		widget.NewCheck("Keyboard navigation", func(b bool) {
			tabs.SetKeyboardNavigation(b)
		}),
	)

	tabs.AddTabWithIconAndTooltip("Settings", settingsContent, theme.SettingsIcon(), "Configure application")

	// Logs tab
	logs := widget.NewMultiLineEntry()
	logs.SetText("Application started...\nEnhanced MultiTab initialized.\n")
	logs.Disable() // Make it read-only

	tabs.AddTabWithIconAndTooltip("Logs", logs, theme.InfoIcon(), "Application logs")

	// Demo tab with interactive content
	demoContent := container.NewVBox(
		widget.NewLabel("Interactive Demo"),
		widget.NewButton("Show Tab Info", func() {
			info := fmt.Sprintf("Active tab: %d/%d\n", tabs.ActiveIndex()+1, tabs.TabCount())
			for i := 0; i < tabs.TabCount(); i++ {
				info += fmt.Sprintf("Tab %d: %s\n", i+1, tabs.GetTabTitle(i))
			}
			widget.NewPopUp(widget.NewLabel(info), w.Canvas()).Show()
		}),
		widget.NewButton("Insert Tab at Position 1", func() {
			tabs.InsertTab(1, "Inserted", widget.NewLabel("This tab was inserted"))
		}),
		widget.NewButton("Move Current Tab Left", func() {
			active := tabs.ActiveIndex()
			if active > 0 {
				tabs.MoveTab(active, active-1)
			}
		}),
		widget.NewButton("Move Current Tab Right", func() {
			active := tabs.ActiveIndex()
			if active < tabs.TabCount()-1 {
				tabs.MoveTab(active, active+1)
			}
		}),
	)

	tabs.AddTabWithIconAndTooltip("Demo", demoContent, theme.ComputerIcon(), "Interactive features demo")

	w.SetContent(tabs)
	w.ShowAndRun()
}
