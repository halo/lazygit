package filetree

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestFileTreeViewModel_ToggleSubtreeExpansion_Recursive(t *testing.T) {
	files := []*models.File{
		{Path: "a/b/c/file1"},
		{Path: "a/b/d/file2"},
		{Path: "a/e/file3"},
	}

	cmn := common.NewDummyCommon()
	vm := NewFileTreeViewModel(func() []*models.File { return files }, cmn, true)
	vm.SetTree()

	// ensure initially nothing is collapsed
	if vm.IsCollapsed("./a") || vm.IsCollapsed("./a/b") || vm.IsCollapsed("./a/b/c") {
		t.Fatalf("expected no collapsed paths initially")
	}

	// collapse recursively at ./a
	vm.ToggleSubtreeExpansion("./a")

	// now a and descendants should be collapsed
	assert.True(t, vm.IsCollapsed("./a"))
	assert.True(t, vm.IsCollapsed("./a/b"))
	assert.True(t, vm.IsCollapsed("./a/b/c"))

	// toggle again should expand them
	vm.ToggleSubtreeExpansion("./a")
	assert.False(t, vm.IsCollapsed("./a"))
	assert.False(t, vm.IsCollapsed("./a/b"))
	assert.False(t, vm.IsCollapsed("./a/b/c"))
}
