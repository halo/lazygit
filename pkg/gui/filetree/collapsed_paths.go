package filetree

import (
	"strings"

	"github.com/jesseduffield/generics/set"
)

type CollapsedPaths struct {
	collapsedPaths *set.Set[string]
}

func NewCollapsedPaths() *CollapsedPaths {
	return &CollapsedPaths{
		collapsedPaths: set.New[string](),
	}
}

func (self *CollapsedPaths) ExpandToPath(path string) {
	// need every directory along the way
	splitPath := split(path)
	for i := range splitPath {
		dir := join(splitPath[0 : i+1])
		self.collapsedPaths.Remove(dir)
	}
}

func (self *CollapsedPaths) IsCollapsed(path string) bool {
	return self.collapsedPaths.Includes(path)
}

func (self *CollapsedPaths) Collapse(path string) {
	self.collapsedPaths.Add(path)
}

func (self *CollapsedPaths) ToggleCollapsed(path string) {
	if self.collapsedPaths.Includes(path) {
		self.collapsedPaths.Remove(path)
	} else {
		self.collapsedPaths.Add(path)
	}
}

// Expand removes the given path from the collapsed set (i.e. mark expanded)
func (self *CollapsedPaths) Expand(path string) {
	self.collapsedPaths.Remove(path)
}

func (self *CollapsedPaths) ExpandAll() {
	// Could be cleaner if Set had a Clear() method...
	self.collapsedPaths.RemoveSlice(self.collapsedPaths.ToSlice())
}

// ToggleSubtreeExpansionOn toggles the collapsed state for the given path and all
// descendant directories within the provided tree root. It's a generic helper
// shared by different tree implementations to avoid code duplication.
func ToggleSubtreeExpansionOn[T any](self *CollapsedPaths, root *Node[T], path string) {
	if root == nil {
		return
	}

	if self.IsCollapsed(path) {
		expandSubtree(self, root, path)
	} else {
		collapseSubtree(self, root, path)
	}
}

// expandSubtree expands the given path and all its descendant directories
func expandSubtree[T any](self *CollapsedPaths, root *Node[T], path string) {
	collectAndApply(self, root, path, func(self *CollapsedPaths, p string) {
		self.Expand(p)
	})
}

// collapseSubtree collapses the given path and all its descendant directories
func collapseSubtree[T any](self *CollapsedPaths, root *Node[T], path string) {
	collectAndApply(self, root, path, func(self *CollapsedPaths, p string) {
		self.Collapse(p)
	})
}

// collectAndApply traverses the tree and applies the given operation to all
// directories that are the target path or its descendants
func collectAndApply[T any](self *CollapsedPaths, root *Node[T], path string, apply func(*CollapsedPaths, string)) {
	var collect func(n *Node[T])
	collect = func(n *Node[T]) {
		if n == nil {
			return
		}
		if !n.IsFile() {
			p := n.GetInternalPath()
			if p == path || strings.HasPrefix(p, path+"/") {
				apply(self, p)
			}
		}
		for _, ch := range n.Children {
			collect(ch)
		}
	}
	collect(root)
}
