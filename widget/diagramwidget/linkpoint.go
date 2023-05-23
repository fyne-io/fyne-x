package diagramwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// LinkPoint identifies the point at which a link end is connected to another diagram element's connection pad
type LinkPoint struct {
	widget.BaseWidget
	link DiagramLink
}

func NewLinkPoint(link DiagramLink) *LinkPoint {
	lp := &LinkPoint{}
	lp.BaseWidget.ExtendBaseWidget(lp)
	lp.link = link
	return lp
}

func (lp *LinkPoint) CreateRenderer() fyne.WidgetRenderer {
	lpr := &linkPointRenderer{}
	return lpr
}

func (lp *LinkPoint) GetLink() DiagramLink {
	return lp.link
}

func (lp *LinkPoint) IsConnectionAllowed(connectionPad ConnectionPad) bool {
	return lp.link.IsConnectionAllowed(lp, connectionPad)
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
