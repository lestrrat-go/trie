package trie_test

import (
	"context"
	"fmt"

	"github.com/lestrrat-go/trie"
)

func ExampleStringKey() {
	// An example where a string is used as the keys
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tree := trie.New()

	// Put values in the trie
	tree.Put(ctx, trie.StringKey("foo"), "one")
	tree.Put(ctx, trie.StringKey("bar"), 2)
	tree.Put(ctx, trie.StringKey("baz"), 3.0)
	tree.Put(ctx, trie.StringKey("日本語"), []byte{'f', 'o', 'u', 'r'})

	// Get a value from the trie
	v, ok := tree.Get(ctx, trie.StringKey("日本語"))
	if !ok {
		fmt.Printf("failed to find key '日本語'\n")
		return
	}
	_ = v

	// Delete a key from the trie
	if !tree.Delete(ctx, trie.StringKey("日本語")) {
		fmt.Printf("failed to delete key '日本語'\n")
		return
	}

	// This time Get() should fail
	v, ok = tree.Get(ctx, trie.StringKey("日本語"))
	if ok {
		fmt.Printf("key '日本語' should not exist\n")
		return
	}
	_ = v

	// Or, walk the entire trie
	for p := range tree.Walk(ctx) {
		// Do something with the values...
		_ = p
	}

	// OUTPUT:
}
