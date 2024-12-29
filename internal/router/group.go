package router

import (
	"fmt"
	"net/url"
	"reflect"

	"go.uber.org/fx"
)

type WebRouter struct {
	prefix string
	router func(*Router)
	name   []string
}

type ApiRouter struct {
	version string
	prefix  string
	router  func(*Router)
	name    []string
}

type WebRouters []WebRouter
type ApiRouters []ApiRouter

func AsWebRouter(f any) any {
	validateRouterFunc(f, "WebRouter", "AsWebRouter")
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"web_router"`),
	)
}

func AsApiRouter(f any) any {
	validateRouterFunc(f, "ApiRouter", "AsApiRouter")
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"api_router"`),
	)
}

func NewWebRouter(prefix string, router func(*Router), names ...string) WebRouter {
	return WebRouter{
		prefix: prefix,
		router: router,
		name:   names,
	}
}

func NewApiRouter(version, prefix string, router func(*Router), names ...string) ApiRouter {
	return ApiRouter{
		version: version,
		prefix:  prefix,
		router:  router,
		name:    names,
	}
}

func (wr WebRouters) Register(router *Router) {
	for _, r := range wr {
		router.Route(r.prefix, r.router, r.name...)
	}
}

func (ar ApiRouters) Register(router *Router) {
	for _, r := range ar {
		if prefix, err := url.JoinPath(r.version, r.prefix); err == nil {
			router.Route(prefix, r.router, r.name...)
		}
	}
}

func validateRouterFunc(f any, expectedType, description string) {
	fnType := reflect.TypeOf(f)
	if fnType == nil {
		panic(fmt.Sprintf("%s expects a construction function, got nil", description))
	}

	if fnType.Kind() != reflect.Func {
		panic(fmt.Sprintf("%s expects a construction function, got %s", description, fnType.Kind()))
	}

	if fnType.NumOut() != 1 {
		panic(fmt.Sprintf("%s function must return exactly one value, got %d", description, fnType.NumOut()))
	}

	outType := fnType.Out(0)
	if outType.Kind() == reflect.Ptr {
		outType = outType.Elem()
	}

	if outType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%s function must return a struct, got %s", description, outType.Kind()))
	}

	if outType.Name() != expectedType {
		panic(fmt.Sprintf("%s function must return %s, got %s", description, expectedType, outType.Name()))
	}
}
