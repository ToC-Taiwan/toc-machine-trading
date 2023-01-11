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

func NewKey(category, index string) *Key {
	return &Key{
		category: category,
		index:    index,
	}
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
	k.index = fmt.Sprintf("%s:%s", k.index, strings.Join(opt, ":"))
	return k
}
