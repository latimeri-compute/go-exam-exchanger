package cache

import "time"

type cacheItem[V any] struct {
	value  V
	expiry time.Time
}

func (i cacheItem[V]) isExpired() bool {
	return time.Now().After(i.expiry)
}

func newItem[V any](value V, delta time.Duration) cacheItem[V] {
	return cacheItem[V]{
		value:  value,
		expiry: time.Now().Add(delta),
	}
}

func (c *cacheItem[V]) getValue() V {
	return c.value
}
