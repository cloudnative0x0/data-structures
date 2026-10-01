package bfs

import (
	"errors"
	"slices"
	"testing"
)

func TestTraverseDistancesAndShortestPath(t *testing.T) {
	g := mustGraph(t, 7, false)
	addEdges(t, g, [][2]int{{0, 1}, {0, 2}, {1, 3}, {1, 4}, {2, 5}, {4, 5}})

	order, err := g.Traverse(0)
	if err != nil || !slices.Equal(order, []int{0, 1, 2, 3, 4, 5}) {
		t.Fatalf("Traverse(0) = (%v, %v)", order, err)
	}

	distance, err := g.Distances(0)
	if err != nil || !slices.Equal(distance, []int{0, 1, 1, 2, 2, 2, -1}) {
		t.Fatalf("Distances(0) = (%v, %v)", distance, err)
	}

	path, found, err := g.ShortestPath(0, 5)
	if err != nil || !found || !slices.Equal(path, []int{0, 2, 5}) {
		t.Fatalf("ShortestPath(0, 5) = (%v, %v, %v)", path, found, err)
	}

	path, found, err = g.ShortestPath(0, 6)
	if err != nil || found || path != nil {
		t.Fatalf("unreachable path = (%v, %v, %v)", path, found, err)
	}
}

func TestDirectedGraphAndSelfPath(t *testing.T) {
	g := mustGraph(t, 3, true)
	addEdges(t, g, [][2]int{{0, 1}, {1, 2}})

	order, err := g.Traverse(2)
	if err != nil || !slices.Equal(order, []int{2}) {
		t.Fatalf("Traverse(2) = (%v, %v)", order, err)
	}
	path, found, err := g.ShortestPath(1, 1)
	if err != nil || !found || !slices.Equal(path, []int{1}) {
		t.Fatalf("ShortestPath(1, 1) = (%v, %v, %v)", path, found, err)
	}
}

func TestValidation(t *testing.T) {
	if _, err := New(-1, false); !errors.Is(err, ErrInvalidVertexCount) {
		t.Fatalf("New(-1) error = %v", err)
	}

	g := mustGraph(t, 2, false)
	for _, edge := range [][2]int{{-1, 0}, {0, 2}} {
		if err := g.AddEdge(edge[0], edge[1]); !errors.Is(err, ErrVertexOutOfRange) {
			t.Fatalf("AddEdge%v error = %v", edge, err)
		}
	}
	if _, err := g.Traverse(2); !errors.Is(err, ErrVertexOutOfRange) {
		t.Fatalf("Traverse(2) error = %v", err)
	}
	if _, _, err := g.ShortestPath(0, -1); !errors.Is(err, ErrVertexOutOfRange) {
		t.Fatalf("ShortestPath error = %v", err)
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
