package trie

// Tree implemnets ternary trie-tree.
type Tree struct {
	// Root is root of the tree. Only Child is valid.
	Root Node

	// nc means node counts
	nc int
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

// NodeProc provides procedure for nodes.
type NodeProc func(*Node) bool

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


