package widget

import (
	"testing"

	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestNewMap_WithDefaults(t *testing.T) {
	// arrange
	w := test.NewApp().NewWindow("TestMap")
	m := NewMap()
	// action
	w.SetContent(m)
	// verify
	assert.Equal(t, "https://tile.openstreetmap.org/%d/%d/%d.png", m.tileSource)
	assert.Equal(t, "OpenStreetMap", m.attributionLabel)
	assert.Equal(t, "https://openstreetmap.org", m.attributionURL)
	assert.False(t, m.hideAttribution)
	assert.False(t, m.hideMoveButtons)
	assert.False(t, m.hideZoomButtons)
}

func TestNewMap_WithOptions(t *testing.T) {
	// arrange
	w := test.NewApp().NewWindow("TestMap")
	m := NewMapWithOptions(
		WithScrollButtons(false),
		WithZoomButtons(false),
		WithAttribution(true, "test", "http://test.org"),
	)
	// action
	w.SetContent(m)
	// verify
	assert.Equal(t, "test", m.attributionLabel)
	assert.Equal(t, "http://test.org", m.attributionURL)
	assert.False(t, m.hideAttribution)
	assert.True(t, m.hideMoveButtons)
	assert.True(t, m.hideZoomButtons)
}
