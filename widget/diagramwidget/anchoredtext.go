package diagramwidget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"
)

// AnchoredText provides a text annotation for a path that is anchored to one
// of the path's reference points (e.g. end or middle). The anchored text may
// be moved independently, but it keeps track of its position relative to the
// reference point. If the reference point moves, the AnchoredText will also
// move by the same amount
type AnchoredText struct {
	widget.BaseWidget
	diagramElement       DiagramElement
	ID                   string
	offset               r2.Vec2
	referencePosition    fyne.Position
	displayedTextBinding binding.String
	ForegroundColor      color.Color
	textEntry            *widget.Entry
}

// NewAnchoredText creates an textual annotation for a link. After it is created, one of the
// three Add<position>AnchoredText methods must be called on the link to actually associate the
// anchored text with the appropriate reference point on the link.
func NewAnchoredText(text string, id ...string) *AnchoredText {
	at := &AnchoredText{
		offset:            r2.MakeVec2(0, 0),
		ForegroundColor:   theme.Color(theme.ColorNameForeground),
		referencePosition: fyne.Position{X: 0, Y: 0},
	}
	if len(id) > 0 {
		at.ID = id[0]
	}
	at.displayedTextBinding = binding.NewString()
	at.displayedTextBinding.Set(text)
	at.textEntry = widget.NewEntryWithData(at.displayedTextBinding)
	at.displayedTextBinding.AddListener(at)
	at.textEntry.Wrapping = fyne.TextWrapOff
	at.textEntry.Scroll = container.ScrollNone
	at.textEntry.Validator = nil
	at.ExtendBaseWidget(at)
	return at
}

// CreateRenderer is the required method for a widget extension
func (at *AnchoredText) CreateRenderer() fyne.WidgetRenderer {
	atr := &anchoredTextRenderer{
		widget: at,
	}
	atr.Refresh()

	return atr
}

// DataChanged is the callback function for the displayedTextBinding.
func (at *AnchoredText) DataChanged() {
	callback := at.diagramElement.GetDiagram().AnchoredTextChangedCallback
	if callback != nil {
		callback(at)
	}
	at.Refresh()
}

// Displace moves the anchored text relative to its reference position.
func (at *AnchoredText) Displace(delta fyne.Position) {
	at.Move(at.Position().Add(delta))
	callback := at.diagramElement.GetDiagram().AnchoredTextChangedCallback
	if callback != nil {
		callback(at)
	}
}

// DragEnd is one of the required methods for a draggable widget. It just refreshes the widget.
func (at *AnchoredText) DragEnd() {
	callback := at.diagramElement.GetDiagram().AnchoredTextChangedCallback
	if callback != nil {
		callback(at)
	}
	at.Refresh()
}

// Dragged is the required method for a draggable widget. It moves the anchored text
// relative to its reference position
func (at *AnchoredText) Dragged(event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	at.Move(at.Position().Add(delta))
	callback := at.diagramElement.GetDiagram().AnchoredTextChangedCallback
	if callback != nil {
		callback(at)
	}
	at.Refresh()
}

// GetDiagramElement returns the diagram element to which the anchored text belongs
func (at *AnchoredText) GetDiagramElement() DiagramElement {
	return at.diagramElement
}

// GetDisplayedTextBinding returns the binding for the displayed text
func (at *AnchoredText) GetDisplayedTextBinding() binding.String {
	return at.displayedTextBinding
}

// GetOffset returns the X and Y values of the anchored text's offset from the reference position
func (at *AnchoredText) GetOffset() (X float64, Y float64) {
	return at.offset.X, at.offset.Y
}

// GetReferencePosition returns the X and Y values of the anchored text's reference position
func (at *AnchoredText) GetReferencePosition() (X float64, Y float64) {
	return float64(at.referencePosition.X), float64(at.referencePosition.Y)
}

// GetTextEntry returns the entry widget
func (at *AnchoredText) GetTextEntry() *widget.Entry {
	return at.textEntry
}

// MinSize returns the size of the entry widget plus a one-pixel border
func (at *AnchoredText) MinSize() fyne.Size {
	textEntryMinSize := at.textEntry.MinSize()
	minSize := fyne.NewSize(textEntryMinSize.Width+10, textEntryMinSize.Height+10)
	return minSize
}

