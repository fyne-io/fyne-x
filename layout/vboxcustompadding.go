package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Layout = (*VBoxCustomPadding)(nil)

// VBoxCustomPadding is a layout based on Fyne's VBox layout that
// allows for the amount of padding between items to be customized.
// The padding between items is theme.Padding() + ExtraPad (which may be negative).
// If DisableThemePad is set, the padding is just ExtraPad. This is
// useful for placing items that directly touch each other (with ExtraPad=0).
type VBoxCustomPadding struct {
	DisableThemePad bool
	ExtraPad        float32
}

func (*VBoxCustomPadding) isSpacer(obj fyne.CanvasObject) bool {
	if !obj.Visible() {
		return false
	}
	if spacer, ok := obj.(layout.SpacerObject); ok {
		return spacer.ExpandVertical()
	}

	return false
}

func (v *VBoxCustomPadding) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false
	padding := v.themePad() + v.ExtraPad
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if v.isSpacer(child) {
			continue
		}

		minSize.Width = fyne.Max(child.MinSize().Width, minSize.Width)
		minSize.Height += child.MinSize().Height
		if addPadding {
			minSize.Height += padding
		}
		addPadding = true
	}
	return minSize
}

func (v *VBoxCustomPadding) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := 0
	total := float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if v.isSpacer(child) {
			spacers++
			continue
		}
		total += child.MinSize().Height
	}

	y := float32(0)
	padding := v.themePad() + v.ExtraPad
	extra := size.Height - total - (padding * float32(len(objects)-spacers-1))
	extraCell := float32(0)
	if spacers > 0 {
		extraCell = extra / float32(spacers)
	}
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if v.isSpacer(child) {
			y += extraCell
		}
		height := child.MinSize().Height
		child.Move(fyne.NewPos(0, y))
		child.Resize(fyne.NewSize(size.Width, height))
		y += padding + height
	}
}

func (v *VBoxCustomPadding) themePad() float32 {
	if v.DisableThemePad {
		return 0
	}
	return theme.Padding()
}
