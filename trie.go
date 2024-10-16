// Package trie implements a trie that allows a generic key to point to
// stored data. Whereas many trie implementations are optimized for string
// based keys, this implementation allows keys to be a sequence of "labels".
package trie

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"
)

// Tokenizer is an object that tokenize a L into individual keys.
// For example, a string tokenizer would split a string into individual runes.
type Tokenizer[L any, K cmp.Ordered] interface {
	Tokenize(L) ([]K, error)
}

// TokenizeFunc is a function that implements the Tokenizer interface
type TokenizeFunc[L any, K cmp.Ordered] func(L) ([]K, error)

func (f TokenizeFunc[L, K]) Tokenize(in L) ([]K, error) {
	return f(in)
}

// Trie is a trie that accepts arbitrary Key types as its input.
//
// L represents the "label", the input that is used to Get/Set/Delete
// a value from the trie.
//
// K represents the "key", the individual components that are associated
// with the nodes in the trie.
//
// V represents the "value", the data that is stored in the trie.
// Data is stored at the leaf nodes of the trie.
type Trie[L any, K cmp.Ordered, V any] struct {
	mu        sync.RWMutex
	root      *node[K, V]
	tokenizer Tokenizer[L, K]
}

// Node represents an individual node in the trie.
type Node[K cmp.Ordered, V any] interface {
	// Key returns the key associated with this node
	Key() K

	// Value returns the value associated with this node
	Value() V

	// Children returns the immediate children nodes of this node
	Children() []Node[K, V]

	// First returns the first child of this node
	First() Node[K, V]

	// AddChild adds a child to this node
	AddChild(Node[K, V])

	// Parent returns the parent of this node
	Parent() Node[K, V]

	// Ancestors returns a sequence of ancestors of this node.
	// The first element is the root element, progressing all the way
	// up to the parent of this node.
	Ancestors() []Node[K, V]
}

// New creates a new Trie object.
func New[L any, K cmp.Ordered, V any](tokenizer Tokenizer[L, K]) *Trie[L, K, V] {
	node := newNode[K, V]()
	node.isRoot = true
	return &Trie[L, K, V]{
		root:      node,
		tokenizer: tokenizer,
	}
}

// Get returns the value associated with `key`. The second return value
// indicates if the value was found.
func (t *Trie[L, K, V]) Get(key L) (V, bool) {
	var zero V
	tokens, err := t.tokenizer.Tokenize(key)
	if err != nil {
		return zero, false
	}

	t.mu.RLock()
	defer t.mu.RUnlock()
	node, ok := getNode(t.root, tokens)
	if !ok {
		return zero, false
	}
	return node.Value(), true
}

func (t *Trie[L, K, V]) GetNode(key L) (Node[K, V], bool) {
	tokens, err := t.tokenizer.Tokenize(key)
	if err != nil {
		return nil, false
	}

	t.mu.RLock()
	defer t.mu.RUnlock()
	return getNode(t.root, tokens)
}

func getNode[K cmp.Ordered, V any](root Node[K, V], tokens []K) (Node[K, V], bool) {
	if len(tokens) > 0 {
		for _, child := range root.Children() {
			if child.Key() == tokens[0] {
				// found the current token in the children.
				if len(tokens) == 1 {
					// this is the node we're looking for
					return child, true
				}
				// we need to traverse down the trie
				return getNode[K, V](child, tokens[1:])
			}
		}
	}

	// if we got here, that means we couldn't find a common ancestor
	return nil, false
}

