package m

import (
	"github.com/PandaTtttt/go-assembly/util"
	"sort"
)

// M is a shortcut for map[string]interface{}
type M map[string]interface{}

// L is a shortcut for []interface{}
type L []interface{}

func (m M) SortedKeys() []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (m M) Contains(keys []string) bool {
	var existKeys []string
	for k := range m {
		existKeys = append(existKeys, k)
	}
	for _, k := range keys {
		if !util.InArrayString(k, existKeys) {
			return false
		}
	}
	return true
}

func (m M) ExactlyContains(keys []string) bool {
	var existKeys []string
	for k := range m {
		if !util.InArrayString(k, keys) {
			return false
		}
		existKeys = append(existKeys, k)
	}
	for _, k := range keys {
		if !util.InArrayString(k, existKeys) {
			return false
		}
	}
	return true
}
