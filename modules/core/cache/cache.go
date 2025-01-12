package cache

import (
	"context"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/mrrizkin/pohara/app/config"
)

type Cache struct {
	cache  *ristretto.Cache[string, interface{}]
	config *config.App
}

func NewRessetto(config *config.App) (*Cache, error) {
	ristrettoConfig := ristretto.Config[string, interface{}]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,
	}
	cache, err := ristretto.NewCache(&ristrettoConfig)
	if err != nil {
		return nil, err
	}

	return &Cache{
		cache:  cache,
		config: config,
	}, nil
}

func (c *Cache) Set(context context.Context, key string, value interface{}) error {
	c.cache.SetWithTTL(key, value, 0, c.config.CacheTTLSecond())
	return nil
}

func (c *Cache) Get(context context.Context, key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *Cache) Delete(context context.Context, key string) error {
	c.cache.Del(key)
	return nil
}

func (c *Cache) Close() error {
	c.cache.Close()

	return nil
}
