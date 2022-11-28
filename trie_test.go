package trie_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/lestrrat-go/trie"
	"github.com/stretchr/testify/assert"
)

func TestTrie(t *testing.T) {
	t.Parallel()

	tree := trie.New()
	tree.Put(trie.StringKey("foo"), 1)
	tree.Put(trie.StringKey("bar"), 2)
	tree.Put(trie.StringKey("baz"), 3)
	tree.Put(trie.StringKey("日本語"), 4)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for p := range tree.Walk(ctx) {
		t.Logf("%#v", p)
	}

	testcases := []struct {
		Key      trie.Key
		Expected interface{}
		Missing  bool
	}{
		{
			Key:      trie.StringKey("foo"),
			Expected: 1,
		},
		{
			Key:      trie.StringKey("日本語"),
			Expected: 4,
		},
		{
			Key:     trie.StringKey("hoge"),
			Missing: true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("%s", tc.Key), func(t *testing.T) {
			t.Parallel()
			v, ok := tree.Get(tc.Key)
			if tc.Missing {
				if !assert.False(t, ok, `tree.Get should return false`) {
					return
				}
			} else {
				if !assert.True(t, ok, `tree.Get should return true`) {
					return
				}

				if !assert.Equal(t, tc.Expected, v, `tree.Get should return expected value`) {
					return
				}
			}
		})
	}
}
