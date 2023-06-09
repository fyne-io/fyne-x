package widget

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestNewPagination(t *testing.T) {
	p := NewPagination(&PaginationConfig{
		Page:     1,
		PageSize: 10,
	})
	_ = test.WidgetRenderer(p)
	e, _ := p.Objects[1].(*widget.Entry)
	assert.Equal(t, "1", e.Text)
	e, _ = p.Objects[4].(*widget.Entry)
	assert.Equal(t, "10", e.Text)
}

func TestNewPagination_Previous(t *testing.T) {
	p := NewPagination(&PaginationConfig{
		Page:     2,
		PageSize: 10,
	})
	p.SetTotalRows(100)
	_ = test.WidgetRenderer(p)
	btn, _ := p.Objects[0].(*widget.Button)

	assert.Equal(t, 2, p.GetPage())
	test.Tap(btn)
	assert.Equal(t, 1, p.GetPage())
	test.Tap(btn)
	assert.Equal(t, 1, p.GetPage())
}

func TestNewPagination_Next(t *testing.T) {
	p := NewPagination(&PaginationConfig{
		Page:     2,
		PageSize: 10,
	})
	p.SetTotalRows(25)

	_ = test.WidgetRenderer(p)
	btn, _ := p.Objects[2].(*widget.Button)

	assert.Equal(t, 2, p.GetPage())
	test.Tap(btn)
	assert.Equal(t, 3, p.GetPage())
	test.Tap(btn)
	assert.Equal(t, 3, p.GetPage())
}

func TestNewPagination_Page(t *testing.T) {
	p := NewPagination(&PaginationConfig{
		Page:     1,
		PageSize: 10,
	})
	p.SetTotalRows(5)

	_ = test.WidgetRenderer(p)
	entry, _ := p.Objects[1].(*widget.Entry)
	entry.SetText("2")

	assert.NotNil(t, entry.Validate())
}

func TestNewPagination_PageSize(t *testing.T) {
	p := NewPagination(&PaginationConfig{
		Page:     1,
		PageSize: 10,
	})
	p.SetTotalRows(15)

	_ = test.WidgetRenderer(p)
	pageSize, _ := p.Objects[4].(*widget.Entry)
	pageSize.SetText("20")

	assert.Equal(t, 1, p.GetTotalPages())
}