// MouseIn is one of the required methods for a mouseable widget.
func (at *AnchoredText) MouseIn(event *desktop.MouseEvent) {
}

// MouseMoved is one of the required methods for a mouseable widget
func (at *AnchoredText) MouseMoved(event *desktop.MouseEvent) {

}

// MouseOut is one of the required methods for a mouseable widget
func (at *AnchoredText) MouseOut() {
}

// Move overrides the BaseWidget's Move method. It updates the anchored text's offset
// and then calls the normal BaseWidget.Move method.
func (at *AnchoredText) Move(position fyne.Position) {
	delta := r2.MakeVec2(float64(position.X-at.Position().X), float64(position.Y-at.Position().Y))
	newOffset := at.offset.Add(delta)
	at.SetOffsetNoCallback(newOffset.X, newOffset.Y)
	callback := at.diagramElement.GetDiagram().AnchoredTextChangedCallback
	if callback != nil {
		callback(at)
	}
}

// SetForegroundColor sets the text color
func (at *AnchoredText) SetForegroundColor(fc color.Color) {
	at.ForegroundColor = fc
	at.Refresh()
}

// SetOffset sets the X and Y values of the anchored text's offset from the reference position
func (at *AnchoredText) SetOffset(X float64, Y float64) {
	if at.offset.X != X || at.offset.Y != Y {
		at.SetOffsetNoCallback(X, Y)
		callback := at.diagramElement.GetDiagram().AnchoredTextChangedCallback
		if callback != nil {
			callback(at)
		}
	}
}

// SetOffsetNoCallback sets the X and Y values of the anchored text's offset from the reference position
func (at *AnchoredText) SetOffsetNoCallback(X float64, Y float64) {
	if at.offset.X != X || at.offset.Y != Y {
		at.offset.X = X
		at.offset.Y = Y
		delta := fyne.Delta{DX: float32(X), DY: float32(Y)}
		at.BaseWidget.Move(at.referencePosition.Add(delta))
	}
}

// SetReferencePosition sets the reference position of the anchored text and calls
// the BaseWidget.Move() method to actually move the displayed text
func (at *AnchoredText) SetReferencePosition(position fyne.Position) {
	if position != at.referencePosition {
		at.SetReferencePositionNoCallback(position)
		callback := at.diagramElement.GetDiagram().AnchoredTextChangedCallback
		if callback != nil {
			callback(at)
		}
	}
}

// SetReferencePositionNoCallback sets the reference position of the anchored text and calls
// the BaseWidget.Move() method to actually move the displayed text
func (at *AnchoredText) SetReferencePositionNoCallback(position fyne.Position) {
	if position != at.referencePosition {
		delta := fyne.Delta{DX: float32(position.X - at.referencePosition.X), DY: float32(position.Y - at.referencePosition.Y)}
		// We don't want to change the offset here, so we call the BaseWidget.Move directly
		at.BaseWidget.Move(at.Position().Add(delta))
		at.referencePosition = position
	}
}

// SetText sets the anchored text without invoking the callback
func (at *AnchoredText) SetText(text string) {
	at.displayedTextBinding.Set(text)
}

// anchoredTextRenderer
type anchoredTextRenderer struct {
	widget *AnchoredText
}

func (atr *anchoredTextRenderer) Destroy() {

}

func (atr *anchoredTextRenderer) Layout(size fyne.Size) {
}

func (atr *anchoredTextRenderer) MinSize() fyne.Size {
	return atr.widget.textEntry.MinSize()
}

func (atr *anchoredTextRenderer) Objects() []fyne.CanvasObject {
	canvasObjects := []fyne.CanvasObject{
		atr.widget.textEntry,
	}
	return canvasObjects
}

func (atr *anchoredTextRenderer) Refresh() {
	atr.widget.Resize(atr.widget.MinSize())
	atr.widget.textEntry.Resize(atr.widget.textEntry.MinSize())
	atr.widget.textEntry.Move(fyne.NewPos(5, 5))
	atr.widget.textEntry.Refresh()
}
