package diagramwidget

import "fyne.io/fyne/v2"

// A DiagramElement is a widget that can be placed directly in a diagram. The most common
// elements are Node and Link widgets.
type DiagramElement interface {
	fyne.Widget
	GetDiagram() *DiagramWidget
	GetDiagramElementID() string
	handleDragged(handle *Handle, event *fyne.DragEvent)
	HideHandles()
	ShowHandles()
}

type diagramElement struct {
	diagram *DiagramWidget
	id      string
	handles map[string]*Handle
}

func (de *diagramElement) GetDiagram() *DiagramWidget {
	return de.diagram
}

func (de *diagramElement) GetDiagramElementID() string {
	return de.id
}

func (de *diagramElement) HideHandles() {
	for _, handle := range de.handles {
		handle.Hide()
	}
}

func (de *diagramElement) initialize(diagram *DiagramWidget, id string) {
	de.diagram = diagram
	de.id = id
	de.handles = make(map[string]*Handle)
}

func (de *diagramElement) ShowHandles() {
	for _, handle := range de.handles {
		handle.Show()
	}
}
