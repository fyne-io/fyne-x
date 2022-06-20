package widget

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestNewMap_WithDefaults(t *testing.T) {
	// arrange
	w := test.NewApp().NewWindow("TestMap")
	m := NewMap()
	// action
	w.SetContent(m)
	// verify
	equals(t, "https://tile.openstreetmap.org/%d/%d/%d.png", m.tileSource)
	equals(t, "OpenStreetMap", m.disclaimerLabel)
	equals(t, "https://openstreetmap.org", m.disclaimerUrl)
	equals(t, false, m.hideAttribution)
	equals(t, false, m.hideMoveButtons)
	equals(t, false, m.hideZoomButtons)
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
	equals(t, "test", m.disclaimerLabel)
	equals(t, "http://test.org", m.disclaimerUrl)
	equals(t, false, m.hideAttribution)
	equals(t, true, m.hideMoveButtons)
	equals(t, true, m.hideZoomButtons)
}

func equals(t *testing.T, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}
