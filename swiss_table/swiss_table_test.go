package swiss_table

import (
	"fmt"
	"math/bits"
	"math/rand"
	"testing"
	"time"
)

func hashInt(k int) uint64 {
	h := uint64(k)
	h ^= h >> 33
	h *= 0xff51afd7ed558ccd
	h ^= h >> 33
	h *= 0xc4ceb9fe1a85ec53
	h ^= h >> 33
	return h
}

func TestEmptyMap(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	if m.Len() != 0 {
		t.Fatalf("Len() = %d, want 0", m.Len())
	}
	if _, ok := m.Get(42); ok {
		t.Fatalf("Get on empty map returned ok=true")
	}
	if m.Delete(42) {
		t.Fatalf("Delete on empty map returned true")
	}
}

func TestPutGetSingle(t *testing.T) {
	m := NewFastSwissMap[int, string](hashInt)
	m.Put(1, "one")
	v, ok := m.Get(1)
	if !ok || v != "one" {
		t.Fatalf("Get(1) = (%q,%v), want (\"one\",true)", v, ok)
	}
	if m.Len() != 1 {
		t.Fatalf("Len() = %d, want 1", m.Len())
	}
}

func TestPutGetMany(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	const n = 1000
	for i := 0; i < n; i++ {
		m.Put(i, i*i)
	}
	if m.Len() != n {
		t.Fatalf("Len() = %d, want %d", m.Len(), n)
	}
	for i := 0; i < n; i++ {
		v, ok := m.Get(i)
		if !ok || v != i*i {
			t.Fatalf("Get(%d) = (%d,%v), want (%d,true)", i, v, ok, i*i)
		}
	}
}

func TestOverwriteDoesNotGrowLen(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	m.Put(42, 1)
	m.Put(42, 2)
	m.Put(42, 3)
	if v, ok := m.Get(42); !ok || v != 3 {
		t.Fatalf("Get(42) = (%d,%v), want (3,true)", v, ok)
	}
	if m.Len() != 1 {
		t.Fatalf("Len() = %d, want 1 (overwrite must not add a new entry)", m.Len())
	}
}

func TestGetMissingKey(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	for i := 0; i < 100; i++ {
		m.Put(i, i)
	}
	for _, k := range []int{-1, 100, 100000, -100000} {
		if _, ok := m.Get(k); ok {
			t.Fatalf("Get(%d) = ok=true, want false (key was never inserted)", k)
		}
	}
}

func TestDeleteReducesLen(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	for i := 0; i < 10; i++ {
		m.Put(i, i)
	}
	if !m.Delete(5) {
		t.Fatalf("Delete(5) = false, want true")
	}
	if m.Len() != 9 {
		t.Fatalf("Len() = %d, want 9", m.Len())
	}
	if _, ok := m.Get(5); ok {
		t.Fatalf("Get(5) after Delete returned ok=true")
	}

	for _, k := range []int{4, 6} {
		if v, ok := m.Get(k); !ok || v != k {
			t.Fatalf("Get(%d) = (%d,%v), want (%d,true)", k, v, ok, k)
		}
	}
}

func TestDeleteMissingKeyIsNoop(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	m.Put(1, 1)
	if m.Delete(999) {
		t.Fatalf("Delete(999) = true, want false (key was never present)")
	}
	if m.Len() != 1 {
		t.Fatalf("Len() = %d, want 1 (unaffected by no-op delete)", m.Len())
	}
}

func TestDeleteThenReinsert(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	const n = 500
	for i := 0; i < n; i++ {
		m.Put(i, i)
	}
	for i := 0; i < n; i += 2 {
		if !m.Delete(i) {
			t.Fatalf("Delete(%d) = false, want true", i)
		}
	}
	if m.Len() != n/2 {
		t.Fatalf("Len() = %d, want %d", m.Len(), n/2)
	}

	for i := 0; i < n; i += 2 {
		m.Put(i, i*10)
	}
	if m.Len() != n {
		t.Fatalf("Len() = %d, want %d", m.Len(), n)
	}
	for i := 0; i < n; i++ {
		want := i
		if i%2 == 0 {
			want = i * 10
		}
		v, ok := m.Get(i)
		if !ok || v != want {
			t.Fatalf("Get(%d) = (%d,%v), want (%d,true)", i, v, ok, want)
		}
	}
}

