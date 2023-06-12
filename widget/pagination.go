package widget

import (
	"errors"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type paginationLayout struct{}

var _ fyne.Layout = (*paginationLayout)(nil)

var firstPage = 1

// MinSize set the minimal size for pagination widget
func (p paginationLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	maxHeight := float32(0)
	totalWidth := float32(0)
	for _, o := range objects {
		maxHeight = fyne.Max(o.MinSize().Height, maxHeight)
	}
	return fyne.NewSize(totalWidth, maxHeight)
}

// Layout is called to put all child objects into same row
func (p paginationLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	x := float32(0.0)
	y := float32(0.0)

	for _, o := range objects {
		o.Move(fyne.NewPos(x, y))
		width := o.MinSize().Width
		o.Resize(fyne.NewSize(width, size.Height))
		x += width
	}
}

// Pagination create a pagination widget
type Pagination struct {
	widget.BaseWidget
	page            binding.Int
	pageSize        binding.Int
	totalPages      binding.Int
	totalRows       binding.Int
	OnChange        func(page, pageSize int)
	Objects         []fyne.CanvasObject
	defaultPageSize int
}

// Reset will set page and pageSize to initial config
func (p *Pagination) Reset() {
	p.page.Set(firstPage)
	p.pageSize.Set(p.defaultPageSize)
}

func (p *Pagination) SetTotalRows(total int) {
	p.Reset()
	p.totalRows.Set(total)
	p.setTotalPages()
	p.onSubmit()
	p.Refresh()
}

func (p *Pagination) GetPage() int {
	page, _ := p.page.Get()
	return page
}

func (p *Pagination) GetPageSize() int {
	size, _ := p.pageSize.Get()
	return size
}

func (p *Pagination) GetTotalPages() int {
	pages, _ := p.totalPages.Get()
	return pages
}

func (p *Pagination) setTotalPages() {
	size, _ := p.pageSize.Get()
	rows, _ := p.totalRows.Get()
	pages := int(math.Ceil(float64(rows) / float64(size)))
	p.totalPages.Set(pages)
}

func (p *Pagination) handlePreClick() {
	current, _ := p.page.Get()
	if current == 1 {
		return
	}
	current -= 1
	p.page.Set(current)
	p.onSubmit()
}

func (p *Pagination) handleNextClick() {
	current, _ := p.page.Get()
	total, _ := p.totalPages.Get()
	if current >= total {
		return
	}
	current += 1
	p.page.Set(current)
	p.onSubmit()
}

func (p *Pagination) onSubmit() {
	if p.OnChange != nil {
		page, _ := p.page.Get()
		pSize, _ := p.pageSize.Get()
		p.OnChange(page, pSize)
	}
}

func (p *Pagination) CreateRenderer() fyne.WidgetRenderer {
	pre := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), p.handlePreClick)
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), p.handleNextClick)

	curPageStr := binding.IntToString(p.page)
	curPage := widget.NewEntryWithData(curPageStr)
	curPage.Validator = func(s string) error {
		targetPage, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		totalPages := p.GetTotalPages()
		if targetPage < 1 {
			return errors.New("page should not be smaller than 1")
		}
		if targetPage > totalPages {
			return errors.New("page should not bigger than total page")
		}
		return nil
	}
	curPage.OnSubmitted = func(s string) {
		if err := curPage.Validate(); err == nil {
			p.onSubmit()
		}
	}

	pageSizeText := widget.NewLabel("Size")

	pageSizeStr := binding.IntToString(p.pageSize)
	pageSizeValue := widget.NewEntryWithData(pageSizeStr)
	pageSizeValue.Validator = func(s string) error {
		target, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if target < 1 {
			return errors.New("page size should be bigger than 1")
		}
		return nil
	}
	pageSizeValue.OnSubmitted = func(_ string) {
		if err := pageSizeValue.Validate(); err == nil {
			p.setTotalPages()
			p.page.Set(firstPage)
			p.onSubmit()
		}
	}

	totalPages := widget.NewLabel("TotalPages")

	totalPagesStr := binding.IntToString(p.totalPages)
	totalPageValue := widget.NewLabelWithData(totalPagesStr)

	totalRowsText := widget.NewLabel("TotalRows")

	totalRowsStr := binding.IntToString(p.totalRows)
	totalRowsValue := widget.NewLabelWithData(totalRowsStr)
	objects := []fyne.CanvasObject{
		pre, curPage, next,
		pageSizeText, pageSizeValue,
		totalPages, totalPageValue,
		totalRowsText, totalRowsValue,
	}
	p.Objects = objects

	c := container.New(&paginationLayout{}, objects...)
	return widget.NewSimpleRenderer(c)
}

func NewPagination(pageSize int) *Pagination {
	p := &Pagination{
		page:            binding.NewInt(),
		pageSize:        binding.NewInt(),
		totalPages:      binding.NewInt(),
		totalRows:       binding.NewInt(),
		defaultPageSize: pageSize,
	}
	p.Reset()
	p.ExtendBaseWidget(p)
	return p
}
