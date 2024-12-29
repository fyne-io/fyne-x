package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// This is essentially a VBox layout, but glued to the bottom instead of
// the top. Intended for mobile devices with thumb control.
func NewBottomUp(objects ...fyne.CanvasObject) *fyne.Container {
	return &fyne.Container{
		Layout:  NewBottomUpLayout(),
		Objects: objects,
	}
}

type bottomUpLayout struct {
	*fyne.Container
}

func NewBottomUpLayout() fyne.Layout {
	return &bottomUpLayout{}
}

func (c *bottomUpLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	var offset float32
	padding := theme.Padding()
	for i := len(objects) - 1; i >= 0; i-- {
		child := objects[i]
		childMin := child.MinSize()

		child.Resize(fyne.NewSize(size.Width, childMin.Height))
		child.Move(fyne.NewPos(0, size.Height-childMin.Height-offset))

		offset += childMin.Height + padding
	}
}

func (c *bottomUpLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	padding := theme.Padding()
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		ms := child.MinSize()
		//minSize.Width = max(minSize.Width, ms.Width) // fixme: swap once go version is >=1.22
		if ms.Width > minSize.Width {
			minSize.Width = ms.Width
		}
		minSize.Height += ms.Height + padding
	}
	minSize.Height -= padding
	if minSize.Height < 0 {
		minSize.Height = 0
	}

	return minSize
}
