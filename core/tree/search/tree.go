package search

import (
	"errors"
	"fmt"
)

const (
	colon = ':'
	slash = '/'
)

var (
	// errDupItem means adding duplicated item.
	errDupItem = errors.New("duplicated item")
	// errDupSlash means item is started with more than one slash.
	errDupSlash = errors.New("duplicated slash")
	// errEmptyItem means adding empty item.
	errEmptyItem = errors.New("empty item")
	// errInvalidState means search tree is in an invalid state.
	errInvalidState = errors.New("search tree is in an invalid state")
	// errNotFromRoot means path is not starting with slash.
	errNotFromRoot = errors.New("path should start with /")

	// NotFound is used to hold the not found result.
	NotFound Result
)

type (
	innerResult struct {
		key   string
		value string
		named bool
		found bool
	}

	node struct {
		item     interface{}
		children [2]map[string]*node
	}

	Tree struct {
		root *node
	}

	Result struct {
		Item   interface{}
		Params map[string]string
	}
)

func NewTree() *Tree {
	return &Tree{
		root: newNode(nil),
	}
}

func (t *Tree) Add(route string, item interface{}) (err error) {
	// * route不是空字符串 && 首字符必须是 /
	if len(route) == 0 || route[0] != slash {
		return errNotFromRoot
	}

	if item == nil {
		return errEmptyItem
	}

	err = add(t.root, route[1:], item)

	switch err {
	case errDupItem:
		return duplicatedItem(route)
	case errDupSlash:
		return duplicatedSlash(route)
	default:
		return
	}
}

// * 搜索路径
func (t *Tree) Search(route string) (result Result, is bool) {
	if len(route) == 0 || route[0] != slash {
		return NotFound, false
	}
	ok := t.next(t.root, route[1:], &result)
	return result, ok
}

func (t *Tree) next(n *node, route string, result *Result) bool {
	if len(route) == 0 && n.item != nil {
		result.Item = n.item
		return true
	}

	for i := range route {
		if route[i] != slash {
			continue
		}

		token := route[:i]
		return n.forEach(func(k string, v *node) bool {
			r := match(k, token)
			// * 继续匹配
			if !r.found || !t.next(v, route[i+1:], result) {
				return false
			}
			// * 说明找到了
			if r.named {
				addParam(result, r.key, r.value)
			}

			return true
		})
	}

	return n.forEach(func(k string, v *node) bool {
		if r := match(k, route); r.found && v.item != nil {
			result.Item = v.item
			if r.named {
				addParam(result, r.key, r.value)
			}

			return true
		}

		return false
	})
}

func duplicatedItem(item string) error {
	return fmt.Errorf("duplicated item for %s", item)
}

func duplicatedSlash(item string) error {
	return fmt.Errorf("duplicated slash for %s", item)
}

func add(prev *node, route string, item interface{}) (err error) {

	// * 如果开头是 /
	if route[0] == slash {
		return errDupSlash
	}

	// * route 为空时说明已经到底了,item放入prev
	if len(route) == 0 {
		if prev.item != nil {
			return errDupItem
		}
		prev.item = item
	}

	// * 处理后续节点
	for i := range route {
		if route[i] != slash {
			continue
		}
		// * 找到 slash 结尾
		token := route[:i]
		// * 上一个节点是否有这个token路径.直接继续递归
		children := prev.getChildren(token)
		if child, ok := children[token]; ok {
			if child != nil {
				return add(child, route[i+1:], item)
			}
			return errInvalidState
		}

		// * 创建新的阶段路径并递归
		child := newNode(nil)
		children[token] = child
		return add(child, route[i+1:], item)
	}

	// * 匹配结束都没有遇见 slash
	children := prev.getChildren(route)
	// * 如果上一个节点存在这个子节点,插入数值(如果节点已经有数值了则报错)
	if child, ok := children[route]; ok {
		if child.item != nil {
			return errDupItem
		}

		child.item = item
	} else {
		// * 直接新增子节点
		children[route] = newNode(item)
	}

	return
}

func (nd *node) getChildren(route string) map[string]*node {
	if len(route) > 0 && route[0] == colon {
		return nd.children[1]
	}

	return nd.children[0]
}

func (nd *node) forEach(fn func(string, *node) bool) bool {
	for _, children := range nd.children {
		for k, v := range children {
			if fn(k, v) {
				return true
			}
		}
	}

	return false
}

// * pat 节点key 如果是:开头数值匹配
func match(pat, token string) innerResult {
	if pat[0] == colon {
		return innerResult{
			key:   pat[1:],
			value: token,
			named: true,
			found: true,
		}
	}

	return innerResult{
		found: pat == token,
	}
}

func addParam(result *Result, k, v string) {
	if result.Params == nil {
		result.Params = make(map[string]string)
	}

	result.Params[k] = v
}

// * 创建一个节点
func newNode(item interface{}) *node {
	return &node{
		item: item,
		children: [2]map[string]*node{
			make(map[string]*node),
			make(map[string]*node),
		},
	}
}
