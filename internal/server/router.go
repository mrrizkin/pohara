package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/internal/common/sql"
	"github.com/mrrizkin/pohara/internal/common/validator"
	"github.com/mrrizkin/pohara/internal/ports"
	"github.com/mrrizkin/pohara/internal/web/inertia"
	"github.com/mrrizkin/pohara/internal/web/template"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"
)

type WebRouter struct {
	prefix string
	router func(*Router)
	name   []string
}

func NewWebRouter(prefix string, router func(*Router), names ...string) WebRouter {
	return WebRouter{
		prefix: prefix,
		router: router,
		name:   names,
	}
}

type WebRouters []WebRouter

func (wr WebRouters) Register(router *Router) {
	for _, r := range wr {
		router.Route(r.prefix, r.router, r.name...)
	}
}

func AsWebRouter(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"web_router"`),
	)
}

type ApiRouter struct {
	version string
	prefix  string
	router  func(*Router)
	name    []string
}

func NewApiRouter(version, prefix string, router func(*Router), names ...string) ApiRouter {
	return ApiRouter{
		prefix:  prefix,
		router:  router,
		name:    names,
		version: version,
	}
}

func AsApiRouter(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"api_router"`),
	)
}

type ApiRouters []ApiRouter

func (ar ApiRouters) Register(router *Router) {
	for _, r := range ar {
		prefix, err := url.JoinPath(r.version, r.prefix)
		if err != nil {
			continue
		}

		router.Route(prefix, r.router, r.name...)
	}
}

type Router struct {
	core      fiber.Router
	template  *template.Template
	inertia   *inertia.Inertia
	validator *validator.Validator
	config    *config.App
	cache     ports.Cache
	log       ports.Logger
}

type Ctx struct {
	*fiber.Ctx

	router *Router
}

func (c *Ctx) Render(template string, bind Map) error {
	cacheKey := c.cacheKey(template, bind)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if cacheKey.Valid {
		if value, ok := c.router.cache.Get(ctx, cacheKey.String); ok {
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
		c.router.cache.Set(ctx, cacheKey.String, html)
	}

	return c.Type("html").Send(html)
}

func (c *Ctx) InertiaRender(template string, props ...gonertia.Props) error {
	return c.router.inertia.Render(c.Ctx, template, props...)
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

func (c *Ctx) cacheKey(template string, data Map) sql.StringNullable {
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
	inertia *inertia.Inertia,
	config *config.App,
	cache ports.Cache,
	log ports.Logger,
	validator *validator.Validator,
) *Router {
	return &Router{
		core:      r,
		template:  template,
		config:    config,
		cache:     cache,
		log:       log,
		inertia:   inertia,
		validator: validator,
	}
}

func (pr *Router) Inherit(router fiber.Router) *Router {
	r := *pr
	r.core = router
	return &r
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
		fn(pr.Inherit(router))
	}, name...)
}

func (pr *Router) Mount(prefix string, fiber *fiber.App) fiber.Router {
	return pr.core.Mount(prefix, fiber)
}

func (pr *Router) Name(name string) fiber.Router {
	return pr.core.Name(name)
}
