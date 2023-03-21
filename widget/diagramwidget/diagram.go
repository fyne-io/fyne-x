package diagramwidget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var Globaldiagram *DiagramWidget

func ForceRefresh() {
	Globaldiagram.DummyBox.Refresh()
}

type DiagramWidget struct {
	widget.BaseWidget

	// Diagrams may want to use a different theme and variant than the application. The default value is the
	// applicaton's theme and variant
	DiagramTheme fyne.Theme
	ThemeVariant fyne.ThemeVariant
	Offset       fyne.Position

	// DesiredSize specifies the size which the graph widget should take
	// up, defaults to 800 x 600
	DesiredSize fyne.Size

	Nodes    map[string]*DiagramNode
	Links    map[string]*DiagramLink
	DummyBox *canvas.Rectangle
}

func NewDiagramWidget() *DiagramWidget {
	d := &DiagramWidget{
		DiagramTheme: fyne.CurrentApp().Settings().Theme(),
		ThemeVariant: fyne.CurrentApp().Settings().ThemeVariant(),
		DesiredSize:  fyne.Size{Width: 800, Height: 600},
		Offset:       fyne.Position{X: 0, Y: 0},
		Nodes:        map[string]*DiagramNode{},
		Links:        map[string]*DiagramLink{},
		DummyBox:     canvas.NewRectangle(color.Transparent),
	}
	d.DummyBox.SetMinSize(fyne.Size{Height: 50, Width: 50})
	d.DummyBox.Move(fyne.Position{X: 50, Y: 50})

	d.ExtendBaseWidget(d)

	return d
}

func (dw *DiagramWidget) CreateRenderer() fyne.WidgetRenderer {
	r := diagramWidgetRenderer{
		diagramWidget: dw,
	}

	return &r
}

func (dw *DiagramWidget) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (dw *DiagramWidget) DragEnd() {
	dw.Refresh()
}

func (dw *DiagramWidget) Dragged(event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	for _, n := range dw.Nodes {
		n.Displace(delta)
	}
	dw.Refresh()
}

func (dw *DiagramWidget) MouseIn(event *desktop.MouseEvent) {
}

func (dw *DiagramWidget) MouseOut() {
}

func (dw *DiagramWidget) MouseMoved(event *desktop.MouseEvent) {
}

func (d *DiagramWidget) GetEdges(n *DiagramNode) []*DiagramLink {
	links := []*DiagramLink{}

	for _, link := range d.Links {
		if link.Origin == n {
			links = append(links, link)
		} else if link.Target == n {
			links = append(links, link)
		}
	}

	return links
}

type diagramWidgetRenderer struct {
	diagramWidget *DiagramWidget
}

func (r *diagramWidgetRenderer) MinSize() fyne.Size {
	return r.diagramWidget.DesiredSize
}

func (r *diagramWidgetRenderer) Layout(size fyne.Size) {
	// r.diagramWidget.at.Move(fyne.Position{X: 100, Y: 100})
}

func (r *diagramWidgetRenderer) ApplyTheme(size fyne.Size) {
}

func (r *diagramWidgetRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *diagramWidgetRenderer) Destroy() {
}

func (r *diagramWidgetRenderer) Objects() []fyne.CanvasObject {
	obj := make([]fyne.CanvasObject, 0)
	for _, n := range r.diagramWidget.Nodes {
		obj = append(obj, n)
	}
	for _, e := range r.diagramWidget.Links {
		obj = append(obj, e)
	}
	obj = append(obj, r.diagramWidget.DummyBox)
	return obj
}

func (r *diagramWidgetRenderer) Refresh() {
	for _, e := range r.diagramWidget.Links {
		e.Refresh()
	}
	for _, n := range r.diagramWidget.Nodes {
		n.Refresh()
	}
}
