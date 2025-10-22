package filetree

import "testing"

func TestCollapsedPathsBasic(t *testing.T) {
	cp := NewCollapsedPaths()

	p := "./a"

	// initially not collapsed
	if cp.IsCollapsed(p) {
		t.Fatalf("expected %s to be expanded", p)
	}

	// collapse
	cp.Collapse(p)
	if !cp.IsCollapsed(p) {
		t.Fatalf("expected %s to be collapsed", p)
	}

	// toggle -> expand
	cp.ToggleCollapsed(p)
	if cp.IsCollapsed(p) {
		t.Fatalf("expected %s to be expanded after toggle", p)
	}

	// toggle -> collapse
	cp.ToggleCollapsed(p)
	if !cp.IsCollapsed(p) {
		t.Fatalf("expected %s to be collapsed after toggle", p)
	}

	// expand
	cp.Expand(p)
	if cp.IsCollapsed(p) {
		t.Fatalf("expected %s to be expanded after Expand", p)
	}
}
