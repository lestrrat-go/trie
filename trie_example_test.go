package trie_test

import (
	"fmt"

	"github.com/lestrrat-go/trie/v2"
)

func Example() {
	tree := trie.New[string, rune, any](trie.String())

	// Put values in the trie
	_ = tree.Put("foo", "one")
	_ = tree.Put("bar", 2)
	_ = tree.Put("baz", 3.0)
	_ = tree.Put("日本語", []byte{'f', 'o', 'u', 'r'})

	// Get a value from the trie
	v, ok := tree.Get("日本語")
	if !ok {
		fmt.Printf("failed to find key '日本語'\n")
		return
	}
	_ = v

	// Delete a key from the trie
	if !tree.Delete("日本語") {
		fmt.Printf("failed to delete key '日本語'\n")
		return
	}

	// This time Get() should fail
	v, ok = tree.Get("日本語")
	if ok {
		fmt.Printf("key '日本語' should not exist\n")
		return
	}
	_ = v

	// OUTPUT:
}
