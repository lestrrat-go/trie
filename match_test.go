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
	tr.Put("ab", 2)
	tr.Put("bc", 4)
	tr.Put("bab", 6)
	tr.Put("d", 7)
	tr.Put("abcde", 10)
	mt := Compile(tr)

	// Check tree.
	f := func(s string, exp []int) {
		act := toInts(mt.MatchAll(s, nil))
		assert.Equal(t, act, exp, "not match for key=%q", s)
	}
	f("ab", []int{2})
	f("bc", []int{4})
	f("d", []int{7})
	f("abcde", []int{2, 4, 7, 10})
	f("babc", []int{6, 2, 4})
}
