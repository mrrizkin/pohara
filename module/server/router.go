package server

import (
	"net/url"

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
		fx.ResultTags(`group:"router"`),
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
