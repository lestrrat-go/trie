// Package trie implements a trie that allows a generic key to point to
// stored data. Whereas many trie implementations are optimized for string
// based keys, this implementation allows keys to be a sequence of "labels".
package trie

import (
	"context"
	"sync"
)

// Trie is a simple trie that accepts arbitrary Key types as its input.
type Trie struct {
	children map[Label]*Trie
	hasValue bool
	mu       sync.RWMutex
	value    interface{}
}

// New creates a new Trie
func New() *Trie {
	return &Trie{
		children: make(map[Label]*Trie),
	}
}

// Get returns the value associated with `key`. The second return value
// is true if the value exists, false otherwise
func (t *Trie) Get(ctx context.Context, key Key) (interface{}, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	gctx, cancel := context.WithCancel(ctx)
	defer cancel()

	node := t
	for l := range key.Iterate(gctx) {
		node = node.children[l]
		if node == nil {
			return nil, false
		}
	}
	return node.value, true
}

// Put sets `key` to point to data `value`. The return value is true
// if the value was set anew. If this was an update operation, the return
// value would be false
func (t *Trie) Put(ctx context.Context, key Key, value interface{}) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	pctx, cancel := context.WithCancel(ctx)
	defer cancel()

	node := t
	for l := range key.Iterate(pctx) {
		child := node.children[l]
		if child == nil {
			child = New()
			node.children[l] = child
		}
		node = child
	}

	isNewVal := node.hasValue
	node.hasValue = true
	node.value = value
	return isNewVal
}

func (t *Trie) isLeaf() bool {
	return len(t.children) == 0
}

type ancestor struct {
	Label Label
	Node  *Trie
}

// Delete removes data associated with `key`. It returns true if the value
// was found and deleted, false otherwise
func (t *Trie) Delete(ctx context.Context, key Key) bool {
	var ancestors []ancestor
	node := t
	for l := range key.Iterate(ctx) {
		ancestors = append(ancestors, ancestor{
			Label: l,
			Node: node,
		})
		node = node.children[l]
		if node == nil {
			// node does not exist
			return false
		}
	}

	// delete the node value
	node.value = nil

	// if leaf, remove it from its parent's children map. Repeat for ancestors.
	if !node.isLeaf() {
		return true
	}
	// iterate backwards over the ancestors
	for i := len(ancestors) - 1; i >= 0; i-- {
		ancestor := ancestors[i]
		parent := ancestor.Node
		delete(parent.children, ancestor.Label)

		if !parent.isLeaf() {
			// parent has other children, stop
			break
		}
		parent.children = nil
		if parent.hasValue {
			// parent has a value, stop
			break
		}
	}
	return true
}

// WalkPair is what you get when you call `Walk()` on a trie.
type WalkPair struct {
	// Because we have a generic "Label" type, we unfortunately cannot
	// provide a re-constructed Key object for the user to handle.
	// Instead we provide this value as a slice of Labels
	Labels []Label

	// Value is the value associated with the Labels
	Value interface{}
}

// Walk returns a channel that you can read from to access all data
// that is stored within this trie.
func (t *Trie) Walk(ctx context.Context) <-chan WalkPair {
	ch := make(chan WalkPair)
	go t.walk(ctx, ch, nil)
	return ch
}

func (t *Trie) walk(ctx context.Context, dst chan WalkPair, labels []Label) {
	if labels == nil {
		t.mu.RLock()
		defer t.mu.RUnlock()
		defer close(dst)
	}

	if t.hasValue {
		p := WalkPair{
			Labels: labels,
			Value:  t.value,
		}
		select {
		case <-ctx.Done():
			return
		case dst <- p:
		}
	}

	for l, child := range t.children {
		child.walk(ctx, dst, append(labels, l))
	}
}
