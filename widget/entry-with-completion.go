package widget

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CompletionEntry is an Entry with options displayed in a PopUpMenu.
type CompletionEntry struct {
	widget.Entry
	popupMenu     *widget.PopUp
	navigableList *navigableList
	Items         []string
	pause         bool
}

// NewCompletionEntry creates a new CompletionEntry which creates a popup menu that responds to keystrokes to navigate through the items without losing the editing ability of the text input.
func NewCompletionEntry(options []string) *CompletionEntry {
	c := &CompletionEntry{Items: options}
	c.ExtendBaseWidget(c)
	return c
}

// ShowCompletion displays the completion menu
func (c *CompletionEntry) ShowCompletion() {
	if c.pause {
		return
	}

	if c.navigableList == nil {
		c.navigableList = newNavigableList(c.Items, &c.Entry, c.setTextFromMenu, c.HideCompletion)
	}
	holder := fyne.CurrentApp().Driver().CanvasForObject(c)
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(c)

	if c.popupMenu == nil {
		c.popupMenu = widget.NewPopUp(c.navigableList, holder)
	}
	max := fyne.Min(10, float32(len(c.Items)))
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

// SetOptions set the completion list with itemList.
func (c *CompletionEntry) SetOptions(itemList []string) {
	c.Items = itemList
	if c.navigableList != nil {
		c.navigableList.SetOptions(itemList)
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
	data            binding.StringList
	entry           *widget.Entry
	selected        int
	setTextFromMenu func(string)
	hide            func()
	navigating      bool
}

func newNavigableList(items []string, entry *widget.Entry, setTextFromMenu func(string), hide func()) *navigableList {
	n := &navigableList{entry: entry, selected: -1, setTextFromMenu: setTextFromMenu, hide: hide}
	n.data = binding.BindStringList(&items)

	n.List = widget.List{
		Length: func() int {
			return n.data.Length()
		},
		CreateItem: func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		UpdateItem: func(i widget.ListItemID, o fyne.CanvasObject) {
			item, err := n.data.GetItem(i)
			if err != nil {
				fyne.LogError(fmt.Sprintf("Error getting data item %d", i), err)
				return
			}
			o.(*widget.Label).Bind(item.(binding.String))
		},
	}

	n.ExtendBaseWidget(n)

	n.OnSelected = func(id int) {
		if !n.navigating {
			item, _ := n.data.GetItem(id)
			val, _ := item.(binding.String).Get()
			setTextFromMenu(val)
		}
		n.navigating = false
	}

	return n
}

func (n *navigableList) SetOptions(items []string) {
	n.Unselect(n.selected)
	n.data.Set(items)
	n.selected = -1
}

func (n *navigableList) TypedKey(event *fyne.KeyEvent) {
	n.entry.TypedKey(event)
	switch event.Name {
	case fyne.KeyDown:
		if n.selected < n.data.Length()-1 {
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
			n.selected = n.data.Length() - 1
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
