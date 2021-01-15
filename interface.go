package trie

import "context"

// Key is a sequence of Labels, which is associated with a value.
// The Key interface represents any type that allows the Trie to
// iterate over its elements, known as Labels.
//
// For example, a string can be thought of as a key consisting of
// runes as its labels.
//
// This allows the user maximum flexibility in terms of the input to
// use for our trie.
type Key interface {
	Iterate(context.Context) <-chan Label
}

// Label is a single entry in a Key. It can be anything, really
type Label interface{}
