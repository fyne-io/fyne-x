package diagramwidget

import "fyne.io/fyne/v2"

// A DiagramElement is a widget that can be placed directly in a diagram. The most common
// elements are Node and Link widgets.
type DiagramElement interface {
	fyne.Widget
	GetDiagram() *DiagramWidget
	handleDragged(handle *Handle, event *fyne.DragEvent)
}

type diagramElement struct {
	diagram *DiagramWidget
}

func (de *diagramElement) GetDiagram() *DiagramWidget {
	return de.diagram
}
