package queue

import (
	"testing"
)

func TestNewQueueInvalidCapacity(t *testing.T) {
	for _, c := range []int{0, -1, -100} {
		q, err := NewQueue[int](c)
		if err != ErrInvalidCapacity {
			t.Errorf("capacity %d: expected ErrInvalidCapacity, got %v", c, err)
		}
		if q != nil {
			t.Errorf("capacity %d: expected nil queue, got %v", c, q)
		}
	}
}

func TestNewQueueValid(t *testing.T) {
	q, err := NewQueue[int](3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !q.IsEmpty() {
		t.Error("new queue should be empty")
	}
	if q.IsFull() {
		t.Error("new queue should not be full")
	}
}

func TestEnqueueDequeue(t *testing.T) {
	q, _ := NewQueue[int](5)
	values := []int{10, 20, 30}
	for _, v := range values {
		if err := q.Enqueue(v); err != nil {
			t.Fatalf("Enqueue(%d) unexpected error: %v", v, err)
		}
	}
	for _, want := range values {
		got, err := q.Dequeue()
		if err != nil {
			t.Fatalf("Dequeue() unexpected error: %v", err)
		}
		if got != want {
			t.Errorf("Dequeue() = %d, want %d", got, want)
		}
	}
}

func TestEnqueueUntilFull(t *testing.T) {
	q, _ := NewQueue[int](2)
	_ = q.Enqueue(1)
	_ = q.Enqueue(2)

	if !q.IsFull() {
		t.Error("queue should be full")
	}
	if err := q.Enqueue(3); err != ErrFull {
		t.Errorf("expected ErrFull, got %v", err)
	}
}

func TestDequeueFromEmpty(t *testing.T) {
	q, _ := NewQueue[int](2)
	val, err := q.Dequeue()
	if err != ErrEmpty {
		t.Errorf("expected ErrEmpty, got %v", err)
	}
	if val != 0 {
		t.Errorf("expected zero value, got %v", val)
	}
}

func TestPeek(t *testing.T) {
	q, _ := NewQueue[int](3)
	_, err := q.Peek()
	if err != ErrEmpty {
		t.Errorf("Peek on empty: expected ErrEmpty, got %v", err)
	}

	if err := q.Enqueue(42); err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}
	val, err := q.Peek()
	if err != nil {
		t.Fatalf("Peek unexpected error: %v", err)
	}
	if val != 42 {
		t.Errorf("Peek = %d, want 42", val)
	}
	if q.IsEmpty() {
		t.Error("queue should not be empty after Peek")
	}

	val2, _ := q.Dequeue()
	if val2 != 42 {
		t.Errorf("Dequeue after Peek = %d, want 42", val2)
	}
}

func TestCircularWrapAround(t *testing.T) {
	q, _ := NewQueue[int](3)
	if err := q.Enqueue(1); err != nil { // head=0, tail=1
		t.Fatalf("Enqueue failed: %v", err)
	}
	if err := q.Enqueue(2); err != nil { // head=0, tail=2
		t.Fatalf("Enqueue failed: %v", err)
	}

	v, err := q.Dequeue()
	if err != nil || v != 1 {
		t.Fatalf("expected 1, got %v", v)
	}
	// head=1, tail=2
	if err := q.Enqueue(3); err != nil { // tail=0
		t.Fatalf("Enqueue failed: %v", err)
	}
	if err := q.Enqueue(4); err != nil { // tail=1
		t.Fatalf("Enqueue failed: %v", err)
	}

	want := []int{2, 3, 4}
	for _, w := range want {
		got, err := q.Dequeue()
		if err != nil {
			t.Fatalf("Dequeue unexpected error: %v", err)
		}
		if got != w {
			t.Errorf("Dequeue = %d, want %d", got, w)
		}
	}
	if !q.IsEmpty() {
		t.Error("queue should be empty after extracting all")
	}
}

func TestQueueWithStrings(t *testing.T) {
	q, _ := NewQueue[string](2)
	_ = q.Enqueue("hello")
	_ = q.Enqueue("world")
	if err := q.Enqueue("extra"); err != ErrFull {
		t.Errorf("expected ErrFull, got %v", err)
	}
	val1, _ := q.Dequeue()
	val2, _ := q.Dequeue()
	if val1 != "hello" || val2 != "world" {
		t.Errorf("Dequeue: got %q %q, want %q %q", val1, val2, "hello", "world")
	}
}

func TestStress(t *testing.T) {
	capacity := 100
	q, _ := NewQueue[int](capacity)
	var reference []int

	ops := 10000
	for i := 0; i < ops; i++ {
		if len(reference) == 0 || (len(reference) < capacity && i%2 == 0) {
			// enqueue
			val := i
			refErr := error(nil)
			if len(reference) == capacity {
				refErr = ErrFull
			}
			qErr := q.Enqueue(val)
			if refErr != qErr {
				t.Fatalf("Enqueue error mismatch: expected %v, got %v", refErr, qErr)
			}
			if qErr == nil {
				reference = append(reference, val)
			}
		} else {
			// dequeue
			var wantVal int
			var refErr error
			if len(reference) == 0 {
				wantVal = 0
				refErr = ErrEmpty
			} else {
				wantVal = reference[0]
				reference = reference[1:]
			}
			qVal, qErr := q.Dequeue()
			if qErr != refErr {
				t.Fatalf("Dequeue error mismatch: expected %v, got %v", refErr, qErr)
			}
			if qErr == nil && qVal != wantVal {
				t.Fatalf("Dequeue value: got %d, want %d", qVal, wantVal)
			}
		}
	}
}
