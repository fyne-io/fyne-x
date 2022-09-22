package binding

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type disableableBinding struct {
	binding.Bool

	inverted bool
	widgets  []fyne.Disableable
}

// NewDisableableBinding returns a `Bool` binding which accepts Disableable widgets.
// When the Bool changes, the widgets Enable or Disable method will be executed.
func NewDisableableBinding(widgets ...fyne.Disableable) *disableableBinding {
	newBinding := &disableableBinding{
		Bool:    binding.NewBool(),
		widgets: widgets,
	}

	// Add default listener
	newBinding.AddListener(binding.NewDataListener(newBinding.update))

	return newBinding
}

// Adding widgets to the binding.
// This will update the Disable status of the widgets immediately.
func (b *disableableBinding) AddWidgets(widgets ...fyne.Disableable) {
	b.widgets = append(b.widgets, widgets...)
	b.update()
}

// SetInverted will switch the behavior of when the widgets will be Enabled or Disabled.
// This will update the Disable status of the widgets immediately.
func (b *disableableBinding) SetInverted(inverted bool) {
	b.inverted = inverted
	b.update()
}

func (b *disableableBinding) update() {
	val, err := b.Get()
	if err != nil {
		log.Println(err)
		return
	}

	if (!b.inverted && val) || (b.inverted && !val) {
		for _, widget := range b.widgets {
			widget.Enable()
		}
	} else {
		for _, widget := range b.widgets {
			widget.Disable()
		}
	}
}
