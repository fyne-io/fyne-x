package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TwoStateToolbarAction is a push button style of ToolbarItem that displays a different
// icon depending on its state.
//
// state is a boolean indicating off and on. The actual meaning of the boolean depends on how it is used. For
// example, in a media play app, false might indicate that the medium file is not being played, and true might
// indicate that the file is being played.
// Similarly, the two states could be used to indicate that a panel is being hidden or shown.
type TwoStateToolbarAction struct {
	on          bool
	offIcon     fyne.Resource
	onIcon      fyne.Resource
	OnActivated func(bool) `json:"-"`

	button widget.Button
}

// NewTwoStateToolbarAction returns a new push button style of Toolbar item that displays
// a different icon for each of its two states
func NewTwoStateToolbarAction(offStateIcon fyne.Resource,
	onStateIcon fyne.Resource,
	onTapped func(bool)) *TwoStateToolbarAction {
	t := &TwoStateToolbarAction{offIcon: offStateIcon, onIcon: onStateIcon, OnActivated: onTapped}
	t.button.SetIcon(t.offIcon)
	t.button.OnTapped = t.activated
	return t
}

// GetOn returns the current state of the toolbaraction
func (t *TwoStateToolbarAction) GetOn() bool {
	return t.on
}

// SetOn sets the state of the toolbaraction
func (t *TwoStateToolbarAction) SetOn(on bool) {
	t.on = on
	if t.OnActivated != nil {
		t.OnActivated(t.on)
	}
	t.setButtonIcon()
	t.button.Refresh()
}

// SetOffIcon sets the icon that is displayed when the state is false
func (t *TwoStateToolbarAction) SetOffIcon(icon fyne.Resource) {
	t.offIcon = icon
	t.setButtonIcon()
	t.button.Refresh()
}

// SetOnIcon sets the icon that is displayed when the state is true
func (t *TwoStateToolbarAction) SetOnIcon(icon fyne.Resource) {
	t.onIcon = icon
	t.setButtonIcon()
	t.button.Refresh()
}

// ToolbarObject gets a button to render this ToolbarAction
func (t *TwoStateToolbarAction) ToolbarObject() fyne.CanvasObject {
	t.button.Importance = widget.LowImportance

	// synchronize properties
	t.setButtonIcon()
	t.button.OnTapped = t.activated
	return &t.button
}

func (t *TwoStateToolbarAction) activated() {
	if !t.on {
		t.on = true
	} else {
		t.on = false
	}
	t.setButtonIcon()
	if t.OnActivated != nil {
		t.OnActivated(t.on)
	}
	t.button.Refresh()
}

func (t *TwoStateToolbarAction) setButtonIcon() {
	if !t.on {
		t.button.Icon = t.offIcon
	} else {
		t.button.Icon = t.onIcon
	}
}
