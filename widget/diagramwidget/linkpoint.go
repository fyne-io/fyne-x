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

// NewLinkPoint creates an instance of a LinkPoint for a specific DiagramLink
func NewLinkPoint(link DiagramLink) *LinkPoint {
	lp := &LinkPoint{}
	lp.BaseWidget.ExtendBaseWidget(lp)
	lp.link = link
	return lp
}

// CreateRenderer creates the renderere for a LinkPoint
func (lp *LinkPoint) CreateRenderer() fyne.WidgetRenderer {
	lpr := &linkPointRenderer{}
	return lpr
}

// GetLink returns the Link to which the LinkPoint belongs
func (lp *LinkPoint) GetLink() DiagramLink {
	return lp.link
}

// IsConnectionAllowed returns true if a connection is permitted with the indicated pad. The
// question is passed to the owning link
func (lp *LinkPoint) IsConnectionAllowed(connectionPad ConnectionPad) bool {
	return lp.link.isConnectionAllowed(lp, connectionPad)
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
