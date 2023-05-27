package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CompletionEntry is an Entry with options displayed in a PopUpMenu.
type CompletionEntry struct {
	widget.Entry

	Options           []string
	OnCompleted       func(string) string // Called when you select an item from the list ; return different string to override the completion
	SubmitOnCompleted bool                // True will call Entry.OnSubmited when you select a completion item

	CustomCreate func() fyne.CanvasObject
	CustomUpdate func(id widget.ListItemID, object fyne.CanvasObject)

	popup    *widget.PopUp
	list     *completionEntryList
	selected widget.ListItemID
	pause    bool
}

// NewCompletionEntry creates a new CompletionEntry which creates a popup menu that responds to keystrokes to navigate through the items without losing the editing ability of the text input.
func NewCompletionEntry(options []string) *CompletionEntry {
	c := &CompletionEntry{Options: options}
	c.ExtendBaseWidget(c)
	return c
}

// HideCompletion hides the completion menu.
func (c *CompletionEntry) HideCompletion() {
	if c.popup != nil {
		c.list.UnselectAll()
		c.popup.Hide()
	}
}

// Move changes the relative position of the select entry.
//
// Implements: fyne.Widget
func (c *CompletionEntry) Move(pos fyne.Position) {
	c.Entry.Move(pos)
	if c.popup != nil && c.popup.Visible() {
		c.popup.Move(c.popUpPos())
		c.popup.Resize(c.maxSize())
	}
}

// Refresh the list to update the options to display.
func (c *CompletionEntry) Refresh() {
	c.Entry.Refresh()
	if c.list != nil && c.list.Visible() {
		c.list.Refresh()
	}
}

// SetOptions set the completion list with itemList and update the view.
func (c *CompletionEntry) SetOptions(itemList []string) {
	c.Options = itemList
	c.Refresh()
}

// ShowCompletion displays the completion menu
func (c *CompletionEntry) ShowCompletion() {
	if c.pause {
		return
	}
	if len(c.Options) <= 0 {
		c.HideCompletion()
		return
	}

	cnv := fyne.CurrentApp().Driver().CanvasForObject(c)
	if cnv == nil {
		return // canvas acquisiton failed (widget not showed yet?)
	}

	if c.list == nil {
		c.list = newCompletionEntryList(c)
	}
	if c.popup == nil {
		c.popup = widget.NewPopUp(c.list, cnv)
	}

	c.popup.ShowAtPosition(c.popUpPos())
	c.popup.Resize(c.maxSize())

	c.list.Select(0)
	cnv.Focus(c.list)
}

// calculate the max size to make the popup to cover everything below the entry
func (c *CompletionEntry) maxSize() fyne.Size {
	cnv := fyne.CurrentApp().Driver().CanvasForObject(c)

	// return empty size if cannot get canvas (widget not showed yet?)
	if cnv == nil {
		return fyne.Size{}
	}

	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(c)

	// define size boundaries
	minWidth := c.Size().Width
	maxWidth := cnv.Size().Width - pos.X - theme.Padding()
	maxHeight := cnv.Size().Height - pos.Y - c.MinSize().Height - 2*theme.Padding()

	// iterating items until the end or we rech maxHeight
	var width, height float32
	for i := 0; i < len(c.Options); i++ {
		item := c.list.CreateItem()
		c.list.UpdateItem(i, item)
		sz := item.MinSize()
		if sz.Width > width {
			width = sz.Width
		}
		height += sz.Height + theme.Padding()
		if height > maxHeight {
			height = maxHeight
			break
		}
	}
	height += theme.Padding() // popup padding

	width += 2 * theme.Padding() // let some padding on the trailing end of the longest item
	if width < minWidth {
		width = minWidth
	}
	if width > maxWidth {
		width = maxWidth
	}

	return fyne.NewSize(width, height)
}

// calculate where the popup should appear
func (c *CompletionEntry) popUpPos() fyne.Position {
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(c)
	return pos.Add(fyne.NewPos(0, c.Size().Height+theme.Padding()))
}

