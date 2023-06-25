package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type HBoxRatioStruct struct {
	WidgetSizes []float32 // size in percentage of total width
}

// top and bottom gaps
func (*HBoxRatioStruct) HeightGaps() float32 {
	return float32(10)
}

// total width and height of all widgets/objects
func (d *HBoxRatioStruct) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()

		// total width of all childs
		w += childSize.Width

		// get longest child height
		if h < childSize.Height {
			h = childSize.Height
		}
	}
	return fyne.NewSize(w, h+d.HeightGaps())
}

// layout all the child(objects) by give containrSize horizontally
// align center vertically
func (d *HBoxRatioStruct) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	var newSize fyne.Size

	cont_width := containerSize.Width - (theme.Padding() * 2)
	pos := fyne.NewPos(theme.Padding(), d.HeightGaps()/2)
	for i, o := range objects {
		newSize = fyne.NewSize(cont_width*(d.WidgetSizes[i]/float32(100)), o.MinSize().Height)
		if newSize.Width < o.MinSize().Width {
			newSize.Width = o.MinSize().Width
		}
		o.Resize(newSize)
		o.Move(pos)
		pos = pos.Add(fyne.NewPos(newSize.Width, 0))
	}
}

func NewHBoxRatioLayout(widths []float32) *HBoxRatioStruct {
	t := new(HBoxRatioStruct)
	t.WidgetSizes = widths
	return t
}
