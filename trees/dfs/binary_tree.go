package dfs

// Node is a node of a binary tree.
type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// New builds a complete binary tree from values given in level order.
// For example, []int{1, 2, 3} creates a root 1 with children 2 and 3.
func New(values []int) *Node {
	if len(values) == 0 {
		return nil
	}

	nodes := make([]*Node, len(values))
	for i, value := range values {
		nodes[i] = &Node{Value: value}
	}

	for i := range nodes {
		left := 2*i + 1
		right := 2*i + 2
		if left < len(nodes) {
			nodes[i].Left = nodes[left]
		}
		if right < len(nodes) {
			nodes[i].Right = nodes[right]
		}
	}

	return nodes[0]
}

// Search performs iterative depth-first search in pre-order and returns the
// first node whose value matches. The left subtree is explored before the
// right subtree.
func (n *Node) Search(value int) *Node {
	if n == nil {
		return nil
	}

	stack := []*Node{n}
	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]

		if current.Value == value {
			return current
		}

		// The stack is LIFO, so push right first to visit left first.
		if current.Right != nil {
			stack = append(stack, current.Right)
		}
		if current.Left != nil {
			stack = append(stack, current.Left)
		}
	}

	return nil
}

// TraverseDFS performs iterative pre-order DFS: root, left, right.
// A nil visit function is treated as a no-op.
func (n *Node) TraverseDFS(visit func(int)) {
	if n == nil || visit == nil {
		return
	}

	stack := []*Node{n}
	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]

		visit(current.Value)
		if current.Right != nil {
			stack = append(stack, current.Right)
		}
		if current.Left != nil {
			stack = append(stack, current.Left)
		}
	}
}

// TraversePreOrder visits root, left subtree, then right subtree.
func (n *Node) TraversePreOrder(visit func(int)) {
	if n == nil || visit == nil {
		return
	}

	visit(n.Value)
	n.Left.TraversePreOrder(visit)
	n.Right.TraversePreOrder(visit)
}

// TraverseInOrder visits left subtree, root, then right subtree.
func (n *Node) TraverseInOrder(visit func(int)) {
	if n == nil || visit == nil {
		return
	}

	n.Left.TraverseInOrder(visit)
	visit(n.Value)
	n.Right.TraverseInOrder(visit)
}

// TraversePostOrder visits left subtree, right subtree, then root.
func (n *Node) TraversePostOrder(visit func(int)) {
	if n == nil || visit == nil {
		return
	}

	n.Left.TraversePostOrder(visit)
	n.Right.TraversePostOrder(visit)
	visit(n.Value)
}
