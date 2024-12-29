package router

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/common/sql"
	"github.com/romsar/gonertia"
)

func (c *Ctx) Render(template string, bind Map) error {
	cacheKey := c.generateCacheKey(template, bind)
	if html, ok := c.tryGetFromCache(cacheKey); ok {
		return c.sendHTML(html)
	}

	html, err := c.router.template.Render(template, bind)
	if err != nil {
		return err
	}

	if c.router.config.IsCacheView() && cacheKey.Valid {
		c.cacheHTML(cacheKey, html)
	}

	return c.sendHTML(html)
}

func (c *Ctx) InertiaRender(template string, props ...gonertia.Props) error {
	return c.router.inertia.Render(c.Ctx, template, props...)
}

func (c *Ctx) ParseBodyAndValidate(out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		return &fiber.Error{Code: 400, Message: "payload not valid"}
	}
	return c.router.validator.MustValidate(out)
}

// Private helper methods
func (c *Ctx) generateCacheKey(template string, data Map) sql.StringNullable {
	if !c.router.config.IsCacheView() {
		return sql.StringNullable{Valid: false}
	}

	var hash [16]byte
	if data != nil {
		if encodedData, err := json.Marshal(data); err == nil {
			hash = md5.Sum(append([]byte(template), encodedData...))
		} else {
			return sql.StringNullable{Valid: false}
		}
	} else {
		hash = md5.Sum([]byte(template))
	}

	return sql.StringNullable{
		Valid:  true,
		String: hex.EncodeToString(hash[:]),
	}
}

func (c *Ctx) tryGetFromCache(cacheKey sql.StringNullable) ([]byte, bool) {
	if !cacheKey.Valid {
		return nil, false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if value, ok := c.router.cache.Get(ctx, cacheKey.String); ok {
		if html, ok := value.([]byte); ok {
			return html, true
		}
		c.router.log.Warn("cached template invalid type")
	}
	return nil, false
}

func (c *Ctx) cacheHTML(cacheKey sql.StringNullable, html []byte) {
	if cacheKey.Valid {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		c.router.cache.Set(ctx, cacheKey.String, html)
	}
}

func (c *Ctx) sendHTML(html []byte) error {
	return c.Type("html").Send(html)
}
