package cache

import (
	"github.com/bluele/gcache"
	"log"
	"time"
)

type Cache struct {
	cache gcache.Cache
}

func NewCache() Cache {
	gc := Cache{
		cache: gcache.New(10).LRU().Expiration(30 * time.Minute).
			EvictedFunc(func(key interface{}, value interface{}) {
				log.Printf("Evicted key: %s, val: %s", key, value)
			}).
			PurgeVisitorFunc(func(key interface{}, value interface{}) {
				log.Printf("Purged key: %s, val: %s", key, value)
			}).
			AddedFunc(func(key interface{}, value interface{}) {
				log.Printf("Added key: %s, val: %s", key, value)
			}).
			Build(),
	}
	return gc
}

func (receiver Cache) Set(key string, value interface{}) error {
	return receiver.cache.Set(key, value)
}

func (receiver Cache) Get(key string) (interface{}, error) {
	return receiver.cache.Get(key)
}