// Delete removes data associated with `key`. It returns true if the value
// was found and deleted, false otherwise
func (t *Trie[L, K, V]) Delete(key L) bool {
	tokens, err := t.tokenizer.Tokenize(key)
	if err != nil {
		return false
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	return delete[K, V](t.root, tokens)
}

func delete[K cmp.Ordered, V any](root *node[K, V], tokens []K) bool {
	if len(tokens) <= 0 {
		return false
	}

	for i, child := range root.children {
		if child.Key() == tokens[0] {
			if len(tokens) == 1 {
				// this is the node we're looking for
				root.children = slices.Delete(root.children, i, i+1)
				return true
			}

			// we need to traverse down the trie
			if delete[K, V](child, tokens[1:]) {
				if len(child.children) == 0 {
					root.children = slices.Delete(root.children, i, i+1)
				}
				return true
			}
			return false
		}
	}

	return false
}

// Put sets `key` to point to data `value`.
func (t *Trie[L, K, V]) Put(key L, value V) error {
	tokens, err := t.tokenizer.Tokenize(key)
	if err != nil {
		return fmt.Errorf(`failed to tokenize key: %w`, err)
	}

	node := t.root
	t.mu.Lock()
	defer t.mu.Unlock()
	put[K, V](node, tokens, value)
	return nil
}

func put[K cmp.Ordered, V any](root Node[K, V], tokens []K, value V) {
	if len(tokens) == 0 {
		return
	}

	for _, token := range tokens {
		for _, child := range root.Children() {
			if child.Key() == token {
				// found the current token in the children.
				// we need to traverse down the trie
				put[K, V](child, tokens[1:], value)
				return
			}
		}
	}

	// if we got here, that means we couldn't find a common ancestor

	// the first token has already been consumed, create a new node,
	var newRoot *node[K, V]
	var cur *node[K, V]
	for _, token := range tokens { // duplicate token?
		newNode := newNode[K, V]()
		newNode.key = token
		if cur == nil {
			newRoot = newNode
		} else {
			cur.AddChild(newNode)
		}
		cur = newNode
	}
	// cur holds the last element.
	cur.value = value

	root.AddChild(newRoot)
}

type node[K cmp.Ordered, V any] struct {
	mu       sync.RWMutex
	key      K
	value    V
	isRoot   bool
	children []*node[K, V]
	parent   *node[K, V]
}

func newNode[K cmp.Ordered, V any]() *node[K, V] {
	return &node[K, V]{}
}

func (n *node[K, V]) Key() K {
	return n.key
}

func (n *node[K, V]) Value() V {
	return n.value
}

func (n *node[K, V]) Parent() Node[K, V] {
	return n.parent
}

func (n *node[K, V]) Ancestors() []Node[K, V] {
	var ancestors []Node[K, V]
	for {
		n = n.parent
		if n == nil {
			break
		}
		ancestors = append(ancestors, n)
	}
	return ancestors
}

func (n *node[K, V]) Children() []Node[K, V] {
	n.mu.RLock()
	children := make([]Node[K, V], 0, len(n.children))
	for _, child := range n.children {
		children = append(children, child)
	}
	n.mu.RUnlock()
	return children
}

func (n *node[K, V]) First() Node[K, V] {
	if len(n.children) == 0 {
		return nil
	}
	return n.children[0]
}

func (n *node[K, V]) AddChild(child Node[K, V]) {
	n.mu.Lock()
	// This is kind of gross, but we're only covering *node[T] with
	// Node[T] interface because we don't want the users to instantiate
	// their own nodes... so this type conversion is safe.
	//nolint:forcetypeassert
	raw := child.(*node[K, V])
	raw.parent = n
	n.children = append(n.children, raw)
	sort.Slice(n.children, func(i, j int) bool {
		return n.children[i].Key() < n.children[j].Key()
	})
	n.mu.Unlock()
}

type VisitMetadata struct {
	Depth int
}

type Visitor[K cmp.Ordered, V any] interface {
	Visit(Node[K, V], VisitMetadata) bool
}

type VisitFunc[K cmp.Ordered, V any] func(Node[K, V], VisitMetadata) bool

func (f VisitFunc[K, V]) Visit(n Node[K, V], m VisitMetadata) bool {
	return f(n, m)
}

func Walk[L any, K cmp.Ordered, V any](trie *Trie[L, K, V], v Visitor[K, V]) {
	var meta VisitMetadata
	meta.Depth = 1
	walk(trie.root, v, meta)
}

func walk[K cmp.Ordered, V any](node Node[K, V], v Visitor[K, V], meta VisitMetadata) {
	for _, child := range node.Children() {
		if !v.Visit(child, meta) {
			break
		}
		walk(child, v, VisitMetadata{Depth: meta.Depth + 1})
	}
}

type dumper[K cmp.Ordered, V any] struct{}

func (dumper[K, V]) Visit(n Node[K, V], meta VisitMetadata) bool {
	var sb strings.Builder
	for i := 0; i < meta.Depth; i++ {
		sb.WriteString("  ")
	}

	fmt.Fprintf(&sb, "%q: %v", fmt.Sprintf("%v", n.Key()), n.Value())
	fmt.Println(sb.String())
	return true
}

func Dumper[K cmp.Ordered, V any]() Visitor[K, V] {
	return dumper[K, V]{}
}
