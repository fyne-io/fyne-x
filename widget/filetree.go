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
	Filter storage.FileFilter
	Sorter func(fyne.URI, fyne.URI) bool

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
		listableCache: make(map[widget.TreeNodeID]fyne.ListableURI),
		uriCache:      make(map[widget.TreeNodeID]fyne.URI),
	}
	tree.IsBranch = func(id widget.TreeNodeID) bool {
		if _, ok := tree.listableCache[id]; ok {
			return true
		}
		var err error
		uri, ok := tree.uriCache[id]
		if !ok {
			uri, err = storage.ParseURI(id)
			if err != nil {
				fyne.LogError("Unable to parse URI", err)
				return false
			}
			tree.uriCache[id] = uri
		}
		if luri, err := storage.ListerForURI(uri); err == nil {
			tree.listableCache[id] = luri
			return true
		}
		return false
	}
	tree.ChildUIDs = func(id widget.TreeNodeID) (c []string) {
		var err error
		luri, ok := tree.listableCache[id]
		if !ok {
			uri, ok := tree.uriCache[id]
			if !ok {
				uri, err = storage.ParseURI(id)
				if err != nil {
					fyne.LogError("Unable to parse URI", err)
					return
				}
				tree.uriCache[id] = uri
			}

			l, err := storage.ListerForURI(uri)
			if err != nil {
				fyne.LogError("Unable to get lister for "+id, err)
				return
			}
			luri = l
			tree.listableCache[id] = l
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
				if filter.Matches(u) {
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
		var err error
		uri, ok := tree.uriCache[id]
		if !ok {
			uri, err = storage.ParseURI(id)
			if err != nil {
				fyne.LogError("Unable to parse URI", err)
				return
			}
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
			c.Objects[1].(*widget.Icon).SetResource(r)
		} else {
			// Set file uri to update icon
			c.Objects[1].(*widget.FileIcon).SetURI(uri)
		}

		l := c.Objects[0].(*widget.Label)
		if tree.Root == id {
			l.SetText(id)
		} else {
			l.SetText(uri.Name())
		}
	}
	tree.ExtendBaseWidget(tree)
	return tree
}