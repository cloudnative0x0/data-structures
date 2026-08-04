package stack

import (
	"math/rand"
	"testing"
)

func TestPushPop(t *testing.T) {
	s := NewStack(5)

	if err := s.Push(10); err != nil {
		t.Fatalf("unexpected error on push: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("unexpected error on pop: %v", err)
	}
	if val != 10 {
		t.Errorf("expected 10, got %d", val)
	}
}

func TestLIFOOrder(t *testing.T) {
	s := NewStack(3)

	if err := s.Push(1); err != nil {
		t.Fatalf("unexpected error on push: %v", err)
	}
	if err := s.Push(2); err != nil {
		t.Fatalf("unexpected error on push: %v", err)
	}
	if err := s.Push(3); err != nil {
		t.Fatalf("unexpected error on push: %v", err)
	}

	expected := []int{3, 2, 1}

	for i, want := range expected {
		got, err := s.Pop()
		if err != nil {
			t.Fatalf("unexpected error on pop #%d: %v", i, err)
		}
		if got != want {
			t.Errorf("pop #%d: expected %d, got %d", i, want, got)
		}
	}
}

func TestEmptyStack(t *testing.T) {
	s := NewStack(3)

	if !s.IsEmpty() {
		t.Error("new stack should be empty")
	}

	_, err := s.Pop()
	if err == nil {
		t.Error("expected underflow error, got nil")
	}
}

func TestFullStack(t *testing.T) {
	s := NewStack(2)

	if err := s.Push(1); err != nil {
		t.Fatalf("unexpected error on push: %v", err)
	}
	if err := s.Push(2); err != nil {
		t.Fatalf("unexpected error on push: %v", err)
	}

	if !s.IsFull() {
		t.Error("stack should be full")
	}

	err := s.Push(3)
	if err == nil {
		t.Error("expected overflow error, got nil")
	}
}

func TestSize(t *testing.T) {
	s := NewStack(5)

	if s.Size() != 0 {
		t.Errorf("expected size 0, got %d", s.Size())
	}

	if err := s.Push(1); err != nil {
		t.Fatalf("unexpected error on push: %v", err)
	}
	if err := s.Push(2); err != nil {
		t.Fatalf("unexpected error on push: %v", err)
	}

	if s.Size() != 2 {
		t.Errorf("expected size 2, got %d", s.Size())
	}

	if _, err := s.Pop(); err != nil {
		t.Fatalf("unexpected error on pop: %v", err)
	}

	if s.Size() != 1 {
		t.Errorf("expected size 1, got %d", s.Size())
	}
}

func TestStress(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	const iterations = 1000
	const opsPerIteration = 50
	const capacity = 30

	for iter := 0; iter < iterations; iter++ {
		mine := NewStack(capacity)
		var ref []int

		for op := 0; op < opsPerIteration; op++ {
			switch rng.Intn(3) {
			case 0: // push
				if !mine.IsFull() {
					val := rng.Intn(1000)
					if err := mine.Push(val); err != nil {
						t.Fatalf("iter %d: unexpected push error: %v", iter, err)
					}
					ref = append(ref, val)
				}

			case 1: // pop
				if len(ref) > 0 {
					want := ref[len(ref)-1]
					ref = ref[:len(ref)-1]

					got, err := mine.Pop()
					if err != nil {
						t.Fatalf("iter %d: unexpected pop error: %v", iter, err)
					}
					if got != want {
						t.Fatalf("iter %d: mismatch on pop: expected %d, got %d", iter, want, got)
					}
				}

			case 2:
				if mine.Size() != len(ref) {
					t.Fatalf("iter %d: size mismatch: expected %d, got %d", iter, len(ref), mine.Size())
				}
				if mine.IsEmpty() != (len(ref) == 0) {
					t.Fatalf("iter %d: isEmpty mismatch", iter)
				}
			}
		}
	}
}
