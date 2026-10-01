package binary_heap

import (
	"errors"
	"math/rand"
	"slices"
	"testing"
)

func TestMinHeapPopsAscending(t *testing.T) {
	h := NewMinHeap(7, -2, 5, 5, 0, 11, -9)
	want := []int{-9, -2, 0, 5, 5, 7, 11}

	got := popAll(t, h)
	if !slices.Equal(got, want) {
		t.Fatalf("Pop order = %v, want %v", got, want)
	}
}

func TestMaxHeapPopsDescending(t *testing.T) {
	h := NewMaxHeap(7, -2, 5, 5, 0, 11, -9)
	want := []int{11, 7, 5, 5, 0, -2, -9}

	got := popAll(t, h)
	if !slices.Equal(got, want) {
		t.Fatalf("Pop order = %v, want %v", got, want)
	}
}

func TestPushPeekLenAndIsEmpty(t *testing.T) {
	h := NewMinHeap()
	if !h.IsEmpty() || h.Len() != 0 {
		t.Fatalf("new heap: IsEmpty=%v Len=%d, want true and 0", h.IsEmpty(), h.Len())
	}

	for _, value := range []int{8, 3, 10, 1} {
		h.Push(value)
	}

	if h.IsEmpty() || h.Len() != 4 {
		t.Fatalf("after pushes: IsEmpty=%v Len=%d, want false and 4", h.IsEmpty(), h.Len())
	}
	if got, err := h.Peek(); err != nil || got != 1 {
		t.Fatalf("Peek() = (%d, %v), want (1, nil)", got, err)
	}
	if h.Len() != 4 {
		t.Fatalf("Peek changed length to %d, want 4", h.Len())
	}
}

func TestEmptyHeapErrors(t *testing.T) {
	h := NewMinHeap()

	if got, err := h.Peek(); got != 0 || !errors.Is(err, ErrEmpty) {
		t.Fatalf("Peek() = (%d, %v), want (0, ErrEmpty)", got, err)
	}
	if got, err := h.Pop(); got != 0 || !errors.Is(err, ErrEmpty) {
		t.Fatalf("Pop() = (%d, %v), want (0, ErrEmpty)", got, err)
	}
}

func TestNewCopiesInput(t *testing.T) {
	values := []int{3, 1, 2}
	h := NewMinHeap(values...)

	values[0] = -100
	if got, err := h.Peek(); err != nil || got != 1 {
		t.Fatalf("heap aliases constructor input: Peek() = (%d, %v)", got, err)
	}
}

func TestGenericHeapWithCustomPriority(t *testing.T) {
	type job struct {
		name     string
		priority int
	}

	h := New(func(a, b job) bool {
		return a.priority > b.priority
	},
		job{name: "report", priority: 2},
		job{name: "incident", priority: 10},
		job{name: "email", priority: 1},
	)

	got, err := h.Pop()
	if err != nil {
		t.Fatalf("Pop() error = %v", err)
	}
	if got.name != "incident" {
		t.Fatalf("first job = %q, want incident", got.name)
	}
}

func TestNilComparatorPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("New with nil comparator did not panic")
		}
	}()

	New[int](nil)
}

func TestRandomOperationsAgainstReference(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	h := NewMinHeap()
	var reference []int

	for step := 0; step < 20_000; step++ {
		if len(reference) == 0 || rng.Intn(3) != 0 {
			value := rng.Intn(20_001) - 10_000
			h.Push(value)
			reference = append(reference, value)
		} else {
			minIndex := 0
			for i := 1; i < len(reference); i++ {
				if reference[i] < reference[minIndex] {
					minIndex = i
				}
			}

			want := reference[minIndex]
			reference[minIndex] = reference[len(reference)-1]
			reference = reference[:len(reference)-1]

			got, err := h.Pop()
			if err != nil || got != want {
				t.Fatalf("step %d: Pop() = (%d, %v), want (%d, nil)", step, got, err, want)
			}
		}

		if h.Len() != len(reference) {
			t.Fatalf("step %d: Len() = %d, want %d", step, h.Len(), len(reference))
		}
		assertInvariant(t, h)
	}
}

func popAll(t *testing.T, h *Heap[int]) []int {
	t.Helper()

	result := make([]int, 0, h.Len())
	for !h.IsEmpty() {
		value, err := h.Pop()
		if err != nil {
			t.Fatalf("Pop() error = %v", err)
		}
		result = append(result, value)
	}

	return result
}

func assertInvariant[T any](t *testing.T, h *Heap[T]) {
	t.Helper()

	for child := 1; child < len(h.data); child++ {
		parent := (child - 1) / 2
		if h.less(h.data[child], h.data[parent]) {
			t.Fatalf("heap invariant broken between parent %d and child %d", parent, child)
		}
	}
}
