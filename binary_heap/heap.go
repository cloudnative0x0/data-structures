package binary_heap

import "errors"

var ErrEmpty = errors.New("binary heap is empty")

type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

func New[T any](less func(a, b T) bool, values ...T) *Heap[T] {
	if less == nil {
		panic("binary_heap: nil comparator")
	}

	h := &Heap[T]{
		data: append([]T(nil), values...),
		less: less,
	}

	for i := len(h.data)/2 - 1; i >= 0; i-- {
		h.siftDown(i)
	}

	return h
}

func NewMinHeap(values ...int) *Heap[int] {
	return New(func(a, b int) bool { return a < b }, values...)
}

func NewMaxHeap(values ...int) *Heap[int] {
	return New(func(a, b int) bool { return a > b }, values...)
}

func (h *Heap[T]) Len() int {
	return len(h.data)
}

func (h *Heap[T]) IsEmpty() bool {
	return len(h.data) == 0
}

func (h *Heap[T]) Peek() (T, error) {
	if h.IsEmpty() {
		var zero T
		return zero, ErrEmpty
	}

	return h.data[0], nil
}

func (h *Heap[T]) Push(value T) {
	h.data = append(h.data, value)
	h.siftUp(len(h.data) - 1)
}

func (h *Heap[T]) Pop() (T, error) {
	if h.IsEmpty() {
		var zero T
		return zero, ErrEmpty
	}

	root := h.data[0]
	last := len(h.data) - 1
	h.data[0] = h.data[last]

	var zero T
	h.data[last] = zero
	h.data = h.data[:last]

	if len(h.data) > 0 {
		h.siftDown(0)
	}

	return root, nil
}

func (h *Heap[T]) siftUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if !h.less(h.data[index], h.data[parent]) {
			return
		}

		h.data[index], h.data[parent] = h.data[parent], h.data[index]
		index = parent
	}
}

func (h *Heap[T]) siftDown(index int) {
	for {
		left := 2*index + 1
		if left >= len(h.data) {
			return
		}

		best := left
		right := left + 1
		if right < len(h.data) && h.less(h.data[right], h.data[left]) {
			best = right
		}

		if !h.less(h.data[best], h.data[index]) {
			return
		}

		h.data[index], h.data[best] = h.data[best], h.data[index]
		index = best
	}
}
