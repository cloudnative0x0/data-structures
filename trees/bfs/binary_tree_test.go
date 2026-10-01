package bfs

import "testing"

func TestBinaryTree(t *testing.T) {
	root := New([]int{1, 2, 3, 4, 5, 6, 7})

	var bfs []int
	root.TraverseBFS(func(v int) { bfs = append(bfs, v) })
	want := []int{1, 2, 3, 4, 5, 6, 7}
	if !equal(bfs, want) {
		t.Fatalf("BFS после вставки = %v, хотел %v", bfs, want)
	}

	if root.Left.Value != 2 || root.Right.Value != 3 {
		t.Fatalf("первый уровень собран неверно: left=%d right=%d", root.Left.Value, root.Right.Value)
	}
	if root.Left.Left.Value != 4 || root.Left.Right.Value != 5 {
		t.Fatalf("второй уровень (слева) собран неверно")
	}

	if found := root.Search(6); found == nil || found.Value != 6 {
		t.Fatalf("Search(6) = %v", found)
	}
	if found := root.Search(42); found != nil {
		t.Fatalf("Search(42) = %v, ожидал nil", found)
	}

	root = Delete(root, 3)

	if root.Right.Value != 7 {
		t.Fatalf("после Delete(3) на месте узла должно быть значение 7, получили %d", root.Right.Value)
	}
	if root.Right.Right != nil {
		t.Fatalf("последний узел (7) должен был отрезаться, но он всё ещё на месте")
	}
	if root.Search(3) != nil {
		t.Fatalf("значение 3 всё ещё находится после удаления")
	}

	for _, v := range []int{7, 6, 5, 4, 2} {
		root = Delete(root, v)
	}
	if root == nil || root.Value != 1 || root.Left != nil || root.Right != nil {
		t.Fatalf("ожидал единственный узел со значением 1, получили %+v", root)
	}
	root = Delete(root, 1)
	if root != nil {
		t.Fatalf("после удаления последнего узла дерево должно быть nil, получили %+v", root)
	}
}

func TestBinaryEmptyTree(t *testing.T) {
	if root := New(nil); root != nil {
		t.Fatalf("New(nil) = %v, ожидал nil", root)
	}
	if root := New([]int{}); root != nil {
		t.Fatalf("New([]int{}) = %v, ожидал nil", root)
	}

	var root *Node

	if found := root.Search(1); found != nil {
		t.Fatalf("Search на пустом дереве = %v, ожидал nil", found)
	}
	if got := Delete(root, 1); got != nil {
		t.Fatalf("Delete на пустом дереве = %v, ожидал nil", got)
	}

	var visited []int
	visit := func(v int) { visited = append(visited, v) }
	root.TraversePreOrder(visit)
	root.TraverseInOrder(visit)
	root.TraversePostOrder(visit)
	root.TraverseBFS(visit)
	if len(visited) != 0 {
		t.Fatalf("обходы пустого дерева что-то посетили: %v", visited)
	}
}

func TestBinaryInsertOnNilRoot(t *testing.T) {
	got := Insert(nil, 5)
	if got == nil || got.Value != 5 {
		t.Fatalf("Insert(nil, 5) = %v, ожидал узел со значением 5", got)
	}
}

func TestBinarySingleElement(t *testing.T) {
	root := New([]int{7})
	if root == nil || root.Value != 7 || root.Left != nil || root.Right != nil {
		t.Fatalf("New([]int{7}) = %+v, ожидал единственный узел 7", root)
	}

	root = Delete(root, 7)
	if root != nil {
		t.Fatalf("Delete последнего узла = %v, ожидал nil", root)
	}
}

func TestBinaryDeleteNonExistent(t *testing.T) {
	root := New([]int{1, 2, 3, 4, 5})

	got := Delete(root, 999)
	if got != root {
		t.Fatalf("Delete несуществующего значения должен вернуть то же дерево без изменений")
	}

	var bfs []int
	got.TraverseBFS(func(v int) { bfs = append(bfs, v) })
	if !equal(bfs, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("дерево изменилось после неудачного Delete: %v", bfs)
	}
}

func manualTree() *Node {
	n4 := &Node{Value: 4}
	n2 := &Node{Value: 2, Left: n4}
	n3 := &Node{Value: 3}
	return &Node{Value: 1, Left: n2, Right: n3}
}

func TestBinaryTraversalOrders(t *testing.T) {
	root := manualTree()

	cases := []struct {
		name string
		run  func(func(int))
		want []int
	}{
		{"PreOrder", root.TraversePreOrder, []int{1, 2, 4, 3}},
		{"InOrder", root.TraverseInOrder, []int{4, 2, 1, 3}},
		{"PostOrder", root.TraversePostOrder, []int{4, 2, 3, 1}},
		{"BFS", root.TraverseBFS, []int{1, 2, 3, 4}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got []int
			c.run(func(v int) { got = append(got, v) })
			if !equal(got, c.want) {
				t.Fatalf("%s = %v, хотел %v", c.name, got, c.want)
			}
		})
	}
}

func isComplete(root *Node) bool {
	queue := []*Node{root}
	seenGap := false

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]

		if n == nil {
			seenGap = true
			continue
		}
		if seenGap {
			return false
		}
		queue = append(queue, n.Left, n.Right)
	}

	return true
}

func TestBinaryStaysCompleteAfterInsertsAndDeletes(t *testing.T) {
	var root *Node
	for v := 1; v <= 15; v++ {
		root = Insert(root, v)
		if !isComplete(root) {
			t.Fatalf("дерево перестало быть полным после вставки %d", v)
		}
	}

	for _, v := range []int{8, 1, 15, 7, 3, 12} {
		root = Delete(root, v)
		if !isComplete(root) {
			t.Fatalf("дерево перестало быть полным после удаления %d", v)
		}
		if root.Search(v) != nil {
			t.Fatalf("значение %d всё ещё найдено после удаления", v)
		}
	}

	var remaining []int
	root.TraverseBFS(func(v int) { remaining = append(remaining, v) })
	if len(remaining) != 9 {
		t.Fatalf("осталось %d значений, ожидал 9: %v", len(remaining), remaining)
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
