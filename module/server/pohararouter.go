package server

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/module/cache"
	"github.com/mrrizkin/pohara/module/database/sql"
	"github.com/mrrizkin/pohara/module/logger"
	"github.com/mrrizkin/pohara/module/template"
	"github.com/mrrizkin/pohara/module/validator"
)

type Router struct {
	core      fiber.Router
	template  *template.Template
	validator *validator.Validator
	config    *config.App
	cache     *cache.Cache
	log       *logger.Logger
}

type Ctx struct {
	*fiber.Ctx

	router *Router
}

func (c *Ctx) Render(template string, bind fiber.Map) error {
	cacheKey := c.cacheKey(template, bind)
	if cacheKey.Valid {
		if value, ok := c.router.cache.Get(cacheKey.String); ok {
			if html, ok := value.([]byte); ok {
				return c.Type("html").Send(html)
			}

			c.router.log.Warn("cached template invalid type", "template", template)
		}
	}

	html, err := c.router.template.Render(template, bind)
	if err != nil {
		return err
	}
	if c.router.config.IsCacheView() && cacheKey.Valid {
		c.router.cache.Set(cacheKey.String, html)
	}

	return c.Type("html").Send(html)
}

func (c *Ctx) ParseBodyAndValidate(out interface{}) error {
	err := c.BodyParser(out)
	if err != nil {
		return &fiber.Error{
			Code:    400,
			Message: "payload not valid",
		}
	}

	err = c.router.validator.MustValidate(out)
	if err != nil {
		return err
	}

	return nil
}

func (c *Ctx) cacheKey(template string, data map[string]interface{}) sql.StringNullable {
	if !c.router.config.IsCacheView() {
		return sql.StringNullable{
			Valid: false,
		}
	}

	var cacheKey string
	var hash [16]byte
	if data != nil {
		encodedData, err := json.Marshal(data)
		if err != nil {
			return sql.StringNullable{
				Valid: false,
			}
		}

		hash = md5.Sum(append([]byte(template), encodedData...))
	} else {
		hash = md5.Sum([]byte(template))
	}

	cacheKey = hex.EncodeToString(hash[:])

	return sql.StringNullable{
		Valid:  true,
		String: cacheKey,
	}
}

type Handler = func(*Ctx) error

type Handlers []Handler

func (h Handlers) Convert(pr *Router) []fiber.Handler {
	fiberHandler := make([]fiber.Handler, len(h))

	for i, handler := range h {
		fiberHandler[i] = func(c *fiber.Ctx) error {
			return handler(&Ctx{c, pr})
		}
	}

	return fiberHandler
}

func NewRouter(
	r fiber.Router,
	template *template.Template,
	config *config.App,
	cache *cache.Cache,
	log *logger.Logger,
	validator *validator.Validator,
) *Router {
	return &Router{
		core:      r,
		template:  template,
		config:    config,
		cache:     cache,
		log:       log,
		validator: validator,
	}
}

func (pr *Router) Use(args ...interface{}) fiber.Router {
	return pr.core.Use(args...)
}

func (pr *Router) Get(path string, handlers ...Handler) fiber.Router {
	return pr.core.Get(path, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Head(path string, handlers ...Handler) fiber.Router {
	return pr.core.Head(path, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Post(path string, handlers ...Handler) fiber.Router {
	return pr.core.Post(path, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Put(path string, handlers ...Handler) fiber.Router {
	return pr.core.Put(path, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Delete(path string, handlers ...Handler) fiber.Router {
	return pr.core.Delete(path, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Connect(path string, handlers ...Handler) fiber.Router {
	return pr.core.Connect(path, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Options(path string, handlers ...Handler) fiber.Router {
	return pr.core.Options(path, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Trace(path string, handlers ...Handler) fiber.Router {
	return pr.core.Trace(path, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Patch(path string, handlers ...Handler) fiber.Router {
	return pr.core.Patch(path, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Add(method, path string, handlers ...Handler) fiber.Router {
	return pr.core.Add(method, path, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Static(prefix, root string, config ...fiber.Static) fiber.Router {
	return pr.core.Static(prefix, root, config...)
}

func (pr *Router) All(path string, handlers ...Handler) fiber.Router {
	return pr.core.All(path, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Group(prefix string, handlers ...Handler) fiber.Router {
	return pr.core.Group(prefix, Handlers(handlers).Convert(pr)...)
}

func (pr *Router) Route(prefix string, fn func(router *Router), name ...string) fiber.Router {
	return pr.core.Route(prefix, func(router fiber.Router) {
		fn(NewRouter(
			router,
			pr.template,
			pr.config,
			pr.cache,
			pr.log,
			pr.validator,
		))
	}, name...)
}

func (pr *Router) Mount(prefix string, fiber *fiber.App) fiber.Router {
	return pr.core.Mount(prefix, fiber)
}

func (pr *Router) Name(name string) fiber.Router {
	return pr.core.Name(name)
}
