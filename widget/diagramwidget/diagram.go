package diagramwidget

import (
	"image/color"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ForceRepaint is a workaround for a Fyne bug (Issue #2205) in which moving a canvas object does not
// trigger repainting. When the issue is resolved, this function and all references to it should be
// removed. The DummyBox on the GlobalDiagram should also be removed.
// The conditionals here are required during initialization.
func (dw *DiagramWidget) ForceRepaint() {
	if dw != nil && dw.dummyBox != nil {
		dw.dummyBox.Refresh()
	}
}

// Verify that interfaces are fully implemented
var _ fyne.Tappable = (*DiagramWidget)(nil)

// Default values
var defaultPadColor = color.RGBA{121, 237, 119, 255}

type linkPadPair struct {
	link *BaseDiagramLink
	pad  ConnectionPad
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

	Nodes                          map[string]DiagramNode
	Links                          map[string]DiagramLink
	primarySelection               DiagramElement
	selection                      map[string]DiagramElement
	diagramElementLinkDependencies map[string][]linkPadPair
	connectionTransaction          *connectionTransaction
	padColor                       color.Color
	// IsConnectionAllowedCallback is called to determine whether a particular connection between a link and a pad is allowed
	IsConnectionAllowedCallback func(DiagramLink, LinkEnd, ConnectionPad) bool
	// LinkConnectionChangedCallback is called when a link connection changes. The string can either be
	// "source" or "target". The first pad is the old pad, the second one is the new pad
	LinkConnectionChangedCallback func(DiagramLink, string, ConnectionPad, ConnectionPad)
	// OnTappedCallback is called when the diagram background is tapped. If present, it overrides the default
	// diagram behavior for Tapped()
	OnTappedCallback func(*DiagramWidget, *fyne.PointEvent)
	// PrimaryDiagramElementSelectionChangedCallback is called when the primary element selection changes
	PrimaryDiagramElementSelectionChangedCallback func(string)
	// ElementTappedExtendsSelection determines the behavior when one or more elements are already selected and
	// an element that is not currently selected is tapped. When true, the new element is added to the selection.
	// When false, the selection is cleared and the new element is made the only selected element.
	ElementTappedExtendsSelection bool

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
		Nodes:                          map[string]DiagramNode{},
		Links:                          map[string]DiagramLink{},
		dummyBox:                       canvas.NewRectangle(color.Transparent),
		selection:                      map[string]DiagramElement{},
		diagramElementLinkDependencies: map[string][]linkPadPair{},
		padColor:                       defaultPadColor,
	}
	dw.dummyBox.SetMinSize(fyne.Size{Height: 50, Width: 50})
	dw.dummyBox.Move(fyne.Position{X: 50, Y: 50})

	dw.ExtendBaseWidget(dw)

	return dw
}

func (dw *DiagramWidget) addElementToSelection(de DiagramElement) {
	if !dw.IsSelected(de) {
		if dw.primarySelection == nil {
			dw.primarySelection = de
			if dw.PrimaryDiagramElementSelectionChangedCallback != nil {
				dw.PrimaryDiagramElementSelectionChangedCallback(de.GetDiagramElementID())
			}
		}
		dw.selection[de.GetDiagramElementID()] = de
		de.ShowHandles()
	}
}

// addLink adds a link to the diagram
func (dw *DiagramWidget) addLink(link DiagramLink) {
	dw.Links[link.GetDiagramElementID()] = link
	link.Refresh()
	// TODO add logic to rezise diagram if necessary
}

func (dw *DiagramWidget) addLinkDependency(diagramElement DiagramElement, link *BaseDiagramLink, pad ConnectionPad) {
	deID := diagramElement.GetDiagramElementID()
	currentDependencies := dw.diagramElementLinkDependencies[deID]
	if currentDependencies == nil {
		dw.diagramElementLinkDependencies[deID] = []linkPadPair{{link, pad}}
	} else {
		for _, pair := range currentDependencies {
			if pair.link == link && pair.pad == pad {
				// it's already there
				return
			}
		}
		dw.diagramElementLinkDependencies[deID] = append(currentDependencies, linkPadPair{link, pad})
	}
}

// addNode adds a node to the diagram
func (dw *DiagramWidget) addNode(node DiagramNode) {
	dw.Nodes[node.GetDiagramElementID()] = node
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

// ClearSelection clears the selection and invokes the PrimaryDiagramElementSelectionChangedCallback
func (dw *DiagramWidget) ClearSelection() {
	for _, de := range dw.selection {
		de.HideHandles()
		dw.removeElementFromSelection(de)
	}
}

// ClearSelectionNoCallback clears the selection. It does not invoke the PrimaryDiagramElementSelectionChangedCallback
func (dw *DiagramWidget) ClearSelectionNoCallback() {
	dw.primarySelection = nil
	for _, element := range dw.selection {
		element.HideHandles()
	}
	dw.selection = map[string]DiagramElement{}
}

// Cursor returns the default cursor
func (dw *DiagramWidget) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// DiagramElementTapped adds the element to the selection when the element is tapped
func (dw *DiagramWidget) DiagramElementTapped(de DiagramElement, event *fyne.PointEvent) {
	if !dw.ElementTappedExtendsSelection {
		dw.ClearSelectionNoCallback()
	}
	if !dw.IsSelected(de) {
		dw.addElementToSelection(de)
	}
	dw.ForceRepaint()
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
	if de == nil || reflect.ValueOf(de).IsNil() {
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
func (dw *DiagramWidget) DiagramNodeDragged(node *BaseDiagramNode, event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	dw.DisplaceNode(node, delta)
	dw.ForceRepaint()
}

// DisplaceNode moves the indicated node and refreshes any links that may be attached
// to it
func (dw *DiagramWidget) DisplaceNode(node DiagramNode, delta fyne.Position) {
	node.Move(node.Position().Add(delta))
	dw.refreshDependentLinks(node)
	dw.ForceRepaint()
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

// GetDiagramElements returns a map of all of the diagram's DiagramElements
func (dw *DiagramWidget) GetDiagramElements() map[string]DiagramElement {
	diagramElements := map[string]DiagramElement{}
	for key, node := range dw.Nodes {
		diagramElements[key] = node
	}
	for key, link := range dw.Links {
		diagramElements[key] = link
	}
	return diagramElements
}

// hideAllPads is a work-around for fyne Issue #3906 in which a child's Hoverable interface
// (i.e. the pad) masks the parent's Tappable interface. This function (and all references to
// it) should be removed when this issue has been resolved
func (dw *DiagramWidget) hideAllPads() {
	for _, node := range dw.Nodes {
		for _, pad := range node.GetConnectionPads() {
			pad.Hide()
		}
	}
	for _, link := range dw.Links {
		for _, pad := range link.GetConnectionPads() {
			pad.Hide()
		}
	}
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

// removeDependenciesInvolvingLink re-creates the diagram's dependencies, omitting any
// that involve the indicated link. This is a convoluted way of removing any entries
// involving the link.
func (dw *DiagramWidget) removeDependenciesInvolvingLink(linkID string) {
	newMap := map[string][]linkPadPair{}
	for elementID, dependencies := range dw.diagramElementLinkDependencies {
		newDependencies := []linkPadPair{}
		for _, pair := range dependencies {
			if pair.link.id != linkID {
				newDependencies = append(newDependencies, pair)
			}
		}
		if len(newDependencies) > 0 {
			newMap[elementID] = newDependencies
		}
	}
	dw.diagramElementLinkDependencies = newMap
}

func (dw *DiagramWidget) removeElementFromSelection(de DiagramElement) {
	if dw.IsSelected(de) {
		delete(dw.selection, de.GetDiagramElementID())
		if dw.primarySelection == de {
			dw.primarySelection = nil
			if dw.PrimaryDiagramElementSelectionChangedCallback != nil {
				dw.PrimaryDiagramElementSelectionChangedCallback("")
			}
		}
		de.HideHandles()
	}
}

func (dw *DiagramWidget) removeLinkDependency(diagramElement DiagramElement, link *BaseDiagramLink, pad ConnectionPad) {
	deID := diagramElement.GetDiagramElementID()
	currentDependencies := dw.diagramElementLinkDependencies[deID]
	if currentDependencies == nil {
		return
	}
	for i, pair := range currentDependencies {
		if pair.link == link && pair.pad == pad {
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

// RemoveElement removes the element from the diagram. It also removes any linkss to the element
func (dw *DiagramWidget) RemoveElement(elementID string) {
	element := dw.GetDiagramElement(elementID)
	// We make a copy of the dependencies because the array can get modified during the iteration
	currentDependencies := append([]linkPadPair(nil), dw.diagramElementLinkDependencies[elementID]...)
	for _, pair := range currentDependencies {
		dw.RemoveElement(pair.link.id)
	}
	delete(dw.diagramElementLinkDependencies, elementID)
	if element.IsNode() {
		delete(dw.Nodes, elementID)
	} else if element.IsLink() {
		delete(dw.Links, elementID)
		dw.removeDependenciesInvolvingLink(elementID)
	}
	dw.Refresh()
}

// SelectDiagramElement clears the selection, makes the indicated element the primary selection, and invokes
// the PrimaryDiagramElementSelectionChangedCallback
func (dw *DiagramWidget) SelectDiagramElement(element DiagramElement) {
	dw.ClearSelection()
	dw.addElementToSelection(element)
}

// SelectDiagramElementNoCallback makes the indicated element the PrimarySelection. It does not invoke the
// PrimaryDiagramElementSelectionChangedCallback
func (dw *DiagramWidget) SelectDiagramElementNoCallback(id string) {
	dw.ClearSelectionNoCallback()
	element := dw.GetDiagramElement(id)
	if element != nil {
		dw.primarySelection = element
		dw.selection[id] = element
		element.ShowHandles()
	}
}

// showAllPads is a work-around for fyne Issue #3906 in which a child's Hoverable interface
// (i.e. the pad) masks the parent's Tappable interface. This function (and all references to
// it) should be removed when this issue has been resolved
func (dw *DiagramWidget) showAllPads() {
	for _, node := range dw.Nodes {
		for _, pad := range node.GetConnectionPads() {
			pad.Show()
		}
	}
	for _, link := range dw.Links {
		for _, pad := range link.GetConnectionPads() {
			pad.Show()
		}
	}
}

// StartNewLinkConnectionTransaction starts the process of adding a link, setting up for the source connection
func (dw *DiagramWidget) StartNewLinkConnectionTransaction(link DiagramLink) {
	dw.connectionTransaction = NewConnectionTransaction(link.getBaseDiagramLink().linkPoints[0], link, nil, fyne.NewPos(0, 0))
	dw.showAllPads()
}

// Tapped  respondss to taps in the diagram background. It removes all diagram elements
// from the selection
func (dw *DiagramWidget) Tapped(event *fyne.PointEvent) {
	if dw.OnTappedCallback != nil {
		dw.OnTappedCallback(dw, event)
	} else {
		dw.ClearSelection()
	}
	dw.ForceRepaint()
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