// Prevent the menu to open when the user validate value from the menu.
func (c *CompletionEntry) setTextFromMenu(s string) {
	c.popup.Hide()
	c.pause = true
	if c.OnCompleted != nil {
		s = c.OnCompleted(s)
	}
	c.Entry.Text = s
	c.Entry.CursorColumn = len([]rune(s))
	c.Entry.Refresh()
	c.pause = false
	if c.SubmitOnCompleted && c.OnSubmitted != nil {
		c.OnSubmitted(c.Entry.Text)
	}
}

type completionEntryList struct {
	widget.List

	// just hold a reference to the "parent" CompletionEntry that already holds all what the list needs to operate
	parent *CompletionEntry
}

func newCompletionEntryList(parent *CompletionEntry) *completionEntryList {
	list := &completionEntryList{parent: parent}
	list.ExtendBaseWidget(list)

	list.List.Length = func() int { return len(parent.Options) }
	list.List.CreateItem = func() fyne.CanvasObject {
		var item *completionEntryListItem
		if parent.CustomCreate != nil {
			item = newCompletionEntryListItem(parent, parent.CustomCreate())
		} else {
			item = newCompletionEntryListItem(parent, widget.NewLabel(""))
		}
		return item
	}
	list.List.UpdateItem = func(id widget.ListItemID, co fyne.CanvasObject) {
		if parent.CustomUpdate != nil {
			parent.CustomUpdate(id, co.(*completionEntryListItem).co)
		} else {
			co.(*completionEntryListItem).co.(*widget.Label).Text = parent.Options[id]
		}
		co.(*completionEntryListItem).id = id
		list.SetItemHeight(id, co.MinSize().Height)
		co.Refresh()
	}
	list.List.OnSelected = func(id widget.ListItemID) {
		parent.selected = id
	}
	list.List.OnUnselected = func(_ widget.ListItemID) {
		parent.selected = -1
	}
	return list
}

// Implements: fyne.Focusable
func (list *completionEntryList) FocusGained() {}

// Implements: fyne.Focusable
func (list *completionEntryList) FocusLost() {}

// Implements: fyne.Focusable
func (list *completionEntryList) TypedKey(ke *fyne.KeyEvent) {
	switch ke.Name {
	case fyne.KeyDown:
		if list.parent.selected < len(list.parent.Options)-1 {
			list.parent.list.Select(list.parent.selected + 1)
		} else {
			list.parent.list.Select(0)
		}
	case fyne.KeyUp:
		if list.parent.selected > 0 {
			list.parent.list.Select(list.parent.selected - 1)
		} else {
			list.parent.list.Select(len(list.parent.Options) - 1)
		}
	case fyne.KeyReturn, fyne.KeyEnter:
		if list.parent.selected >= 0 {
			list.parent.setTextFromMenu(list.parent.Options[list.parent.selected])
		} else {
			list.parent.HideCompletion()
			list.parent.Entry.TypedKey(ke)
		}
	case fyne.KeyTab, fyne.KeyEscape:
		list.parent.HideCompletion()
	default:
		list.parent.TypedKey(ke)
	}
}

// Implements: fyne.Focusable
func (list *completionEntryList) TypedRune(r rune) {
	list.parent.TypedRune(r)
}

// Implements: fyne.Shortcutable
func (list *completionEntryList) TypedShortcut(s fyne.Shortcut) {
	list.parent.TypedShortcut(s)
}

type completionEntryListItem struct {
	widget.BaseWidget
	parent *CompletionEntry
	id     widget.ListItemID // each list item knows where it belongs in parent's data slice
	co     fyne.CanvasObject
}

func newCompletionEntryListItem(parent *CompletionEntry, co fyne.CanvasObject) *completionEntryListItem {
	item := &completionEntryListItem{parent: parent, id: -1, co: co}
	item.ExtendBaseWidget(item)
	return item
}

func (item *completionEntryListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(item.co)
}

func (item *completionEntryListItem) Tapped(_ *fyne.PointEvent) {
	item.parent.setTextFromMenu(item.parent.Options[item.id])
}

func (item *completionEntryListItem) MouseIn(_ *desktop.MouseEvent)    { item.parent.list.Select(item.id) }
func (item *completionEntryListItem) MouseMoved(_ *desktop.MouseEvent) {}
func (item *completionEntryListItem) MouseOut()                        {}
