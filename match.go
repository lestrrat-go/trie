package trie

// Match is matched data.
type Match struct {
	Index   int
	Pattern string
	Value   interface{}
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

// Match matches text and return all matched data.
func (mt *MatchTree) Match(text string) []Match {
	// TODO:
	return nil
}
