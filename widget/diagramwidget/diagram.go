package diagramwidget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// forceRepaint is a workaround for a Fyne bug (Issue #2205) in which moving a canvas object does not
// trigger repainting. When the issue is resolved, this function and all references to it should be
// removed. The DummyBox on the GlobalDiagram should also be removed.
// The conditionals here are required during initialization.
func (dw *DiagramWidget) forceRepaint() {
	if dw != nil && dw.dummyBox != nil {
		dw.dummyBox.Refresh()
	}
}

// Verify that interfaces are fully implemented
var _ fyne.Tappable = (*DiagramWidget)(nil)

type linkPinPair struct {
	link *DiagramLink
	pin  *LinkPoint
}

// DiagramWidget maintains a diagram consisting of DiagramNodes and DiagramLinks. The layout of
// the nodes and links does not change when the DiagramWidget is resized: they are either positioned
// manually (interactively) or programmatically.
type DiagramWidget struct {
	widget.BaseWidget

	// ID is expected to be unique across all DiagramWidgets in the application.
	ID string

	// Diagrams may want to use a different theme and variant than the application. The default value is the
	// applicaton's theme and variant
	DiagramTheme fyne.Theme
	ThemeVariant fyne.ThemeVariant
	Offset       fyne.Position

	// DesiredSize specifies the size of the displayed diagram. Defaults to 800 x 600
	DesiredSize fyne.Size

	Nodes                          map[string]*DiagramNode
	Links                          map[string]*DiagramLink
	selection                      map[string]DiagramElement
	diagramElementLinkDependencies map[string][]linkPinPair

	// TODO Remove dummyBox when fyne rendering issue is resolved
	dummyBox *canvas.Rectangle
}

// NewDiagramWidget creates a DiagramWidget. The user-supplied ID can be used to map the diagram
// to data structures within the of the application. It is expected to  be unique within the application
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

// AddLink adds a link to the diagram
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

// AddNode adds a node to the diagram
func (dw *DiagramWidget) AddNode(node *DiagramNode) {
	dw.Nodes[node.id] = node
	node.Refresh()
	// TODO add logic to rezise diagram if necessary
}

// CreateRenderer creates the renderer for the diagram
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

// Cursor returns the default cursor
func (dw *DiagramWidget) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// DiagramElementTapped adds the element to the selection when the element is tapped
func (dw *DiagramWidget) DiagramElementTapped(de DiagramElement, event *fyne.PointEvent) {
	if !dw.IsSelected(de) {
		dw.addElementToSelection(de)
	}
	dw.forceRepaint()
}

// DragEnd is called when the drag comes to an end. It refreshes the widget
func (dw *DiagramWidget) DragEnd() {
	dw.Refresh()
}

// GetBackgroundColor returns the background color for the widget from the diagram's theme, which
// may be different from the application's theme.
func (dw *DiagramWidget) GetBackgroundColor() color.Color {
	return dw.DiagramTheme.Color(theme.ColorNameBackground, dw.ThemeVariant)
}

// GetDiagramElement returns the diagram element with the specified ID, whether
// it is a node or a link
func (dw *DiagramWidget) GetDiagramElement(elementID string) DiagramElement {
	var de DiagramElement
	de = dw.Nodes[elementID]
	if de == nil {
		de = dw.Links[elementID]
	}
	return de
}

// GetForegroundColor returns the foreground color from the diagram's theme, which may
// be different from the application's theme
func (dw *DiagramWidget) GetForegroundColor() color.Color {
	return dw.DiagramTheme.Color(theme.ColorNameForeground, dw.ThemeVariant)
}

// GetHoverColor returns the hover color from the diagram's theme, which may may
// be different from the application's  theme
func (dw *DiagramWidget) GetHoverColor() color.Color {
	return dw.DiagramTheme.Color(theme.ColorNameHover, dw.ThemeVariant)
}

// DiagramNodeDragged moves the indicated node and refreshes any links that may be attached
// to it
func (dw *DiagramWidget) DiagramNodeDragged(node *DiagramNode, event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	dw.DisplaceNode(node, delta)
	dw.forceRepaint()
}

// DisplaceNode moves the indicated node and refreshes any links that may be attached
// to it
func (dw *DiagramWidget) DisplaceNode(node *DiagramNode, delta fyne.Position) {
	node.Move(node.Position().Add(delta))
	dw.refreshDependentLinks(node)
	dw.forceRepaint()
}

// Dragged responds to a drag movement in the background of the diagram. It moves all nodes
// in the diagram and refreshes all links.
func (dw *DiagramWidget) Dragged(event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	for _, n := range dw.Nodes {
		dw.DisplaceNode(n, delta)
	}
	dw.Refresh()
}

// IsSelected returns true if the indicated element is currently part of the selection
func (dw *DiagramWidget) IsSelected(de DiagramElement) bool {
	return dw.selection[de.GetDiagramElementID()] != nil
}

// MouseIn responds to the mouse moving into the diagram. It presently is a noop
func (dw *DiagramWidget) MouseIn(event *desktop.MouseEvent) {
}

// MouseOut responds to the mouse leaving the diagram. It presently is a noop
func (dw *DiagramWidget) MouseOut() {
}

// MouseMoved responds to mouse movements in the diagram. It presently is a noop
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

// Tapped  respondss to taps in the diagram background. It removes all diagram elements
// from the selection
func (dw *DiagramWidget) Tapped(event *fyne.PointEvent) {
	for _, de := range dw.selection {
		dw.removeElementFromSelection(de)
	}
	dw.forceRepaint()
}

// diagramWidgetRenderer
type diagramWidgetRenderer struct {
	diagramWidget *DiagramWidget
}

func (r *diagramWidgetRenderer) Destroy() {
}

func (r *diagramWidgetRenderer) Layout(size fyne.Size) {
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
