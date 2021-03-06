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
  return &ProximaTableCache{cache: cache.New(cacheExpiration, 1*time.Minute), cacheExpiration: cacheExpiration, capacity: -1}
}

func (cache *ProximaTableCache) Get(key interface{}) (*ProximaDBResult, bool) {

  cached, found := cache.cache.Get(key.(string))
  if found {
    cache.cache.SetDefault(key.(string), cached);
    return cached.(*ProximaDBResult), found
  }

  return nil, found
}

func (cache *ProximaTableCache) Remove(key interface{}) {
  cache.cache.Delete(key.(string))
  // if cache.cache.ItemCount() >= cache.capacity {
  //   cache.cache.DeleteExpired()
  // }
}

func (cache *ProximaTableCache) Set(key interface{}, value interface{}) {
  // if cache.cache.ItemCount() < cache.capacity {
    cache.cache.SetDefault(key.(string), value)
  // }
}