func TestZeroValueKeysAndValues(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	m.Put(0, 0)
	v, ok := m.Get(0)
	if !ok || v != 0 {
		t.Fatalf("Get(0) = (%d,%v), want (0,true) — zero key/value must be distinguishable from absence", v, ok)
	}
	if !m.Delete(0) {
		t.Fatalf("Delete(0) = false, want true")
	}
	if _, ok := m.Get(0); ok {
		t.Fatalf("Get(0) after delete = ok=true, want false")
	}
}

func TestStringKeys(t *testing.T) {
	hashStr := func(s string) uint64 {
		var h uint64 = 14695981039346656037
		for i := 0; i < len(s); i++ {
			h ^= uint64(s[i])
			h *= 1099511628211
		}
		return h
	}
	m := NewFastSwissMap[string, int](hashStr)
	words := []string{"", "a", "ab", "abc", "swiss", "table", "quadratic", "probing"}
	for i, w := range words {
		m.Put(w, i)
	}
	for i, w := range words {
		v, ok := m.Get(w)
		if !ok || v != i {
			t.Fatalf("Get(%q) = (%d,%v), want (%d,true)", w, v, ok, i)
		}
	}
	if m.Len() != len(words) {
		t.Fatalf("Len() = %d, want %d", m.Len(), len(words))
	}
}

func TestRehashInPlaceBelowCapacityLimit(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)

	const n = 40
	for i := 0; i < n; i++ {
		m.Put(i, i)
	}
	if m.globalDepth != 0 {
		t.Fatalf("globalDepth = %d, want 0 (must not have split yet)", m.globalDepth)
	}
	if len(m.dir) != 1 {
		t.Fatalf("len(dir) = %d, want 1", len(m.dir))
	}
	for i := 0; i < n; i++ {
		if v, ok := m.Get(i); !ok || v != i {
			t.Fatalf("Get(%d) = (%d,%v), want (%d,true)", i, v, ok, i)
		}
	}
}

func TestBucketSplitTriggers(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	const n = 50_000
	for i := 0; i < n; i++ {
		m.Put(i, i)
	}
	if m.globalDepth == 0 {
		t.Fatalf("globalDepth = 0, want > 0 after %d inserts (split should have triggered)", n)
	}
	if len(m.dir) != 1<<m.globalDepth {
		t.Fatalf("len(dir) = %d, want %d (2^globalDepth)", len(m.dir), 1<<m.globalDepth)
	}
	if m.Len() != n {
		t.Fatalf("Len() = %d, want %d", m.Len(), n)
	}
	for i := 0; i < n; i++ {
		if v, ok := m.Get(i); !ok || v != i {
			t.Fatalf("Get(%d) = (%d,%v), want (%d,true)", i, v, ok, i)
		}
	}
}

func TestDirectoryBucketSharing(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	for i := 0; i < 20_000; i++ {
		m.Put(i, i)
	}
	if m.globalDepth == 0 {
		t.Fatalf("test requires at least one split to have happened")
	}
	seen := 0
	for i, b := range m.dir {
		step := uint64(1) << (m.globalDepth - b.localDepth)
		start := (uint64(i) / step) * step
		for j := start; j < start+step; j++ {
			if m.dir[j] != b {
				t.Fatalf("dir[%d] and dir[%d] should share the same bucket (localDepth=%d, globalDepth=%d)",
					i, j, b.localDepth, m.globalDepth)
			}
		}
		seen++
	}
	if seen != len(m.dir) {
		t.Fatalf("iterated %d directory slots, want %d", seen, len(m.dir))
	}
}

func TestGrowthLeftNeverNegative(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	rng := rand.New(rand.NewSource(7))
	for i := 0; i < 20_000; i++ {
		k := rng.Intn(500)
		if rng.Intn(2) == 0 {
			m.Put(k, k)
		} else {
			m.Delete(k)
		}
		for _, b := range m.dir {
			if b.growthLeft < 0 {
				t.Fatalf("iter %d: bucket growthLeft = %d, must never be negative", i, b.growthLeft)
			}
		}
	}
}

func TestNoInfiniteLoopOnTombstoneChurn(t *testing.T) {
	m := NewFastSwissMap[int, int](hashInt)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 100_000; i++ {
			k := i % 8
			m.Put(k, k)
			m.Delete(k)
		}
		_, _ = m.Get(999999)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf("Put/Delete/Get churn did not finish in time — suspected infinite loop")
	}
}

