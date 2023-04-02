package widget

import (
	"fmt"
	"math"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Declare conformity with Widget interface.
var _ fyne.Widget = (*GridWrapList)(nil)

// GridWrapListItemID is the ID of an individual item in the GridWrapList widget.
type GridWrapListItemID int

// GridWrapList is a widget with an API very similar to widget.List,
// that lays out items in a wrapping grid similar to container.NewGridWrap.
// It caches and reuses widgets for performance.
type GridWrapList struct {
	widget.BaseWidget

	Length     func() int                                          `json:"-"`
	CreateItem func() fyne.CanvasObject                            `json:"-"`
	UpdateItem func(id GridWrapListItemID, item fyne.CanvasObject) `json:"-"`

	scroller      *container.Scroll
	itemMin       fyne.Size
	offsetY       float32
	offsetUpdated func(fyne.Position)
}

// NewGridWrapList creates and returns a GridWrapList widget for displaying items in
// a wrapping grid layout with scrolling and caching for performance.
func NewGridWrapList(length func() int, createItem func() fyne.CanvasObject, updateItem func(GridWrapListItemID, fyne.CanvasObject)) *GridWrapList {
	gwList := &GridWrapList{BaseWidget: widget.BaseWidget{}, Length: length, CreateItem: createItem, UpdateItem: updateItem}
	gwList.ExtendBaseWidget(gwList)
	return gwList
}

// NewGridWrapListWithData creates a new GridWrapList widget that will display the contents of the provided data.
func NewGridWrapListWithData(data binding.DataList, createItem func() fyne.CanvasObject, updateItem func(binding.DataItem, fyne.CanvasObject)) *GridWrapList {
	gwList := NewGridWrapList(
		data.Length,
		createItem,
		func(i GridWrapListItemID, o fyne.CanvasObject) {
			item, err := data.GetItem(int(i))
			if err != nil {
				fyne.LogError(fmt.Sprintf("Error getting data item %d", i), err)
				return
			}
			updateItem(item, o)
		})

	data.AddListener(binding.NewDataListener(gwList.Refresh))
	return gwList
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (l *GridWrapList) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)

	if f := l.CreateItem; f != nil {
		if l.itemMin.IsZero() {
			l.itemMin = f().MinSize()
		}
	}
	layout := &fyne.Container{}
	l.scroller = container.NewVScroll(layout)
	layout.Layout = newGridWrapListLayout(l)
	layout.Resize(layout.MinSize())
	objects := []fyne.CanvasObject{l.scroller}
	lr := newGridWrapListRenderer(objects, l, l.scroller, layout)
	return lr
}

// MinSize returns the size that this widget should not shrink below.
func (l *GridWrapList) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)

	return l.BaseWidget.MinSize()
}

func (l *GridWrapList) scrollTo(id GridWrapListItemID) {
	if l.scroller == nil {
		return
	}
	row := math.Floor(float64(id) / float64(l.getColCount()))
	y := (float32(row) * l.itemMin.Height) + (float32(row) * theme.Padding())
	if y < l.scroller.Offset.Y {
		l.scroller.Offset.Y = y
	} else if y+l.itemMin.Height > l.scroller.Offset.Y+l.scroller.Size().Height {
		l.scroller.Offset.Y = y + l.itemMin.Height - l.scroller.Size().Height
	}
	l.offsetUpdated(l.scroller.Offset)
}

// Resize is called when this GridWrapList should change size. We refresh to ensure invisible items are drawn.
func (l *GridWrapList) Resize(s fyne.Size) {
	l.BaseWidget.Resize(s)
	l.offsetUpdated(l.scroller.Offset)
	l.scroller.Content.(*fyne.Container).Layout.(*gridWrapListLayout).updateList(true)
}

// ScrollTo scrolls to the item represented by id
func (l *GridWrapList) ScrollTo(id GridWrapListItemID) {
	length := 0
	if f := l.Length; f != nil {
		length = f()
	}
	if id < 0 || int(id) >= length {
		return
	}
	l.scrollTo(id)
	l.Refresh()
}

// ScrollToBottom scrolls to the end of the list
func (l *GridWrapList) ScrollToBottom() {
	length := 0
	if f := l.Length; f != nil {
		length = f()
	}
	if length > 0 {
		length--
	}
	l.scrollTo(GridWrapListItemID(length))
	l.Refresh()
}

// ScrollToTop scrolls to the start of the list
func (l *GridWrapList) ScrollToTop() {
	l.scrollTo(0)
	l.Refresh()
}

// ScrollToOffset scrolls the list to the given offset position
func (l *GridWrapList) ScrollToOffset(offset float32) {
	l.scroller.Offset.Y = offset
	l.offsetUpdated(l.scroller.Offset)
}

// GetScrollOffset returns the current scroll offset position
func (l *GridWrapList) GetScrollOffset() float32 {
	return l.offsetY
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*gridWrapListRenderer)(nil)

type gridWrapListRenderer struct {
	objects []fyne.CanvasObject

	list     *GridWrapList
	scroller *container.Scroll
	layout   *fyne.Container
}

func newGridWrapListRenderer(objects []fyne.CanvasObject, l *GridWrapList, scroller *container.Scroll, layout *fyne.Container) *gridWrapListRenderer {
	lr := &gridWrapListRenderer{objects: objects, list: l, scroller: scroller, layout: layout}
	lr.scroller.OnScrolled = l.offsetUpdated
	return lr
}

func (l *gridWrapListRenderer) Layout(size fyne.Size) {
	l.scroller.Resize(size)
}

func (l *gridWrapListRenderer) MinSize() fyne.Size {
	return l.scroller.MinSize().Max(l.list.itemMin)
}

