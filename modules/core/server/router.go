package server

import (
	"github.com/gofiber/fiber/v2"
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
	globRouters := []WebRouter{}
	for _, r := range wr {
		if r.prefix == "*" {
			globRouters = append(globRouters, r)
			continue
		}

		router.Route(r.prefix, r.router, r.name...)
	}

	for _, r := range globRouters {
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
	router  func(fiber.Router)
	name    []string
}

func NewApiRouter(version string, router func(fiber.Router), names ...string) ApiRouter {
	return ApiRouter{
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
		router.Route(r.version, r.router, r.name...)
	}
}
