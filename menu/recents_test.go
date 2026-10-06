package menu

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestNewRecents_Empty(t *testing.T) {
	test.NewTempApp(t)

	r := NewRecents("Recent", nil)
	assert.Equal(t, "Recent", r.MenuItem().Label)
	assert.True(t, r.MenuItem().Disabled)
	assert.Empty(t, r.MenuItem().ChildMenu.Items)
	assert.Empty(t, r.URIs())
}

func TestRecents_Add(t *testing.T) {
	test.NewTempApp(t)

	r := NewRecents("Recent", nil)
	r.Add(storage.NewFileURI("/tmp/one.txt"))
	r.Add(storage.NewFileURI("/tmp/two.txt"))

	assert.False(t, r.MenuItem().Disabled)
	assert.Equal(t, []string{"two.txt", "one.txt"}, labels(r))
}

func TestRecents_Add_Duplicate(t *testing.T) {
	test.NewTempApp(t)

	r := NewRecents("Recent", nil)
	r.Add(storage.NewFileURI("/tmp/one.txt"))
	r.Add(storage.NewFileURI("/tmp/two.txt"))
	r.Add(storage.NewFileURI("/tmp/three.txt"))
	r.Add(storage.NewFileURI("/tmp/one.txt"))

	assert.Equal(t, []string{"one.txt", "three.txt", "two.txt"}, labels(r))
}

func TestRecents_Add_Limit(t *testing.T) {
	test.NewTempApp(t)

	r := NewRecents("Recent", nil)
	for i := 1; i <= 7; i++ {
		r.Add(storage.NewFileURI(fmt.Sprintf("/tmp/%d.txt", i)))
	}

	assert.Equal(t, []string{"7.txt", "6.txt", "5.txt", "4.txt", "3.txt"}, labels(r))
	assert.Len(t, fyne.CurrentApp().Preferences().StringList(recentsKey), 5)
}

func TestRecents_Add_Refresh(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)

	r := NewRecents("Recent", nil)
	w.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File", fyne.NewMenuItem("Open", nil), r.MenuItem())))

	r.Add(storage.NewFileURI("/tmp/one.txt"))
	recent := w.MainMenu().Items[0].Items[1]
	assert.False(t, recent.Disabled)
	assert.Len(t, recent.ChildMenu.Items, 1)
}

func TestRecents_Open(t *testing.T) {
	test.NewTempApp(t)

	var opened fyne.URI
	r := NewRecents("Recent", func(u fyne.URI) {
		opened = u
	})
	r.Add(storage.NewFileURI("/tmp/one.txt"))
	r.Add(storage.NewFileURI("/tmp/two.txt"))

	r.MenuItem().ChildMenu.Items[1].Action()
	assert.Equal(t, "file:///tmp/one.txt", opened.String())
}

func TestRecents_Persist(t *testing.T) {
	test.NewTempApp(t)

	r := NewRecents("Recent", nil)
	r.Add(storage.NewFileURI("/tmp/one.txt"))
	r.Add(storage.NewFileURI("/tmp/two.txt"))

	assert.Equal(t, []string{"file:///tmp/two.txt", "file:///tmp/one.txt"},
		fyne.CurrentApp().Preferences().StringList(recentsKey))

	r = NewRecents("Recent", nil)
	assert.False(t, r.MenuItem().Disabled)
	assert.Equal(t, []string{"two.txt", "one.txt"}, labels(r))
}

func TestRecents_Add_SameName(t *testing.T) {
	test.NewTempApp(t)

	r := NewRecents("Recent", nil)
	r.Add(storage.NewFileURI("/tmp/first/pres.md"))
	r.Add(storage.NewFileURI("/tmp/other.md"))
	r.Add(storage.NewFileURI("/tmp/second/pres.md"))
	assert.Equal(t, []string{"second/pres.md", "other.md", "first/pres.md"}, labels(r))

	// parent dirs share a name too, so go up another level
	r.Add(storage.NewFileURI("/home/work/docs/pres.md"))
	r.Add(storage.NewFileURI("/home/play/docs/pres.md"))
	assert.Equal(t, []string{"play/docs/pres.md", "work/docs/pres.md", "second/pres.md", "other.md", "first/pres.md"}, labels(r))
}

func TestRecents_SetItemLabel(t *testing.T) {
	test.NewTempApp(t)

	r := NewRecents("Recent", nil)
	r.Add(storage.NewFileURI("/tmp/first/pres.md"))
	r.Add(storage.NewFileURI("/tmp/second/pres.md"))
	assert.Equal(t, []string{"second/pres.md", "first/pres.md"}, labels(r))

	r.SetItemLabel(func(u fyne.URI) string {
		parent, _ := storage.Parent(u)
		return parent.Name()
	})
	assert.Equal(t, []string{"second", "first"}, labels(r))
}

func labels(r *Recents) []string {
	var ret []string
	for _, i := range r.MenuItem().ChildMenu.Items {
		ret = append(ret, i.Label)
	}
	return ret
}
