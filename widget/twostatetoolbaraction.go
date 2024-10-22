package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TwoStateState defines the type for the state of a TwoStateToolbarAction.
type TwoStateState bool

const (
	OffState TwoStateState = false
	OnState  TwoStateState = true
)

// TwoStateToolbarAction is a push button style of ToolbarItem that displays a different
// icon depending on its state
type TwoStateToolbarAction struct {
	state       TwoStateState
	offIcon     fyne.Resource
	onIcon      fyne.Resource
	OnActivated func(TwoStateState) `json:"-"`

	button widget.Button
}

// NewTwoStateToolbarAction returns a new push button style of Toolbar item that displays
// a different icon for each of its two states
func NewTwoStateToolbarAction(offStateIcon fyne.Resource,
	onStateIcon fyne.Resource,
	onTapped func(TwoStateState)) *TwoStateToolbarAction {
	t := &TwoStateToolbarAction{offIcon: offStateIcon, onIcon: onStateIcon, OnActivated: onTapped}
	t.button.SetIcon(t.offIcon)
	t.button.OnTapped = t.activated
	return t
}

// GetState returns the current state of the toolbaraction
func (t *TwoStateToolbarAction) GetState() TwoStateState {
	return t.state
}

// SetState sets the state of the toolbaraction
func (t *TwoStateToolbarAction) SetState(state TwoStateState) {
	t.state = state
	if t.OnActivated != nil {
		t.OnActivated(t.state)
	}
	t.setButtonIcon()
	t.button.Refresh()
}

// SetOffStateIcon sets the icon that is displayed when the state is OffState
func (t *TwoStateToolbarAction) SetOffStateIcon(icon fyne.Resource) {
	t.offIcon = icon
	t.setButtonIcon()
	t.button.Refresh()
}

// SetOnStateIcon sets the icon that is displayed when the state is OnState
func (t *TwoStateToolbarAction) SetOnStateIcon(icon fyne.Resource) {
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
	if t.state == OffState {
		t.state = OnState
	} else {
		t.state = OffState
	}
	t.setButtonIcon()
	if t.OnActivated != nil {
		t.OnActivated(t.state)
	}
	t.button.Refresh()
}

func (t *TwoStateToolbarAction) setButtonIcon() {
	if t.state == OffState {
		t.button.Icon = t.offIcon
	} else {
		t.button.Icon = t.onIcon
	}
}
