# github.com/lestrrat-go/trie ![](https://github.com/lestrrat-go/trie/workflows/CI/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/lestrrat-go/trie.svg)](https://pkg.go.dev/github.com/lestrrat-go/trie)

This trie is implemented such that generic Key types can be used. 
Most other trie implementations are optimized for string based keys, but my use
case is to match certain numeric opcodes to arbitrary data.

<!-- INCLUDE(trie_example_test.go) -->
<!-- END INCLUDE -->