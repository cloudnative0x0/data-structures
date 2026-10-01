package dfs

import (
	"errors"
	"slices"
	"testing"
)

func TestIterativeAndRecursiveTraversal(t *testing.T) {
	g := mustGraph(t, 7, false)
	addEdges(t, g, [][2]int{{0, 1}, {0, 2}, {1, 3}, {1, 4}, {2, 5}, {4, 5}})
	want := []int{0, 1, 3, 4, 5, 2}

	iterative, err := g.Traverse(0)
	if err != nil || !slices.Equal(iterative, want) {
		t.Fatalf("Traverse(0) = (%v, %v), want %v", iterative, err, want)
	}
	recursive, err := g.TraverseRecursive(0)
	if err != nil || !slices.Equal(recursive, want) {
		t.Fatalf("TraverseRecursive(0) = (%v, %v), want %v", recursive, err, want)
	}
}

func TestDirectedReachability(t *testing.T) {
	g := mustGraph(t, 4, true)
	addEdges(t, g, [][2]int{{0, 1}, {1, 2}, {2, 0}})

	for _, tc := range []struct {
		from, to int
		want     bool
	}{{0, 2, true}, {2, 1, true}, {0, 3, false}, {3, 3, true}} {
		got, err := g.HasPath(tc.from, tc.to)
		if err != nil || got != tc.want {
			t.Fatalf("HasPath(%d, %d) = (%v, %v), want %v", tc.from, tc.to, got, err, tc.want)
		}
	}
}

func TestIterativeTraversalHandlesDeepGraph(t *testing.T) {
	const vertices = 100_000
	g := mustGraph(t, vertices, true)
	for vertex := 0; vertex < vertices-1; vertex++ {
		if err := g.AddEdge(vertex, vertex+1); err != nil {
			t.Fatalf("AddEdge(%d, %d) error = %v", vertex, vertex+1, err)
		}
	}

	order, err := g.Traverse(0)
	if err != nil || len(order) != vertices || order[vertices-1] != vertices-1 {
		t.Fatalf("deep Traverse: len=%d last=%d err=%v", len(order), order[len(order)-1], err)
	}
}

func TestValidation(t *testing.T) {
	if _, err := New(-1, false); !errors.Is(err, ErrInvalidVertexCount) {
		t.Fatalf("New(-1) error = %v", err)
	}
	g := mustGraph(t, 2, false)
	if err := g.AddEdge(0, 2); !errors.Is(err, ErrVertexOutOfRange) {
		t.Fatalf("AddEdge error = %v", err)
	}
	if _, err := g.Traverse(-1); !errors.Is(err, ErrVertexOutOfRange) {
		t.Fatalf("Traverse error = %v", err)
	}
	if _, err := g.HasPath(0, 2); !errors.Is(err, ErrVertexOutOfRange) {
		t.Fatalf("HasPath error = %v", err)
	}
}

func mustGraph(t *testing.T, vertices int, directed bool) *Graph {
	t.Helper()
	g, err := New(vertices, directed)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return g
}

func addEdges(t *testing.T, g *Graph, edges [][2]int) {
	t.Helper()
	for _, edge := range edges {
		if err := g.AddEdge(edge[0], edge[1]); err != nil {
			t.Fatalf("AddEdge%v error = %v", edge, err)
		}
	}
}
