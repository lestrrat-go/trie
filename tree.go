package trie

import "container/list"

// New creates a Tree.
func New() *Tree {
	return &Tree{
		root: &Node{},
	}
}

// Get retrieves a value for key.
func (tr *Tree) Get(key Key) *Node {
	n := tr.root
	for label := range key.Iterate() {
		n = n.Get(label)
		if n == nil {
			return nil
		}
	}
	return n
}

// Put stores a pair of key and value.
func (tr *Tree) Put(key Key, value interface{}) *Node {
	n := tr.root
	for label := range key.Iterate() {
		var f bool
		n, f = n.Dig(label)
		if f {
			tr.nc++
		}
	}
	n.Value = value
	return n
}

// Each processes all nodes in width first.
func (tr *Tree) Each(proc NodeProc) {
	q := list.New()
	q.PushBack(tr.root)
	for q.Len() != 0 {
		f := q.Front()
		q.Remove(f)
		curr := f.Value.(*Node)
		if !proc(curr) {
			break
		}
		curr.Child.Each(func(n *Node) bool {
			q.PushBack(n)
			return true
		})
	}
}

// Get finds a child node which Label matches r.
func (n *Node) Get(l Label) *Node {
	n = n.Child
	for n != nil {
		switch l.Compare(n.label) {
		case 0:
			return n
		case -1:
			n = n.Low
		default:
			n = n.High
		}
	}
	return nil
}

// Dig finds a child node which Label matches r. Or create a new one when there
// are no nodes.
func (n *Node) Dig(l Label) (node *Node, isNew bool) {
	if n.Child == nil {
		n.Child = &Node{label: l}
		n.cc = 1
		return n.Child, true
	}
	m := n
	n = n.Child
	for {
		switch l.Compare(n.label) {
		case 0:
			return n, false
		case -1:
			if n.Low == nil {
				n.Low = &Node{label: l}
				m.cc++
				return n.Low, true
			}
			n = n.Low
		default:
			if n.High == nil {
				n.High = &Node{label: l}
				m.cc++
				return n.High, true
			}
			n = n.High
		}
	}
}

// Balance balances children nodes.
func (n *Node) Balance() {
	if n.Child == nil {
		return
	}
	nodes := make([]*Node, 0, n.cc)
	n.Child.Each(func(m *Node) bool {
		nodes = append(nodes, m)
		return true
	})
	n.Child = balanceNodes(nodes, 0, len(nodes))
}

// Each processes all sibiling nodes with proc.
func (n *Node) Each(proc NodeProc) bool {
	if n == nil {
		return true
	}
	if !n.Low.Each(proc) || !proc(n) || !n.High.Each(proc) {
		return false
	}
	return true
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
