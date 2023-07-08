package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Layout = (*HBoxCustomPadding)(nil)

// HBoxCustomPadding is a layout based on Fyne's HBox layout that
// allows for the amount of padding between items to be customized.
// The padding between items is theme.Padding() + ExtraPad (which may be negative).
// If DisableThemePad is set, the padding is just ExtraPad. This is
// useful for placing items that directly touch each other (with ExtraPad=0).
type HBoxCustomPadding struct {
	DisableThemePad bool
	ExtraPad        float32
}

func (*HBoxCustomPadding) isSpacer(obj fyne.CanvasObject) bool {
	if !obj.Visible() {
		return false
	}
	if spacer, ok := obj.(layout.SpacerObject); ok {
		return spacer.ExpandHorizontal()
	}

	return false
}

func (h *HBoxCustomPadding) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false
	padding := h.themePad() + h.ExtraPad
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if h.isSpacer(child) {
			continue
		}

		minSize.Height = fyne.Max(child.MinSize().Height, minSize.Height)
		minSize.Width += child.MinSize().Width
		if addPadding {
			minSize.Width += padding
		}
		addPadding = true
	}
	return minSize
}

func (h *HBoxCustomPadding) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := 0
	total := float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if h.isSpacer(child) {
			spacers++
			continue
		}
		total += child.MinSize().Width
	}

	x := float32(0)
	padding := h.themePad() + h.ExtraPad
	extra := size.Width - total - (padding * float32(len(objects)-spacers-1))
	extraCell := float32(0)
	if spacers > 0 {
		extraCell = extra / float32(spacers)
	}
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if h.isSpacer(child) {
			x += extraCell
		}
		width := child.MinSize().Width
		child.Move(fyne.NewPos(x, 0))
		child.Resize(fyne.NewSize(width, size.Height))
		x += padding + width
	}
}

func (h *HBoxCustomPadding) themePad() float32 {
	if h.DisableThemePad {
		return 0
	}
	return theme.Padding()
}
