package queue

import "errors"

var (
	ErrEmpty           = errors.New("queue is empty")
	ErrFull            = errors.New("queue is full")
	ErrInvalidCapacity = errors.New("capacity must be greater than zero")
)

type Queue[T any] struct {
	arr   []T
	head  int
	tail  int
	count int
}

func NewQueue[T any](capacity int) (*Queue[T], error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}
	return &Queue[T]{
		arr: make([]T, capacity),
	}, nil
}

func (q *Queue[T]) Enqueue(value T) error {
	if q.IsFull() {
		return ErrFull
	}

	q.arr[q.tail] = value
	q.tail = (q.tail + 1) % len(q.arr)
	q.count++

	return nil
}

func (q *Queue[T]) Dequeue() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, ErrEmpty
	}

	oldHead := q.arr[q.head]
	q.head = (q.head + 1) % len(q.arr)
	q.count--

	return oldHead, nil
}

func (q *Queue[T]) Peek() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, ErrEmpty
	}

	return q.arr[q.head], nil
}

func (q *Queue[T]) IsEmpty() bool {
	return q.count == 0
}

func (q *Queue[T]) IsFull() bool {
	return q.count == len(q.arr)
}
