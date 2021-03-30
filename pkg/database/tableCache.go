package database

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

type ProximaTableCache struct {
	cache           *cache.Cache
	cacheExpiration time.Duration //might not need a cache expiration
	capacity        int
}

func NewTableCache(cacheExpiration time.Duration) *ProximaTableCache {
	return &ProximaTableCache{cache: cache.New(cacheExpiration, 5*time.Minute), cacheExpiration: cacheExpiration, capacity: -1}
}

func (cache *ProximaTableCache) Get(key interface{}) (*ProximaDBResult, bool) {

	cached, found := cache.cache.Get(key.(string))
	if found && cached != nil {
		cache.cache.SetDefault(key.(string), cached)
		return cached.(*ProximaDBResult), found
	}

	return nil, false
}

//keySlice
//offset-ordering-sliceNumber

func (cache *ProximaTableCache) GetSlice(keySlice interface{}) ([]*ProximaDBResult, bool) {
	//key slice
	//update, ...
	//keySlice, //
	cached, found := cache.cache.Get(keySlice.(string))
	if found && cached != nil {
		cache.cache.SetDefault(keySlice.(string), cached)
		return cached.([]*ProximaDBResult), found
	}
	return make([]*ProximaDBResult, 0), false
}

func (cache *ProximaTableCache) SetSlice(keySlice interface{}, values interface{}) {
	cache.cache.SetDefault(keySlice.(string), values)
}

func (cache *ProximaTableCache) Remove(key interface{}) {
	cache.cache.Delete(key.(string))
	// if cache.cache.ItemCount() >= cache.capacity {
	cache.cache.DeleteExpired()
	// }
}

func (cache *ProximaTableCache) Set(key interface{}, value interface{}) {
	// if cache.cache.ItemCount() < cache.capacity {
	cache.cache.SetDefault(key.(string), value)
	// }
}
