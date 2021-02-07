package widget_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/x/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestFileTree(t *testing.T) {
	tree := &widget.FileTree{}
	tree.Refresh() // Should not crash
}

func TestFileTree_Layout(t *testing.T) {
	test.NewApp()

	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	root := storage.NewURI("file://" + tempDir)
	tree := widget.NewFileTree(root)
	tree.OpenAllBranches()

	branch, err := storage.Child(root, "B")
	assert.NoError(t, err)
	leaf, err := storage.Child(branch, "C")
	assert.NoError(t, err)
	tree.Select(leaf.String())

	window := test.NewWindow(tree)
	defer window.Close()
	window.Resize(fyne.NewSize(300, 100))
	test.AssertImageMatches(t, "filetree/selected.png", window.Canvas().Capture())
}

func Test_NewFileTree(t *testing.T) {
	test.NewApp()

	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	root := storage.NewURI("file://" + tempDir)
	tree := widget.NewFileTree(root)
	tree.OpenAllBranches()

	assert.True(t, tree.IsBranchOpen(root.String()))
	branch1, err := storage.Child(root, "A")
	assert.NoError(t, err)
	assert.True(t, tree.IsBranchOpen(branch1.String()))
	branch2, err := storage.Child(root, "B")
	assert.NoError(t, err)
	assert.True(t, tree.IsBranchOpen(branch2.String()))
	leaf, err := storage.Child(branch2, "C")
	assert.NoError(t, err)
	assert.False(t, tree.IsBranchOpen(leaf.String()))
}

func createTempDir(t *testing.T) string {
	t.Helper()
	tempDir, err := ioutil.TempDir("", "test")
	assert.NoError(t, err)
	err = os.MkdirAll(path.Join(tempDir, "A"), os.ModePerm)
	assert.NoError(t, err)
	err = os.MkdirAll(path.Join(tempDir, "B"), os.ModePerm)
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join(tempDir, "B", "C"), []byte("c"), os.ModePerm)
	assert.NoError(t, err)
	return tempDir
}
