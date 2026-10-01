package bfs

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

	return &Graph{
		adjacency: make([][]int, vertexCount),
		directed:  directed,
	}, nil
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

// Traverse returns vertices reachable from start in breadth-first order.
// Neighbours on the same level are visited in edge insertion order.
func (g *Graph) Traverse(start int) ([]int, error) {
	if err := g.validateVertex(start); err != nil {
		return nil, err
	}

	visited := make([]bool, len(g.adjacency))
	visited[start] = true
	queue := []int{start}
	order := make([]int, 0, len(g.adjacency))

	for head := 0; head < len(queue); head++ {
		vertex := queue[head]
		order = append(order, vertex)

		for _, neighbour := range g.adjacency[vertex] {
			if visited[neighbour] {
				continue
			}
			visited[neighbour] = true
			queue = append(queue, neighbour)
		}
	}

	return order, nil
}

// Distances returns the minimum number of edges from start to every vertex.
// An unreachable vertex has distance -1.
func (g *Graph) Distances(start int) ([]int, error) {
	if err := g.validateVertex(start); err != nil {
		return nil, err
	}

	distance := make([]int, len(g.adjacency))
	for i := range distance {
		distance[i] = -1
	}
	distance[start] = 0
	queue := []int{start}

	for head := 0; head < len(queue); head++ {
		vertex := queue[head]
		for _, neighbour := range g.adjacency[vertex] {
			if distance[neighbour] != -1 {
				continue
			}
			distance[neighbour] = distance[vertex] + 1
			queue = append(queue, neighbour)
		}
	}

	return distance, nil
}

// ShortestPath returns a path with the minimum number of edges. found is false
// when target is unreachable from start.
func (g *Graph) ShortestPath(start, target int) (path []int, found bool, err error) {
	if err := g.validateVertex(start); err != nil {
		return nil, false, err
	}
	if err := g.validateVertex(target); err != nil {
		return nil, false, err
	}

	parent := make([]int, len(g.adjacency))
	for i := range parent {
		parent[i] = -1
	}
	parent[start] = start
	queue := []int{start}

	for head := 0; head < len(queue) && parent[target] == -1; head++ {
		vertex := queue[head]
		for _, neighbour := range g.adjacency[vertex] {
			if parent[neighbour] != -1 {
				continue
			}
			parent[neighbour] = vertex
			queue = append(queue, neighbour)
		}
	}

	if parent[target] == -1 {
		return nil, false, nil
	}

	for vertex := target; ; vertex = parent[vertex] {
		path = append(path, vertex)
		if vertex == start {
			break
		}
	}
	reverse(path)

	return path, true, nil
}

func (g *Graph) validateVertex(vertex int) error {
	if vertex < 0 || vertex >= len(g.adjacency) {
		return ErrVertexOutOfRange
	}
	return nil
}

func reverse(values []int) {
	for left, right := 0, len(values)-1; left < right; left, right = left+1, right-1 {
		values[left], values[right] = values[right], values[left]
	}
}
