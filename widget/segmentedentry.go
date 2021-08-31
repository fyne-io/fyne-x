package widget

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Focusable = (*SegmentedEntry)(nil)

// SegmentedEntry shows a series of entries where each one only
// accepts one value. You can optionally place a separator
// between all or select fields.
type SegmentedEntry struct {
	widget.DisableableWidget

	// A function to run when the grid is full
	OnCompletion func(string)
	// Text to use to separate fields. It will be
	// forced to monospace to ensure consistent
	// sizing.
	Delimiter string
	// Optionally specify which fields should have
	// a separator after them.
	DelimitAt []int

	allowedRunes []rune
	length       int

	entries  []*SingleRuneEntry
	selected int
}

// NewSegmentedEntry returns a new SegmentedEntry of the given length allowing
// the given runes.
func NewSegmentedEntry(allowedRunes []rune, length int) *SegmentedEntry {
	// Initialize the input
	input := &SegmentedEntry{
		allowedRunes: allowedRunes,
		length:       length,
		entries:      make([]*SingleRuneEntry, length),
		selected:     0,
	}
	input.ExtendBaseWidget(input)
	return input
}

// GetValue returns the value of all the entries in the grid combined.
func (g *SegmentedEntry) GetValue() string {
	var out string
	for _, e := range g.entries {
		out += strings.TrimSpace(e.Text)
	}
	return out
}

// FocusGained implements fyne.Focusable
func (g *SegmentedEntry) FocusGained() {
}

// FocusLost implements fyne.Focusable
func (g *SegmentedEntry) FocusLost() {
}

// FocusIndex will focus the entry at the given index.
func (g *SegmentedEntry) FocusIndex(idx int) {
	fyne.CurrentApp().Driver().CanvasForObject(g.entries[idx]).Focus(g.entries[idx])
	g.selected = idx
}

// TypedKey forwards the event to the focussed entry
func (g *SegmentedEntry) TypedKey(key *fyne.KeyEvent) {
	g.entries[g.selected].TypedKey(key)
}

// TypedRune forwards the rune to the focussed entry
func (g *SegmentedEntry) TypedRune(r rune) {
	g.entries[g.selected].TypedRune(r)
}

// Disbale disables all entries in the grid.
func (g *SegmentedEntry) Disable() {
	for _, e := range g.entries {
		e.Disable()
	}
}

// Enable enables all entries in the grid.
func (g *SegmentedEntry) Enable() {
	for _, e := range g.entries {
		e.Enable()
	}
}

// CreateRenderer will initialize the entries and create a renderer
// for them.
func (g *SegmentedEntry) CreateRenderer() fyne.WidgetRenderer {
	g.ExtendBaseWidget(g)

	// Set up the entries
	for i := 0; i < g.length; i++ {
		// Create the entry
		g.entries[i] = NewSingleRuneEntry(g.allowedRunes)

		// Take a reference to this index for the closures
		this := i

		// Navigation for this entry
		g.entries[i].OnNavigate = func(key *fyne.KeyEvent) {
			// This was a regular input
			if key == nil {
				if len(g.GetValue()) == g.length && g.OnCompletion != nil {
					g.OnCompletion(g.GetValue())
					return
				}
				if this != g.length-1 {
					g.FocusIndex(this + 1)
				}
				return
			}
			switch key.Name {
			case fyne.KeyBackspace:
				if this == 0 {
					return
				}
				if g.entries[this-1].Text != "" {
					g.entries[this-1].SetText("")
					g.FocusIndex(this - 1)
				}
			case fyne.KeyLeft:
				if this == 0 {
					return
				}
				g.FocusIndex(this - 1)
			case fyne.KeyRight:
				if this == g.length-1 {
					return
				}
				g.FocusIndex(this + 1)
			}
		}
	}

	return &segmentedEntryRenderer{input: g}
}

func (g *SegmentedEntry) objects() []fyne.CanvasObject {
	out := make([]fyne.CanvasObject, len(g.entries))
	for i, e := range g.entries {
		out[i] = e
	}
	return out
}

type segmentedEntryRenderer struct {
	input      *SegmentedEntry
	separators []fyne.CanvasObject
}

// Layout implements fyne.WidgetRenderer and produces the layout for the widget.
func (g *segmentedEntryRenderer) Layout(size fyne.Size) {
	basePos := fyne.NewPos(0, 0)

	fieldSize := widget.NewEntry().MinSize()
	fieldSize = fyne.NewSize(fieldSize.Width-theme.Padding(), fieldSize.Height)

	if g.input.Delimiter != "" {
		style := fyne.TextStyle{Monospace: true}
		align := fyne.TextAlignCenter
		sepSize := widget.NewLabelWithStyle(g.input.Delimiter, align, style).Size()

		fieldSize = fieldSize.Subtract(fyne.NewSize(sepSize.Width+theme.Padding()*1.8, 0))
		g.separators = make([]fyne.CanvasObject, g.input.length-1)
		for i := 0; i < g.input.length-1; i++ {
			if g.input.DelimitAt != nil && !intSliceContains(g.input.DelimitAt, i) {
				continue
			}
			g.separators[i] = widget.NewLabelWithStyle(g.input.Delimiter, align, style)
		}
	}

	xpos := basePos.X
	for i, entry := range g.input.entries {
		if i != 0 && g.input.Delimiter != "" {
			sep := g.separators[i-1]
			if sep != nil {
				sep.Move(fyne.NewPos(xpos+fieldSize.Width/2, 0))
				xpos = sep.Position().X + fieldSize.Width/2
			}
		}
		entry.Move(fyne.NewPos(xpos, basePos.Y))
		entry.Resize(fieldSize)
		xpos = xpos + fieldSize.Width + theme.Padding()
	}

	g.Refresh()
}

// MinSize implements fyne.WidgetRenderer and returns the minimum size for the layout.
func (g *segmentedEntryRenderer) MinSize() fyne.Size {
	numEntries := float32(len(g.input.entries))
	fieldsSizes := fyne.NewSize(numEntries*theme.InputBorderSize(), theme.TextSize()+theme.InputBorderSize())
	return fieldsSizes
}

// Objects implements fyne.WidgetRenderer and returns the objects to render.
func (g *segmentedEntryRenderer) Objects() []fyne.CanvasObject {
	return append(g.input.objects(), g.separators...)
}

// Refresh implements fyne.WidgetRenderer and refreshes all the entries.
func (g *segmentedEntryRenderer) Refresh() {
	for _, entry := range g.input.entries {
		entry.Refresh()
	}
}

// Destroy implements fyne.WidgetRenderer.
func (g *segmentedEntryRenderer) Destroy() {}

func intSliceContains(ii []int, i int) bool {
	for _, x := range ii {
		if x == i {
			return true
		}
	}
	return false
}
