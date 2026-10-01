package bfs

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func New(slice []int) *Node {
	var root *Node

	for _, v := range slice {
		root = Insert(root, v)
	}

	return root
}

func Insert(root *Node, value int) *Node {
	newNode := &Node{Value: value}
	if root == nil {
		return newNode
	}

	queue := []*Node{root}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.Left == nil {
			cur.Left = newNode

			return root
		}
		queue = append(queue, cur.Left)

		if cur.Right == nil {
			cur.Right = newNode

			return root
		}
		queue = append(queue, cur.Right)
	}

	return root
}

func (n *Node) Search(value int) *Node {
	if n == nil {
		return nil
	}

	queue := []*Node{n}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.Value == value {
			return cur
		}
		if cur.Left != nil {
			queue = append(queue, cur.Left)
		}
		if cur.Right != nil {
			queue = append(queue, cur.Right)
		}
	}

	return nil
}

func Delete(root *Node, value int) *Node {
	if root == nil {
		return nil
	}

	target := root.Search(value)
	if target == nil {
		return root
	}

	var lastParent, last *Node
	queue := []*Node{root}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.Left != nil {
			lastParent, last = cur, cur.Left
			queue = append(queue, cur.Left)
		}
		if cur.Right != nil {
			lastParent, last = cur, cur.Right
			queue = append(queue, cur.Right)
		}
	}

	if last == nil {
		return nil
	}

	target.Value = last.Value

	if lastParent.Right == last {
		lastParent.Right = nil
	} else {
		lastParent.Left = nil
	}

	return root
}

func (n *Node) TraversePreOrder(visit func(int)) {
	if n == nil {
		return
	}

	visit(n.Value)
	n.Left.TraversePreOrder(visit)
	n.Right.TraversePreOrder(visit)
}

func (n *Node) TraverseInOrder(visit func(int)) {
	if n == nil {
		return
	}

	n.Left.TraverseInOrder(visit)
	visit(n.Value)
	n.Right.TraverseInOrder(visit)
}

func (n *Node) TraversePostOrder(visit func(int)) {
	if n == nil {
		return
	}

	n.Left.TraversePostOrder(visit)
	n.Right.TraversePostOrder(visit)
	visit(n.Value)
}

func (n *Node) TraverseBFS(visit func(int)) {
	if n == nil {
		return
	}

	queue := []*Node{n}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		visit(cur.Value)
		if cur.Left != nil {
			queue = append(queue, cur.Left)
		}
		if cur.Right != nil {
			queue = append(queue, cur.Right)
		}
	}
}
