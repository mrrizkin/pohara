package server

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type WebRouter struct {
	Prefix string
	Router func(fiber.Router)
	Name   string
}

func (w *WebRouter) SetPrefix(prefix string) *WebRouter {
	w.Prefix = prefix
	return w
}

func (w *WebRouter) SetRouter(router func(fiber.Router)) *WebRouter {
	w.Router = router
	return w
}

func (w *WebRouter) SetName(name string) *WebRouter {
	w.Name = name
	return w
}

type WebRouters []WebRouter

func (wr WebRouters) Register(r fiber.Router) {
	for _, router := range wr {
		if router.Name != "" {
			r.Route(router.Prefix, router.Router, router.Name)
		} else {
			r.Route(router.Prefix, router.Router)
		}
	}
}

func AsWebRouter(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"router"`),
	)
}

type ApiRouter struct {
	Version string
	Prefix  string
	Router  func(fiber.Router)
	Name    string
}

func AsApiRouter(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"api_router"`),
	)
}

type ApiRouters []ApiRouter

func (ar ApiRouters) Register(r fiber.Router) {
	for _, router := range ar {
		prefix, err := url.JoinPath(router.Version, router.Prefix)
		if err != nil {
			continue
		}

		if router.Name != "" {
			r.Route(prefix, router.Router, router.Name)
		} else {
			r.Route(prefix, router.Router)
		}
	}

}
