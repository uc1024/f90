package search

import "testing"

type mock struct {
	route string
	value int
}

func TestSearch(t *testing.T) {
	tree := NewTree()
	_ = tree
	routes := []mock{
		{"/api/users", 1},
		{"/api/:layer", 19},
		{"/api/:layer1", 3},
		{"/:layer1/:layer2/::layer3", 4},
	}
	for i := range routes {
		item := routes[i]
		err := tree.Add(item.route, item.value)
		if err != nil {
			t.Log(err)
		}
	}
	query := "/api/:layer2"
	result, ok := tree.Search(query)
	if ok {
		t.Logf("%v", result.Item)
		t.Logf("%v", result.Params)
		for k, v := range result.Params {
			t.Log(k, v)
		}
	} else {
		t.Log(ok)
	}
	query = "/api/users"
	result, ok = tree.Search(query)
	if ok {
		t.Logf("%v", result.Item)
		t.Logf("%v", result.Params)
		for k, v := range result.Params {
			t.Log(k, v)
		}
	} else {
		t.Log(ok)
	}
	query = "/o/k/x"
	result, ok = tree.Search(query)
	if ok {
		t.Logf("%v", result.Item)
		t.Logf("%v", result.Params)
		for k, v := range result.Params {
			t.Log(k, v)
		}
	} else {
		t.Log(ok)
	}
}
