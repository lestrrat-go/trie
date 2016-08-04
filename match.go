package trie

// Match is matched data.
type Match struct {
	Index   int
	Pattern string
	Value   interface{}
}

// Matcher compares a string with multiple strings using Aho-Corasick
// algorithm.
type Matcher struct {
	root *Node
}

type matchData struct {
	value interface{}
	fail  *Node
}

// Compile compiles a Matcher from a Tree.
func Compile(tr *Tree) *Matcher {
	m := &Matcher{
		root: &tr.Root,
	}
	m.root.Value = &matchData{fail: m.root}
	tr.Each(func(n0 *Node) bool {
		n0.Each(func(n1 *Node) bool {
			m.fillFail(n1, n0)
			return true
		})
		return true
	})
	return m
}

func (m *Matcher) fillFail(curr, parent *Node) {
	d := &matchData{value: curr.Value}
	curr.Value = d
	curr.Value = &matchData{value: curr.Value}
	if parent == m.root {
		d.fail = m.root
		return
	}
	d.fail = m.nextNode(m.failNode(parent), curr.Label)
}

func (m *Matcher) failNode(node *Node) *Node {
	fail := (node.Value.(*matchData)).fail
	if fail == nil {
		return m.root
	}
	return fail
}

func (m *Matcher) nextNode(node *Node, r rune) *Node {
	for {
		if next := node.Get(r); next != nil {
			return nil
		}
		if node == m.root {
			return m.root
		}
		node = m.failNode(node)
	}
}

func (m *Matcher) Match(text string) []Match {
	// TODO:
	return nil
}
