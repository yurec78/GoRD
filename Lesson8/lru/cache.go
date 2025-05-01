package lru

import (
	"container/list"
)

type LRUCache interface {
	Put(key, value string)
	Get(key string) (string, bool)
}

type lruCacheImpl struct {
	capacity int
	items    map[string]*list.Element
	order    *list.List
}

func NewLruCache(capacity int) LRUCache {
	return &lruCacheImpl{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

func (c *lruCacheImpl) Put(key, value string) {
	if el, exists := c.items[key]; exists {
		el.Value.(*entry).value = value
		c.order.MoveToFront(el)
		return
	}
	if c.order.Len() == c.capacity {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.items, oldest.Value.(*entry).key)
		}
	}
	newEntry := &entry{key: key, value: value}
	element := c.order.PushFront(newEntry)
	c.items[key] = element
}

func (c *lruCacheImpl) Get(key string) (string, bool) {
	if el, exists := c.items[key]; exists {
		c.order.MoveToFront(el)
		return el.Value.(*entry).value, true
	}
	return "", false
}

type entry struct {
	key   string
	value string
}
