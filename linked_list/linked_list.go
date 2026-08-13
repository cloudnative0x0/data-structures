package linked_list

import (
	"errors"
	"strings"
)

var (
	ErrEmptyList   = errors.New("list is empty")
	ErrOutOfBounds = errors.New("index out of bounds")
	ErrNotFound    = errors.New("value not found")
)

type Node[T any] struct {
	Value T
	Next  *Node[T]
}

type LinkedList[T comparable] struct {
	head *Node[T]
	tail *Node[T]
	size int
}

func New[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{}
}

func (ll *LinkedList[T]) Len() int { return ll.size }

func (ll *LinkedList[T]) IsEmpty() bool { return ll.size == 0 }

func (ll *LinkedList[T]) Prepend(value T) {
	node := &Node[T]{Value: value, Next: ll.head}
	ll.head = node
	if ll.tail == nil {
		ll.tail = node
	}

	ll.size++
}

func (ll *LinkedList[T]) Append(value T) {
	node := &Node[T]{Value: value}
	if ll.tail == nil {
		ll.head = node
		ll.tail = node
	} else {
		ll.tail.Next = node
		ll.tail = node
	}

	ll.size++
}

func (ll *LinkedList[T]) InsertAt(position int, value T) error {
	if position < 0 || position > ll.size {
		return ErrOutOfBounds
	}
	if position == 0 {
		ll.Prepend(value)
		return nil
	}
	if position == ll.size {
		ll.Append(value)
		return nil
	}

	prev := ll.head
	for i := 0; i < position-1; i++ {
		prev = prev.Next
	}

	node := &Node[T]{Value: value, Next: prev.Next}
	prev.Next = node
	ll.size++

	return nil
}

func (ll *LinkedList[T]) RemoveFirst() (T, error) {
	var zero T
	if ll.head == nil {
		return zero, ErrEmptyList
	}

	val := ll.head.Value

	ll.head = ll.head.Next
	if ll.head == nil {
		ll.tail = nil
	}

	ll.size--

	return val, nil
}

func (ll *LinkedList[T]) RemoveLast() (T, error) {
	var zero T
	if ll.head == nil {
		return zero, ErrEmptyList
	}
	if ll.head.Next == nil {
		val := ll.head.Value
		ll.head = nil
		ll.tail = nil
		ll.size--
		return val, nil
	}

	prev := ll.head
	for prev.Next.Next != nil {
		prev = prev.Next
	}

	val := prev.Next.Value
	prev.Next = nil
	ll.tail = prev
	ll.size--

	return val, nil
}

func (ll *LinkedList[T]) Remove(value T) error {
	if ll.head == nil {
		return ErrEmptyList
	}
	if ll.head.Value == value {
		ll.head = ll.head.Next

		if ll.head == nil {
			ll.tail = nil
		}
		ll.size--

		return nil
	}

	prev := ll.head
	for prev.Next != nil {
		if prev.Next.Value == value {
			if prev.Next == ll.tail {
				ll.tail = prev
			}

			prev.Next = prev.Next.Next
			ll.size--

			return nil
		}

		prev = prev.Next
	}

	return ErrNotFound
}

func (ll *LinkedList[T]) Search(value T) bool {
	for cur := ll.head; cur != nil; cur = cur.Next {
		if cur.Value == value {
			return true
		}
	}

	return false
}

func (ll *LinkedList[T]) Get(position int) (T, error) {
	var zero T
	if position < 0 || position >= ll.size {
		return zero, ErrOutOfBounds
	}
	cur := ll.head
	for i := 0; i < position; i++ {
		cur = cur.Next
	}

	return cur.Value, nil
}

func (ll *LinkedList[T]) Traverse(callback func(value T)) {
	for cur := ll.head; cur != nil; cur = cur.Next {
		callback(cur.Value)
	}
}

func (ll *LinkedList[T]) Reverse() {
	var prev *Node[T]

	cur := ll.head
	ll.tail = ll.head

	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}

	ll.head = prev
}

func (ll *LinkedList[T]) ToSlice() []T {
	slice := make([]T, 0, ll.size)
	for cur := ll.head; cur != nil; cur = cur.Next {
		slice = append(slice, cur.Value)
	}

	return slice
}

func (ll *LinkedList[T]) String() string {
	var b strings.Builder
	b.WriteString("[")
	for cur := ll.head; cur != nil; cur = cur.Next {

		if cur.Next != nil {
			b.WriteString(" -> ")
		}
	}
	b.WriteString("]")
	
	return b.String()
}
