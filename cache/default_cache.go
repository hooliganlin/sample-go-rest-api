package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type DefaultCache struct {
	underlying *cache.Cache
	expirationDuration time.Duration
}

func NewDefaultCache(expirationDuration time.Duration, cleanupInterval time.Duration) DefaultCache {
	c := cache.New(expirationDuration, cleanupInterval)
	return DefaultCache{
		underlying: c,
		expirationDuration: expirationDuration,
	}
}

func (c DefaultCache) Set(key string, value interface{}) {
	c.underlying.Set(key, value, c.expirationDuration)
}

func (c DefaultCache) Get(key string) (interface{}, bool) {
	return c.underlying.Get(key)
}
