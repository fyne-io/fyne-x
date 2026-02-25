package multitab

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TabChangeCallback func(index int)
type TabRemovedCallback func()

type Tabs struct {
	widget.BaseWidget

	tabs         []Tab
	active       int
	onTabChange  TabChangeCallback
	onTabRemoved TabRemovedCallback
	allowClose   bool
	allowReorder bool
	keyboardNav  bool
}

func New() *Tabs {
	t := &Tabs{
		active:       -1,
		allowClose:   true,
		allowReorder: false,
		keyboardNav:  true,
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *Tabs) AddTab(title string, content fyne.CanvasObject) {
	t.AddTabWithIcon(title, content, nil)
}

func (t *Tabs) AddTabWithIcon(title string, content fyne.CanvasObject, icon fyne.Resource) {
	t.AddTabWithIconAndTooltip(title, content, icon, "")
}

func (t *Tabs) AddTabWithIconAndTooltip(title string, content fyne.CanvasObject, icon fyne.Resource, tooltip string) {
	t.tabs = append(t.tabs, Tab{
		title:   title,
		content: content,
		icon:    icon,
		tooltip: tooltip,
	})

	t.active = len(t.tabs) - 1
	t.Refresh()
	if t.onTabChange != nil {
		t.onTabChange(t.active)
	}
}

func (t *Tabs) RemoveTab(index int) {
	if index < 0 || index >= len(t.tabs) {
		return
	}

	t.tabs = append(t.tabs[:index], t.tabs[index+1:]...)

	if t.active >= len(t.tabs) {
		t.active = len(t.tabs) - 1
	}

	t.Refresh()
	if t.onTabRemoved != nil {
		t.onTabRemoved()
	}
}

func (t *Tabs) SetActive(index int) {
	if index < 0 || index >= len(t.tabs) {
		return
	}
	t.active = index
	t.Refresh()
}

func (t *Tabs) ActiveIndex() int {
	return t.active
}

func (t *Tabs) TabCount() int {
	return len(t.tabs)
}

// Configuration methods
func (t *Tabs) SetAllowClose(allow bool) {
	t.allowClose = allow
	t.Refresh()
}

func (t *Tabs) SetAllowReorder(allow bool) {
	t.allowReorder = allow
	t.Refresh()
}

func (t *Tabs) SetKeyboardNavigation(enable bool) {
	t.keyboardNav = enable
}

// Event callbacks
func (t *Tabs) OnTabChange(callback TabChangeCallback) {
	t.onTabChange = callback
}

func (t *Tabs) OnTabRemoved(callback TabRemovedCallback) {
	t.onTabRemoved = callback
}

// Get tab information
func (t *Tabs) GetTabTitle(index int) string {
	if index < 0 || index >= len(t.tabs) {
		return ""
	}
	return t.tabs[index].title
}

func (t *Tabs) GetTabIcon(index int) fyne.Resource {
	if index < 0 || index >= len(t.tabs) {
		return nil
	}
	return t.tabs[index].icon
}

func (t *Tabs) GetTabTooltip(index int) string {
	if index < 0 || index >= len(t.tabs) {
		return ""
	}
	return t.tabs[index].tooltip
}

// Tab manipulation
func (t *Tabs) InsertTab(index int, title string, content fyne.CanvasObject) {
	if index < 0 || index > len(t.tabs) {
		return
	}

	tab := Tab{
		title:   title,
		content: content,
	}

	if index == len(t.tabs) {
		t.tabs = append(t.tabs, tab)
	} else {
		t.tabs = append(t.tabs[:index], append([]Tab{tab}, t.tabs[index:]...)...)
	}

	t.active = index
	t.Refresh()
	if t.onTabChange != nil {
		t.onTabChange(t.active)
	}
}

func (t *Tabs) MoveTab(from, to int) {
	if from < 0 || from >= len(t.tabs) || to < 0 || to >= len(t.tabs) {
		return
	}

	tab := t.tabs[from]
	t.tabs = append(t.tabs[:from], t.tabs[from+1:]...)
	t.tabs = append(t.tabs[:to], append([]Tab{tab}, t.tabs[to:]...)...)

	if t.active == from {
		t.active = to
	} else if t.active > from && t.active <= to {
		t.active--
	} else if t.active < from && t.active >= to {
		t.active++
	}

	t.Refresh()
}

// Keyboard navigation
func (t *Tabs) KeyDown(key *fyne.KeyEvent) {
	if !t.keyboardNav {
		return
	}

	// Simple keyboard shortcuts without complex modifier checking
	switch key.Name {
	case fyne.KeyT:
		// T: New tab (could be extended to show a dialog)
		t.AddTab("New Tab", widget.NewLabel("New tab content"))
	case fyne.KeyW:
		// W: Close current tab
		if t.active >= 0 {
			t.RemoveTab(t.active)
		}
	case fyne.KeyTab:
		// Tab: Next tab
		next := t.active + 1
		if next >= len(t.tabs) {
			next = 0
		}
		t.SetActive(next)
	}
}

func (t *Tabs) TypedRune(r rune) {
	// Handle typed runes if needed
}

func (t *Tabs) FocusGained() {
	// Handle focus gained
}

func (t *Tabs) FocusLost() {
	// Handle focus lost
}

func (t *Tabs) TypedShortcut(shortcut fyne.Shortcut) {
	// Handle shortcuts
}
