package trie_test

import (
	"testing"

	"github.com/lestrrat-go/trie/v2"
	"github.com/stretchr/testify/require"
)

func TestTrie(t *testing.T) {
	t.Parallel()

	tree := trie.New[string, rune, int](trie.String())

	testcases := []struct {
		Key   string
		Value int
	}{
		{"foo", 1},
		{"far", 2},
		{"for", 3},
		{"bar", 4},
		{"baz", 5},
	}

	for _, tc := range testcases {
		tree.Put(tc.Key, tc.Value)
	}

	for _, tc := range testcases {
		t.Run(tc.Key, func(t *testing.T) {
			v, ok := tree.Get(tc.Key)
			require.True(t, ok, `tree.Get should return true`)
			require.Equal(t, tc.Value, v, `tree.Get should return expected value`)
		})
	}

	require.True(t, tree.Delete("foo"), `tree.Delete should return true`)
	_, ok := tree.Get("foo")
	require.False(t, ok, `tree.Get should return false`)
}
