package widget

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// FileTree extends widget.Tree to display a file system hierarchy.
type FileTree struct {
	widget.Tree
	Filter       storage.FileFilter
	ShowRootPath bool
	Sorter       func(fyne.URI, fyne.URI) bool

	listCache     map[widget.TreeNodeID][]widget.TreeNodeID
	listableCache map[widget.TreeNodeID]fyne.ListableURI
	uriCache      map[widget.TreeNodeID]fyne.URI
}

// NewFileTree creates a new FileTree from the given root URI.
func NewFileTree(root fyne.URI) *FileTree {
	tree := &FileTree{
		Tree: widget.Tree{
			Root: root.String(),
			CreateNode: func(branch bool) fyne.CanvasObject {
				var icon fyne.CanvasObject
				if branch {
					icon = widget.NewIcon(nil)
				} else {
					icon = widget.NewFileIcon(nil)
				}
				return container.NewBorder(nil, nil, icon, nil, widget.NewLabel("Template Object"))
			},
		},
		listCache:     make(map[widget.TreeNodeID][]widget.TreeNodeID),
		listableCache: make(map[widget.TreeNodeID]fyne.ListableURI),
		uriCache:      make(map[widget.TreeNodeID]fyne.URI),
	}
	tree.IsBranch = func(id widget.TreeNodeID) bool {
		_, err := tree.toListable(id)
		return err == nil
	}
	tree.ChildUIDs = func(id widget.TreeNodeID) (c []string) {
		listable, err := tree.toListable(id)
		if err != nil {
			fyne.LogError("Unable to get lister for "+id, err)
			return
		}

		ids, ok := tree.listCache[id]
		if ok {
			return ids
		}

		uris, err := listable.List()
		if err != nil {
			fyne.LogError("Unable to list "+listable.String(), err)
			return
		}

		for _, u := range tree.sort(tree.filter(uris)) {
			// Convert to String
			c = append(c, u.String())
		}

		tree.listCache[id] = c
		return
	}
	tree.UpdateNode = func(id widget.TreeNodeID, branch bool, node fyne.CanvasObject) {
		uri, err := tree.toURI(id)
		if err != nil {
			fyne.LogError("Unable to parse URI", err)
			return
		}

		c := node.(*fyne.Container)
		if branch {
			var r fyne.Resource
			if tree.IsBranchOpen(id) {
				// Set open folder icon
				r = theme.FolderOpenIcon()
			} else {
				// Set folder icon
				r = theme.FolderIcon()
			}
			c.Objects[1].(*widget.Icon).SetResource(r)
		} else {
			// Set file uri to update icon
			c.Objects[1].(*widget.FileIcon).SetURI(uri)
		}

		var l string
		if tree.Root == id && tree.ShowRootPath {
			l = id
		} else {
			l = uri.Name()
		}
		c.Objects[0].(*widget.Label).SetText(l)
	}

	// reset sorted child ID cache if the branch is closed - in the future we do FS watch
	tree.OnBranchClosed = func(id widget.TreeNodeID) {
		delete(tree.listCache, id)
	}

	tree.ExtendBaseWidget(tree)
	return tree
}

func (t *FileTree) filter(uris []fyne.URI) []fyne.URI {
	filter := t.Filter
	if filter == nil {
		return uris
	}
	var filtered []fyne.URI
	for _, u := range uris {
		if filter.Matches(u) {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

func (t *FileTree) sort(uris []fyne.URI) []fyne.URI {
	if sorter := t.Sorter; sorter != nil {
		sort.Slice(uris, func(i, j int) bool {
			return sorter(uris[i], uris[j])
		})
	}
	return uris
}

func (t *FileTree) toListable(id widget.TreeNodeID) (fyne.ListableURI, error) {
	listable, ok := t.listableCache[id]
	if ok {
		return listable, nil
	}
	uri, err := t.toURI(id)
	if err != nil {
		return nil, err
	}

	listable, err = storage.ListerForURI(uri)
	if err != nil {
		return nil, err
	}
	t.listableCache[id] = listable
	return listable, nil
}

func (t *FileTree) toURI(id widget.TreeNodeID) (fyne.URI, error) {
	uri, ok := t.uriCache[id]
	if ok {
		return uri, nil
	}

	uri, err := storage.ParseURI(id)
	if err != nil {
		return nil, err
	}
	t.uriCache[id] = uri
	return uri, nil
}
