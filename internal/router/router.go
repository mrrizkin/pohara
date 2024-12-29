package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/internal/common/validator"
	"github.com/mrrizkin/pohara/internal/ports"
	"github.com/mrrizkin/pohara/internal/web/inertia"
	"github.com/mrrizkin/pohara/internal/web/template"
)

func New(
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

func (r *Router) Inherit(router fiber.Router) *Router {
	clone := *r
	clone.core = router
	return &clone
}

func (r *Router) Use(args ...interface{}) fiber.Router {
	return r.core.Use(args...)
}

func (r *Router) Get(path string, handlers ...Handler) fiber.Router {
	return r.core.Get(path, r.convertHandlers(handlers...)...)
}

func (r *Router) Head(path string, handlers ...Handler) fiber.Router {
	return r.core.Head(path, r.convertHandlers(handlers...)...)
}

func (r *Router) Post(path string, handlers ...Handler) fiber.Router {
	return r.core.Post(path, r.convertHandlers(handlers...)...)
}

func (r *Router) Put(path string, handlers ...Handler) fiber.Router {
	return r.core.Put(path, r.convertHandlers(handlers...)...)
}

func (r *Router) Delete(path string, handlers ...Handler) fiber.Router {
	return r.core.Delete(path, r.convertHandlers(handlers...)...)
}

func (r *Router) Connect(path string, handlers ...Handler) fiber.Router {
	return r.core.Connect(path, r.convertHandlers(handlers...)...)
}

func (r *Router) Options(path string, handlers ...Handler) fiber.Router {
	return r.core.Options(path, r.convertHandlers(handlers...)...)
}

func (r *Router) Trace(path string, handlers ...Handler) fiber.Router {
	return r.core.Trace(path, r.convertHandlers(handlers...)...)
}

func (r *Router) Patch(path string, handlers ...Handler) fiber.Router {
	return r.core.Patch(path, r.convertHandlers(handlers...)...)
}

func (r *Router) Add(method, path string, handlers ...Handler) fiber.Router {
	return r.core.Add(method, path, r.convertHandlers(handlers...)...)
}

func (r *Router) Static(prefix, root string, config ...fiber.Static) fiber.Router {
	return r.core.Static(prefix, root, config...)
}

func (r *Router) All(path string, handlers ...Handler) fiber.Router {
	return r.core.All(path, r.convertHandlers(handlers...)...)
}

func (r *Router) Group(prefix string, handlers ...Handler) fiber.Router {
	return r.core.Group(prefix, r.convertHandlers(handlers...)...)
}

func (r *Router) Mount(prefix string, fiber *fiber.App) fiber.Router {
	return r.core.Mount(prefix, fiber)
}

func (r *Router) Name(name string) fiber.Router {
	return r.core.Name(name)
}

func (r *Router) Route(prefix string, fn func(router *Router), name ...string) fiber.Router {
	return r.core.Route(prefix, func(router fiber.Router) {
		fn(r.Inherit(router))
	}, name...)
}

func (r *Router) convertHandlers(handlers ...Handler) []fiber.Handler {
	fiberHandlers := make([]fiber.Handler, len(handlers))
	for i, handler := range handlers {
		h := handler // Create a new variable to avoid closure problems
		fiberHandlers[i] = func(c *fiber.Ctx) error {
			return h(&Ctx{Ctx: c, router: r})
		}
	}
	return fiberHandlers
}
