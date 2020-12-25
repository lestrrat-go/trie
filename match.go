package trie

import "context"

// Compile compiles a MatchTree from a Tree.
func Compile(tr *Tree) *MatchTree {
	mt := &MatchTree{
		root: tr.root,
	}
	mt.root.Value = &matchData{fail: mt.root}
	for n0 := range tr.Iterate(context.TODO()) {
		for n1 := range n0.Child.Iterate(context.TODO()) {
			mt.fillFail(n1, n0)
		}
	}
	return mt
}

func (mt *MatchTree) fillFail(curr, parent *Node) {
	d := &matchData{value: curr.Value}
	curr.Value = d
	if parent == mt.root {
		d.fail = mt.root
		return
	}
	d.fail = mt.nextNode(mt.failNode(parent), curr.label)
}

func (mt *MatchTree) failNode(node *Node) *Node {
	fail := (node.Value.(*matchData)).fail
	if fail == nil {
		return mt.root
	}
	return fail
}

func (mt *MatchTree) nextNode(node *Node, l Label) *Node {
	for {
		if next := node.Get(l); next != nil {
			return next
		}
		if node == mt.root {
			return mt.root
		}
		node = mt.failNode(node)
	}
}

// MatchAll matches text and return all matched data.I
func (mt *MatchTree) MatchAll(key Key, matches []Match) []Match {
	m := mt.Matcher()
	for l := range key.Iterate() {
		matches = m.Next(l, matches)
	}
	return matches
}

// Matcher implements an iterator to match.
type Matcher struct {
	mt   *MatchTree
	curr *Node
}

// Matcher creates a new Matcher which is matching context.
func (mt *MatchTree) Matcher() *Matcher {
	return (&Matcher{mt: mt}).Reset()
}

// Reset resets matchin context.
func (m *Matcher) Reset() *Matcher {
	m.curr = m.mt.root
	return m
}

// Next appends a rune to match string, then get matches.
func (m *Matcher) Next(l Label, matches []Match) []Match {
	m.curr = m.mt.nextNode(m.curr, l)
	if m.curr == m.mt.root {
		return nil
	}
	for n := m.curr; n != m.mt.root; {
		d := n.Value.(*matchData)
		if d.value != nil {
			matches = append(matches, Match{
				Value: d.value,
			})
		}
		n = d.fail
	}
	return matches
}
