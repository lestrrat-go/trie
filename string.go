package trie

type stringTokenizer struct{}

func (stringTokenizer) Tokenize(s string) ([]rune, error) {
	var list []rune
	for _, r := range s {
		list = append(list, r)
	}
	return list, nil
}

// String returns a Tokenizer that tokenizes a string into individual runes.
func String() Tokenizer[string, rune] {
	return stringTokenizer{}
}
