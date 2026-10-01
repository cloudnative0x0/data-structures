package dfs

import (
	"slices"
	"testing"
)

func TestNewBuildsCompleteTree(t *testing.T) {
	root := New([]int{1, 2, 3, 4, 5, 6, 7})

	if root == nil || root.Value != 1 {
		t.Fatalf("root = %+v, want node 1", root)
	}
	if root.Left.Value != 2 || root.Right.Value != 3 {
		t.Fatalf("first level = (%d, %d), want (2, 3)", root.Left.Value, root.Right.Value)
	}
	if root.Left.Left.Value != 4 || root.Left.Right.Value != 5 {
		t.Fatalf("children of node 2 are incorrect")
	}
	if root.Right.Left.Value != 6 || root.Right.Right.Value != 7 {
		t.Fatalf("children of node 3 are incorrect")
	}
}

func TestNewEmpty(t *testing.T) {
	if got := New(nil); got != nil {
		t.Fatalf("New(nil) = %+v, want nil", got)
	}
	if got := New([]int{}); got != nil {
		t.Fatalf("New(empty) = %+v, want nil", got)
	}
}

func TestTraverseOrders(t *testing.T) {
	root := New([]int{1, 2, 3, 4, 5, 6, 7})

	tests := []struct {
		name string
		run  func(func(int))
		want []int
	}{
		{name: "DFS", run: root.TraverseDFS, want: []int{1, 2, 4, 5, 3, 6, 7}},
		{name: "pre-order", run: root.TraversePreOrder, want: []int{1, 2, 4, 5, 3, 6, 7}},
		{name: "in-order", run: root.TraverseInOrder, want: []int{4, 2, 5, 1, 6, 3, 7}},
		{name: "post-order", run: root.TraversePostOrder, want: []int{4, 5, 2, 6, 7, 3, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int
			tt.run(func(value int) { got = append(got, value) })
			if !slices.Equal(got, tt.want) {
				t.Fatalf("order = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	root := New([]int{1, 2, 3, 4, 5, 6, 7})

	if got := root.Search(6); got == nil || got.Value != 6 {
		t.Fatalf("Search(6) = %+v, want node 6", got)
	}
	if got := root.Search(42); got != nil {
		t.Fatalf("Search(42) = %+v, want nil", got)
	}
}

func TestSearchReturnsFirstMatchInDFSOrder(t *testing.T) {
	leftMatch := &Node{Value: 9}
	rightMatch := &Node{Value: 9}
	root := &Node{
		Value: 1,
		Left:  &Node{Value: 2, Left: leftMatch},
		Right: rightMatch,
	}

	if got := root.Search(9); got != leftMatch {
		t.Fatalf("Search(9) = %p, want left DFS match %p", got, leftMatch)
	}
}

func TestNilReceiverAndNilVisitor(t *testing.T) {
	var root *Node
	if got := root.Search(1); got != nil {
		t.Fatalf("Search on nil tree = %+v, want nil", got)
	}

	root.TraverseDFS(nil)
	root.TraversePreOrder(nil)
	root.TraverseInOrder(nil)
	root.TraversePostOrder(nil)

	nonEmpty := New([]int{1, 2, 3})
	nonEmpty.TraverseDFS(nil)
	nonEmpty.TraversePreOrder(nil)
	nonEmpty.TraverseInOrder(nil)
	nonEmpty.TraversePostOrder(nil)
}

func TestIterativeDFSHandlesDeepTree(t *testing.T) {
	const depth = 100_000
	root := &Node{Value: 0}
	current := root
	for i := 1; i < depth; i++ {
		current.Left = &Node{Value: i}
		current = current.Left
	}

	visited := 0
	root.TraverseDFS(func(int) { visited++ })
	if visited != depth {
		t.Fatalf("visited %d nodes, want %d", visited, depth)
	}
	if got := root.Search(depth - 1); got == nil || got.Value != depth-1 {
		t.Fatalf("failed to find deepest node: %+v", got)
	}
}
