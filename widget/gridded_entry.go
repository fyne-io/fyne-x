package widget

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// GriddedEntry shows a series of entries where each one only
// accepts one value. You can optionally place a separator
// between all or select fields.
type GriddedEntry struct {
	widget.DisableableWidget

	// A function to run when the grid is full
	OnCompletion func(string)
	// Text to use to separate fields. It will be
	// forced to monospace to ensure consistent
	// sizing.
	Separator *canvas.Text
	// Optionally specify which fields should have
	// a separator after them.
	SeparateAt []int

	allowedRunes []rune
	length       int

	box      *canvas.Rectangle
	entries  []*SingleRuneEntry
	selected int
}

// NewGriddedEntry returns a new GriddedEntry of the given length allowing
// the given runes.
func NewGriddedEntry(allowedRunes []rune, length int) *GriddedEntry {
	// Initialize the input
	input := &GriddedEntry{
		allowedRunes: allowedRunes,
		length:       length,
		box:          canvas.NewRectangle(theme.HoverColor()),
		entries:      make([]*SingleRuneEntry, length),
		selected:     0,
	}
	input.ExtendBaseWidget(input)
	return input
}

// GetValue returns the value of all the entries in the grid combined.
func (g *GriddedEntry) GetValue() string {
	var out string
	for _, e := range g.entries {
		out += strings.TrimSpace(e.Text)
	}
	return out
}

// FocusIndex will focus the entry at the given index.
func (g *GriddedEntry) FocusIndex(idx int) {
	fyne.CurrentApp().Driver().CanvasForObject(g.entries[idx]).Focus(g.entries[idx])
	g.selected = idx
}

// Disbale disables all entries in the grid.
func (g *GriddedEntry) Disable() {
	for _, e := range g.entries {
		e.Disable()
	}
}

// Enable enables all entries in the grid.
func (g *GriddedEntry) Enable() {
	for _, e := range g.entries {
		e.Enable()
	}
}

// CreateRenderer will initialize the entries and create a renderer
// for them.
func (g *GriddedEntry) CreateRenderer() fyne.WidgetRenderer {
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

	return &griddedEntryRenderer{input: g}
}

func (g *GriddedEntry) objects() []fyne.CanvasObject {
	out := make([]fyne.CanvasObject, len(g.entries)+1)
	out[0] = g.box
	for i, e := range g.entries {
		out[i+1] = e
	}
	return out
}

type griddedEntryRenderer struct {
	input      *GriddedEntry
	separators []fyne.CanvasObject
}

func (g *griddedEntryRenderer) Layout(size fyne.Size) {
	// Grab the size of an entry
	entryMinSize := g.input.entries[0].MinSize()

	// Size and position the border
	boxSize := size.Subtract(fyne.NewSize(theme.InputBorderSize(), 0))
	boxSize.Height = entryMinSize.Height + theme.Padding()
	g.input.box.Resize(boxSize)
	boxPos := fyne.NewPos(0, theme.InputBorderSize()*2)
	g.input.box.Move(boxPos)

	// Size and position the fields
	fieldSize := fyne.NewSize(boxSize.Width/float32(len(g.input.entries))-theme.InputBorderSize()*2.25, entryMinSize.Height)

	if g.input.Separator != nil {
		sepSize := g.input.Separator.Size()
		fieldSize = fieldSize.Subtract(fyne.NewSize(sepSize.Width+theme.Padding()*1.8, 0))
		g.separators = make([]fyne.CanvasObject, g.input.length-1)
		for i := 0; i < g.input.length-1; i++ {
			if g.input.SeparateAt != nil && !intSliceContains(g.input.SeparateAt, i) {
				g.input.box.Resize(g.input.box.Size().Subtract(fyne.NewSize(sepSize.Width+theme.Padding()*2, 0)))
				continue
			}
			g.separators[i] = canvas.NewText(g.input.Separator.Text, g.input.Separator.Color)
		}
	}

	xpos := boxPos.X + 4
	for i, entry := range g.input.entries {
		if i != 0 && g.input.Separator != nil {
			sep := g.separators[i-1]
			if sep != nil {
				sep.(*canvas.Text).TextSize = fieldSize.Width / 2
				sep.(*canvas.Text).TextStyle.Monospace = true
				sep.Move(fyne.NewPos(xpos-sep.Size().Width*2-theme.Padding()/2, fieldSize.Height/2-theme.Padding()))
				xpos = xpos + sep.Size().Width + theme.Padding()*2
			}
		}
		entry.Move(fyne.NewPos(xpos, boxPos.Y+2))
		entry.Resize(fieldSize)
		xpos = xpos + fieldSize.Width + theme.Padding()
	}

	g.Refresh()
}

func (g *griddedEntryRenderer) MinSize() fyne.Size {
	numEntries := float32(len(g.input.entries))
	fieldsSize := fyne.NewSize(numEntries*1.5, theme.TextSize()*numEntries/2)
	min := g.input.box.MinSize().Add(fieldsSize)
	return min
}

func (g *griddedEntryRenderer) Objects() []fyne.CanvasObject {
	return append(g.input.objects(), g.separators...)
}

func (g *griddedEntryRenderer) Refresh() {
	g.input.box.Refresh()
	for _, entry := range g.input.entries {
		entry.Refresh()
	}
}

func (g *griddedEntryRenderer) Destroy() {}

func intSliceContains(ii []int, i int) bool {
	for _, x := range ii {
		if x == i {
			return true
		}
	}
	return false
}
