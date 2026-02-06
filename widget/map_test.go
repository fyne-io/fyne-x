package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestMap_Pan(t *testing.T) {
	m := NewMap()
	m.Resize(fyne.NewSize(200, 200))
	m.Zoom(3)
	assert.Equal(t, 0, m.x)
	assert.Equal(t, 0, m.y)

	m.PanSouth()
	m.PanEast()
	assert.Equal(t, 1, m.x)
	assert.Equal(t, 1, m.y)

	m.PanNorth()
	m.PanWest()
	assert.Equal(t, 0, m.x)
	assert.Equal(t, 0, m.y)
}

func TestMap_Drag(t *testing.T) {
	m := NewMap()
	m.Resize(fyne.NewSize(200, 200))
	m.Zoom(3)
	assert.Equal(t, 0, m.x)
	assert.Equal(t, 0, m.y)

	m.Dragged(&fyne.DragEvent{Dragged: fyne.Delta{
		DX: 300,
		DY: 72,
	}})
	assert.Equal(t, float32(-212), m.offsetX)
	assert.Equal(t, -2, m.x)
	assert.Equal(t, float32(72), m.offsetY)
	assert.Equal(t, 0, m.y)

	m.Dragged(&fyne.DragEvent{Dragged: fyne.Delta{
		DX: -564,
		DY: 0,
	}})
	assert.Equal(t, float32(-264), m.offsetX)
	assert.Equal(t, 0, m.x)
	assert.Equal(t, float32(72), m.offsetY)
	assert.Equal(t, 0, m.y)
}

func TestMap_Zoom(t *testing.T) {
	m := NewMap()
	m.Resize(fyne.NewSize(200, 200))
	assert.Equal(t, 0, m.zoom)
	m.ZoomIn()
	assert.Equal(t, 1, m.zoom)
	m.ZoomOut()
	assert.Equal(t, 0, m.zoom)

	m.Zoom(5)
	assert.Equal(t, 5, m.zoom)
	m.Zoom(55) // invalid
	assert.Equal(t, 5, m.zoom)
}

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