func TestInsertOrderIndependence(t *testing.T) {
	const n = 2000
	keys := make([]int, n)
	for i := range keys {
		keys[i] = i
	}

	build := func(order []int) map[int]int {
		m := NewFastSwissMap[int, int](hashInt)
		for _, k := range order {
			m.Put(k, k*2)
		}
		out := make(map[int]int, n)
		for _, k := range keys {
			v, ok := m.Get(k)
			if !ok {
				t.Fatalf("key %d missing after building with a shuffled order", k)
			}
			out[k] = v
		}
		return out
	}

	forward := append([]int(nil), keys...)
	shuffled := append([]int(nil), keys...)
	rand.New(rand.NewSource(3)).Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	a := build(forward)
	b := build(shuffled)

	if len(a) != len(b) {
		t.Fatalf("result sizes differ: %d vs %d", len(a), len(b))
	}
	for k, v := range a {
		if b[k] != v {
			t.Fatalf("key %d: forward-order value %d != shuffled-order value %d", k, v, b[k])
		}
	}
}

func TestMatchByteFindsAllPositions(t *testing.T) {
	var word uint64
	target := byte(0x42)
	wantPositions := map[int]bool{1: true, 4: true, 7: true}
	for pos := 0; pos < 8; pos++ {
		b := byte(0x11 * (pos + 1))
		if wantPositions[pos] {
			b = target
		}
		word = setByte(word, pos, b)
	}

	got := map[int]bool{}
	for m := matchByte(word, target); m != 0; m &= m - 1 {
		pos := bits.TrailingZeros64(m) / 8
		got[pos] = true
	}
	if len(got) != len(wantPositions) {
		t.Fatalf("matchByte found positions %v, want %v", got, wantPositions)
	}
	for p := range wantPositions {
		if !got[p] {
			t.Fatalf("matchByte missed position %d", p)
		}
	}
}

func TestSetByteAndByteAtRoundtrip(t *testing.T) {
	var word uint64
	for pos := 0; pos < 8; pos++ {
		word = setByte(word, pos, byte(pos+1))
	}
	for pos := 0; pos < 8; pos++ {
		if got := byteAt(word, pos); got != byte(pos+1) {
			t.Fatalf("byteAt(%d) = %d, want %d", pos, got, pos+1)
		}
	}
}

func TestProbeSeqVisitsEveryGroupExactlyOnce(t *testing.T) {
	for _, groups := range []int{1, 2, 4, 8, 16, 64} {
		mask := uint64(groups - 1)
		for h1 := uint64(0); h1 < 32; h1++ {
			seq := makeProbeSeq(h1, mask)
			seen := make(map[uint64]bool, groups)
			for i := 0; i < groups; i++ {
				if seen[seq.offset] {
					t.Fatalf("groups=%d h1=%d: offset %d visited twice within %d steps",
						groups, h1, seq.offset, groups)
				}
				seen[seq.offset] = true
				seq = seq.next()
			}
			if len(seen) != groups {
				t.Fatalf("groups=%d h1=%d: visited %d distinct offsets, want %d",
					groups, h1, len(seen), groups)
			}
		}
	}
}

func TestAgainstReferenceMap(t *testing.T) {
	for seed := int64(0); seed < 6; seed++ {
		seed := seed
		t.Run(fmt.Sprintf("seed=%d", seed), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			ref := make(map[int]int)
			m := NewFastSwissMap[int, int](hashInt)

			const ops = 60_000
			const keySpace = 4000

			for i := 0; i < ops; i++ {
				k := rng.Intn(keySpace)
				switch rng.Intn(3) {
				case 0, 1:
					v := rng.Int()
					ref[k] = v
					m.Put(k, v)
				case 2:
					delete(ref, k)
					m.Delete(k)
				}
			}

			if len(ref) != m.Len() {
				t.Fatalf("Len mismatch: ref=%d got=%d", len(ref), m.Len())
			}
			for k, want := range ref {
				got, ok := m.Get(k)
				if !ok || got != want {
					t.Fatalf("key %d: got (%d,%v), want (%d,true)", k, got, ok, want)
				}
			}
			for k := 0; k < keySpace; k++ {
				if _, inRef := ref[k]; !inRef {
					if _, ok := m.Get(k); ok {
						t.Fatalf("key %d: present in map, should be absent", k)
					}
				}
			}
		})
	}
}

func BenchmarkPut(b *testing.B) {
	m := NewFastSwissMap[int, int](hashInt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func BenchmarkGetHit(b *testing.B) {
	m := NewFastSwissMap[int, int](hashInt)
	const n = 100_000
	for i := 0; i < n; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(i % n)
	}
}

func BenchmarkGetMiss(b *testing.B) {
	m := NewFastSwissMap[int, int](hashInt)
	const n = 100_000
	for i := 0; i < n; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(n + i%n)
	}
}
