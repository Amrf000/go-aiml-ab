package ab

import "fmt"

type NodemapperOperator struct{}

func Size(node *Nodemapper) int {
	set := make(map[string]bool)
	if node.ShortCut {
		set["<THAT>"] = true
	}
	if node.Key != "" {
		set[node.Key] = true
	}
	if node.Map != nil {
		for k := range node.Map {
			set[k] = true
		}
	}
	return len(set)
}

func Put(node *Nodemapper, key string, value *Nodemapper) {
	if node.Map != nil {
		node.Map[key] = value
	} else {
		node.Key = key
		node.Value = value
	}
}

func Get(node *Nodemapper, key string) *Nodemapper {
	if node.Map != nil {
		return node.Map[key]
	}
	if key == node.Key {
		return node.Value
	}
	return nil
}

func ContainsKey(node *Nodemapper, key string) bool {
	if node.Map != nil {
		_, ok := node.Map[key]
		return ok
	}
	return key == node.Key
}

func PrintKeys(node *Nodemapper) {
	keySet := KeySet(node)
	for _, k := range keySet {
		fmt.Println(k)
	}
}

func KeySet(node *Nodemapper) []string {
	setObj := make(map[string]bool)
	if node.Map != nil {
		for k := range node.Map {
			setObj[k] = true
		}
	} else {
		if node.Key != "" {
			setObj[node.Key] = true
		}
	}
	ret := []string{}
	for k := range setObj {
		ret = append(ret, k)
	}
	return ret
}

func IsLeaf(node *Nodemapper) bool {
	return node.Category != nil
}

func Upgrade(node *Nodemapper) {
	node.Map = map[string]*Nodemapper{}
	node.Map[node.Key] = node.Value
	node.Key = ""
	node.Value = nil
}
