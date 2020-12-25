package trie

import (
	"context"
)

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

func iterateTreeBFS(ctx context.Context, n *Node, ch chan *Node) {
	defer close(ch)
	nodes := []*Node{n}
	for len(nodes) > 0 {
		q := nodes[0]
		nodes = nodes[1:]
		select {
		case <-ctx.Done():
			return
		case ch <- q:
		}

		if child := q.Child; child == nil {
			continue
		}

		for c := range q.Child.Iterate(ctx) {
			nodes = append(nodes, c)
		}
	}
}

func (tr *Tree) Iterate(ctx context.Context) <-chan *Node {
	ch := make(chan *Node)
	go iterateTreeBFS(ctx, tr.root, ch)
	return ch
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
	for m := range n.Child.Iterate(context.TODO()) {
		nodes = append(nodes, m)
	}
	n.Child = balanceNodes(nodes, 0, len(nodes))
}

func iterateNodesBFS(ctx context.Context, n, root *Node, ch chan *Node, reverse bool) {
	//	if root != nil && n == root {
	defer close(ch)
	//	}

	if n == nil {
		return
	}

	nodes := []*Node{n}
	for len(nodes) > 0 {
		q := nodes[0]
		nodes = nodes[1:]
		select {
		case <-ctx.Done():
			return
		case ch <- q:
		}

		if reverse {
			if next := q.High; next != nil {
				nodes = append(nodes, next)
			}
			if next := q.Low; next != nil {
				nodes = append(nodes, next)
			}
		} else {
			if next := q.Low; next != nil {
				nodes = append(nodes, next)
			}
			if next := q.High; next != nil {
				nodes = append(nodes, next)
			}
		}
	}
}

func iterateNodesDFS(ctx context.Context, n, root *Node, ch chan *Node) {
	if root != nil && n == root {
		defer close(ch)
	}

	if n == nil {
		return
	}

	select {
	case <-ctx.Done():
		return
	default:
		if next := n.Low; next != nil {
			iterateNodesDFS(ctx, next, root, ch)
		}
	}

	select {
	case <-ctx.Done():
		return
	case ch <- n:
	}

	select {
	case <-ctx.Done():
		return
	default:
		if next := n.High; next != nil {
			iterateNodesDFS(ctx, next, root, ch)
		}
	}
}

type iterationStrategy struct{}

const (
	iterateStrategyDFS = iota
	iterateStrategyBFS
	iterateStrategyBFSReverse
)

func WithBFS(ctx context.Context) context.Context {
	return context.WithValue(ctx, iterationStrategy{}, iterateStrategyBFS)
}

func WithBFSReverse(ctx context.Context) context.Context {
	return context.WithValue(ctx, iterationStrategy{}, iterateStrategyBFSReverse)
}

func (n *Node) Iterate(ctx context.Context) <-chan *Node {
	ch := make(chan *Node)
	if n == nil {
		close(ch)
		return ch
	}

	switch st := ctx.Value(iterationStrategy{}); st {
	case iterateStrategyDFS, nil:
		go iterateNodesDFS(ctx, n, n, ch)
	case iterateStrategyBFS:
		go iterateNodesBFS(ctx, n, n, ch, false)
	case iterateStrategyBFSReverse:
		go iterateNodesBFS(ctx, n, n, ch, true)
	default:
		panic("unknown iteration strategy")
	}

	return ch
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
