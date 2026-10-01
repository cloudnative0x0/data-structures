package floyd_warshal

import "errors"

const Infinity = int(^uint(0)>>1) / 4

var (
	ErrInvalidVertexCount = errors.New("vertex count must be non-negative")
	ErrVertexOutOfRange   = errors.New("vertex is out of range")
	ErrWeightOutOfRange   = errors.New("edge weight is out of supported range")
	ErrNegativeCycle      = errors.New("graph contains a negative cycle")
)

type Graph struct {
	weights [][]int
}

type Result struct {
	distance [][]int
	next     [][]int
}

func New(vertexCount int) (*Graph, error) {
	if vertexCount < 0 {
		return nil, ErrInvalidVertexCount
	}

	weights := make([][]int, vertexCount)
	for from := range weights {
		weights[from] = make([]int, vertexCount)
		for to := range weights[from] {
			weights[from][to] = Infinity
		}
		weights[from][from] = 0
	}
	return &Graph{weights: weights}, nil
}

func (g *Graph) AddEdge(from, to, weight int) error {
	if err := g.validateVertex(from); err != nil {
		return err
	}
	if err := g.validateVertex(to); err != nil {
		return err
	}
	if weight <= -Infinity || weight >= Infinity {
		return ErrWeightOutOfRange
	}
	if weight < g.weights[from][to] {
		g.weights[from][to] = weight
	}
	return nil
}

// AddUndirectedEdge adds the same weighted edge in both directions.
func (g *Graph) AddUndirectedEdge(a, b, weight int) error {
	if err := g.AddEdge(a, b, weight); err != nil {
		return err
	}
	return g.AddEdge(b, a, weight)
}

func (g *Graph) ShortestPaths() (*Result, error) {
	n := len(g.weights)
	distance := make([][]int, n)
	next := make([][]int, n)
	for from := 0; from < n; from++ {
		distance[from] = append([]int(nil), g.weights[from]...)
		next[from] = make([]int, n)
		for to := 0; to < n; to++ {
			next[from][to] = -1
			if distance[from][to] != Infinity {
				next[from][to] = to
			}
		}
	}

	for via := 0; via < n; via++ {
		for from := 0; from < n; from++ {
			if distance[from][via] == Infinity {
				continue
			}
			for to := 0; to < n; to++ {
				if distance[via][to] == Infinity {
					continue
				}
				candidate := boundedAdd(distance[from][via], distance[via][to])
				if candidate < distance[from][to] {
					distance[from][to] = candidate
					next[from][to] = next[from][via]
				}
			}
		}
	}

	for vertex := 0; vertex < n; vertex++ {
		if distance[vertex][vertex] < 0 {
			return nil, ErrNegativeCycle
		}
	}
	return &Result{distance: distance, next: next}, nil
}

// Distance returns a shortest-path distance. found is false when to is not
// reachable from from.
func (r *Result) Distance(from, to int) (distance int, found bool, err error) {
	if err := r.validateVertex(from); err != nil {
		return 0, false, err
	}
	if err := r.validateVertex(to); err != nil {
		return 0, false, err
	}
	if r.distance[from][to] == Infinity {
		return 0, false, nil
	}
	return r.distance[from][to], true, nil
}

// Path reconstructs one shortest path. found is false when no path exists.
func (r *Result) Path(from, to int) (path []int, found bool, err error) {
	if err := r.validateVertex(from); err != nil {
		return nil, false, err
	}
	if err := r.validateVertex(to); err != nil {
		return nil, false, err
	}
	if r.next[from][to] == -1 {
		return nil, false, nil
	}

	path = append(path, from)
	for from != to {
		from = r.next[from][to]
		path = append(path, from)
	}
	return path, true, nil
}

func (g *Graph) validateVertex(vertex int) error {
	if vertex < 0 || vertex >= len(g.weights) {
		return ErrVertexOutOfRange
	}
	return nil
}

func (r *Result) validateVertex(vertex int) error {
	if vertex < 0 || vertex >= len(r.distance) {
		return ErrVertexOutOfRange
	}
	return nil
}

func boundedAdd(a, b int) int {
	if b > 0 && a > Infinity-b {
		return Infinity
	}
	if b < 0 && a < -Infinity-b {
		return -Infinity
	}
	return a + b
}
