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

// ForceRepaint is a workaround for a Fyne bug (Issue #2205) in which moving a canvas object does not
// trigger repainting. When the issue is resolved, this function and all references to it should be
// removed. The DummyBox on the GlobalDiagram should also be removed.
// The conditionals here are required during initialization.
func ForceRepaint() {
	if Globaldiagram != nil {
		if Globaldiagram.dummyBox != nil {
			Globaldiagram.dummyBox.Refresh()
		}
	}
}

// Verify that interfaces are fully implemented
var _ fyne.Tappable = (*DiagramWidget)(nil)

type linkPinPair struct {
	link *DiagramLink
	pin  *LinkPoint
}

type DiagramWidget struct {
	widget.BaseWidget

	// ID is expected to be unique across all DiagramWidgets in the application.
	ID string

	// Diagrams may want to use a different theme and variant than the application. The default value is the
	// applicaton's theme and variant
	DiagramTheme fyne.Theme
	ThemeVariant fyne.ThemeVariant
	Offset       fyne.Position

	// DesiredSize specifies the size which the graph widget should take
	// up, defaults to 800 x 600
	DesiredSize fyne.Size

	Nodes                          map[string]*DiagramNode
	Links                          map[string]*DiagramLink
	selection                      map[string]DiagramElement
	diagramElementLinkDependencies map[string][]linkPinPair

	// TODO Remove dummyBox when fyne rendering issue is resolved
	dummyBox *canvas.Rectangle
}

func NewDiagramWidget(id string) *DiagramWidget {
	dw := &DiagramWidget{
		ID:                             id,
		DiagramTheme:                   fyne.CurrentApp().Settings().Theme(),
		ThemeVariant:                   fyne.CurrentApp().Settings().ThemeVariant(),
		DesiredSize:                    fyne.Size{Width: 800, Height: 600},
		Offset:                         fyne.Position{X: 0, Y: 0},
		Nodes:                          map[string]*DiagramNode{},
		Links:                          map[string]*DiagramLink{},
		dummyBox:                       canvas.NewRectangle(color.Transparent),
		selection:                      map[string]DiagramElement{},
		diagramElementLinkDependencies: map[string][]linkPinPair{},
	}
	dw.dummyBox.SetMinSize(fyne.Size{Height: 50, Width: 50})
	dw.dummyBox.Move(fyne.Position{X: 50, Y: 50})

	dw.ExtendBaseWidget(dw)

	return dw
}

func (dw *DiagramWidget) AddLink(link *DiagramLink) {
	dw.Links[link.id] = link
	link.Refresh()
	// TODO add logic to rezise diagram if necessary
}

func (dw *DiagramWidget) addLinkDependency(diagramElement DiagramElement, link *DiagramLink, pin *LinkPoint) {
	deID := diagramElement.GetDiagramElementID()
	currentDependencies := dw.diagramElementLinkDependencies[deID]
	if currentDependencies == nil {
		dw.diagramElementLinkDependencies[deID] = []linkPinPair{{link, pin}}
	} else {
		for _, pair := range currentDependencies {
			if pair.link == link && pair.pin == pin {
				// it's already there
				return
			}
		}
		dw.diagramElementLinkDependencies[deID] = append(currentDependencies, linkPinPair{link, pin})
	}
}

func (dw *DiagramWidget) AddNode(node *DiagramNode) {
	dw.Nodes[node.id] = node
	node.Refresh()
	// TODO add logic to rezise diagram if necessary
}

func (dw *DiagramWidget) CreateRenderer() fyne.WidgetRenderer {
	r := diagramWidgetRenderer{
		diagramWidget: dw,
	}
	return &r
}

func (dw *DiagramWidget) addElementToSelection(de DiagramElement) {
	if !dw.IsSelected(de) {
		dw.selection[de.GetDiagramElementID()] = de
		de.ShowHandles()
	}
}

func (dw *DiagramWidget) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (dw *DiagramWidget) DiagramElementTapped(de DiagramElement, event *fyne.PointEvent) {
	if !dw.IsSelected(de) {
		dw.addElementToSelection(de)
	}
	ForceRepaint()
}

func (dw *DiagramWidget) DragEnd() {
	dw.Refresh()
}

func (dw *DiagramWidget) GetBackgroundColor() color.Color {
	return dw.DiagramTheme.Color(theme.ColorNameBackground, dw.ThemeVariant)
}

func (dw *DiagramWidget) GetForegroundColor() color.Color {
	return dw.DiagramTheme.Color(theme.ColorNameForeground, dw.ThemeVariant)
}

func (dw *DiagramWidget) GetHoverColor() color.Color {
	return dw.DiagramTheme.Color(theme.ColorNameHover, dw.ThemeVariant)
}

func (dw *DiagramWidget) DiagramNodeDragged(node *DiagramNode, event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	dw.DisplaceNode(node, delta)
	ForceRepaint()
}

func (dw *DiagramWidget) DisplaceNode(node *DiagramNode, delta fyne.Position) {
	node.Move(node.Position().Add(delta))
	dw.refreshDependentLinks(node)
	ForceRepaint()
}

func (dw *DiagramWidget) Dragged(event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	for _, n := range dw.Nodes {
		dw.DisplaceNode(n, delta)
	}
	dw.Refresh()
}

func (dw *DiagramWidget) IsSelected(de DiagramElement) bool {
	return dw.selection[de.GetDiagramElementID()] != nil
}

func (dw *DiagramWidget) MouseIn(event *desktop.MouseEvent) {
}

func (dw *DiagramWidget) MouseOut() {
}

func (dw *DiagramWidget) MouseMoved(event *desktop.MouseEvent) {
}

func (dw *DiagramWidget) removeElementFromSelection(de DiagramElement) {
	delete(dw.selection, de.GetDiagramElementID())
	de.HideHandles()
}

func (dw *DiagramWidget) removeLinkDependency(diagramElement DiagramElement, link *DiagramLink, pin *LinkPoint) {
	deID := diagramElement.GetDiagramElementID()
	currentDependencies := dw.diagramElementLinkDependencies[deID]
	if currentDependencies == nil {
		return
	}
	for i, pair := range currentDependencies {
		if pair.link == link && pair.pin == pin {
			dw.diagramElementLinkDependencies[deID] = append(currentDependencies[:i], currentDependencies[i+1:]...)
			return
		}
	}
}

func (dw *DiagramWidget) refreshDependentLinks(de DiagramElement) {
	dependencies := dw.diagramElementLinkDependencies[de.GetDiagramElementID()]
	for _, pair := range dependencies {
		pair.link.Refresh()
	}
}

func (dw *DiagramWidget) Tapped(event *fyne.PointEvent) {
	for _, de := range dw.selection {
		dw.removeElementFromSelection(de)
	}
	ForceRepaint()
}

// diagramWidgetRenderer
type diagramWidgetRenderer struct {
	diagramWidget *DiagramWidget
}

func (r *diagramWidgetRenderer) Destroy() {
}

func (r *diagramWidgetRenderer) Layout(size fyne.Size) {
	// r.diagramWidget.at.Move(fyne.Position{X: 100, Y: 100})
}

func (r *diagramWidgetRenderer) MinSize() fyne.Size {
	return r.diagramWidget.DesiredSize
}

func (r *diagramWidgetRenderer) Objects() []fyne.CanvasObject {
	obj := make([]fyne.CanvasObject, 0)
	for _, n := range r.diagramWidget.Nodes {
		obj = append(obj, n)
	}
	for _, e := range r.diagramWidget.Links {
		obj = append(obj, e)
	}
	obj = append(obj, r.diagramWidget.dummyBox)
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
