package cache

import (
	"sync"
)

type Cache struct {
	mu    *sync.RWMutex
	cache map[string]string
}

func (c Cache) Add(key, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = val
}

func (c Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.cache[key]
	if !ok {
		return "", false
	}

	return val, true
}

func (c Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	clear(c.cache)
}

func NewCache() Cache {
	c := Cache{
		cache: make(map[string]string),
		mu:    &sync.RWMutex{},
	}
	return c
}
