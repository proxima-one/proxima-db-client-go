package database

import (
  cache "github.com/patrickmn/go-cache"
  "time"
)

type ProximaTableCache struct {
  cache *cache.Cache
  cacheExpiration time.Duration //might not need a cache expiration
  capacity int
}

func NewTableCache(cacheExpiration time.Duration) (*ProximaTableCache) {
  return &ProximaTableCache{cache: cache.New(cacheExpiration, 5*time.Minute), cacheExpiration: cacheExpiration, capacity: -1}
}

func (cache *ProximaTableCache) Get(key string) (*ProximaDBResult, bool) {

  cached, found := cache.cache.Get(key)
  if found {
    cache.cache.SetDefault(key, cached);
    return cached.(*ProximaDBResult), found
  }

  return nil, found
}

func (cache *ProximaTableCache) Remove(key string) {
  cache.cache.Delete(key)
  // if cache.cache.ItemCount() >= cache.capacity {
  //   cache.cache.DeleteExpired()
  // }
}

func (cache *ProximaTableCache) Set(key string, value interface{}) {
  // if cache.cache.ItemCount() < cache.capacity {
    cache.cache.SetDefault(key, value)
  // }
}
