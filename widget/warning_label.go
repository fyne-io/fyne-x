package widget

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// WarningLabel is a generic label for showing error
// messages alongside other widgets. It starts hidden
// and will reveal itself with a given message when
// Warn() is called.
type WarningLabel struct {
	widget.BaseWidget

	canvas *canvas.Text
}

// NewWarningLabel returns a new WarningLabel.
func NewWarningLabel() *WarningLabel {
	w := &WarningLabel{}

	w.canvas = canvas.NewText("", theme.ErrorColor())
	w.canvas.Alignment = fyne.TextAlignTrailing
	w.canvas.TextStyle = fyne.TextStyle{Italic: true}
	w.canvas.TextSize = theme.CaptionTextSize()
	w.canvas.Hide()

	w.ExtendBaseWidget(w)
	return w
}

// Warn shows the label with the given message
func (w *WarningLabel) Warn(msg string) {
	w.ExtendBaseWidget(w)
	// Capitalize beginning of error
	spl := strings.Split(msg, " ")
	if len(spl) > 0 {
		spl[0] = strings.Title(spl[0])
		msg = strings.Join(spl, " ")
	}
	w.canvas.Text = msg
	w.canvas.Show()
}

// Reset will clear the warning message and hide the label.
func (w *WarningLabel) Reset() {
	w.ExtendBaseWidget(w)
	w.canvas.Text = ""
	w.canvas.Hide()
}

func (w *WarningLabel) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	return &warningLabelRenderer{label: w}
}

type warningLabelRenderer struct {
	label *WarningLabel
}

func (w *warningLabelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{w.label.canvas}
}

func (w *warningLabelRenderer) MinSize() fyne.Size {
	return fyne.NewSize(float32(len(w.label.canvas.Text)), w.label.canvas.TextSize+theme.Padding()*2)
}

func (w *warningLabelRenderer) Layout(size fyne.Size) { w.label.canvas.Resize(size) }
func (w *warningLabelRenderer) Refresh()              { w.label.canvas.Refresh() }
func (w *warningLabelRenderer) Destroy()              {}
