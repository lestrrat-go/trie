package trie

// Match is matched data.
type Match struct {
	Index   int
	Pattern string
	Value   interface{}
}

// Matcher compares a string with multiple strings using Aho-Corasick
// algorithm.
type Matcher struct {
	// TODO:
}

// Compile compiles a Matcher from a Tree.
func Compile(tr *Tree) *Matcher {
	// TODO:
	return nil
}

func (m *Matcher) Match(text string) []Match {
	// TODO:
	return nil
}
