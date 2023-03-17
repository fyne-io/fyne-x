package diagramwidget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type diagramRenderer struct {
	graph *DiagramWidget
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

	Nodes map[string]*DiagramNode
	Links map[string]*DiagramLink
}

func (r *diagramRenderer) MinSize() fyne.Size {
	return r.graph.DesiredSize
}

func (r *diagramRenderer) Layout(size fyne.Size) {
}

func (r *diagramRenderer) ApplyTheme(size fyne.Size) {
}

func (r *diagramRenderer) Refresh() {
	for _, e := range r.graph.Links {
		e.Refresh()
	}
	for _, n := range r.graph.Nodes {
		n.Refresh()
	}
}

func (r *diagramRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *diagramRenderer) Destroy() {
}

func (r *diagramRenderer) Objects() []fyne.CanvasObject {
	obj := make([]fyne.CanvasObject, 0)
	for _, n := range r.graph.Nodes {
		obj = append(obj, n)
	}
	for _, e := range r.graph.Links {
		obj = append(obj, e)
		for _, sourceDecoration := range e.SourceDecorations {
			if sourceDecoration != nil {
				obj = append(obj, sourceDecoration)
			}
		}
		for _, midpointDecoration := range e.MidpointDecorations {
			if midpointDecoration != nil {
				obj = append(obj, midpointDecoration)
			}
		}
		for _, targetDecoration := range e.TargetDecorations {
			if targetDecoration != nil {
				obj = append(obj, targetDecoration)
			}
		}
	}

	return obj
}

func (g *DiagramWidget) CreateRenderer() fyne.WidgetRenderer {
	r := diagramRenderer{
		graph: g,
	}

	return &r
}

func (g *DiagramWidget) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (g *DiagramWidget) DragEnd() {
	g.Refresh()
}

func (g *DiagramWidget) Dragged(event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	for _, n := range g.Nodes {
		n.Displace(delta)
	}
	g.Refresh()
}

func (g *DiagramWidget) MouseIn(event *desktop.MouseEvent) {
}

func (g *DiagramWidget) MouseOut() {
}

func (g *DiagramWidget) MouseMoved(event *desktop.MouseEvent) {
}

func NewDiagram() *DiagramWidget {
	d := &DiagramWidget{
		DiagramTheme: fyne.CurrentApp().Settings().Theme(),
		ThemeVariant: fyne.CurrentApp().Settings().ThemeVariant(),
		DesiredSize:  fyne.Size{Width: 800, Height: 600},
		Offset:       fyne.Position{X: 0, Y: 0},
		Nodes:        map[string]*DiagramNode{},
		Links:        map[string]*DiagramLink{},
	}

	d.ExtendBaseWidget(d)

	return d
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
