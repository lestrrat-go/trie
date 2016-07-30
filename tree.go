package trie

// Tree implemnets ternary trie-tree.
type Tree struct {
	// Root is root of the tree. Only Child is valid.
	Root Node
}

// Get retrieves a value for key.
func (tr *Tree) Get(key string) *Node {
	n := &tr.Root
	for _, r := range key {
		n = n.Get(r)
		if n == nil {
			return nil
		}
	}
	return n
}

// Put stores a pair of key and value.
func (tr *Tree) Put(key string, value interface{}) *Node {
	n := &tr.Root
	for _, r := range key {
		n, _ = n.Dig(r)
	}
	n.Value = value
	return n
}

// Node implemnets node of ternary trie-tree.
type Node struct {
	Label rune
	Value interface{}
	Low   *Node
	High  *Node
	Child *Node
}

// Get finds a child node which Label matches r.
func (n *Node) Get(r rune) *Node {
	n = n.Child
	for n != nil {
		switch {
		case r == n.Label:
			return n
		case r < n.Label:
			n = n.Low
		default:
			n = n.High
		}
	}
	return nil
}

// Dig finds a child node which Label matches r. Or create a new one when there
// are no nodes.
func (n *Node) Dig(r rune) (node *Node, isNew bool) {
	if n.Child == nil {
		n.Child = &Node{Label: r}
		return n.Child, true
	}
	n = n.Child
	for {
		switch {
		case r == n.Label:
			return n, false
		case r < n.Label:
			if n.Low == nil {
				n.Low = &Node{Label: r}
				return n.Low, true
			}
			n = n.Low
		default:
			if n.High == nil {
				n.High = &Node{Label: r}
				return n.High, true
			}
			n = n.High
		}
	}
}
