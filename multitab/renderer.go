package multitab

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type tabsRenderer struct {
	tabs        *Tabs
	tabBar      *fyne.Container
	contentArea *fyne.Container
	root        *fyne.Container
}

func (t *Tabs) CreateRenderer() fyne.WidgetRenderer {
	r := &tabsRenderer{
		tabs:        t,
		tabBar:      container.NewHBox(),
		contentArea: container.NewMax(),
	}

	r.root = container.NewBorder(
		r.tabBar,
		nil,
		nil,
		nil,
		r.contentArea,
	)

	r.refresh()
	return r
}

func (r *tabsRenderer) refresh() {
	r.tabBar.Objects = nil

	for i := range r.tabs.tabs {
		index := i
		tab := r.tabs.tabs[i]

		// Create tab button with icon if available
		var btn *widget.Button
		if tab.icon != nil {
			btn = widget.NewButtonWithIcon(tab.title, tab.icon, func() {
				r.tabs.SetActive(index)
			})
		} else {
			btn = widget.NewButton(tab.title, func() {
				r.tabs.SetActive(index)
			})
		}
		btn.Importance = widget.LowImportance

		// Set tooltip if available
		if tab.tooltip != "" {
			if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
				if hover, ok := drv.(interface {
					SetHoverable(*widget.Button, func(bool))
				}); ok {
					hover.SetHoverable(btn, func(over bool) {
						if over {
							// Show tooltip - this is a simplified implementation
							// In a real implementation, you'd want proper tooltip support
						}
					})
				}
			}
		}

		// Create close button if closing is allowed
		if r.tabs.allowClose {
			closeBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				r.tabs.RemoveTab(index)
			})
			closeBtn.Importance = widget.LowImportance

			r.tabBar.Add(container.NewHBox(btn, closeBtn))
		} else {
			r.tabBar.Add(btn)
		}
	}

	// Add "+" button for creating new tabs
	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		r.tabs.AddTab("New Tab", widget.NewLabel(""))
	})
	r.tabBar.Add(addBtn)

	if r.tabs.active >= 0 && r.tabs.active < len(r.tabs.tabs) {
		r.contentArea.Objects = []fyne.CanvasObject{
			container.NewPadded(r.tabs.tabs[r.tabs.active].content),
		}
	} else {
		r.contentArea.Objects = nil
	}

	r.tabBar.Refresh()
	r.contentArea.Refresh()
}

func (r *tabsRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
}

func (r *tabsRenderer) MinSize() fyne.Size {
	return r.root.MinSize()
}

func (r *tabsRenderer) Refresh() {
	r.refresh()
	canvas := fyne.CurrentApp().Driver().CanvasForObject(r.tabs)
	if canvas != nil {
		canvas.Refresh(r.tabs)
	}
}

func (r *tabsRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *tabsRenderer) Destroy() {}

func (r *tabsRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}
