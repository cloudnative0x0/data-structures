package hash_map

import (
	"fmt"
	"testing"
)

func fnv1aString(s string) uint64 {
	var h uint64 = 14695981039346656037
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}
	return h
}

func djb2String(s string) uint64 {
	var h uint64 = 5381
	for i := 0; i < len(s); i++ {
		h = ((h << 5) + h) + uint64(s[i])
	}
	return h
}

func identInt(x int) uint64 { return uint64(x) }
func mulInt(x int) uint64   { return uint64(x)*2654435761 + 1 }

func zeroHash[K comparable](K) uint64 { return 0 }

func oneHash[K comparable](K) uint64 { return 1 }

func newStringTable(size int, strategy ProbeStrategy) *HashTable[string, int] {
	return NewHashMap[string, int](size, strategy, fnv1aString, djb2String)
}

func strategyName(s ProbeStrategy) string {
	switch s {
	case Linear:
		return "Linear"
	case Quadratic:
		return "Quadratic"
	case Double:
		return "Double"
	default:
		return "Unknown"
	}
}

func TestInsertAndSearch(t *testing.T) {
	ht := newStringTable(16, Linear)
	if ok := ht.Insert("a", 1); !ok {
		t.Fatalf("Insert(a) = false, want true")
	}
	v, found := ht.Search("a")
	if !found || v != 1 {
		t.Fatalf("Search(a) = (%d, %v), want (1, true)", v, found)
	}
}

func TestSearchMissingKey(t *testing.T) {
	ht := newStringTable(16, Linear)
	ht.Insert("a", 1)
	if _, found := ht.Search("nope"); found {
		t.Fatalf("Search(nope) found = true, want false")
	}
}

func TestSearchOnEmptyTable(t *testing.T) {
	ht := newStringTable(16, Linear)
	if _, found := ht.Search("a"); found {
		t.Fatalf("Search on empty table found = true, want false")
	}
}

func TestInsertUpdatesExistingKey(t *testing.T) {
	ht := newStringTable(16, Linear)
	ht.Insert("a", 1)
	if ok := ht.Insert("a", 2); !ok {
		t.Fatalf("Insert(a, 2) = false, want true")
	}
	v, found := ht.Search("a")
	if !found || v != 2 {
		t.Fatalf("Search(a) = (%d, %v), want (2, true)", v, found)
	}
	for i := 0; i < 15; i++ {
		key := fmt.Sprintf("k%d", i)
		if ok := ht.Insert(key, i); !ok {
			t.Fatalf("Insert(%s) = false, want true (table should have room)", key)
		}
	}
}

func TestDeleteThenSearch(t *testing.T) {
	ht := newStringTable(16, Linear)
	ht.Insert("a", 1)
	if ok := ht.Delete("a"); !ok {
		t.Fatalf("Delete(a) = false, want true")
	}
	if _, found := ht.Search("a"); found {
		t.Fatalf("Search(a) after delete found = true, want false")
	}
}

func TestDeleteNonExistentKey(t *testing.T) {
	ht := newStringTable(16, Linear)
	ht.Insert("a", 1)
	if ok := ht.Delete("nope"); ok {
		t.Fatalf("Delete(nope) = true, want false")
	}
}

func TestDeleteOnEmptyTable(t *testing.T) {
	ht := newStringTable(16, Linear)
	if ok := ht.Delete("a"); ok {
		t.Fatalf("Delete on empty table = true, want false")
	}
}

func TestReinsertAfterDelete(t *testing.T) {
	ht := newStringTable(16, Linear)
	ht.Insert("a", 1)
	ht.Delete("a")
	if ok := ht.Insert("a", 42); !ok {
		t.Fatalf("Insert(a, 42) after delete = false, want true")
	}
	v, found := ht.Search("a")
	if !found || v != 42 {
		t.Fatalf("Search(a) = (%d, %v), want (42, true)", v, found)
	}
}

func TestTombstoneIsReusedOnInsert(t *testing.T) {
	size := 4
	ht := newStringTable(size, Linear)
	keys := []string{"a", "b", "c", "d"}
	for _, k := range keys {
		if ok := ht.Insert(k, 0); !ok {
			t.Fatalf("Insert(%s) = false, want true", k)
		}
	}
	if ok := ht.Insert("e", 0); ok {
		t.Fatalf("Insert(e) into full table = true, want false")
	}
	ht.Delete("a")
	if ok := ht.Insert("e", 99); !ok {
		t.Fatalf("Insert(e) after delete = false, want true")
	}
	v, found := ht.Search("e")
	if !found || v != 99 {
		t.Fatalf("Search(e) = (%d, %v), want (99, true)", v, found)
	}
	if _, found := ht.Search("a"); found {
		t.Fatalf("Search(a) after delete found = true, want false")
	}
}

