package imdb

import (
	"container/list"
	"time"
)

type LRUCache[K comparable, V any] struct {
	capacity int
	cache    map[K]*list.Element
	list     *list.List
	keys     map[K]struct{}
}

type entry[K comparable, V any] struct {
	key   K
	value V
	ttl   time.Time
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[K]*list.Element),
		list:     list.New(),
		keys:     make(map[K]struct{}),
	}
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	if elem, found := c.cache[key]; found {
		c.list.MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, true
	}
	var zero V
	return zero, false
}

func (c *LRUCache[K, V]) Put(key K, value V) {
	if elem, found := c.cache[key]; found && elem.Value.(*entry[K, V]).ttl.Before(time.Now()) {
		c.list.MoveToFront(elem)
		elem.Value.(*entry[K, V]).value = value
		return
	}

	if c.list.Len() >= c.capacity {
		back := c.list.Back()
		if back != nil {
			c.list.Remove(back)
			delete(c.cache, back.Value.(*entry[K, V]).key)
		}
	}

	elem := c.list.PushFront(&entry[K, V]{key, value, time.Now().Add(5 * time.Minute)})
	c.cache[key] = elem
	c.keys[key] = struct{}{}
}

func (c *LRUCache[K, V]) Delete(key K) {
	if elem, found := c.cache[key]; found {
		c.list.Remove(elem)
		delete(c.cache, key)
		delete(c.keys, key)
	}
}

func (c *LRUCache[K, V]) Len() int {
	return c.list.Len()
}
func (c *LRUCache[K, V]) GetKeys() []K {
	keys := make([]K, 0, len(c.keys))
	for k := range c.keys {
		keys = append(keys, k)
	}
	return keys
}
func (c *LRUCache[K, V]) DeleteAll() {
	c.cache = make(map[K]*list.Element)
	c.keys = make(map[K]struct{})
	c.list = list.New()

}
