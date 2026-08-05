package stack

import "errors"

var (
	ErrUnderflow       = errors.New("stack underflow: cannot pop from empty stack")
	ErrOverflow        = errors.New("stack overflow: cannot push onto full stack")
	ErrInvalidCapacity = errors.New("capacity must be greater than zero")
)

type Stack[T any] struct {
	arr []T
	top int
	n   int
}

func NewStack[T any](n int) (*Stack[T], error) {
	if n <= 0 {
		return nil, ErrInvalidCapacity
	}

	return &Stack[T]{
		arr: make([]T, n+1),
		top: 0,
		n:   n,
	}, nil
}

func (s *Stack[T]) Pop() (T, error) {
	if s.top == 0 {
		var zero T
		return zero, ErrUnderflow
	}

	s.top--

	return s.arr[s.top+1], nil
}

func (s *Stack[T]) Push(x T) error {
	if s.IsFull() {
		return ErrOverflow
	}

	s.top++
	s.arr[s.top] = x

	return nil
}

func (s *Stack[T]) IsEmpty() bool {
	return s.top == 0
}

func (s *Stack[T]) IsFull() bool {
	return s.top == s.n
}

func (s *Stack[T]) Size() int {
	return s.top
}
