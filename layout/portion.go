package layout

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Layout = (*HPortion)(nil)

type HPortion struct {
	Portions []float32
}

func (p *HPortion) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(p.Portions) != len(objects) {
		log.Println("Mismatch between partitions and objects")
		return
	}

	padding := theme.Padding()
	xpos := padding

	for i, child := range objects {
		width := p.Portions[i] * size.Width
		child.Resize(fyne.NewSize(width, size.Height))
		child.Move(fyne.NewPos(xpos, 0))

		xpos += width + padding
	}
}

func (p *HPortion) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(p.Portions) != len(objects) {
		log.Println("Mismatch between partitions and objects")
		return fyne.NewSize(0, 0)
	}

	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}

	maxMinWidth := float32(0)
	maxIndex := -1
	height := float32(0)

	for i := 0; i < len(objects); i++ {
		min := objects[i].MinSize()
		height = fyne.Max(height, min.Height)

		if min.Width > maxMinWidth {
			maxMinWidth = min.Width
			maxIndex = i
		}
	}

	totalPadding := float32(len(objects)-1) * theme.Padding()
	return fyne.NewSize(maxMinWidth/p.Portions[maxIndex]+totalPadding, height)
}

// NewHPortion creates a layout that partitions objects verticaly taking up
// as large of a portion of the space as defined by the given slice.
// The portions should be between 0 and 1 but not equal to.
// The length of the Portions slice needs to be equal to the amount of objects.
func NewHPortion(Portions []float32) *HPortion {
	return &HPortion{Portions: Portions}
}

var _ fyne.Layout = (*VPortion)(nil)

type VPortion struct {
	Portions []float32
}

func (p *VPortion) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(p.Portions) != len(objects) {
		log.Println("Mismatch between partitions and objects")
		return
	}

	padding := theme.Padding()
	ypos := padding

	for i, child := range objects {
		height := p.Portions[i] * size.Height
		child.Resize(fyne.NewSize(ypos, height))
		child.Move(fyne.NewPos(ypos, 0))

		ypos += height + padding
	}
}

func (p *VPortion) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(p.Portions) != len(objects) {
		log.Println("Mismatch between partitions and objects")
		return fyne.NewSize(0, 0)
	}

	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}

	maxMinHeight := float32(0)
	maxIndex := -1
	width := float32(0)

	for i := 0; i < len(objects); i++ {
		min := objects[i].MinSize()
		width = fyne.Max(width, min.Width)

		if min.Height > maxMinHeight {
			maxMinHeight = min.Height
			maxIndex = i
		}
	}

	totalPadding := float32(len(objects)-1) * theme.Padding()
	return fyne.NewSize(width, maxMinHeight/p.Portions[maxIndex]+totalPadding)
}

// NewVPortion creates a layout that partitions objects verticaly taking up
// as large of a portion of the space as defined by the given slice.
// The portions should be between 0 and 1 but not equal to.
// The length of the Portions slice needs to be equal to the amount of objects.
func NewVPortion(portion []float32) *VPortion {
	return &VPortion{Portions: portion}
}
