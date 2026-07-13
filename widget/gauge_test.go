package widget

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestGauge_Defaults(t *testing.T) {
	test.NewTempApp(t)

	g := NewGauge()
	g.SetMinSize(fyne.NewSize(150, 150))
	r := test.WidgetRenderer(g).(*gaugeRenderer)

	assert.Equal(t, 0.0, g.Min)
	assert.Equal(t, 100.0, g.Max)
	assert.Len(t, r.ticks, 11)
	assert.Equal(t, "0", r.readout.Text)
	assert.Equal(t, "0", r.labels[0].Text)
	assert.Equal(t, "20", r.labels[2].Text)
	assert.Equal(t, "100", r.labels[10].Text)
	assert.Nil(t, r.labels[1], "minor ticks have no label")
	assert.True(t, g.MinSize().Width >= 150)
}

func TestGauge_ZeroValueStruct(t *testing.T) {
	test.NewTempApp(t)

	g := &Gauge{}
	r := test.WidgetRenderer(g).(*gaugeRenderer)

	assert.Equal(t, 100.0, g.Max)
	assert.Len(t, r.ticks, 11)
}

func TestGauge_SetValue(t *testing.T) {
	test.NewTempApp(t)

	g := NewGauge()
	r := test.WidgetRenderer(g).(*gaugeRenderer)
	r.Layout(fyne.NewSize(200, 200))

	g.SetValue(50)

	assert.Equal(t, "50", r.readout.Text)
	// at half range the needle points straight up
	assert.InDelta(t, 100, r.needle.Position1.X, 0.001)
	assert.InDelta(t, 100, r.needle.Position2.X, 0.001)
	assert.Less(t, r.needle.Position2.Y, r.needle.Position1.Y)
}

func TestGauge_NeedleClamped(t *testing.T) {
	test.NewTempApp(t)

	g := NewGauge()
	r := test.WidgetRenderer(g).(*gaugeRenderer)
	r.Layout(fyne.NewSize(200, 200))

	g.SetValue(100)
	atMax1, atMax2 := r.needle.Position1, r.needle.Position2

	g.SetValue(9000)
	assert.Equal(t, atMax1, r.needle.Position1)
	assert.Equal(t, atMax2, r.needle.Position2)
	assert.Equal(t, "9000", r.readout.Text, "readout shows the value unclamped")

	g.SetValue(-10)
	atMin1 := r.needle.Position1
	g.SetValue(0)
	assert.Equal(t, atMin1, r.needle.Position1)
}

func TestGauge_TextFormatter(t *testing.T) {
	test.NewTempApp(t)

	g := NewGauge()
	g.TextFormatter = func() string {
		return fmt.Sprintf("%.2f bar", g.Value)
	}
	r := test.WidgetRenderer(g).(*gaugeRenderer)

	g.SetValue(1.5)
	assert.Equal(t, "1.50 bar", r.readout.Text)
}

func TestGauge_ChangeRangeAndSteps(t *testing.T) {
	test.NewTempApp(t)

	g := NewGauge()
	r := test.WidgetRenderer(g).(*gaugeRenderer)

	g.Title = "RPM"
	g.Max = 8000
	g.Steps = 4
	g.Refresh()

	assert.Equal(t, "RPM", r.title.Text)
	assert.Len(t, r.ticks, 5)
	assert.Equal(t, "4000", r.labels[2].Text)
	assert.Equal(t, "8000", r.labels[4].Text)
	assert.Len(t, r.Objects(), 5+3+5) // ticks + labels + face, title, center, needle, readout
}

func TestGauge_Layout(t *testing.T) {
	test.NewTempApp(t)

	g := NewGauge()
	r := test.WidgetRenderer(g).(*gaugeRenderer)
	r.Layout(fyne.NewSize(200, 200))

	assert.Equal(t, fyne.NewSquareSize(200), r.face.Size())
	// needle at Min points down-left from the centre offset
	g.SetValue(0)
	assert.Less(t, r.needle.Position2.X, float32(100))
	assert.Greater(t, r.needle.Position2.Y, float32(100))
}
