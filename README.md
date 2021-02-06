# trie - A Generic Trie Implementation

This trie is implemented such that generic Key types can be used. 
Most other trie implementations are optimized for string based keys, but my use
case is to match certain numeric opcodes to arbitrary data.

Within this library Keys are treated as sequence of Labels.
For example, a string can be thought of as Key that is comprised of a sequence
of rune Labels.

Each Key need to be able to break down to Labels via the `Iterate` method.
Each Label in turn becomes the local key in a trie node.
Each Label need to implement a `UniqueID` method to identify itself.

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

Originally based on https://github.com/koron/go-trie
Much code stolen from https://github.com/dghubble/trie
