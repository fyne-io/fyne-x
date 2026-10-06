// Package menu provides community extensions for building Fyne menus.
package menu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

const (
	recentsKey = "fyne-x.menu.recents"
	recentsMax = 5
)

// Recents manages a menu item containing the most recently used items for an app.
// The list is ordered by recency and saved in the app preferences so that it is
// available the next time the app is run.
type Recents struct {
	item     *fyne.MenuItem
	uris     []fyne.URI
	onOpen   func(fyne.URI)
	getLabel func(fyne.URI) string
}

// NewRecents creates a new recent items manager that will show a menu item with the given label.
// The onOpen callback is called with the URI of a recent item when the user chooses it from the menu.
// The list of items is loaded from the preferences of the current app, so the app must have a unique ID.
func NewRecents(label string, onOpen func(fyne.URI)) *Recents {
	r := &Recents{onOpen: onOpen}
	r.item = fyne.NewMenuItem(label, nil)
	r.item.ChildMenu = fyne.NewMenu("")

	for _, s := range fyne.CurrentApp().Preferences().StringList(recentsKey) {
		u, err := storage.ParseURI(s)
		if err != nil {
			continue
		}
		r.uris = append(r.uris, u)
	}
	if len(r.uris) > recentsMax {
		r.uris = r.uris[:recentsMax]
	}

	r.update()
	return r
}

// Add marks the given URI as the most recently used item.
// If the item was already in the list it is moved to the top instead of being duplicated,
// and only the most recent items are kept. The updated list is saved to the app preferences.
func (r *Recents) Add(u fyne.URI) {
	if u == nil {
		return
	}

	uris := []fyne.URI{u}
	for _, old := range r.uris {
		if old.String() == u.String() {
			continue
		}
		uris = append(uris, old)
	}
	if len(uris) > recentsMax {
		uris = uris[:recentsMax]
	}
	r.uris = uris

	saved := make([]string, len(uris))
	for i, u := range uris {
		saved[i] = u.String()
	}
	fyne.CurrentApp().Preferences().SetStringList(recentsKey, saved)

	r.update()
	r.refresh()
}

// MenuItem returns the menu item that should be added to a menu to display the recent items.
// It contains a child menu listing each of the items and is disabled when there are none.
func (r *Recents) MenuItem() *fyne.MenuItem {
	return r.item
}

// SetItemLabel sets a function used to generate the label for each recent item.
// By default the name of the URI is used, with parent directories added when two names match.
func (r *Recents) SetItemLabel(fn func(fyne.URI) string) {
	r.getLabel = fn

	r.update()
	r.refresh()
}

// URIs returns the recent items, with the most recently used first.
func (r *Recents) URIs() []fyne.URI {
	return append([]fyne.URI(nil), r.uris...)
}

func (r *Recents) refresh() {
	for _, w := range fyne.CurrentApp().Driver().AllWindows() {
		main := w.MainMenu()
		if main == nil {
			continue
		}

		for _, m := range main.Items {
			if menuContains(m, r.item) {
				w.SetMainMenu(main)
				break
			}
		}
	}
}

func (r *Recents) update() {
	labels := r.labels()
	items := make([]*fyne.MenuItem, len(r.uris))
	for i, u := range r.uris {
		items[i] = fyne.NewMenuItem(labels[i], func() {
			if r.onOpen != nil {
				r.onOpen(u)
			}
		})
	}

	r.item.ChildMenu.Items = items
	r.item.Disabled = len(items) == 0
}

// labels returns the label to show for each of the recent URIs.
// When the default name labels would be ambiguous the parent directories are
// included, such as "dir/file.txt", until the labels are unique.
func (r *Recents) labels() []string {
	labels := make([]string, len(r.uris))
	if r.getLabel != nil {
		for i, u := range r.uris {
			labels[i] = r.getLabel(u)
		}
		return labels
	}

	parents := make([]fyne.URI, len(r.uris))
	for i, u := range r.uris {
		labels[i] = u.Name()
		parents[i] = u
	}

	for {
		counts := make(map[string]int)
		for _, l := range labels {
			counts[l]++
		}

		changed := false
		for i, l := range labels {
			if counts[l] < 2 || parents[i] == nil {
				continue
			}

			parent, err := storage.Parent(parents[i])
			if err != nil || parent == nil || parent.Name() == "" || parent.Name() == "/" {
				parents[i] = nil
				continue
			}

			labels[i] = parent.Name() + "/" + l
			parents[i] = parent
			changed = true
		}

		if !changed {
			return labels
		}
	}
}

func menuContains(m *fyne.Menu, item *fyne.MenuItem) bool {
	for _, i := range m.Items {
		if i == item {
			return true
		}
		if i.ChildMenu != nil && menuContains(i.ChildMenu, item) {
			return true
		}
	}
	return false
}
