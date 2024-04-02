package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	data map[string]cacheEntry
	mu   *sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	ticker := time.NewTicker(interval)
	cache := Cache{data: make(map[string]cacheEntry), mu: &sync.Mutex{}}
	go func() {
		for range ticker.C {
			cache.reapLoop(interval)
		}
	}()
	return &cache
}

func (self *Cache) Add(key string, val []byte) {
	self.mu.Lock()
	self.data[key] = cacheEntry{createdAt: time.Now(), val: val}
	self.mu.Unlock()
}

func (self *Cache) Get(key string) ([]byte, bool) {
	self.mu.Lock()
	entry, ok := self.data[key]
	self.mu.Unlock()
	return entry.val, ok
}

func (self *Cache) reapLoop(interval time.Duration) {
	self.mu.Lock()
	for key, val := range self.data {
		if time.Now().Sub(val.createdAt) >= interval {
			delete(self.data, key)
		}
	}
	self.mu.Unlock()
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}
