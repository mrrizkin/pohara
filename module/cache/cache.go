package cache

import (
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/module/cache/provider"
)

type CacheProvider interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Close() error
}

type Cache struct {
	provider CacheProvider
}

func New(config *config.App) (*Cache, error) {
	cache, err := provider.NewRessetto(config)
	if err != nil {
		return nil, err
	}

	return &Cache{
		provider: cache,
	}, nil
}

func (c *Cache) Has(key string) bool {
	_, ok := c.provider.Get(key)
	return ok
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.provider.Get(key)
}

func (c *Cache) Set(key string, value interface{}) {
	c.provider.Set(key, value)
}
