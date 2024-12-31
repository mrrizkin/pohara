package server

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/internal/common/validator"
	"github.com/mrrizkin/pohara/internal/ports"
	"github.com/mrrizkin/pohara/internal/web/inertia"
	"github.com/mrrizkin/pohara/internal/web/template"
	"go.uber.org/fx"
)

type WebRouter struct {
	prefix string
	router func(fiber.Router)
	name   []string
}

func NewWebRouter(prefix string, router func(fiber.Router), names ...string) WebRouter {
	return WebRouter{
		prefix: prefix,
		router: router,
		name:   names,
	}
}

type WebRouters []WebRouter

func (wr WebRouters) Register(router fiber.Router) {
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
	router  func(fiber.Router)
	name    []string
}

func NewApiRouter(version, prefix string, router func(fiber.Router), names ...string) ApiRouter {
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

func (ar ApiRouters) Register(router fiber.Router) {
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
