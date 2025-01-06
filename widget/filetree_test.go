package widget

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestFileTree(t *testing.T) {
	tree := &FileTree{}
	tree.Refresh() // Should not crash
}

func TestFileTree_Layout(t *testing.T) {
	test.NewApp()

	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	root, err := storage.ParseURI("file://" + tempDir)
	assert.NoError(t, err)
	tree := NewFileTree(root)
	tree.OpenAllBranches()

	window := test.NewWindow(tree)
	defer window.Close()
	window.Resize(fyne.NewSize(300, 100))

	branch, err := storage.Child(root, "B")
	assert.NoError(t, err)
	leaf, err := storage.Child(branch, "C.txt")
	assert.NoError(t, err)
	tree.Select(leaf.String())

	test.AssertImageMatches(t, "filetree/selected.png", window.Canvas().Capture())
}

func TestFileTree_filter(t *testing.T) {
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	root, err := storage.ParseURI("file://" + tempDir)
	assert.NoError(t, err)
	tree := NewFileTree(root)
	tree.Filter = storage.NewExtensionFileFilter([]string{".txt"})

	branch1, err := storage.Child(root, "A")
	assert.NoError(t, err)
	branch2, err := storage.Child(root, "B")
	assert.NoError(t, err)
	leaf1, err := storage.Child(branch2, "C.txt")
	assert.NoError(t, err)
	leaf2, err := storage.Child(branch2, "D.txt")
	assert.NoError(t, err)

	given := []fyne.URI{
		branch1,
		branch2,
		leaf1,
		leaf2,
	}

	expected := []fyne.URI{
		leaf1,
		leaf2,
	}

	assert.Equal(t, expected, tree.filter(given))
}

func TestFileTree_ShowRootPath(t *testing.T) {
	testPath, _ := filepath.Abs("./testdata")
	root, err := storage.ParseURI("file://" + testPath)
	assert.NoError(t, err)

	tree := NewFileTree(root)
	firstNodeContent := func() *widget.Label {
		renderer := tree.CreateRenderer()
		assert.Equal(t, 1, len(renderer.Objects()))
		content := renderer.Objects()[0].(*container.Scroll).Content.(fyne.Widget).CreateRenderer()
		content.Layout(fyne.NewSize(100, 100))

		node := content.Objects()[0].(fyne.Widget).CreateRenderer()
		return node.Objects()[1].(*fyne.Container).Objects[0].(*widget.Label)
	}

	assert.Equal(t, "testdata", firstNodeContent().Text)

	tree.ShowRootPath = true
	tree.Refresh()
	assert.Equal(t, "file://", firstNodeContent().Text[:7])
}

func TestFileTree_sort(t *testing.T) {
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	root, err := storage.ParseURI("file://" + tempDir)
	assert.NoError(t, err)
	tree := NewFileTree(root)
	tree.Sorter = func(u1, u2 fyne.URI) bool {
		return u2.String() < u1.String() // Reverse alphabetical
	}

	branch1, err := storage.Child(root, "A")
	assert.NoError(t, err)
	branch2, err := storage.Child(root, "B")
	assert.NoError(t, err)
	leaf1, err := storage.Child(branch2, "C.txt")
	assert.NoError(t, err)
	leaf2, err := storage.Child(branch2, "D.txt")
	assert.NoError(t, err)

	given := []fyne.URI{
		branch1,
		branch2,
		leaf1,
		leaf2,
	}

	expected := []fyne.URI{
		leaf2,
		leaf1,
		branch2,
		branch1,
	}

	assert.Equal(t, expected, tree.sort(given))
}

func Test_NewFileTree(t *testing.T) {
	test.NewApp()

	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	root, err := storage.ParseURI("file://" + tempDir)
	assert.NoError(t, err)
	tree := NewFileTree(root)
	tree.OpenAllBranches()

	assert.True(t, tree.IsBranchOpen(root.String()))
	branch1, err := storage.Child(root, "A")
	assert.NoError(t, err)
	assert.True(t, tree.IsBranchOpen(branch1.String()))
	branch2, err := storage.Child(root, "B")
	assert.NoError(t, err)
	assert.True(t, tree.IsBranchOpen(branch2.String()))
	leaf, err := storage.Child(branch2, "C.txt")
	assert.NoError(t, err)
	assert.False(t, tree.IsBranchOpen(leaf.String()))
}

func createTempDir(t *testing.T) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	err = os.MkdirAll(path.Join(tempDir, "A"), os.ModePerm)
	assert.NoError(t, err)
	err = os.MkdirAll(path.Join(tempDir, "B"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(path.Join(tempDir, "B", "C.txt"), []byte("c"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(path.Join(tempDir, "B", "D.txt"), []byte("d"), os.ModePerm)
	assert.NoError(t, err)
	return tempDir
}
