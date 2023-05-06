package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestGridWrap_New(t *testing.T) {
	g := createGridWrap(1000)
	template := widget.NewIcon(theme.AccountIcon())

	assert.Equal(t, 1000, g.Length())
	assert.GreaterOrEqual(t, g.MinSize().Width, template.MinSize().Width)
	assert.Equal(t, float32(0), g.offsetY)
}

func TestGridWrap_OffsetChange(t *testing.T) {
	g := createGridWrap(1000)

	assert.Equal(t, float32(0), g.offsetY)

	g.scroller.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -280)})

	assert.NotEqual(t, 0, g.offsetY)
}

func TestGridWrap_ScrollTo(t *testing.T) {
	g := createGridWrap(1000)

	// override update item to keep track of greatest item rendered
	oldUpdateFunc := g.UpdateItem
	var greatest GridWrapItemID = -1
	g.UpdateItem = func(id GridWrapItemID, item fyne.CanvasObject) {
		if id > greatest {
			greatest = id
		}
		oldUpdateFunc(id, item)
	}

	g.ScrollTo(650)
	assert.GreaterOrEqual(t, greatest, 650)

	g.ScrollTo(800)
	assert.GreaterOrEqual(t, greatest, 800)

	g.ScrollToBottom()
	assert.Equal(t, greatest, GridWrapItemID(999))
}

func TestGridWrap_ScrollToTop(t *testing.T) {
	g := createGridWrap(1000)
	g.ScrollTo(750)
	assert.NotEqual(t, g.offsetY, float32(0))
	g.ScrollToTop()
	assert.Equal(t, g.offsetY, float32(0))
}

func createGridWrap(items int) *GridWrap {
	data := make([]fyne.Resource, items)
	for i := 0; i < items; i++ {
		switch i % 10 {
		case 0:
			data[i] = theme.AccountIcon()
		case 1:
			data[i] = theme.CancelIcon()
		case 2:
			data[i] = theme.CheckButtonIcon()
		case 3:
			data[i] = theme.FileApplicationIcon()
		case 4:
			data[i] = theme.FileVideoIcon()
		case 5:
			data[i] = theme.DocumentIcon()
		case 6:
			data[i] = theme.MediaPlayIcon()
		case 7:
			data[i] = theme.MediaRecordIcon()
		case 8:
			data[i] = theme.FolderIcon()
		case 9:
			data[i] = theme.FolderOpenIcon()
		}
	}

	list := NewGridWrap(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			icon := widget.NewIcon(theme.DocumentIcon())
			return icon
		},
		func(id GridWrapItemID, item fyne.CanvasObject) {
			item.(*widget.Icon).SetResource(data[id])
		},
	)
	list.Resize(fyne.NewSize(200, 400))
	return list
}
