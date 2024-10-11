package trie

import (
	"iter"
)

// String returns a Tokenizer that tokenizes a string into individual runes.
func String() Tokenizer[string, rune] {
	return TokenizeFunc[string, rune](func(s string) (iter.Seq[rune], error) {
		return func(yield func(rune) bool) {
			for _, r := range s {
				if !yield(r) {
					break
				}
			}
		}, nil
	})
}
