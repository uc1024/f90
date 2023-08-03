package stringx

import (
	"strings"
)

const defaultMask = '*'

type (
	// TrieOption defines the method to customize a Trie.
	TrieOption func(trie *trieNode)

	// A Trie is a tree implementation that used to find elements rapidly.
	Trie interface {
		Filter(text string) (string, []string, bool)
		FindKeywords(text string) []string
	}

	trieNode struct {
		node
		mask rune
		skip []rune
	}

	scope struct {
		start int
		stop  int
	}
)

// NewTrie returns a Trie.
func NewTrie(words []string, opts ...TrieOption) Trie {
	n := new(trieNode)

	for _, opt := range opts {
		opt(n)
	}
	if n.mask == 0 {
		n.mask = defaultMask
	}
	for _, word := range words {
		n.add(word)
	}

	n.build()

	return n
}

func (n *trieNode) Filter(text string) (sentence string, keywords []string, found bool) {
	chars := []rune(text)
	if len(chars) == 0 {
		return text, nil, false
	}

	scopes := n.find(chars, n.skip...)
	keywords = n.collectKeywords(chars, scopes)

	for _, match := range scopes {
		// we don't care about overlaps, not bringing a performance improvement
		n.replaceWithAsterisk(chars, match.start, match.stop)
	}

	return string(chars), keywords, len(keywords) > 0
}

func (n *trieNode) FindKeywords(text string) []string {
	chars := []rune(text)
	if len(chars) == 0 {
		return nil
	}

	scopes := n.find(chars,n.skip...)
	return n.collectKeywords(chars, scopes)
}

func (n *trieNode) collectKeywords(chars []rune, scopes []scope) []string {
	set := make(map[string]struct{})

	for _, v := range scopes {
		// fmt.Println(string(chars[v.start:v.stop]))
		set[strings.ReplaceAll(string(chars[v.start:v.stop]), " ", "")] = struct{}{}
	}

	var i int
	keywords := make([]string, len(set))
	for k := range set {
		keywords[i] = k
		i++
	}

	return keywords
}

func (n *trieNode) replaceWithAsterisk(chars []rune, start, stop int) {
	for i := start; i < stop; i++ {
		chars[i] = n.mask
	}
}

// WithMask customizes a Trie with keywords masked as given mask char.
func WithMask(mask rune) TrieOption {
	return func(n *trieNode) {
		n.mask = mask
	}
}

func WithSkip(skip []rune) TrieOption {
	return func(n *trieNode) {
		n.skip = skip
	}
}
