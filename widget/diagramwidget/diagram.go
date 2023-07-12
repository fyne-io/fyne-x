package diagramwidget

import (
	"image/color"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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

	Offset fyne.Position

	// DesiredSize specifies the size of the displayed diagram. Defaults to 800 x 600
	DesiredSize fyne.Size

	DefaultDiagramElementProperties DiagramElementProperties
	Nodes                           map[string]DiagramNode
	Links                           map[string]DiagramLink
	primarySelection                DiagramElement
	selection                       map[string]DiagramElement
	diagramElementLinkDependencies  map[string][]linkPadPair
	// ConnectionTransaction holds transient data during the creation of a link. It is public for testing purposes
	ConnectionTransaction *ConnectionTransaction
	padColor              color.Color
	// IsConnectionAllowedCallback is called to determine whether a particular connection between a link and a pad is allowed
	IsConnectionAllowedCallback func(DiagramLink, LinkEnd, ConnectionPad) bool
	// LinkConnectionChangedCallback is called when a link connection changes. The string can either be
	// "source" or "target". The first pad is the old pad, the second one is the new pad
	LinkConnectionChangedCallback func(DiagramLink, string, ConnectionPad, ConnectionPad)
	// MouseDownCallback is called when a MouseDown occurs in the diagram
	MouseDownCallback func(*desktop.MouseEvent)
	// MouseInCallback is called when a MouseIn occurs in the diagram
	MouseInCallback func(*desktop.MouseEvent)
	// MouseMovedCallback is called when a MouseMove occurs in the diagram
	MouseMovedCallback func(*desktop.MouseEvent)
	// MouseOutCallback is called when a MouseOut occurs in the diagram
	MouseOutCallback func()
	// MouseUpCallback is invoked when a MouseUp occurs in the diagram
	MouseUpCallback func(*desktop.MouseEvent)
	// OnTappedCallback is called when the diagram background is tapped. If present, it overrides the default
	// diagram behavior for Tapped()
	OnTappedCallback func(*DiagramWidget, *fyne.PointEvent)
	// PrimaryDiagramElementSelectionChangedCallback is called when the primary element selection changes
	PrimaryDiagramElementSelectionChangedCallback func(string)
	// ElementTappedExtendsSelection determines the behavior when one or more elements are already selected and
	// an element that is not currently selected is tapped. When true, the new element is added to the selection.
	// When false, the selection is cleared and the new element is made the only selected element.
	ElementTappedExtendsSelection bool
}

// NewDiagramWidget creates a DiagramWidget. The user-supplied ID can be used to map the diagram
// to data structures within the of the application. It is expected to  be unique within the application
func NewDiagramWidget(id string) *DiagramWidget {
	dw := &DiagramWidget{
		ID: id,
		// DiagramTheme:                   fyne.CurrentApp().Settings().Theme(),
		// ThemeVariant:                   fyne.CurrentApp().Settings().ThemeVariant(),
		DesiredSize:                    fyne.Size{Width: 800, Height: 600},
		Offset:                         fyne.Position{X: 0, Y: 0},
		Nodes:                          map[string]DiagramNode{},
		Links:                          map[string]DiagramLink{},
		selection:                      map[string]DiagramElement{},
		diagramElementLinkDependencies: map[string][]linkPadPair{},
		padColor:                       defaultPadColor,
	}
	appTheme := fyne.CurrentApp().Settings().Theme()
	appVariant := fyne.CurrentApp().Settings().ThemeVariant()
	dw.DefaultDiagramElementProperties.ForegroundColor = appTheme.Color(theme.ColorNameForeground, appVariant)
	dw.DefaultDiagramElementProperties.HandleColor = appTheme.Color(theme.ColorNameForeground, appVariant)
	dw.DefaultDiagramElementProperties.BackgroundColor = appTheme.Color(theme.ColorNameBackground, appVariant)
	dw.DefaultDiagramElementProperties.TextSize = 12
	dw.DefaultDiagramElementProperties.CaptionTextSize = appTheme.Size(theme.SizeNameCaptionText)
	dw.DefaultDiagramElementProperties.Padding = appTheme.Size(theme.SizeNamePadding)
	dw.DefaultDiagramElementProperties.StrokeWidth = 1
	dw.DefaultDiagramElementProperties.HandleStrokeWidth = 1

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
}

// DragEnd is called when the drag comes to an end. It refreshes the widget
func (dw *DiagramWidget) DragEnd() {
	dw.Refresh()
}

// GetBackgroundColor returns the background color for the widget from the diagram's theme, which
// may be different from the application's theme.
func (dw *DiagramWidget) GetBackgroundColor() color.Color {
	return dw.DefaultDiagramElementProperties.BackgroundColor
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
	return dw.DefaultDiagramElementProperties.ForegroundColor
}

// DiagramNodeDragged moves the indicated node and refreshes any links that may be attached
// to it
func (dw *DiagramWidget) DiagramNodeDragged(node *BaseDiagramNode, event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	dw.DisplaceNode(node, delta)
}

// DisplaceNode moves the indicated node and refreshes any links that may be attached
// to it
func (dw *DiagramWidget) DisplaceNode(node DiagramNode, delta fyne.Position) {
	node.Move(node.Position().Add(delta))
	dw.refreshDependentLinks(node)
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

// GetPrimarySelection returns the diagram element that is currently selected
func (dw *DiagramWidget) GetPrimarySelection() DiagramElement {
	return dw.primarySelection
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

// MouseDown responds to MouseDown events. It invokes the callback, if present
func (dw *DiagramWidget) MouseDown(event *desktop.MouseEvent) {
	if dw.MouseDownCallback != nil {
		dw.MouseDownCallback(event)
	}
}

// MouseIn responds to the mouse moving into the diagram. It presently is a noop
func (dw *DiagramWidget) MouseIn(event *desktop.MouseEvent) {
	if dw.MouseInCallback != nil {
		dw.MouseInCallback(event)
	}
}

// MouseMoved responds to mouse movements in the diagram. It presently is a noop
func (dw *DiagramWidget) MouseMoved(event *desktop.MouseEvent) {
	if dw.MouseMovedCallback != nil {
		dw.MouseMovedCallback(event)
	}
}

// MouseOut responds to the mouse leaving the diagram. It presently is a noop
func (dw *DiagramWidget) MouseOut() {
	if dw.MouseOutCallback != nil {
		dw.MouseOutCallback()
	}
}

// MouseDown responds to MouseDown events. It invokes the callback, if present
func (dw *DiagramWidget) MouseUp(event *desktop.MouseEvent) {
	if dw.MouseUpCallback != nil {
		dw.MouseUpCallback(event)
	}
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
	if element == nil {
		return
	}
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
	dw.ConnectionTransaction = NewConnectionTransaction(link.getBaseDiagramLink().linkPoints[0], link, nil, fyne.NewPos(0, 0))
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
