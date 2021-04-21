package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CompletionEntry is an Entry with options displayed in a PopUpMenu.
type CompletionEntry struct {
	widget.Entry
	popupMenu     *widget.PopUp
	navigableList *navigableList
	Options       []string
	pause         bool
}

// NewCompletionEntry creates a new CompletionEntry which creates a popup menu that responds to keystrokes to navigate through the items without losing the editing ability of the text input.
func NewCompletionEntry(options []string) *CompletionEntry {
	c := &CompletionEntry{Options: options}
	c.ExtendBaseWidget(c)
	return c
}

// ShowCompletion displays the completion menu
func (c *CompletionEntry) ShowCompletion() {
	if c.pause {
		return
	}

	if c.navigableList == nil {
		c.navigableList = newNavigableList(c.Options, &c.Entry, c.setTextFromMenu, c.HideCompletion)
	}
	holder := fyne.CurrentApp().Driver().CanvasForObject(c)
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(c)

	if c.popupMenu == nil {
		c.popupMenu = widget.NewPopUp(c.navigableList, holder)
	}
	max := fyne.Min(10, float32(len(c.Options)))
	c.popupMenu.Resize(fyne.NewSize(c.Entry.Size().Width, c.popupMenu.MinSize().Height*max-3*theme.Padding()))
	c.popupMenu.ShowAtPosition(fyne.Position{X: pos.X, Y: pos.Y + c.Size().Height})
	holder.Focus(c.navigableList)
}

// HideCompletion hides the completion menu.
func (c *CompletionEntry) HideCompletion() {
	if c.popupMenu != nil {
		c.popupMenu.Hide()
	}
}

// SetOptions set the completion list with itemList and update the view.
func (c *CompletionEntry) SetOptions(itemList []string) {
	c.Options = itemList
	c.Refresh()
}

// Update the list to refresh the options to display.
func (c *CompletionEntry) Refresh() {
	c.Entry.Refresh()
	if c.navigableList != nil {
		c.navigableList.SetOptions(c.Options)
	}
}

// Prevent the menu to open when the user validate value from the menu.
func (c *CompletionEntry) setTextFromMenu(s string) {
	c.pause = true
	c.Entry.SetText(s)
	c.Entry.Refresh()
	c.pause = false
	c.popupMenu.Hide()
}

type navigableList struct {
	widget.List
	entry           *widget.Entry
	selected        int
	setTextFromMenu func(string)
	hide            func()
	navigating      bool
	items           []string
}

func newNavigableList(items []string, entry *widget.Entry, setTextFromMenu func(string), hide func()) *navigableList {
	n := &navigableList{
		entry:           entry,
		selected:        -1,
		setTextFromMenu: setTextFromMenu,
		hide:            hide,
		items:           items,
	}

	n.List = widget.List{
		Length: func() int {
			return len(n.items)
		},
		CreateItem: func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		UpdateItem: func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(n.items[i])
		},
		OnSelected: func(id widget.ListItemID) {
			if !n.navigating {
				setTextFromMenu(n.items[id])
			}
			n.navigating = false
		},
	}
	n.ExtendBaseWidget(n)
	return n
}

func (n *navigableList) SetOptions(items []string) {
	n.Unselect(n.selected)
	n.items = items
	n.Refresh()
	n.selected = -1
}

func (n *navigableList) TypedKey(event *fyne.KeyEvent) {
	n.entry.TypedKey(event)
	switch event.Name {
	case fyne.KeyDown:
		if n.selected < len(n.items)-1 {
			n.selected++
		} else {
			n.selected = 0
		}
		n.navigating = true
		n.Select(n.selected)

	case fyne.KeyUp:
		if n.selected > 0 {
			n.selected--
		} else {
			n.selected = len(n.items) - 1
		}
		n.navigating = true
		n.Select(n.selected)
	case fyne.KeyReturn:
		n.navigating = false
		n.OnSelected(n.selected)
	case fyne.KeyEscape:
		n.hide()

	}
}
func (n *navigableList) TypedRune(r rune) {
	n.entry.TypedRune(r)
}
func (n *navigableList) FocusGained() {}
func (n *navigableList) FocusLost()   {}