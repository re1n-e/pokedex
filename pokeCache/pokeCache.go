package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCacheEntry(val []byte) cacheEntry {
	return cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

type Cache struct {
	mu       sync.Mutex
	entry    map[string]cacheEntry
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entry:    make(map[string]cacheEntry),
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.entry[key] = NewCacheEntry(val)
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	if key == "" {
		return nil, false
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()

	val, ok := cache.entry[key]
	if !ok {
		return nil, false
	}
	
	return val.val, true
}

func (cache *Cache) reapLoop() {
	ticker := time.NewTicker(cache.interval)
	defer ticker.Stop()

	for range ticker.C {
		cache.mu.Lock()
		for key, entry := range cache.entry {
			if time.Since(entry.createdAt) > cache.interval {
				delete(cache.entry, key)
			}
		}
		cache.mu.Unlock()
	}
}
