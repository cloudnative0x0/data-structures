package dsu

import "errors"

var (
	ErrInvalidSize       = errors.New("size must be non-negative")
	ErrElementOutOfRange = errors.New("element is out of range")
)

type DSU struct {
	parent     []int
	size       []int
	components int
}

func New(n int) (*DSU, error) {
	if n < 0 {
		return nil, ErrInvalidSize
	}

	d := &DSU{
		parent:     make([]int, n),
		size:       make([]int, n),
		components: n,
	}
	for element := 0; element < n; element++ {
		d.parent[element] = element
		d.size[element] = 1
	}
	return d, nil
}

func (d *DSU) Find(element int) (int, error) {
	if err := d.validate(element); err != nil {
		return 0, err
	}

	root := element
	for root != d.parent[root] {
		root = d.parent[root]
	}
	for element != root {
		next := d.parent[element]
		d.parent[element] = root
		element = next
	}
	return root, nil
}

// Union merges the sets containing a and b. It returns true only when two
// previously separate sets were merged.
func (d *DSU) Union(a, b int) (bool, error) {
	rootA, err := d.Find(a)
	if err != nil {
		return false, err
	}
	rootB, err := d.Find(b)
	if err != nil {
		return false, err
	}
	if rootA == rootB {
		return false, nil
	}

	if d.size[rootA] < d.size[rootB] {
		rootA, rootB = rootB, rootA
	}
	d.parent[rootB] = rootA
	d.size[rootA] += d.size[rootB]
	d.components--
	return true, nil
}

// Connected reports whether a and b belong to the same set.
func (d *DSU) Connected(a, b int) (bool, error) {
	rootA, err := d.Find(a)
	if err != nil {
		return false, err
	}
	rootB, err := d.Find(b)
	if err != nil {
		return false, err
	}
	return rootA == rootB, nil
}

func (d *DSU) Size(element int) (int, error) {
	root, err := d.Find(element)
	if err != nil {
		return 0, err
	}
	return d.size[root], nil
}

func (d *DSU) Components() int {
	return d.components
}

func (d *DSU) validate(element int) error {
	if element < 0 || element >= len(d.parent) {
		return ErrElementOutOfRange
	}
	return nil
}
