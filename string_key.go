package trie

type RuneLabel rune

func (r RuneLabel) UniqueID() interface{} {
	return rune(r)
}

type RuneLabelIterator struct {
	list []rune
	cur  int
}

func (iter *RuneLabelIterator) Next() bool {
	return iter.cur < len(iter.list)
}

func (iter *RuneLabelIterator) Label() Label {
	r := RuneLabel(iter.list[iter.cur])
	iter.cur++
	return r
}

type StringKey string

func (sk StringKey) Labels() LabelIterator {
	var list []rune
	for _, r := range string(sk) {
		list = append(list, r)
	}
	return &RuneLabelIterator{
		list: list,
	}
}