func TestSearchPastTombstone(t *testing.T) {
	ht := NewHashMap[string, int](8, Linear, zeroHash[string], oneHash[string])
	ht.Insert("a", 1) // slot 0
	ht.Insert("b", 2) // slot 1
	ht.Insert("c", 3) // slot 2
	ht.Delete("a")    // slot 0 -> tombstone

	if v, found := ht.Search("b"); !found || v != 2 {
		t.Fatalf("Search(b) = (%d, %v), want (2, true)", v, found)
	}
	if v, found := ht.Search("c"); !found || v != 3 {
		t.Fatalf("Search(c) = (%d, %v), want (3, true)", v, found)
	}
}

func TestInsertWithForcedCollisions(t *testing.T) {
	strategies := []ProbeStrategy{Linear, Quadratic, Double}
	for _, strategy := range strategies {
		strategy := strategy
		t.Run(strategyName(strategy), func(t *testing.T) {
			size := 11
			ht := NewHashMap[string, int](size, strategy, zeroHash[string], oneHash[string])
			keys := []string{"a", "b", "c", "d", "e"}
			for i, k := range keys {
				if ok := ht.Insert(k, i); !ok {
					t.Fatalf("[%s] Insert(%s) = false, want true", strategyName(strategy), k)
				}
			}
			for i, k := range keys {
				v, found := ht.Search(k)
				if !found || v != i {
					t.Fatalf("[%s] Search(%s) = (%d, %v), want (%d, true)", strategyName(strategy), k, v, found, i)
				}
			}
		})
	}
}

func TestInsertFullTableRejectsNewKey(t *testing.T) {
	size := 5
	ht := newStringTable(size, Linear)
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("k%d", i)
		if ok := ht.Insert(key, i); !ok {
			t.Fatalf("Insert(%s) = false, want true", key)
		}
	}
	if ok := ht.Insert("overflow", 0); ok {
		t.Fatalf("Insert into full table = true, want false")
	}
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("k%d", i)
		v, found := ht.Search(key)
		if !found || v != i {
			t.Fatalf("Search(%s) = (%d, %v), want (%d, true)", key, v, found, i)
		}
	}
}

func TestAllStrategiesBasicRoundTrip(t *testing.T) {
	strategies := []ProbeStrategy{Linear, Quadratic, Double}
	for _, strategy := range strategies {
		strategy := strategy
		t.Run(strategyName(strategy), func(t *testing.T) {
			ht := newStringTable(32, strategy)
			keys := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape"}
			for i, k := range keys {
				if ok := ht.Insert(k, i); !ok {
					t.Fatalf("[%s] Insert(%s) = false, want true", strategyName(strategy), k)
				}
			}
			for i, k := range keys {
				v, found := ht.Search(k)
				if !found || v != i {
					t.Fatalf("[%s] Search(%s) = (%d, %v), want (%d, true)", strategyName(strategy), k, v, found, i)
				}
			}
			for i := 0; i < len(keys)/2; i++ {
				if ok := ht.Delete(keys[i]); !ok {
					t.Fatalf("[%s] Delete(%s) = false, want true", strategyName(strategy), keys[i])
				}
			}
			for i := 0; i < len(keys)/2; i++ {
				if _, found := ht.Search(keys[i]); found {
					t.Fatalf("[%s] Search(%s) after delete found = true, want false", strategyName(strategy), keys[i])
				}
			}
			for i := len(keys) / 2; i < len(keys); i++ {
				v, found := ht.Search(keys[i])
				if !found || v != i {
					t.Fatalf("[%s] Search(%s) = (%d, %v), want (%d, true)", strategyName(strategy), keys[i], v, found, i)
				}
			}
		})
	}
}

func TestGenericIntKeys(t *testing.T) {
	ht := NewHashMap[int, string](16, Double, identInt, mulInt)
	for i := 0; i < 10; i++ {
		if ok := ht.Insert(i, fmt.Sprintf("v%d", i)); !ok {
			t.Fatalf("Insert(%d) = false, want true", i)
		}
	}
	for i := 0; i < 10; i++ {
		v, found := ht.Search(i)
		want := fmt.Sprintf("v%d", i)
		if !found || v != want {
			t.Fatalf("Search(%d) = (%q, %v), want (%q, true)", i, v, found, want)
		}
	}
	if ok := ht.Delete(5); !ok {
		t.Fatalf("Delete(5) = false, want true")
	}
	if _, found := ht.Search(5); found {
		t.Fatalf("Search(5) after delete found = true, want false")
	}
}

func TestSearchReturnsZeroValueOnMiss(t *testing.T) {
	ht := newStringTable(8, Linear)
	v, found := ht.Search("missing")
	if found {
		t.Fatalf("found = true, want false")
	}
	if v != 0 {
		t.Fatalf("value on miss = %d, want zero value 0", v)
	}
}
