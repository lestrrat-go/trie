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