func (l *gridWrapListRenderer) Refresh() {
	if f := l.list.CreateItem; f != nil {
		l.list.itemMin = f().MinSize()
	}
	l.Layout(l.list.Size())
	l.scroller.Refresh()
	l.layout.Layout.(*gridWrapListLayout).updateList(true)
	canvas.Refresh(l.list)
}

func (l *gridWrapListRenderer) Destroy() {
}

func (l *gridWrapListRenderer) Objects() []fyne.CanvasObject {
	return l.objects
}

func (r *gridWrapListRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}

// Declare conformity with Layout interface.
var _ fyne.Layout = (*gridWrapListLayout)(nil)

type gridWrapListLayout struct {
	list     *GridWrapList
	children []fyne.CanvasObject

	itemPool   *syncPool
	visible    map[GridWrapListItemID]fyne.CanvasObject
	renderLock sync.Mutex
}

func newGridWrapListLayout(list *GridWrapList) fyne.Layout {
	l := &gridWrapListLayout{list: list, itemPool: &syncPool{}, visible: make(map[GridWrapListItemID]fyne.CanvasObject)}
	list.offsetUpdated = l.offsetUpdated
	return l
}

func (l *gridWrapListLayout) Layout([]fyne.CanvasObject, fyne.Size) {
	l.updateList(true)
}

func (l *gridWrapListLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	if lenF := l.list.Length; lenF != nil {
		cols := l.list.getColCount()
		rows := float32(math.Ceil(float64(lenF()) / float64(cols)))
		return fyne.NewSize(l.list.itemMin.Width,
			(l.list.itemMin.Height+theme.Padding())*rows-theme.Padding())
	}
	return fyne.NewSize(0, 0)
}

func (l *gridWrapListLayout) getItem() fyne.CanvasObject {
	item := l.itemPool.Obtain()
	if item == nil {
		if f := l.list.CreateItem; f != nil {
			item = f()
		}
	}
	return item
}

func (l *gridWrapListLayout) offsetUpdated(pos fyne.Position) {
	if l.list.offsetY == pos.Y {
		return
	}
	l.list.offsetY = pos.Y
	l.updateList(false)
}

func (l *gridWrapListLayout) setupListItem(li fyne.CanvasObject, id GridWrapListItemID) {
	if f := l.list.UpdateItem; f != nil {
		f(id, li)
	}
}

func (l *GridWrapList) getColCount() int {
	colCount := 1
	width := l.Size().Width
	if width > l.itemMin.Width {
		colCount = int(math.Floor(float64(width+theme.Padding()) / float64(l.itemMin.Width+theme.Padding())))
	}
	return colCount
}

func (l *gridWrapListLayout) updateList(refresh bool) {
	// code here is a mashup of listLayout.updateList and gridWrapLayout.Layout

	l.renderLock.Lock()
	defer l.renderLock.Unlock()
	length := 0
	if f := l.list.Length; f != nil {
		length = f()
	}

	colCount := l.list.getColCount()
	visibleRowsCount := int(math.Ceil(float64(l.list.scroller.Size().Height)/float64(l.list.itemMin.Height+theme.Padding()))) + 1

	offY := l.list.offsetY - float32(math.Mod(float64(l.list.offsetY), float64(l.list.itemMin.Height+theme.Padding())))
	minRow := int(offY / (l.list.itemMin.Height + theme.Padding()))
	minItem := GridWrapListItemID(minRow * colCount)
	maxRow := int(math.Min(float64(minRow+visibleRowsCount), math.Ceil(float64(length)/float64(colCount))))
	maxItem := GridWrapListItemID(math.Min(float64(maxRow*colCount), float64(length-1)))

	if l.list.UpdateItem == nil {
		fyne.LogError("Missing UpdateCell callback required for GridWrapList", nil)
	}

	wasVisible := l.visible
	l.visible = make(map[GridWrapListItemID]fyne.CanvasObject)
	var cells []fyne.CanvasObject
	y := offY
	curItem := minItem
	for row := minRow; row <= maxRow && curItem <= maxItem; row++ {
		x := float32(0)
		for col := 0; col < colCount && curItem <= maxItem; col++ {
			c, ok := wasVisible[curItem]
			if !ok {
				c = l.getItem()
				if c == nil {
					continue
				}
				c.Resize(l.list.itemMin)
				l.setupListItem(c, curItem)
			}

			c.Move(fyne.NewPos(x, y))
			if refresh {
				c.Resize(l.list.itemMin)
				if ok { // refresh visible
					l.setupListItem(c, curItem)
				}
			}

			x += l.list.itemMin.Width + theme.Padding()
			l.visible[curItem] = c
			cells = append(cells, c)
			curItem++
		}
		y += l.list.itemMin.Height + theme.Padding()
	}

	for id, old := range wasVisible {
		if _, ok := l.visible[id]; !ok {
			l.itemPool.Release(old)
		}
	}
	l.children = cells

	objects := l.children
	l.list.scroller.Content.(*fyne.Container).Objects = objects
}

type pool interface {
	Obtain() fyne.CanvasObject
	Release(fyne.CanvasObject)
}

var _ pool = (*syncPool)(nil)

type syncPool struct {
	sync.Pool
}

// Obtain returns an item from the pool for use
func (p *syncPool) Obtain() (item fyne.CanvasObject) {
	o := p.Get()
	if o != nil {
		item = o.(fyne.CanvasObject)
	}
	return
}

// Release adds an item into the pool to be used later
func (p *syncPool) Release(item fyne.CanvasObject) {
	p.Put(item)
}
