package linked_list

import (
	"errors"
	"reflect"
	"testing"
)

func TestPrependOnEmptyList(t *testing.T) {
	ll := New[int]()
	ll.Prepend(1) // must not panic on empty list
	if got := ll.ToSlice(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("got %v", got)
	}
}

func TestPrependKeepsRestOfList(t *testing.T) {
	ll := New[int]()
	ll.Append(2)
	ll.Append(3)
	ll.Prepend(1)
	want := []int{1, 2, 3}
	if got := ll.ToSlice(); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestAppendOnEmptyList(t *testing.T) {
	ll := New[int]()
	ll.Append(1) // must not panic on empty list
	if got := ll.ToSlice(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("got %v", got)
	}
}

func TestAppendDoesNotCorruptList(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)
	ll.Append(4)
	want := []int{1, 2, 3, 4}
	if got := ll.ToSlice(); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if ll.Len() != 4 {
		t.Fatalf("len = %d, want 4", ll.Len())
	}
}

func TestInsertAtDoesNotDropTail(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(4)
	if err := ll.InsertAt(2, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{1, 2, 3, 4}
	if got := ll.ToSlice(); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestInsertAtBoundaries(t *testing.T) {
	ll := New[int]()
	ll.Append(2)
	ll.Append(3)

	if err := ll.InsertAt(0, 1); err != nil {
		t.Fatal(err)
	}
	if err := ll.InsertAt(ll.Len(), 4); err != nil {
		t.Fatal(err)
	}
	want := []int{1, 2, 3, 4}
	if got := ll.ToSlice(); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	if err := ll.InsertAt(-1, 99); !errors.Is(err, ErrOutOfBounds) {
		t.Fatalf("want ErrOutOfBounds, got %v", err)
	}
	if err := ll.InsertAt(ll.Len()+1, 99); !errors.Is(err, ErrOutOfBounds) {
		t.Fatalf("want ErrOutOfBounds, got %v", err)
	}
}

func TestRemoveFirstAndLastOnEmpty(t *testing.T) {
	ll := New[int]()
	if _, err := ll.RemoveFirst(); !errors.Is(err, ErrEmptyList) {
		t.Fatalf("want ErrEmptyList, got %v", err)
	}
	if _, err := ll.RemoveLast(); !errors.Is(err, ErrEmptyList) {
		t.Fatalf("want ErrEmptyList, got %v", err)
	}
}

func TestRemoveLastKeepsTailConsistent(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	if _, err := ll.RemoveLast(); err != nil {
		t.Fatal(err)
	}
	ll.Append(30)
	want := []int{1, 2, 30}
	if got := ll.ToSlice(); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRemoveMiddleAndTailKeepsTailPointer(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	if err := ll.Remove(3); err != nil {
		t.Fatal(err)
	}
	ll.Append(30) // must correctly attach to the new tail (2)
	want := []int{1, 2, 30}
	if got := ll.ToSlice(); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRemoveNotFound(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	if err := ll.Remove(99); !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestSearch(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	if !ll.Search(2) {
		t.Fatal("expected to find 2")
	}
	if ll.Search(99) {
		t.Fatal("did not expect to find 99")
	}
}

func TestReverse(t *testing.T) {
	ll := New[int]()
	for _, v := range []int{1, 2, 3, 4} {
		ll.Append(v)
	}
	ll.Reverse()
	want := []int{4, 3, 2, 1}
	if got := ll.ToSlice(); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	ll.Append(0)
	want = []int{4, 3, 2, 1, 0}
	if got := ll.ToSlice(); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestReverseEmptyAndSingle(t *testing.T) {
	ll := New[int]()
	ll.Reverse() // must not panic
	if ll.Len() != 0 {
		t.Fatalf("len = %d, want 0", ll.Len())
	}

	ll.Append(1)
	ll.Reverse()
	if got := ll.ToSlice(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("got %v", got)
	}
}

func TestStringType(t *testing.T) {
	ll := New[string]()
	ll.Append("a")
	ll.Append("b")
	ll.Prepend("z")
	want := []string{"z", "a", "b"}
	if got := ll.ToSlice(); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSizeConsistencyAcrossOps(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Prepend(0)
	err := ll.InsertAt(2, 99)
	if err != nil {
		return
	}
	if ll.Len() != 4 {
		t.Fatalf("len = %d, want 4", ll.Len())
	}
	errR := ll.Remove(99)
	if errR != nil {
		return
	}
	_, errRF := ll.RemoveFirst()
	if errRF != nil {
		return
	}
	_, errRL := ll.RemoveLast()
	if errRL != nil {
		return
	}
	if ll.Len() != 1 {
		t.Fatalf("len = %d, want 1", ll.Len())
	}
}
