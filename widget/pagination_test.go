package widget

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestNewPagination(t *testing.T) {
	p := NewPagination(10)
	_ = test.WidgetRenderer(p)
	assert.Equal(t, 1, p.GetPage())
	assert.Equal(t, 10, p.GetPageSize())
}

func TestPagination_PrevButton(t *testing.T) {
	p := NewPagination(10)
	p.SetTotalRows(100)
	_ = test.WidgetRenderer(p)
	btn, _ := p.Objects[0].(*widget.Button)

	test.Tap(btn)
	assert.Equal(t, 1, p.GetPage())
}

func TestPagination_NextButton(t *testing.T) {
	p := NewPagination(10)
	p.SetTotalRows(15)

	_ = test.WidgetRenderer(p)
	btn, _ := p.Objects[2].(*widget.Button)

	test.Tap(btn)
	assert.Equal(t, 2, p.GetPage())
	test.Tap(btn)
	assert.Equal(t, 2, p.GetPage())
}
