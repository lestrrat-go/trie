package trie

import "testing"

func TestPut(t *testing.T) {
	f := func(tr *Tree, key string, value interface{}) {
		n := tr.Get(key)
		if value == nil {
			assertEquals(t, n, (*Node)(nil), "no nodes for %q", key)
			return
		}
		assertEquals(t, n.Value, value, "value for %q", key)
	}
	tr := New()
	tr.Put("foo", "123")
	tr.Put("bar", "999")
	tr.Put("日本語", "こんにちは")
	f(tr, "foo", "123")
	f(tr, "bar", "999")
	f(tr, "日本語", "こんにちは")
	f(tr, "baz", nil)
	f(tr, "English", nil)
}

func TestNode_cc(t *testing.T) {
	f := func(runes string, cc int) {
		n := new(Node)
		for _, r := range runes {
			n.Dig(r)
		}
		assertEquals(t, n.cc, cc, "runes: %q", runes)
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
	assertEquals(t, string(r1), "84C26AE13579BDF", "should be balanced")
	r2 := collectRunes2(n.Child, n.cc)
	assertEquals(t, string(r2), "8C4EA62FDB97531", "should be balanced")
}
