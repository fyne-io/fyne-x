package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"sort"
)

type FileTree struct {
	widget.Tree
	Filter func(fyne.URI) bool
	Sorter func(fyne.URI, fyne.URI) bool

	luriCache map[widget.TreeNodeID]fyne.ListableURI
	uriCache  map[widget.TreeNodeID]fyne.URI
}

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
				return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), icon, widget.NewLabel("Template Object"))
			},
		},
		luriCache: make(map[widget.TreeNodeID]fyne.ListableURI),
		uriCache:  make(map[widget.TreeNodeID]fyne.URI),
	}
	tree.IsBranch = func(id widget.TreeNodeID) bool {
		if _, ok := tree.luriCache[id]; ok {
			return true
		}
		uri, ok := tree.uriCache[id]
		if !ok {
			uri = storage.NewURI(id)
			tree.uriCache[id] = uri
		}
		if luri, err := storage.ListerForURI(uri); err == nil {
			tree.luriCache[id] = luri
			return true
		}
		return false
	}
	tree.ChildUIDs = func(id widget.TreeNodeID) (c []string) {
		luri, ok := tree.luriCache[id]
		if !ok {
			uri, ok := tree.uriCache[id]
			if !ok {
				uri = storage.NewURI(id)
				tree.uriCache[id] = uri
			}

			l, err := storage.ListerForURI(uri)
			if err != nil {
				fyne.LogError("Unable to get lister for "+id, err)
				return
			} else {
				luri = l
				tree.luriCache[id] = l
			}
		}

		uris, err := luri.List()
		if err != nil {
			fyne.LogError("Unable to list "+luri.String(), err)
			return
		}

		var us []fyne.URI
		// Filter URIs
		if filter := tree.Filter; filter == nil {
			us = uris
		} else {
			for _, u := range uris {
				if filter(u) {
					us = append(us, u)
				}
			}
		}

		// Sort URIs
		if sorter := tree.Sorter; sorter != nil {
			sort.Slice(us, func(i, j int) bool {
				return sorter(us[i], us[j])
			})
		}

		// Convert to Strings
		for _, u := range us {
			c = append(c, u.String())
		}
		return
	}
	tree.UpdateNode = func(id widget.TreeNodeID, branch bool, node fyne.CanvasObject) {
		uri, ok := tree.uriCache[id]
		if !ok {
			uri = storage.NewURI(id)
			tree.uriCache[id] = uri
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
			c.Objects[0].(*widget.Icon).SetResource(r)
		} else {
			// Set file uri to update icon
			c.Objects[0].(*widget.FileIcon).SetURI(uri)
		}

		l := c.Objects[1].(*widget.Label)
		if tree.Root == id {
			l.SetText(id)
		} else {
			l.SetText(uri.Name())
		}
	}
	tree.ExtendBaseWidget(tree)
	return tree
}
