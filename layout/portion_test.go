package layout

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestHPortion(t *testing.T) {
	cont := container.New(&HPortion{Portions: []float64{50, 50}}, widget.NewEntry(), widget.NewEntry())
	cont.Resize(fyne.NewSize(100, 100))

	assert.Equal(t, (100-theme.Padding())/2, cont.Objects[0].Size().Width)
	assert.Equal(t, (100-theme.Padding())/2, cont.Objects[1].Size().Width)

	assert.Equal(t, cont.Objects[0].MinSize().Height, cont.MinSize().Height)
	assert.Equal(t, cont.Objects[1].MinSize().Height, cont.MinSize().Height)

	// Using 0.5 and 0.5 should be the same as 50 and 50.
	cont.Layout = NewHPortion([]float64{0.5, 0.5})

	assert.Equal(t, (100-theme.Padding())/2, cont.Objects[0].Size().Width)
	assert.Equal(t, (100-theme.Padding())/2, cont.Objects[1].Size().Width)

	assert.Equal(t, cont.Objects[0].MinSize().Height, cont.MinSize().Height)
	assert.Equal(t, cont.Objects[1].MinSize().Height, cont.MinSize().Height)

	// Mismatch in length should error out.
	cont.Layout = NewHPortion([]float64{})
	cont.Resize(fyne.NewSize(50, 50))
	assert.Equal(t, fyne.NewSize(0, 0), cont.MinSize())

	// Having no objects should result in zero size.
	cont.Objects = nil
	cont.Resize(fyne.NewSize(100, 100))
	assert.Equal(t, fyne.NewSize(0, 0), cont.MinSize())
}

func TestVPortion(t *testing.T) {
	cont := container.New(&VPortion{Portions: []float64{50, 50}}, widget.NewEntry(), widget.NewEntry())
	cont.Resize(fyne.NewSize(100, 100))

	assert.Equal(t, (100-theme.Padding())/2, cont.Objects[0].Size().Height)
	assert.Equal(t, (100-theme.Padding())/2, cont.Objects[1].Size().Height)

	assert.Equal(t, cont.Objects[0].MinSize().Width, cont.MinSize().Width)
	assert.Equal(t, cont.Objects[1].MinSize().Width, cont.MinSize().Width)

	// Using 0.5 and 0.5 should be the same as 50 and 50.
	cont.Layout = NewVPortion([]float64{0.5, 0.5})

	assert.Equal(t, (100-theme.Padding())/2, cont.Objects[0].Size().Height)
	assert.Equal(t, (100-theme.Padding())/2, cont.Objects[1].Size().Height)

	assert.Equal(t, cont.Objects[0].MinSize().Width, cont.MinSize().Width)
	assert.Equal(t, cont.Objects[1].MinSize().Width, cont.MinSize().Width)

	// Mismatch in length should error out.
	cont.Layout = NewVPortion([]float64{})
	cont.Resize(fyne.NewSize(50, 50))
	assert.Equal(t, fyne.NewSize(0, 0), cont.MinSize())

	// Having no objects should result in zero size.
	cont.Objects = nil
	cont.Resize(fyne.NewSize(100, 100))
	assert.Equal(t, fyne.NewSize(0, 0), cont.MinSize())
}
