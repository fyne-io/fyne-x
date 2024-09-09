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
	State       TwoStateState
	Icon0       fyne.Resource
	Icon1       fyne.Resource
	OnActivated func(TwoStateState) `json:"-"`

	button widget.Button
}

// NewTwoStateToolbarAction returns a new push button style of Toolbar item that displays
// a different icon for each of its two states
func NewTwoStateToolbarAction(icon0 fyne.Resource,
	icon1 fyne.Resource,
	onTapped func(TwoStateState)) *TwoStateToolbarAction {
	t := &TwoStateToolbarAction{Icon0: icon0, Icon1: icon1, OnActivated: onTapped}
	t.button.SetIcon(t.Icon0)
	t.button.OnTapped = t.activated
	return t
}

func (t *TwoStateToolbarAction) ToolbarObject() fyne.CanvasObject {
	t.button.Importance = widget.LowImportance

	// synchronize properties
	if t.State == TwoState0 {
		t.button.Icon = t.Icon0
	} else {
		t.button.Icon = t.Icon1
	}
	t.button.OnTapped = t.activated
	return &t.button
}

func (t *TwoStateToolbarAction) activated() {
	if t.State == TwoState0 {
		t.State = TwoState1
		t.button.Icon = t.Icon1
	} else {
		t.State = TwoState0
		t.button.Icon = t.Icon0
	}
	if t.OnActivated != nil {
		t.OnActivated(t.State)
	}
	t.button.Refresh()
}
