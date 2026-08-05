package stack

import (
	"errors"
)

var (
	ErrOverflow  = errors.New("queue is empty")
	ErrUnderflow = errors.New("queue is full")
)

type Stack struct {
	arr []int
	top int
	n   int
}

func NewStack(n int) *Stack {
	return &Stack{
		arr: make([]int, n+1),
		top: 0,
		n:   n,
	}
}

func (s *Stack) Pop() (int, error) {
	if s.top == 0 {
		return 0, ErrUnderflow
	}

	s.top = s.top - 1

	return s.arr[s.top+1], nil
}

func (s *Stack) Push(x int) error {
	if s.IsFull() {
		return ErrOverflow
	}

	s.top = s.top + 1
	s.arr[s.top] = x

	return nil
}

func (s *Stack) IsEmpty() bool {
	return s.top == 0
}

func (s *Stack) IsFull() bool {
	return s.top == s.n
}

func (s *Stack) Size() int {
	return s.top
}
