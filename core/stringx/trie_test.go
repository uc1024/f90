package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 9

func TestRune(t *testing.T) {
	re := []rune("123456789")
	t.Log(re)
}

func TestTrieSimple(t *testing.T) {
	trie := NewTrie([]string{
		"bc",
		"cd",
	})
	output, keywords, found := trie.Filter("abcd")
	assert.True(t, found)
	assert.Equal(t, "a***", output)
	assert.ElementsMatch(t, []string{"bc", "cd"}, keywords)
}

func TestTrieSimple1(t *testing.T) {
	trie := NewTrie([]string{
		"bc",
		"cd",
		"ab",
		"ac",
	}, WithSkip([]rune(" ")))
	output, keywords, found := trie.Filter("a b c d")
	t.Log(output)
	t.Log(keywords)
	t.Log(found)
}
