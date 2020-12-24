package trie

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEach(t *testing.T) {
	tr := New()
	tr.Put(StringKey("foo"), "123")
	tr.Put(StringKey("bar"), "999")
	tr.Put(StringKey("日本語"), "こんにちは")

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
		t.Logf("%c", n.label)

		var r rune
		if l := n.label; l != nil {
			r = l.(RuneLabel).Rune()
		}
		if !assert.Equal(t, expected[i], r, `labels should match for input %d`, i) {
			return false
		}
		i++
		return true
	}))
}

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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var runes []rune
	for q := range n.Iterate(WithBFS(ctx)) {
		runes = append(runes, q.label.(RuneLabel).Rune())
	}
	return runes
}

// collectRunes2 coolects label runes from sibling nodes in reverse order.
func collectRunes2(n *Node, max int) []rune {
	  ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
  defer cancel()

  var runes []rune
  for q := range n.Iterate(WithBFSReverse(ctx)) {
    runes = append(runes, q.label.(RuneLabel).Rune())
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
	if !assert.Equal(t, "84C26AE13579BDF", string(r1), "should be balanced") {
		return
	}
	r2 := collectRunes2(n.Child, n.cc)
	if !assert.Equal(t, "8C4EA62FDB97531", string(r2), "should be balanced") {
		return
	}
}
