package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func toInts(matches []Match) []int {
	if len(matches) == 0 {
		return nil
	}
	r := make([]int, len(matches))
	for i, m := range matches {
		r[i] = m.Value.(int)
	}
	return r
}

func TestMatch(t *testing.T) {
	// Build tree.
	tr := New()
	tr.Put(StringKey("ab"), 2)
	tr.Put(StringKey("bc"), 4)
	tr.Put(StringKey("bab"), 6)
	tr.Put(StringKey("d"), 7)
	tr.Put(StringKey("abcde"), 10)
	mt := Compile(tr)

	// Check tree.
	f := func(key Key, exp []int) {
		act := toInts(mt.MatchAll(key, nil))
		assert.Equal(t, act, exp, "not match for key=%q", key)
	}
	f(StringKey("ab"), []int{2})
	f(StringKey("bc"), []int{4})
	f(StringKey("d"), []int{7})
	f(StringKey("abcde"), []int{2, 4, 7, 10})
	f(StringKey("babc"), []int{6, 2, 4})
}
