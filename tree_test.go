package trie

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	f := func(t *testing.T, tr *Tree, key Key, value interface{}) {
		t.Helper()

		n := tr.Get(key)
		if value == nil {
			assert.Equal(t, n, (*Node)(nil), "no nodes for %q", key)
			return
		}
		assert.Equal(t, n.Value, value, "value for %q", key)
	}

	testcases := []struct {
		Key   Key
		Value interface{}
	}{
		{Key: StringKey("foo"), Value: "123"},
		{Key: StringKey("bar"), Value: "999"},
		{Key: StringKey("日本語"), Value: "こんにちは"},
		{Key: StringKey("baz")},
		{Key: StringKey("English")},
	}

	tr := New()
	tr.Put(StringKey("foo"), "123")
	tr.Put(StringKey("bar"), "999")
	tr.Put(StringKey("日本語"), "こんにちは")

	for _, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("%s", tc.Key), func(t *testing.T) {
			f(t, tr, tc.Key, tc.Value)
		})
	}
}

func TestTree_nc(t *testing.T) {
	tr := New()
	tr.Put(StringKey("foo"), "123")
	tr.Put(StringKey("bar"), "999")
	tr.Put(StringKey("日本語"), "こんにちは")
	if tr.nc != 9 {
		t.Errorf("nc mismatch: %d", tr.nc)
	}
}

func TestNode_cc(t *testing.T) {
	f := func(key Key, cc int) {
		n := new(Node)
		for l := range key.Iterate() {
			n.Dig(l)
		}
		if !assert.Equal(t, n.cc, cc, "runes: %q", key) {
			return
		}
	}
	f(StringKey(""), 0)
	f(StringKey("a"), 1)
	f(StringKey("bac"), 3)
	f(StringKey("aaa"), 1)
	f(StringKey("bbbaaaccc"), 3)
	f(StringKey("bacbacbac"), 3)
	f(StringKey("日本語こんにちは"), 8)
	f(StringKey("あめんぼあかいなあいうえお"), 10)
}

// collectRunes1 coolects label runes from sibling nodes.
func collectRunes1(n *Node, max int) []rune {
	runes := make([]rune, 0, max)
	q := make([]*Node, 0, max)
	q = append(q, n)
	for len(q) > 0 {
		m := q[0]
		runes = append(runes, m.label.(RuneLabel).Rune())
		if len(runes) > max {
			return []rune("nodes may have infinite loop")
		}
		if m.Low != nil {
			q = append(q, m.Low)
		}
		if m.High != nil {
			q = append(q, m.High)
		}
		q = q[1:]
	}
	return runes
}

// collectRunes1 coolects label runes from sibling nodes in reverse order.
func collectRunes2(n *Node, max int) []rune {
	runes := make([]rune, 0, max)
	q := make([]*Node, 0, max)
	q = append(q, n)
	for len(q) > 0 {
		m := q[0]
		runes = append(runes, m.label.(RuneLabel).Rune())
		if len(runes) > max {
			return []rune("nodes may have infinite loop")
		}
		if m.High != nil {
			q = append(q, m.High)
		}
		if m.Low != nil {
			q = append(q, m.Low)
		}
		q = q[1:]
	}
	return runes
}

func TestNode_Balance(t *testing.T) {
	n := new(Node)
	for l := range StringKey("123456789ABCDEF").Iterate() {
		n.Dig(l)
	}
	n.Balance()
	if n.Child == nil {
		t.Fatal("Child shoud not be nil after balancing")
	}
	r1 := collectRunes1(n.Child, n.cc)
	if !assert.Equal(t, string(r1), "84C26AE13579BDF", "should be balanced") {
		return
	}
	r2 := collectRunes2(n.Child, n.cc)
	if !assert.Equal(t, string(r2), "8C4EA62FDB97531", "should be balanced") {
		return
	}
}
