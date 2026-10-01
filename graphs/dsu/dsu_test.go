package dsu

import (
	"errors"
	"math/rand"
	"testing"
)

func TestUnionConnectedSizeAndComponents(t *testing.T) {
	d := mustDSU(t, 6)
	if d.Components() != 6 {
		t.Fatalf("Components() = %d, want 6", d.Components())
	}

	for _, pair := range [][2]int{{0, 1}, {2, 3}, {1, 2}} {
		merged, err := d.Union(pair[0], pair[1])
		if err != nil || !merged {
			t.Fatalf("Union%v = (%v, %v), want (true, nil)", pair, merged, err)
		}
	}

	connected, err := d.Connected(0, 3)
	if err != nil || !connected {
		t.Fatalf("Connected(0, 3) = (%v, %v)", connected, err)
	}
	size, err := d.Size(2)
	if err != nil || size != 4 {
		t.Fatalf("Size(2) = (%d, %v), want (4, nil)", size, err)
	}
	if d.Components() != 3 {
		t.Fatalf("Components() = %d, want 3", d.Components())
	}

	merged, err := d.Union(0, 3)
	if err != nil || merged {
		t.Fatalf("repeated Union = (%v, %v), want (false, nil)", merged, err)
	}
}

func TestFindCompressesPath(t *testing.T) {
	d := mustDSU(t, 8)
	for _, pair := range [][2]int{{0, 1}, {2, 3}, {0, 2}, {4, 5}, {6, 7}, {4, 6}, {0, 4}} {
		if _, err := d.Union(pair[0], pair[1]); err != nil {
			t.Fatalf("Union%v error = %v", pair, err)
		}
	}

	root, err := d.Find(7)
	if err != nil {
		t.Fatalf("Find(7) error = %v", err)
	}
	if d.parent[7] != root {
		t.Fatalf("path was not compressed: parent[7]=%d root=%d", d.parent[7], root)
	}
}

func TestValidation(t *testing.T) {
	if _, err := New(-1); !errors.Is(err, ErrInvalidSize) {
		t.Fatalf("New(-1) error = %v", err)
	}
	d := mustDSU(t, 2)
	if _, err := d.Find(2); !errors.Is(err, ErrElementOutOfRange) {
		t.Fatalf("Find(2) error = %v", err)
	}
	if _, err := d.Union(-1, 0); !errors.Is(err, ErrElementOutOfRange) {
		t.Fatalf("Union(-1, 0) error = %v", err)
	}
	if _, err := d.Connected(0, 2); !errors.Is(err, ErrElementOutOfRange) {
		t.Fatalf("Connected(0, 2) error = %v", err)
	}
}

func TestRandomOperationsAgainstNaivePartition(t *testing.T) {
	const n = 100
	d := mustDSU(t, n)
	labels := make([]int, n)
	for i := range labels {
		labels[i] = i
	}
	rng := rand.New(rand.NewSource(42))

	for step := 0; step < 10_000; step++ {
		a, b := rng.Intn(n), rng.Intn(n)
		if rng.Intn(2) == 0 {
			oldLabel, newLabel := labels[b], labels[a]
			wantMerged := oldLabel != newLabel
			for i := range labels {
				if labels[i] == oldLabel {
					labels[i] = newLabel
				}
			}
			merged, err := d.Union(a, b)
			if err != nil || merged != wantMerged {
				t.Fatalf("step %d: Union(%d, %d) = (%v, %v), want %v", step, a, b, merged, err, wantMerged)
			}
		} else {
			got, err := d.Connected(a, b)
			want := labels[a] == labels[b]
			if err != nil || got != want {
				t.Fatalf("step %d: Connected(%d, %d) = (%v, %v), want %v", step, a, b, got, err, want)
			}
		}
	}
}

func mustDSU(t *testing.T, n int) *DSU {
	t.Helper()
	d, err := New(n)
	if err != nil {
		t.Fatalf("New(%d) error = %v", n, err)
	}
	return d
}
