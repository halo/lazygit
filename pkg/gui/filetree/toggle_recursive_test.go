package filetree

import (
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestToggleSubtreeExpansion_FileTree(t *testing.T) {
	files := []*models.File{
		{Path: "a/b/c/file1"},
		{Path: "a/b/d/file2"},
		{Path: "a/e/file3"},
	}

	cmn := common.NewDummyCommon()
	ft := NewFileTree(func() []*models.File { return files }, cmn, true)
	ft.SetTree()

	// ensure initially nothing is collapsed
	if ft.IsCollapsed("./a") || ft.IsCollapsed("./a/b") || ft.IsCollapsed("./a/b/c") {
		t.Fatalf("expected no collapsed paths initially")
	}

	// collapse recursively at ./a
	ft.ToggleSubtreeExpansion("./a")

	// now a and descendants should be collapsed
	assert.True(t, ft.IsCollapsed("./a"))
	assert.True(t, ft.IsCollapsed("./a/b"))
	assert.True(t, ft.IsCollapsed("./a/b/c"))

	// toggle again should expand them
	ft.ToggleSubtreeExpansion("./a")
	assert.False(t, ft.IsCollapsed("./a"))
	assert.False(t, ft.IsCollapsed("./a/b"))
	assert.False(t, ft.IsCollapsed("./a/b/c"))
}

func TestFileTreeViewModel_ToggleSubtreeExpansion_PreservesSelection(t *testing.T) {
	files := []*models.File{
		{Path: "a/b/c/file1"},
		{Path: "a/b/d/file2"},
		{Path: "a/e/file3"},
	}

	// create view model and set tree
	cmn := common.NewDummyCommon()
	vm := NewFileTreeViewModel(func() []*models.File { return files }, cmn, true)
	vm.SetTree()

	// select the deepest file
	idx, found := vm.GetIndexForPath("a/b/c/file1")
	if !found {
		// try with root prefix
		idx, found = vm.GetIndexForPath("./a/b/c/file1")
	}
	if !found {
		t.Fatalf("expected to find path")
	}
	vm.SetSelectedLineIdx(idx)

	// toggle recursively on ./a/b
	vm.ToggleSubtreeExpansion("./a/b")

	// selection should either be preserved (if still visible) or be an ancestor
	selected := vm.GetSelected()
	if selected == nil {
		t.Fatalf("expected a selected node")
	}

	// normalize paths by trimming a possible leading "./"
	normalize := func(p string) string { return strings.TrimPrefix(p, "./") }
	selPath := normalize(selected.GetPath())

	// After collapsing ./a/b, the selected node must not be inside the collapsed subtree
	if strings.HasPrefix(selPath, "a/b/") {
		t.Fatalf("expected selected to not be inside the collapsed subtree, got %s", selPath)
	}
}
