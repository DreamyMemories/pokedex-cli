package pokecache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       interface{}
}

type Cache struct {
	Data map[string]cacheEntry
	sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		Data: make(map[string]cacheEntry),
	}
	go newCache.reapLoop(interval)
	log.Printf("Created cache, deleting items after %v", interval)
	return newCache
}

func (cache *Cache) Add(key string, val interface{}) {
	cache.Lock()
	defer cache.Unlock() // Wait function to finish only unlock

	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	cache.Data[key] = entry
}

func (cache *Cache) Get(key string) (interface{}, bool) {
	cache.Lock()
	defer cache.Unlock()

	if data, exist := cache.Data[key]; exist {
		return data.val, exist
	}

	return nil, false
}

func (cache *Cache) reapLoop(t time.Duration) {
	ticker := time.NewTicker(t)

	go func() {
		for range ticker.C {
			cache.Lock()
			for key, item := range cache.Data {
				elapsed := time.Since(item.createdAt)

				if elapsed > t {
					delete(cache.Data, key)
				}
				log.Println("Cache cleared")
				fmt.Print("Pokedex > ")

			}
			cache.Unlock()
		}
	}()
}
