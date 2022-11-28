package trie_test

import (
	"context"
	"fmt"

	"github.com/lestrrat-go/trie"
)

func ExampleStringKey() {
	tree := trie.New()

	// Put values in the trie
	tree.Put(trie.StringKey("foo"), "one")
	tree.Put(trie.StringKey("bar"), 2)
	tree.Put(trie.StringKey("baz"), 3.0)
	tree.Put(trie.StringKey("日本語"), []byte{'f', 'o', 'u', 'r'})

	// Get a value from the trie
	v, ok := tree.Get(trie.StringKey("日本語"))
	if !ok {
		fmt.Printf("failed to find key '日本語'\n")
		return
	}
	_ = v

	// Delete a key from the trie
	if !tree.Delete(trie.StringKey("日本語")) {
		fmt.Printf("failed to delete key '日本語'\n")
		return
	}

	// This time Get() should fail
	v, ok = tree.Get(trie.StringKey("日本語"))
	if ok {
		fmt.Printf("key '日本語' should not exist\n")
		return
	}
	_ = v

	ctx := context.Background()

	// Or, walk the entire trie
	for p := range tree.Walk(ctx) {
		// Do something with the values...
		_ = p
	}

	// OUTPUT:
}
