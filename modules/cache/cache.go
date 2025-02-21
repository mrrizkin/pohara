package cache

import (
	"context"

	"github.com/dgraph-io/ristretto/v2"
	"go.uber.org/fx"
)

type Cache struct {
	core   *ristretto.Cache[string, any]
	config *Config
}

type CacheDeps struct {
	fx.In

	Config *Config
}

func NewCache(lc fx.Lifecycle, deps CacheDeps) (*Cache, error) {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}

	cache := Cache{
		core:   ristrettoCache,
		config: deps.Config,
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			cache.Close()
			return nil
		},
	})

	return &cache, nil
}

func (c *Cache) Set(context context.Context, key string, value any) error {
	c.core.SetWithTTL(key, value, 0, c.config.CacheTTLSecond())
	return nil
}

func (c *Cache) Get(context context.Context, key string) (any, bool) {
	return c.core.Get(key)
}

func (c *Cache) Delete(context context.Context, key string) error {
	c.core.Del(key)
	return nil
}

func (c *Cache) Close() error {
	c.core.Close()

	return nil
}
