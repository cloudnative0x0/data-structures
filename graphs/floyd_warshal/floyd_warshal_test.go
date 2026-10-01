package floyd_warshal

import (
	"errors"
	"slices"
	"testing"
)

func TestShortestPathsAndPathReconstruction(t *testing.T) {
	g := mustGraph(t, 5)
	addEdges(t, g, [][3]int{
		{0, 1, 3}, {0, 2, 10}, {1, 2, 1}, {1, 3, 7},
		{2, 3, 2}, {3, 4, 1}, {0, 1, 5},
	})

	result, err := g.ShortestPaths()
	if err != nil {
		t.Fatalf("ShortestPaths() error = %v", err)
	}

	distance, found, err := result.Distance(0, 4)
	if err != nil || !found || distance != 7 {
		t.Fatalf("Distance(0, 4) = (%d, %v, %v), want (7, true, nil)", distance, found, err)
	}
	path, found, err := result.Path(0, 4)
	if err != nil || !found || !slices.Equal(path, []int{0, 1, 2, 3, 4}) {
		t.Fatalf("Path(0, 4) = (%v, %v, %v)", path, found, err)
	}

	if distance, found, err := result.Distance(4, 0); err != nil || found || distance != 0 {
		t.Fatalf("unreachable Distance = (%d, %v, %v)", distance, found, err)
	}
	if path, found, err := result.Path(4, 0); err != nil || found || path != nil {
		t.Fatalf("unreachable Path = (%v, %v, %v)", path, found, err)
	}
}

func TestNegativeEdgesWithoutNegativeCycle(t *testing.T) {
	g := mustGraph(t, 4)
	addEdges(t, g, [][3]int{{0, 1, 4}, {0, 2, 5}, {1, 2, -2}, {2, 3, 3}})
	result, err := g.ShortestPaths()
	if err != nil {
		t.Fatalf("ShortestPaths() error = %v", err)
	}
	distance, found, err := result.Distance(0, 3)
	if err != nil || !found || distance != 5 {
		t.Fatalf("Distance(0, 3) = (%d, %v, %v), want 5", distance, found, err)
	}
}

func TestUndirectedEdge(t *testing.T) {
	g := mustGraph(t, 2)
	if err := g.AddUndirectedEdge(0, 1, 9); err != nil {
		t.Fatalf("AddUndirectedEdge() error = %v", err)
	}
	result, err := g.ShortestPaths()
	if err != nil {
		t.Fatalf("ShortestPaths() error = %v", err)
	}
	for _, pair := range [][2]int{{0, 1}, {1, 0}} {
		distance, found, err := result.Distance(pair[0], pair[1])
		if err != nil || !found || distance != 9 {
			t.Fatalf("Distance%v = (%d, %v, %v)", pair, distance, found, err)
		}
	}
}

func TestNegativeCycle(t *testing.T) {
	g := mustGraph(t, 3)
	addEdges(t, g, [][3]int{{0, 1, 1}, {1, 2, -3}, {2, 0, 1}})
	if _, err := g.ShortestPaths(); !errors.Is(err, ErrNegativeCycle) {
		t.Fatalf("ShortestPaths() error = %v, want ErrNegativeCycle", err)
	}
}

func TestValidation(t *testing.T) {
	if _, err := New(-1); !errors.Is(err, ErrInvalidVertexCount) {
		t.Fatalf("New(-1) error = %v", err)
	}
	g := mustGraph(t, 2)
	if err := g.AddEdge(0, 2, 1); !errors.Is(err, ErrVertexOutOfRange) {
		t.Fatalf("AddEdge vertex error = %v", err)
	}
	if err := g.AddEdge(0, 1, Infinity); !errors.Is(err, ErrWeightOutOfRange) {
		t.Fatalf("AddEdge weight error = %v", err)
	}
	result, err := g.ShortestPaths()
	if err != nil {
		t.Fatalf("ShortestPaths() error = %v", err)
	}
	if _, _, err := result.Distance(-1, 0); !errors.Is(err, ErrVertexOutOfRange) {
		t.Fatalf("Distance validation error = %v", err)
	}
	if _, _, err := result.Path(0, 2); !errors.Is(err, ErrVertexOutOfRange) {
		t.Fatalf("Path validation error = %v", err)
	}
}

func mustGraph(t *testing.T, vertices int) *Graph {
	t.Helper()
	g, err := New(vertices)
	if err != nil {
		t.Fatalf("New(%d) error = %v", vertices, err)
	}
	return g
}

func addEdges(t *testing.T, g *Graph, edges [][3]int) {
	t.Helper()
	for _, edge := range edges {
		if err := g.AddEdge(edge[0], edge[1], edge[2]); err != nil {
			t.Fatalf("AddEdge%v error = %v", edge, err)
		}
	}
}
