package diagramwidget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	offset            r2.Vec2
	referencePosition fyne.Position
	displayedText     string
	ForegroundColor   color.Color
}

func NewAnchoredText(text string) *AnchoredText {
	at := &AnchoredText{
		displayedText:     text,
		offset:            r2.MakeVec2(0, 0),
		ForegroundColor:   theme.ForegroundColor(),
		referencePosition: fyne.Position{X: 0, Y: 0},
	}
	at.ExtendBaseWidget(at)
	return at
}

func (at *AnchoredText) CreateRenderer() fyne.WidgetRenderer {
	atr := &anchoredTextRenderer{
		widget:     at,
		textObject: canvas.NewText(at.displayedText, color.Black),
	}

	atr.Refresh()

	return atr
}

func (at *AnchoredText) Displace(delta fyne.Position) {
	at.Move(at.Position().Add(delta))
}

func (at *AnchoredText) DragEnd() {
	at.Refresh()
}

func (at *AnchoredText) Dragged(event *fyne.DragEvent) {
	delta := fyne.Position{X: event.Dragged.DX, Y: event.Dragged.DY}
	at.Move(at.Position().Add(delta))
	at.Refresh()
	ForceRepaint()
}

func (at *AnchoredText) MinSize() fyne.Size {
	minSize := fyne.Size{Height: 25, Width: 50}
	return minSize
}

func (at *AnchoredText) MouseIn(event *desktop.MouseEvent) {
	// at.textObject.TextStyle.Bold = true
	at.Refresh()
}

func (at *AnchoredText) MouseMoved(event *desktop.MouseEvent) {

}

func (at *AnchoredText) MouseOut() {
	// at.textObject.TextStyle.Bold = false
	at.Refresh()
}

func (at *AnchoredText) Move(position fyne.Position) {
	delta := r2.MakeVec2(float64(position.X-at.Position().X), float64(position.Y-at.Position().Y))
	at.offset = at.offset.Add(delta)
	at.BaseWidget.Move(position)
}

func (at *AnchoredText) SetForegroundColor(fc color.Color) {
	at.ForegroundColor = fc
	at.Refresh()
}

func (at *AnchoredText) SetReferencePosition(position fyne.Position) {
	delta := fyne.Delta{DX: float32(position.X - at.referencePosition.X), DY: float32(position.Y - at.referencePosition.Y)}
	// We don't want to change the offset here, so we call the BaseWidget.Move directly
	at.BaseWidget.Move(at.Position().Add(delta))
	at.referencePosition = position
}

// anchoredTextRenderer
type anchoredTextRenderer struct {
	widget     *AnchoredText
	textObject *canvas.Text
}

func (atr *anchoredTextRenderer) Destroy() {

}

func (atr *anchoredTextRenderer) Layout(size fyne.Size) {
	atr.widget.Resize(atr.textObject.MinSize())
}

func (atr *anchoredTextRenderer) MinSize() fyne.Size {
	return atr.textObject.MinSize()
}

func (atr *anchoredTextRenderer) Objects() []fyne.CanvasObject {
	canvasObjects := []fyne.CanvasObject{
		atr.textObject,
	}
	return canvasObjects
}

func (atr *anchoredTextRenderer) Refresh() {
	// atr.widget.textObject.Color = atr.widget.ForegroundColor
}
