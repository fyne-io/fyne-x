package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TwoStateState defines the type for the state of a TwoStateToolbarAction.
type TwoStateState bool

const (
	TwoState0 TwoStateState = false
	TwoState1 TwoStateState = true
)

// TwoStateToolbarAction is a push button style of ToolbarItem that displays a different
// icon depending on its state
type TwoStateToolbarAction struct {
	state       TwoStateState
	icon0       fyne.Resource
	Icon1       fyne.Resource
	OnActivated func(TwoStateState) `json:"-"`

	button widget.Button
}

// NewTwoStateToolbarAction returns a new push button style of Toolbar item that displays
// a different icon for each of its two states
func NewTwoStateToolbarAction(state0Icon fyne.Resource,
	state1Icon fyne.Resource,
	onTapped func(TwoStateState)) *TwoStateToolbarAction {
	t := &TwoStateToolbarAction{icon0: state0Icon, Icon1: state1Icon, OnActivated: onTapped}
	t.button.SetIcon(t.icon0)
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

// SetState0Icon sets the icon that is displayed when the state is TwoState0
func (t *TwoStateToolbarAction) SetState0Icon(icon fyne.Resource) {
	t.icon0 = icon
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
	if t.state == TwoState0 {
		t.state = TwoState1
	} else {
		t.state = TwoState0
	}
	t.setButtonIcon()
	if t.OnActivated != nil {
		t.OnActivated(t.state)
	}
	t.button.Refresh()
}

func (t *TwoStateToolbarAction) setButtonIcon() {
	if t.state == TwoState0 {
		t.button.Icon = t.icon0
	} else {
		t.button.Icon = t.Icon1
	}
}
