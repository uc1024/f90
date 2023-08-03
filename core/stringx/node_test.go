package stringx

import (
	"testing"
)

func TestNode(t *testing.T) {
	words := []string{
		"水绿天蓝蓝",
		"水绿火红红",
		"天蓝水蓝蓝",
	}
	nx := new(node)
	for _, word := range words {
		nx.add(word)
	}
	nx.build()

	t.Log(nx.find([]rune("天蓝水蓝蓝 天蓝蓝 ]水绿绿 .天蓝水绿 }天蓝水蓝蓝")))
}

func TestNode1(t *testing.T) {
	words := []string{
		"bc",
		"cd",
	}
	nx := new(node)
	for _, word := range words {
		nx.add(word)
	}
	nx.build()

	t.Log(nx.find([]rune("abcd")))
}
