package pokecache

import (
	"time"
	"sync"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
    mu       sync.Mutex
    entries  map[string]cacheEntry
    interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
    c := &Cache{
        entries:  make(map[string]cacheEntry),
        interval: interval,
    }
    go c.reapLoop()
    return c
}

func (c *Cache)reapLoop() {
ticker := time.NewTicker(c.interval)
for range ticker.C {
	c.mu.Lock()
	for key, entry := range c.entries {
		age := time.Since(entry.createdAt)
		if age > c.interval {
			delete(c.entries, key)
		}
	}
	c.mu.Unlock()
}
}

func (c *Cache)Add(key string, value []byte) {
	c.mu.Lock()
	c.entries[key] = cacheEntry{createdAt: time.Now(), val: value}
	c.mu.Unlock()
}

func (c *Cache)Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	} else {
		return entry.val, true
	}
}