package dfs

import "errors"

var (
	ErrInvalidVertexCount = errors.New("vertex count must be non-negative")
	ErrVertexOutOfRange   = errors.New("vertex is out of range")
)

type Graph struct {
	adjacency [][]int
	directed  bool
}

func New(vertexCount int, directed bool) (*Graph, error) {
	if vertexCount < 0 {
		return nil, ErrInvalidVertexCount
	}
	return &Graph{adjacency: make([][]int, vertexCount), directed: directed}, nil
}

func (g *Graph) AddEdge(from, to int) error {
	if err := g.validateVertex(from); err != nil {
		return err
	}
	if err := g.validateVertex(to); err != nil {
		return err
	}

	g.adjacency[from] = append(g.adjacency[from], to)
	if !g.directed && from != to {
		g.adjacency[to] = append(g.adjacency[to], from)
	}
	return nil
}

// Traverse returns reachable vertices in iterative depth-first order.
func (g *Graph) Traverse(start int) ([]int, error) {
	if err := g.validateVertex(start); err != nil {
		return nil, err
	}

	visited := make([]bool, len(g.adjacency))
	stack := []int{start}
	order := make([]int, 0, len(g.adjacency))

	for len(stack) > 0 {
		last := len(stack) - 1
		vertex := stack[last]
		stack = stack[:last]
		if visited[vertex] {
			continue
		}

		visited[vertex] = true
		order = append(order, vertex)

		// Reverse push preserves adjacency insertion order during traversal.
		for i := len(g.adjacency[vertex]) - 1; i >= 0; i-- {
			neighbour := g.adjacency[vertex][i]
			if !visited[neighbour] {
				stack = append(stack, neighbour)
			}
		}
	}

	return order, nil
}

// TraverseRecursive returns reachable vertices using recursive DFS.
func (g *Graph) TraverseRecursive(start int) ([]int, error) {
	if err := g.validateVertex(start); err != nil {
		return nil, err
	}

	visited := make([]bool, len(g.adjacency))
	order := make([]int, 0, len(g.adjacency))
	var visit func(int)
	visit = func(vertex int) {
		visited[vertex] = true
		order = append(order, vertex)
		for _, neighbour := range g.adjacency[vertex] {
			if !visited[neighbour] {
				visit(neighbour)
			}
		}
	}
	visit(start)

	return order, nil
}

// HasPath reports whether target is reachable from start.
func (g *Graph) HasPath(start, target int) (bool, error) {
	if err := g.validateVertex(start); err != nil {
		return false, err
	}
	if err := g.validateVertex(target); err != nil {
		return false, err
	}

	visited := make([]bool, len(g.adjacency))
	stack := []int{start}
	for len(stack) > 0 {
		last := len(stack) - 1
		vertex := stack[last]
		stack = stack[:last]
		if vertex == target {
			return true, nil
		}
		if visited[vertex] {
			continue
		}
		visited[vertex] = true
		for _, neighbour := range g.adjacency[vertex] {
			if !visited[neighbour] {
				stack = append(stack, neighbour)
			}
		}
	}

	return false, nil
}

func (g *Graph) validateVertex(vertex int) error {
	if vertex < 0 || vertex >= len(g.adjacency) {
		return ErrVertexOutOfRange
	}
	return nil
}
