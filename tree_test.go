package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEach(t *testing.T) {
	tr := New()
	tr.Put("foo", "123")
	tr.Put("bar", "999")
	tr.Put("日本語", "こんにちは")

	expected := []rune{
		0,
		'b',
		'f',
		'日',
		'a',
		'o',
		'本',
		'r',
		'o',
		'語',
	}

	i := 0
	tr.Each(NodeProc(func(n *Node) bool {
		if !assert.Equal(t, expected[i], n.Label, `labels should match for input %d`, i) {
			return false
		}
		i++
		return true
	}))
}

func TestPut(t *testing.T) {
	f := func(t *testing.T, tr *Tree, key string, value interface{}) {
		t.Helper()

		n := tr.Get(key)
		if value == nil {
			assert.Equal(t, n, (*Node)(nil), "no nodes for %q", key)
			return
		}
		assert.Equal(t, n.Value, value, "value for %q", key)
	}

	testcases := []struct {
		Key   string
		Value interface{}
	}{
		{Key: "foo", Value: "123"},
		{Key: "bar", Value: "999"},
		{Key: "日本語", Value: "こんにちは"},
		{Key: "baz"},
		{Key: "English"},
	}

	tr := New()
	tr.Put("foo", "123")
	tr.Put("bar", "999")
	tr.Put("日本語", "こんにちは")

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Key, func(t *testing.T) {
			f(t, tr, tc.Key, tc.Value)
		})
	}
}

func TestTree_nc(t *testing.T) {
	tr := New()
	tr.Put("foo", "123")
	tr.Put("bar", "999")
	tr.Put("日本語", "こんにちは")
	if tr.nc != 9 {
		t.Errorf("nc mismatch: %d", tr.nc)
	}
}

func TestNode_cc(t *testing.T) {
	f := func(runes string, cc int) {
		n := new(Node)
		for _, r := range runes {
			n.Dig(r)
		}
		if !assert.Equal(t, n.cc, cc, "runes: %q", runes) {
			return
		}
	}
	f("", 0)
	f("a", 1)
	f("bac", 3)
	f("aaa", 1)
	f("bbbaaaccc", 3)
	f("bacbacbac", 3)
	f("日本語こんにちは", 8)
	f("あめんぼあかいなあいうえお", 10)
}

// collectRunes1 coolects label runes from sibling nodes.
func collectRunes1(n *Node, max int) []rune {
	runes := make([]rune, 0, max)
	q := make([]*Node, 0, max)
	q = append(q, n)
	for len(q) > 0 {
		m := q[0]
		runes = append(runes, m.Label)
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
		runes = append(runes, m.Label)
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
	for _, r := range "123456789ABCDEF" {
		n.Dig(r)
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
