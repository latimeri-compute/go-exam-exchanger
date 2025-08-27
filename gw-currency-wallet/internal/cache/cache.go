package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	items map[K]cacheItem[V]
	mu    sync.Mutex
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]cacheItem[V]),
	}
}

func (c *Cache[K, V]) Set(key K, value V, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = newItem(value, duration)
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return item.getValue(), false
	}
	if item.isExpired() {
		delete(c.items, key)
		return item.getValue(), false
	}
	return item.getValue(), true
}

func (c *Cache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache[K, V]) IsItemExpired(key K) bool {
	return c.items[key].isExpired()
}
