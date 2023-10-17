package layout

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Layout = (*HPortion)(nil)

// HPortion allows the canvas objects to be divided into portions of the width.
// The length of the Portions slice needs to be equal to the amount of canvas objects.
type HPortion struct {
	Portions []float64
}

// Layout sets the size and position of the canvas objects.
func (p *HPortion) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(p.Portions) != len(objects) {
		log.Println("Mismatch between partitions and objects")
		return
	}

	sum := float64(0)
	for _, child := range p.Portions {
		sum += float64(child)
	}

	padding := theme.Padding()
	width := size.Width - padding*float32(len(objects)-1)
	xpos := float32(0)

	for i, child := range objects {
		width := float32(p.Portions[i]/sum) * width
		child.Resize(fyne.NewSize(width, size.Height))
		child.Move(fyne.NewPos(xpos, 0))

		xpos += width + padding
	}
}

// MinSize calculates the minimum required size to fit all objects.
// It is equal to the largest width MinSize divided by the corresponding portion.
func (p *HPortion) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(p.Portions) != len(objects) {
		log.Println("Mismatch between partitions and objects")
		return fyne.NewSize(0, 0)
	}

	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}

	sum := float64(0)
	for _, child := range p.Portions {
		sum += float64(child)
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
	return fyne.NewSize(maxMinWidth/float32(p.Portions[maxIndex]/sum)+totalPadding, height)
}

// NewHPortion creates a layout that partitions objects horizontally taking up
// as large of a portion of the space as defined by the given slice.
// The length of the Portions slice needs to be equal to the amount of objects.
func NewHPortion(Portions []float64) *HPortion {
	return &HPortion{Portions: Portions}
}

var _ fyne.Layout = (*VPortion)(nil)

// VPortion allows the canvas objects to be divided into portions of the height.
// The length of the Portions slice needs to be equal to the amount of canvas objects.
type VPortion struct {
	Portions []float64
}

// Layout sets the size and position of the canvas objects.
func (p *VPortion) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(p.Portions) != len(objects) {
		log.Println("Mismatch between partitions and objects")
		return
	}

	sum := float64(0)
	for _, child := range p.Portions {
		sum += float64(child)
	}

	padding := theme.Padding()
	height := size.Width - padding*float32(len(objects)-1)
	ypos := float32(0)

	for i, child := range objects {
		height := float32(p.Portions[i]/sum) * height
		child.Resize(fyne.NewSize(ypos, height))
		child.Move(fyne.NewPos(ypos, 0))

		ypos += height + padding
	}
}

// MinSize calculates the minimum required size to fit all objects.
// It is equal to the largest height MinSize divided by the corresponding portion.
func (p *VPortion) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(p.Portions) != len(objects) {
		log.Println("Mismatch between partitions and objects")
		return fyne.NewSize(0, 0)
	}

	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}

	sum := float64(0)
	for _, child := range p.Portions {
		sum += float64(child)
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
	return fyne.NewSize(width, maxMinHeight/float32(p.Portions[maxIndex]/sum)+totalPadding)
}

// NewVPortion creates a layout that partitions objects verticaly taking up
// as large of a portion of the space as defined by the given slice.
// The length of the Portions slice needs to be equal to the amount of objects.
func NewVPortion(portion []float64) *VPortion {
	return &VPortion{Portions: portion}
}
