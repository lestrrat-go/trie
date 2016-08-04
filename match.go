package trie

// Match is matched data.
type Match struct {
	Value interface{}
}

// MatchTree compares a string with multiple strings using Aho-Corasick
// algorithm.
type MatchTree struct {
	root *Node
}

type matchData struct {
	value interface{}
	fail  *Node
}

// Compile compiles a MatchTree from a Tree.
func Compile(tr *Tree) *MatchTree {
	mt := &MatchTree{
		root: &tr.Root,
	}
	mt.root.Value = &matchData{fail: mt.root}
	tr.Each(func(n0 *Node) bool {
		n0.Each(func(n1 *Node) bool {
			mt.fillFail(n1, n0)
			return true
		})
		return true
	})
	return mt
}

func (mt *MatchTree) fillFail(curr, parent *Node) {
	d := &matchData{value: curr.Value}
	curr.Value = d
	curr.Value = &matchData{value: curr.Value}
	if parent == mt.root {
		d.fail = mt.root
		return
	}
	d.fail = mt.nextNode(mt.failNode(parent), curr.Label)
}

func (mt *MatchTree) failNode(node *Node) *Node {
	fail := (node.Value.(*matchData)).fail
	if fail == nil {
		return mt.root
	}
	return fail
}

func (mt *MatchTree) nextNode(node *Node, r rune) *Node {
	for {
		if next := node.Get(r); next != nil {
			return nil
		}
		if node == mt.root {
			return mt.root
		}
		node = mt.failNode(node)
	}
}

// MatchAll matches text and return all matched data.I
func (mt *MatchTree) MatchAll(text string, matches[]Match) []Match {
	m := mt.Matcher()
	for _, r := range text {
		matches = m.Next(r, matches)
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
func (m *Matcher) Next(r rune, matches []Match) []Match {
	m.curr = m.mt.nextNode(m.curr, r)
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
