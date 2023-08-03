package stringx

import "golang.org/x/exp/slices"

type (
	node struct {
		children map[rune]*node
		/*
			如果当前节点没有与字符匹配的子节点，就通过 fail 指针跳转到一个与当前字符匹配的节点，然后继续匹配 (匹配错误时快速回撤)
		*/
		fail  *node
		depth int
		end   bool
	}
)

/*
	add node
*/
func (n *node) add(word string) {

	// * []rune 类型表示 Unicode 字符串
	chars := []rune(word)

	if len(chars) == 0 {
		return
	}

	// * 开始节点
	nd := n
	var depth = int(0)

	// * 循环 unicode数组
	for i, char := range chars {

		if nd.children == nil {
			// * 初始化节点
			child := new(node)
			child.depth = i + 1
			nd.children = map[rune]*node{char: child}
			nd = child
		} else if child, ok := nd.children[char]; ok {
			// * 存在unicode 深度++
			nd = child
			depth = depth + 1
		} else {
			// * 不存在的unicode 记录
			child := new(node)
			child.depth = i + 1
			nd.children[char] = child
			nd = child
		}

	}

	// * 表示词语的结尾
	nd.end = true
}

/*
在这段代码中，节点的结构体类型为 node，其中包含一个 children 字典类型的字段，用于存储子节点。build 方法是用于构建 Trie 树的方法，它首先遍历当前节点的所有子节点，为每个子节点设置 fail 指针，然后将所有子节点加入到一个切片中。

接着，代码进入一个循环，每次取出切片中的第一个节点，并遍历它的所有子节点。对于每个子节点，它首先将其加入到切片中，然后从当前节点开始向上遍历 fail 指针，查找与当前子节点匹配的节点。如果找到了匹配的节点，就将子节点的 fail 指针指向该节点，退出循环。如果遍历到了根节点还没有找到匹配的节点，则将子节点的 fail 指针指向根节点。

通过这个构建过程，可以构建出一个 Trie 树，并为每个节点设置 fail 指针，用于在匹配字符串时快速跳转到下一个节点。具体地，当匹配到一个节点时，如果当前节点没有与字符匹配的子节点，就通过 fail 指针跳转到一个与当前字符匹配的节点，然后继续匹配。这样可以避免重复匹配和回溯，提高字符串匹配的效率。
*/
func (n *node) build() {
	nodes := []*node{}

	for _, child := range n.children {
		child.fail = n
		nodes = append(nodes, child)
	}

	// * 循环直到nodes数组为0
	for len(nodes) > 0 {
		nd := nodes[0]    // * 每次取出第一个几点
		nodes = nodes[1:] // * 剩余节点

		for key, child := range nd.children {
			nodes = append(nodes, child)
			cur := nd
			for cur != nil {
				if cur.fail == nil {
					child.fail = n
					break
				}
				if fail, ok := cur.fail.children[key]; ok {
					child.fail = fail
					break
				}
				cur = cur.fail
			}
		}
	}

}

/*
用于在给定的字符串中查找多个模式串的位置。该方法的参数 chars 是一个字符数组，表示需要进行匹配的字符串。返回值是一个 scope 类型的切片，其中每个元素表示一个匹配到的模式串在字符串中的起始位置和结束位置。

在这个方法中，代码首先初始化了一个空的 scopes 切片用于存储匹配到的模式串的位置。然后，从 Trie 树的根节点开始，按照字符串的顺序依次遍历字符数组 chars。对于每个字符，如果该字符在当前节点的子节点中，则将当前节点指向该子节点。如果该字符在当前节点的子节点中不存在，则需要通过 fail 指针进行跳转，直到找到一个存在该字符的子节点或者跳到了根节点。

如果找到了一个子节点，那么代码会从该子节点开始遍历所有的 fail 指针，直到遍历到了根节点或者某个节点的 end 标志为 true。对于每个遍历到的节点，如果其 end 标志为 true，则将其对应的模式串的位置信息添加到 scopes 切片中。最后返回 scopes 切片即可。

通过这个匹配过程，可以在一个字符串中快速查找多个模式串的位置，并返回位置信息。该算法的时间复杂度为 $O(n)$，其中 $n$ 表示需要匹配的字符串的长度。该算法可以在很短的时间内处理大量的字符串匹配请求，是一种非常实用的字符串匹配算法。
*/
func (n *node) find(chars []rune, skip ...rune) []scope {
	var scopes []scope
	size := len(chars)
	cur := n
	offset := 0

	for i := 0; i < size; i++ {
		// * 过滤
		if slices.Contains[rune](skip, chars[i]) {
			offset++
			continue
		}

		child, ok := cur.children[chars[i]]
		if ok {
			cur = child
		} else {
			for cur != n {
				cur = cur.fail
				if child, ok = cur.children[chars[i]]; ok {
					cur = child
					break
				}
			}

			if child == nil {
				continue
			}
		}

		for child != n {
			if child.end {
				scopes = append(scopes, scope{
					start: (i + 1 - child.depth) - offset,
					stop:  i + 1,
				})
				offset = 0
			}
			child = child.fail
		}
	}

	return scopes
}
