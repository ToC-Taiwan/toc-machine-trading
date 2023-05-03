// Package cache package cache
package cache

import (
	"fmt"
	"strings"
)

// Key Key
type Key struct {
	category string
	index    string
}

func NewKey(category string, index ...string) *Key {
	k := &Key{
		category: category,
	}

	switch len(index) {
	case 0:
		k.index = ""
	case 1:
		k.index = index[0]
	default:
		panic("index length must be 0 or 1")
	}

	return k
}

func (k *Key) String() string {
	return k.category + ":" + k.index
}

func (k *Key) Category() string {
	return k.category
}

func (k *Key) Index() string {
	return k.index
}

func (k *Key) ExtendIndex(opt ...string) *Key {
	if k.index == "" {
		panic("to extend index, original index must not be empty")
	}
	k.index = fmt.Sprintf("%s:%s", k.index, strings.Join(opt, ":"))
	return k
}
