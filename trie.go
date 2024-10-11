// Package trie implements a trie that allows a generic key to point to
// stored data. Whereas many trie implementations are optimized for string
// based keys, this implementation allows keys to be a sequence of "labels".
package trie

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
	"sort"
	"strings"
	"sync"
)

// Tokenizer is an object that tokenize a L into individual keys.
// For example, a string tokenizer would split a string into individual runes.
type Tokenizer[L any, K cmp.Ordered] interface {
	Tokenize(L) (iter.Seq[K], error)
}

// TokenizeFunc is a function that implements the Tokenizer interface
type TokenizeFunc[L any, K cmp.Ordered] func(L) (iter.Seq[K], error)

func (f TokenizeFunc[L, K]) Tokenize(in L) (iter.Seq[K], error) {
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
	Key() K
	Value() V
	Children() iter.Seq[Node[K, V]]
	AddChild(Node[K, V])
}

// New creates a new Trie object.
func New[L any, K cmp.Ordered, V any](tokenizer Tokenizer[L, K]) *Trie[L, K, V] {
	return &Trie[L, K, V]{
		root:      newNode[K, V](),
		tokenizer: tokenizer,
	}
}

// Get returns the value associated with `key`. The second return value
// indicates if the value was found.
func (t *Trie[L, K, V]) Get(key L) (V, bool) {
	var zero V
	iter, err := t.tokenizer.Tokenize(key)
	if err != nil {
		return zero, false
	}

	t.mu.RLock()
	defer t.mu.RUnlock()
	var tokens []K
	for x := range iter {
		tokens = append(tokens, x)
	}
	return get(t.root, tokens)
}

func get[K cmp.Ordered, V any](root Node[K, V], tokens []K) (V, bool) {
	if len(tokens) > 0 {
		for child := range root.Children() {
			if child.Key() == tokens[0] {
				// found the current token in the children.
				if len(tokens) == 1 {
					// this is the node we're looking for
					return child.Value(), true
				}
				// we need to traverse down the trie
				return get[K, V](child, tokens[1:])
			}
		}
	}

	// if we got here, that means we couldn't find a common ancestor
	var zero V
	return zero, false
}

// Delete removes data associated with `key`. It returns true if the value
// was found and deleted, false otherwise
func (t *Trie[L, K, V]) Delete(key L) bool {
	iter, err := t.tokenizer.Tokenize(key)
	if err != nil {
		return false
	}
	var tokens []K
	for x := range iter {
		tokens = append(tokens, x)
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
	iter, err := t.tokenizer.Tokenize(key)
	if err != nil {
		return fmt.Errorf(`failed to tokenize key: %w`, err)
	}
	node := t.root

	var tokens []K
	for x := range iter {
		tokens = append(tokens, x)
	}

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
		for child := range root.Children() {
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
			cur.children = append(cur.children, newNode)
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
	children []*node[K, V]
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

func (n *node[K, V]) Children() iter.Seq[Node[K, V]] {
	n.mu.RLock()
	children := make([]*node[K, V], len(n.children))
	copy(children, n.children)
	n.mu.RUnlock()
	return func(yield func(Node[K, V]) bool) {
		for _, child := range children {
			if !yield(child) {
				break
			}
		}
	}
}

func (n *node[K, V]) AddChild(child Node[K, V]) {
	n.mu.Lock()
	// This is kind of gross, but we're only covering *node[T] with
	// Node[T] interface because we don't want the users to instantiate
	// their own nodes... so this type conversion is safe.
	//nolint:forcetypeassert
	n.children = append(n.children, child.(*node[K, V]))
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

func Walk[L any, K cmp.Ordered, V any](trie *Trie[L, K, V], v Visitor[K, V]) {
	var meta VisitMetadata
	meta.Depth = 1
	walk(trie.root, v, meta)
}

func walk[K cmp.Ordered, V any](node Node[K, V], v Visitor[K, V], meta VisitMetadata) {
	for child := range node.Children() {
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

	fmt.Fprintf(&sb, "%v: %v", n.Key(), n.Value())
	fmt.Println(sb.String())
	return true
}

func Dumper[K cmp.Ordered, V any]() Visitor[K, V] {
	return dumper[K, V]{}
}
