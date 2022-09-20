package binding

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type disableableBinding struct {
	binding.Bool

	inverted bool
	targets  []fyne.Disableable
}

// NewDisableableBinding returns a `Bool` binding which accepts Disableable targets.
// When the Bool changes, the targets Enable or Disable method will be executed.
func NewDisableableBinding(targets ...fyne.Disableable) *disableableBinding {
	newBinding := &disableableBinding{
		Bool:    binding.NewBool(),
		targets: targets,
	}

	// Add default listener
	newBinding.AddListener(binding.NewDataListener(newBinding.update))

	return newBinding
}

// Adding targets to the binding.
// This will update the Disable status of the targets immediately.
func (b *disableableBinding) AddTargets(targets ...fyne.Disableable) {
	b.targets = append(b.targets, targets...)
	b.update()
}

// Invert will switch the behavior of when the targets will be Enabled or Disabled.
// This will update the Disable status of the targets immediately.
func (b *disableableBinding) Invert(inverted bool) {
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
		for _, target := range b.targets {
			target.Enable()
		}
	} else {
		for _, target := range b.targets {
			target.Disable()
		}
	}
}
