package trie

// Tree implemnets ternary trie-tree.
type Tree struct {
	// Root is root of the tree. Only Child is valid.
	Root Node

	// nc means node counts
	nc int
}

// New creates a Tree.
func New() *Tree {
	return new(Tree)
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
		var f bool
		n, f = n.Dig(r)
		if f {
			tr.nc++
		}
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
	cc    int // count of children.
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
		n.cc = 1
		return n.Child, true
	}
	m := n
	n = n.Child
	for {
		switch {
		case r == n.Label:
			return n, false
		case r < n.Label:
			if n.Low == nil {
				n.Low = &Node{Label: r}
				m.cc++
				return n.Low, true
			}
			n = n.Low
		default:
			if n.High == nil {
				n.High = &Node{Label: r}
				m.cc++
				return n.High, true
			}
			n = n.High
		}
	}
}

// balance balances children nodes.
func (n *Node) Balance() {
	if n.Child == nil {
		return
	}
	nodes := make([]*Node, 0, n.cc)
	nodes = enumerateNodes(nodes, n.Child)
	n.Child = balanceNodes(nodes, 0, len(nodes))
}

func enumerateNodes(nodes []*Node, n *Node) []*Node {
	if n == nil {
		return nodes
	}
	nodes = enumerateNodes(nodes, n.Low)
	nodes = append(nodes, n)
	nodes = enumerateNodes(nodes, n.High)
	return nodes
}

func balanceNodes(nodes []*Node, s, e int) *Node {
	c := e - s
	switch {
	case c <= 0:
		return nil
	case c == 1:
		n := nodes[s]
		n.Low = nil
		n.High = nil
		return n
	case c == 2:
		n := nodes[s]
		n.High = nodes[s+1]
		n.Low = nil
		return n
	default:
		m := (s + e) / 2
		n := nodes[m]
		n.Low = balanceNodes(nodes, s, m)
		n.High = balanceNodes(nodes, m+1, e)
		return n
	}
}
