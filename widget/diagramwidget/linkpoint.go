package diagramwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// LinkPoint identifies the point at which a link end is connected to another diagram element's connection pad
type LinkPoint struct {
	widget.BaseWidget
}

func NewLinkPoint(link *BaseDiagramLink) *LinkPoint {
	lp := &LinkPoint{}
	lp.BaseWidget.ExtendBaseWidget(lp)
	return lp
}

func (lp *LinkPoint) CreateRenderer() fyne.WidgetRenderer {
	lpr := &linkPointRenderer{}
	return lpr
}

// linkPointRenderer
type linkPointRenderer struct {
}

func (lpr *linkPointRenderer) Destroy() {

}

func (lpr *linkPointRenderer) Layout(size fyne.Size) {

}

func (lpr *linkPointRenderer) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

func (lp *linkPointRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{}
	return obj
}

func (lp *linkPointRenderer) Refresh() {

}
