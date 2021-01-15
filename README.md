# trie - A Generic Trie Implementation

This trie is implemented such that generic Key types can be used. 
Most other trie implementations are optimized for string based keys, but my use
case is to match certain numeric opcodes to arbitrary.

Within this library Keys are treated as sequence of Labels.
For example, a string can be thought of as Key that is comprised of a sequence
of rune Labels.

# SYNOPSIS

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

t := trie.New()

t.Put(ctx, trie.StringKey("foo"), 1)
v, ok := t.Get(ctx, trie.StringKey("foo"))
ok := t.Delete(ctx, trie.StringKey("foo"))
for p := range t.Walk(ctx) {
	// p.Labels
	// p.Value
}
```

# REFERENCES

Originally based on https://github.com/koron/trie
Much code stolen from https://github.com/dghubble/trie
