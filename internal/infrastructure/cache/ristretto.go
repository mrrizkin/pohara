package cache

import (
	"context"

	"github.com/dgraph-io/ristretto/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/internal/ports"
)

type Ristretto struct {
	cache  *ristretto.Cache[string, interface{}]
	config *config.App
}

var Module = fx.Module("ristretto",
	fx.Provide(NewRessetto),
	fx.Provide(func(cache *Ristretto) ports.Cache { return cache }),
)

func NewRessetto(config *config.App) (*Ristretto, error) {
	ristrettoConfig := ristretto.Config[string, interface{}]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,
	}
	cache, err := ristretto.NewCache(&ristrettoConfig)
	if err != nil {
		return nil, err
	}

	return &Ristretto{
		cache:  cache,
		config: config,
	}, nil
}

func (c *Ristretto) Set(context context.Context, key string, value interface{}) error {
	c.cache.SetWithTTL(key, value, 0, c.config.CacheTTLSecond())
	return nil
}

func (c *Ristretto) Get(context context.Context, key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *Ristretto) Delete(context context.Context, key string) error {
	c.cache.Del(key)
	return nil
}

func (c *Ristretto) Close() error {
	c.cache.Close()

	return nil
}
