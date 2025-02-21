package cache

import "time"

type Config struct {
	TTL int `env:"CACHE_TTL,default=300"`
}

func (c *Config) CacheTTLSecond() time.Duration {
	return time.Duration(c.TTL) * time.Second
}
